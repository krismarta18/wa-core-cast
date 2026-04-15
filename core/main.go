package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wacast/core/appserver"
	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/services/auth"
	"wacast/core/services/billing"
	"wacast/core/services/contact"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/utils"

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
		utils.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	database.DB = db

	utils.Info("Database connection pool configured",
		zap.Int("max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConns),
	)

	// 4. Run migrations
	utils.Info("Running database migrations...")

	migrationRunner := database.NewMigrationRunner(db)
	
	// Load migrations from migrations directory
	migrationsPath := "./migrations"
	err = migrationRunner.LoadMigrationsFromDirectory(migrationsPath)
	if err != nil {
		utils.Error("Failed to load migrations", zap.Error(err))
		// Continue anyway, migrations might not be found during development
	}

	// Run pending migrations
	err = migrationRunner.RunMigrations()
	if err != nil {
		utils.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Print migration status
	err = migrationRunner.PrintMigrationStatus()
	if err != nil {
		utils.Warn("Failed to print migration status", zap.Error(err))
	}

	// 5. Initialize services
	
	// Initialize session service
	utils.Info("Initializing WhatsApp session service...")
	sessionService := session.NewService(
		db,
		cfg.EncryptionKey,
		25, // max sessions
		cfg.SessionTimeout,
	)

	// Start background manager for session cleanup and auto-reconnect
	sessionManager := session.NewManager(sessionService, true, 30*time.Second)
	sessionManager.Start()

	// Attempt to restore previous sessions from database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := sessionService.RestorePreviousSessions(ctx); err != nil {
		utils.Warn("Failed to restore previous sessions",
			zap.Error(err),
		)
	}
	cancel()

	utils.Info("Session service initialized successfully")

	// Initialize message service
	utils.Info("Initializing message service...")
	messageService := message.NewService(
		db,
		sessionService,
		message.DefaultQueueConfig(),
	)

	if err := messageService.Start(); err != nil {
		utils.Fatal("Failed to start message service", zap.Error(err))
	}

	// Register callbacks to handle incoming messages
	messageService.RegisterReceiveCallback(func(rm *message.ReceivedMessage) {
		utils.Debug("Incoming message received",
			zap.String("from_jid", rm.FromJID),
			zap.String("content", rm.Content),
		)
		// TODO: Route to webhook service
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

	// Initialize auth service
	utils.Info("Initializing auth service...")
	authService := auth.NewService(db, cfg.JWTSecret, cfg.JWTExpiryHours, cfg.JWTRefreshExpiryHours)
	utils.Info("Auth service initialized successfully")

	billingService := billing.NewService(db)
	utils.Info("Billing service initialized successfully")

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
