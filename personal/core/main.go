package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	"wacast/core/services/integration"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	err = utils.InitLogger(cfg.LogLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	utils.Info("Starting WACAST Core",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
		zap.String("server", cfg.GetServerAddr()),
	)

	// 3. Connect to database
	db, err := database.InitDatabase(cfg.Database)
	if err != nil {
		utils.Error("Could not connect to database on startup. Please configure via dashboard.", zap.Error(err))
		// We continue anyway so the user can reach the /config API
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	database.DB = db

	utils.Info("Database connection pool initialized",
		zap.Int("max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConns),
	)

	// 4. Run migrations (Only if connected)
	if db != nil && db.HealthCheck() {
		utils.Info("Running database migrations...")
		migrationRunner := database.NewMigrationRunner(db)
		
		// Load migrations from migrations directory
		migrationsPath := "./migrations"
		err = migrationRunner.LoadMigrationsFromDirectory(migrationsPath)
		if err != nil {
			utils.Error("Failed to load migrations", zap.Error(err))
		}

		// Run pending migrations
		err = migrationRunner.RunMigrations()
		if err != nil {
			utils.Error("Failed to run migrations", zap.Error(err))
		}
	} else {
		utils.Warn("Skipping database migrations - no active connection")
	}

	// Initialize services
	
	// Create billing service early so it can be injected into others
	utils.Info("Initializing billing and integration services...")
	billingService := billing.NewService(db)
	analyticService := analytics.NewService(analytics.NewStore(db))
	integrationService := integration.NewService(db)
	utils.Info("Infrastucture services initialized successfully")

	// Initialize session service
	utils.Info("Initializing WhatsApp session service...")
	sessionService := session.NewService(
		db,
		billingService,
		cfg.EncryptionKey,
		25, // max sessions
		cfg.SessionTimeout,
	)

	// Start background manager for session cleanup and auto-reconnect
	sessionManager := session.NewManager(sessionService, true, 30*time.Second)
	sessionManager.Start()

	// Attempt to restore previous sessions from database (Only if connected)
	if db != nil && db.HealthCheck() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := sessionService.RestorePreviousSessions(ctx); err != nil {
			utils.Warn("Failed to restore previous sessions",
				zap.Error(err),
			)
		}
		cancel()
	}

	utils.Info("Session service initialized successfully")

	// Create a dummy context generator for background jobs
	paramContext := func() context.Context { return context.Background() }

	utils.Info("Informing system modules...")

	utils.Info("Initializing autoresponse service...")
	autoresponseStore := autoresponse.NewStore(db)
	autoresponseService := autoresponse.NewService(autoresponseStore)
	utils.Info("Autoresponse service initialized successfully")

	// Initialize message service with Pro configuration
	utils.Info("Initializing message service...")
	msgConfig := message.DefaultQueueConfig()
	
	// Load Pro Settings from DB if available
	if db != nil && db.HealthCheck() {
		if s, err := db.GetSystemSetting("anti_bot_enabled"); err == nil && s != nil {
			msgConfig.AntiBotEnabled = (strings.ToLower(s.Value) == "true")
		}
		if s, err := db.GetSystemSetting("anti_bot_suffix_length"); err == nil && s != nil {
			fmt.Sscanf(s.Value, "%d", &msgConfig.RandomSuffixLength)
		}
		utils.Info("Pro configuration applied from database", 
			zap.Bool("anti_bot_enabled", msgConfig.AntiBotEnabled),
			zap.Int("suffix_length", msgConfig.RandomSuffixLength),
		)
	}

	messageService := message.NewService(
		db,
		sessionService,
		analyticService,
		billingService,
		integrationService,
		msgConfig,
	)

	if err := messageService.Start(); err != nil {
		utils.Fatal("Failed to start message service", zap.Error(err))
	}

	// Register callbacks to handle incoming messages
	messageService.RegisterReceiveCallback(func(rm *message.ReceivedMessage) {
		utils.Debug("Incoming message received",
			zap.String("from_jid", rm.FromJID),
			zap.String("content", rm.Content),
			zap.String("device_id", rm.DeviceID),
		)
		
		// 1. Get UserID from message (or DeviceID if not populated)
		userID := rm.UserID
		parsedDeviceID, _ := uuid.Parse(rm.DeviceID)

		if userID == uuid.Nil {
			utils.Debug("UserID missing in message, looking up device", zap.String("device_id", rm.DeviceID))
			device, err := db.GetDeviceByID(parsedDeviceID)
			if err == nil && device != nil {
				userID = device.UserID
			}
		}

		if userID == uuid.Nil {
			utils.Warn("Could not identify user for incoming message", zap.String("device_id", rm.DeviceID))
			return
		}
		
		// 2. Check AutoResponse
		utils.Debug("Checking auto-response keywords", 
			zap.String("user_id", userID.String()), 
			zap.String("device_id", rm.DeviceID),
			zap.String("content", rm.Content),
		)
		
		replyText := autoresponseService.DetectAndReply(userID, &parsedDeviceID, rm.Content)
		if replyText != "" {
			utils.Info("Auto-reply triggered", 
				zap.String("user_id", userID.String()),
				zap.String("device_id", rm.DeviceID), 
				zap.String("to", rm.FromJID),
				zap.String("reply", replyText),
			)
			// Send reply using messageService
			target := rm.FromJID
			if rm.GroupJID != nil && *rm.GroupJID != "" {
				target = *rm.GroupJID
			}

			msgID, err := messageService.SendMessage(paramContext(), rm.DeviceID, target, replyText, nil, nil)
			if err != nil {
				utils.Error("Failed to send auto-reply", zap.Error(err))
			} else {
				utils.Debug("Auto-reply queued", zap.String("msg_id", msgID))
			}
		} else {
			utils.Debug("No auto-response keywords matched")
		}
	})


	// Register delivery callbacks for status updates
	messageService.RegisterDeliveryCallback(func(msu *message.MessageStatusUpdate) {
		utils.Debug("Message status updated",
			zap.String("message_id", msu.MessageID),
			zap.Int("new_status", int(msu.NewStatus)),
		)
		// TODO: Notify API layer
	})

	utils.Info("Message service initialized successfully")
	
	// Initialize broadcast service
	utils.Info("Initializing broadcast service...")
	broadcastStore := broadcast.NewStore(db)
	broadcastService := broadcast.NewService(broadcastStore, messageService, billingService)
	utils.Info("Broadcast service initialized successfully")

	// Initialize auth service
	utils.Info("Initializing auth service...")
	authService := auth.NewService(db, cfg.JWTSecret, cfg.JWTExpiryHours, cfg.JWTRefreshExpiryHours)
	utils.Info("Auth service initialized successfully")

	utils.Info("Initializing contact service...")
	contactStore := contact.NewStore(db)
	contactService := contact.NewService(contactStore)
	utils.Info("Contact service initialized successfully")

	// 6. Initialize and start HTTP server
	utils.Info("Initializing HTTP server...")
	httpServer := appserver.NewServer(
		authService,
		billingService,
		sessionService,
		messageService,
		contactService,
		analyticService,
		broadcastService,
		autoresponseService,
		db,
		cfg,
		cfg.ServerHost,
		cfg.ServerPort,
	)

	// ✅ Register WebSocket callback for real-time QR updates
	sessionService.RegisterQRUpdateCallback(httpServer.QRUpdateNotifier())

	// Start server in a goroutine so it doesn't block
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- httpServer.Start()
	}()

	utils.Info("WACAST Core started successfully",
		zap.String("server_address", cfg.GetServerAddr()),
	)

	// Open browser automatically for Personal version
	go func() {
		time.Sleep(2 * time.Second) // Give server time to bind
		url := fmt.Sprintf("http://localhost:%d", cfg.ServerPort)
		if cfg.ServerHost != "0.0.0.0" && cfg.ServerHost != "" {
			url = fmt.Sprintf("http://%s:%d", cfg.ServerHost, cfg.ServerPort)
		}
		utils.Info("Opening browser", zap.String("url", url))
		if err := utils.OpenBrowser(url); err != nil {
			utils.Warn("Failed to open browser", zap.Error(err))
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		utils.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case err := <-serverErrors:
		if err != nil {
			utils.Fatal("Server error", zap.Error(err))
		}
	}

	utils.Info("Shutting down gracefully...")

	// Create a context with timeout for shutdown operations
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Cleanup resources
	sessionManager.Stop()
	if err := sessionService.Shutdown(shutdownCtx); err != nil {
		utils.Error("Error shutting down session service", zap.Error(err))
	}

	if err := messageService.Stop(); err != nil {
		utils.Error("Error stopping message service", zap.Error(err))
	}

	if err := messageService.Cleanup(); err != nil {
		utils.Error("Error cleaning up message service", zap.Error(err))
	}

	if err := httpServer.Shutdown(); err != nil {
		utils.Error("Error shutting down HTTP server", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		utils.Error("Error closing database", zap.Error(err))
	}

	utils.Info("WACAST Core stopped")
}
