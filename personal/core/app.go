package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"wacast/core/appserver"
	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/services/analytics"
	"wacast/core/services/auth"
	"wacast/core/services/autoresponse"
	"wacast/core/services/billing"
	"wacast/core/services/broadcast"
	"wacast/core/services/contact"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/utils"
	"wacast/core/services/integration"
	"wacast/core/services/warming"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type App struct {
	cfg            *config.Config
	db             *database.Database
	sessionService *session.Service
	messageService *message.Service
	httpServer     *appserver.Server
	sessionManager *session.Manager
	
	ctx            context.Context
	cancel         context.CancelFunc
	isRunning      bool
}

func NewApp() *App {
	return &App{
		isRunning: false,
	}
}

func (a *App) Start() error {
	if a.isRunning {
		return fmt.Errorf("app is already running")
	}

	// 0. Ensure required files/folders exist
	if _, err := os.Stat("admin.password"); os.IsNotExist(err) {
		_ = os.WriteFile("admin.password", []byte("admin123"), 0600)
		utils.Info("Created default admin.password file with 'admin123'")
	}
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		_ = os.MkdirAll("uploads", 0755)
	}

	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	a.cfg = cfg

	// 2. Initialize logger
	_ = utils.InitLogger(cfg.LogLevel)

	utils.Info("Starting WACAST Core via Control Panel",
		zap.String("version", cfg.ServiceVersion),
	)

	// 3. Connect to database
	db, err := database.InitDatabase(cfg.Database)
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}
	a.db = db
	database.DB = db

	// 4. Run database migrations
	mr := database.NewMigrationRunner(a.db)
	_ = mr.LoadMigrationsFromEmbedded()
	if err := mr.RunMigrations(); err != nil {
		utils.Error("DATABASE MIGRATION FAILED", zap.Error(err))
		return fmt.Errorf("database migration failed: %v", err)
	}

	// 5. Initialize services
	billingService := billing.NewService(db)
	analyticService := analytics.NewService(analytics.NewStore(db))
	integrationService := integration.NewService(db)

	a.sessionService = session.NewService(
		db,
		billingService,
		cfg.EncryptionKey,
		25,
		cfg.SessionTimeout,
	)

	a.sessionManager = session.NewManager(a.sessionService, true, 30*time.Second)
	a.sessionManager.Start()

	// Restore sessions
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	_ = a.sessionService.RestorePreviousSessions(ctx)
	cancel()

	autoresponseStore := autoresponse.NewStore(db)
	autoresponseService := autoresponse.NewService(autoresponseStore)

	msgConfig := message.DefaultQueueConfig()
	if s, err := db.GetSystemSetting("anti_bot_enabled"); err == nil && s != nil {
		msgConfig.AntiBotEnabled = (strings.ToLower(s.Value) == "true")
	}

	a.messageService = message.NewService(
		db,
		a.sessionService,
		analyticService,
		billingService,
		integrationService,
		msgConfig,
	)

	if err := a.messageService.Start(); err != nil {
		return err
	}

	// Setup auto-reply callback
	a.setupCallbacks(autoresponseService)

	broadcastStore := broadcast.NewStore(db)
	broadcastService := broadcast.NewService(broadcastStore, a.messageService, billingService)
	authService := auth.NewService(db, cfg.JWTSecret, cfg.JWTExpiryHours, cfg.JWTRefreshExpiryHours)
	contactStore := contact.NewStore(db)
	contactService := contact.NewService(contactStore)
	warmingService := warming.NewService(db, a.sessionService, a.messageService)

	// 6. Initialize HTTP server
	a.httpServer = appserver.NewServer(
		authService,
		billingService,
		a.sessionService,
		a.messageService,
		contactService,
		analyticService,
		broadcastService,
		autoresponseService,
		warmingService,
		db,
		cfg,
		cfg.ServerHost,
		cfg.ServerPort,
	)

	a.sessionService.RegisterQRUpdateCallback(a.httpServer.QRUpdateNotifier())

	// Start server in a goroutine
	startErr := make(chan error, 1)
	go func() {
		if err := a.httpServer.Start(); err != nil {
			startErr <- err
		}
	}()

	// Wait for a bit to see if it fails early (e.g. port already in use)
	select {
	case err := <-startErr:
		utils.Error("Failed to start HTTP server", zap.Error(err))
		return err
	case <-time.After(12 * time.Second):
		// If no error after 12 seconds, assume it started successfully
		a.isRunning = true
		utils.Info("WACAST Core started successfully", zap.Int("port", cfg.ServerPort))
		return nil
	}
}

func (a *App) Stop() {
	if !a.isRunning {
		return
	}

	utils.Info("Stopping WACAST Core...")
	
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if a.sessionManager != nil {
		a.sessionManager.Stop()
	}
	if a.sessionService != nil {
		_ = a.sessionService.Shutdown(shutdownCtx)
	}
	if a.messageService != nil {
		_ = a.messageService.Stop()
		_ = a.messageService.Cleanup()
	}
	if a.httpServer != nil {
		_ = a.httpServer.Shutdown()
	}
	if a.db != nil {
		_ = a.db.Close()
	}

	a.isRunning = false
	utils.Info("WACAST Core stopped")
}

func (a *App) setupCallbacks(autoresponseService *autoresponse.Service) {
	a.messageService.RegisterReceiveCallback(func(rm *message.ReceivedMessage) {
		userID := rm.UserID
		parsedDeviceID, _ := uuid.Parse(rm.DeviceID)

		if userID == uuid.Nil {
			device, err := a.db.GetDeviceByID(parsedDeviceID)
			if err == nil && device != nil {
				userID = device.UserID
			}
		}

		if userID != uuid.Nil {
			replyText := autoresponseService.DetectAndReply(userID, &parsedDeviceID, rm.Content)
			if replyText != "" {
				target := rm.FromJID
				if rm.GroupJID != nil && *rm.GroupJID != "" {
					target = *rm.GroupJID
				}
				_, _ = a.messageService.SendMessage(context.Background(), rm.DeviceID, target, replyText, nil, nil)
			}
		}
	})
}
