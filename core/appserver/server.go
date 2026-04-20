package appserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/handlers"
	"wacast/core/services/analytics"
	"wacast/core/services/auth"
	"wacast/core/services/autoresponse"
	"wacast/core/services/billing"
	"wacast/core/services/broadcast"
	"wacast/core/services/contact"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/services/license"
	"wacast/core/utils"
)

type Server struct {
	engine           *gin.Engine
	authService      *auth.Service
	billingService   *billing.Service
	sessionService   *session.Service
	messageService   *message.Service
	contactService   *contact.Service
	analyticService  *analytics.Service
	broadcastService *broadcast.Service
	autoresponseService *autoresponse.Service
	db               *database.Database
	config           *config.Config
	licenseService   *license.Service
	port             int
	host             string
	startTime        time.Time
	websocketHandler *handlers.WebSocketHandler
}

func NewServer(
	authService *auth.Service,
	billingService *billing.Service,
	sessionService *session.Service,
	messageService *message.Service,
	contactService *contact.Service,
	analyticService *analytics.Service,
	broadcastService *broadcast.Service,
	autoresponseService *autoresponse.Service,
	db *database.Database,
	cfg *config.Config,
	host string,
	port int,
) *Server {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(corsMiddleware())
	engine.Use(ginLogger())
	engine.Use(gin.Recovery())

	server := &Server{
		engine:           engine,
		authService:      authService,
		billingService:   billingService,
		sessionService:   sessionService,
		messageService:   messageService,
		contactService:   contactService,
		analyticService:  analyticService,
		broadcastService: broadcastService,
		autoresponseService: autoresponseService,
		db:               db,
		config:           cfg,
		port:             port,
		host:             host,
		startTime:        time.Now(),
		licenseService:   license.NewService(),
		websocketHandler: handlers.NewWebSocketHandler(sessionService),
	}

	server.registerRoutes()

	return server
}

func (s *Server) QRUpdateNotifier() func(deviceID, qrCode string, status int) {
	return s.websocketHandler.NotifyQRUpdate
}

func (s *Server) registerRoutes() {
	s.engine.Use(s.licenseMiddleware())

	s.engine.GET("/health", s.HealthCheck)
	s.engine.GET("/health/ready", s.ReadinessCheck)
	s.engine.GET("/health/live", s.LivenessCheck)

	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "WACAST Core",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/docs",
		})
	})

	s.engine.GET("/api/docs", s.ServeSwaggerUI)
	s.engine.GET("/openapi.yaml", s.ServeOpenAPISpec)

	// Serve uploaded files
	s.engine.Static("/uploads", "./uploads")

	v1 := s.engine.Group("/api/v1")
	{
		handlers.RegisterAuthRoutes(v1, s.authService, s.config.JWTSecret)
		handlers.RegisterBillingRoutes(v1, s.billingService, s.config.JWTSecret, s.authService)
		handlers.RegisterSessionRoutes(v1, s.sessionService, s.config.EncryptionKey, s.config.SessionTimeout, s.config.JWTSecret, s.authService)
		handlers.RegisterMessageRoutes(v1, s.messageService, s.config.JWTSecret, s.authService)
		handlers.RegisterContactRoutes(v1, s.contactService, s.config.JWTSecret, s.authService)
		handlers.RegisterAnalyticRoutes(v1, s.analyticService, s.config.JWTSecret, s.authService)
		handlers.RegisterBroadcastRoutes(v1, s.broadcastService, s.config.JWTSecret, s.authService)
		handlers.RegisterAutoResponseRoutes(v1, s.autoresponseService, s.config.JWTSecret, s.authService)

		// Register Database Config Routes
		configHandler := handlers.NewConfigHandler(s.db, s.config)
		configHandler.RegisterRoutes(v1)

		// Register License Routes
		licenseHandler := handlers.NewLicenseHandler(s.licenseService)
		licenseHandler.RegisterRoutes(v1)

		info := v1.Group("/info")
		{
			info.GET("/status", s.ServerStatus)
			info.GET("/stats", s.ServerStats)
		}
	}

	s.engine.GET("/ws/sessions/:device_id/qr", s.websocketHandler.ConnectQR)
	s.engine.GET("/qr/:device_id", s.websocketHandler.ServeQRPage)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	utils.Info("Starting HTTP server",
		zap.String("address", addr),
		zap.String("environment", s.config.Environment),
	)

	return s.engine.Run(addr)
}

func (s *Server) licenseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Always allow license, config, health, and static files
		if strings.Contains(path, "/api/v1/license") || 
		   strings.Contains(path, "/api/v1/config") ||
		   strings.Contains(path, "/health") ||
		   strings.Contains(path, "/uploads") ||
		   path == "/" {
			c.Next()
			return
		}

		status, err := s.licenseService.GetStatus()
		if err != nil || !status.IsActive || status.IsExpired {
			reason := "License required"
			if status != nil && status.IsExpired {
				reason = "License expired. Please contact admin for renewal."
			}

			c.JSON(http.StatusPaymentRequired, gin.H{
				"success": false,
				"message": reason,
				"setup_required": true,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (s *Server) Shutdown() error {
	utils.Info("Shutting down HTTP server")
	return nil
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Uptime    int64                  `json:"uptime_seconds"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func (s *Server) HealthCheck(c *gin.Context) {
	details := map[string]interface{}{
		"database": "healthy",
		"sessions": s.sessionService != nil,
		"messages": s.messageService != nil,
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Unix(),
		Uptime:    int64(time.Since(s.startTime).Seconds()),
		Details:   details,
	})
}

func (s *Server) ReadinessCheck(c *gin.Context) {
	ready := s.db != nil && s.sessionService != nil && s.messageService != nil

	status := http.StatusOK
	statusStr := "READY"

	if !ready {
		status = http.StatusServiceUnavailable
		statusStr = "NOT_READY"
	}

	c.JSON(status, HealthResponse{
		Status:    statusStr,
		Timestamp: time.Now().Unix(),
		Uptime:    int64(time.Since(s.startTime).Seconds()),
	})
}

func (s *Server) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ALIVE",
		Timestamp: time.Now().Unix(),
		Uptime:    int64(time.Since(s.startTime).Seconds()),
	})
}

func (s *Server) ServerStatus(c *gin.Context) {
	activeSessions := len(s.sessionService.GetAllActiveSessions())

	c.JSON(http.StatusOK, gin.H{
		"status":          "running",
		"uptime_seconds":  int64(time.Since(s.startTime).Seconds()),
		"active_sessions": activeSessions,
		"server_address":  fmt.Sprintf("%s:%d", s.host, s.port),
		"environment":     s.config.Environment,
		"timestamp":       time.Now().Unix(),
	})
}

func (s *Server) ServerStats(c *gin.Context) {
	msgStats := s.messageService.GetQueueStats()
	activeSessions := len(s.sessionService.GetAllActiveSessions())

	c.JSON(http.StatusOK, gin.H{
		"sessions": gin.H{
			"active": activeSessions,
			"max":    25,
		},
		"messages":       msgStats,
		"uptime_seconds": int64(time.Since(s.startTime).Seconds()),
		"timestamp":      time.Now().Unix(),
	})
}

func (s *Server) ServeSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>WACAST Core API - Swagger UI</title>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui.css">
	<style>
		html {
			box-sizing: border-box;
			overflow: -moz-scrollbars-vertical;
			overflow-y: scroll;
		}
		*, *:before, *:after {
			box-sizing: inherit;
		}
		body {
			margin: 0;
			padding: 0;
		}
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
	<script>
		window.onload = function() {
			SwaggerUIBundle({
				url: "/openapi.yaml",
				dom_id: '#swagger-ui',
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIBundle.SwaggerUIStandalonePreset
				],
				layout: "BaseLayout",
				deepLinking: true
			})
		}
	</script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

func (s *Server) ServeOpenAPISpec(c *gin.Context) {
	openapi, err := os.ReadFile("openapi.yaml")
	if err != nil {
		utils.Error("Failed to read OpenAPI spec", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load OpenAPI specification",
		})
		return
	}

	c.Header("Content-Type", "application/x-yaml")
	c.String(http.StatusOK, string(openapi))
}

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		if path == "/health" || path == "/health/ready" || path == "/health/live" {
			return
		}

		switch {
		case statusCode >= 400 && statusCode < 500:
			utils.Warn("HTTP request",
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			)
		case statusCode >= 500:
			fields := []zap.Field{
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			}
			if errorMsg, exists := c.Get("error_message"); exists {
				fields = append(fields, zap.String("error", errorMsg.(string)))
			}
			utils.Error("HTTP request", fields...)
		default:
			utils.Debug("HTTP request",
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			)
		}
	}
}