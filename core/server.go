package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/handlers"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/utils"
)

// Server wraps the HTTP server and its dependencies
type Server struct {
	engine           *gin.Engine
	sessionService   *session.Service
	messageService   *message.Service
	db               *database.Database
	config           *config.Config
	port             int
	host             string
	startTime        time.Time
	websocketHandler *handlers.WebSocketHandler
}

// NewServer creates a new HTTP server
func NewServer(
	sessionService *session.Service,
	messageService *message.Service,
	db *database.Database,
	cfg *config.Config,
	host string,
	port int,
) *Server {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	// Add middleware
	// Set a custom logger that uses Zap
	engine.Use(GinLogger())
	engine.Use(gin.Recovery())

	server := &Server{
		engine:           engine,
		sessionService:   sessionService,
		messageService:   messageService,
		db:               db,
		config:           cfg,
		port:             port,
		host:             host,
		startTime:        time.Now(),
		websocketHandler: handlers.NewWebSocketHandler(sessionService),
	}

	// Register routes
	server.registerRoutes()

	return server
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes() {
	// Health check endpoints
	s.engine.GET("/health", s.HealthCheck)
	s.engine.GET("/health/ready", s.ReadinessCheck)
	s.engine.GET("/health/live", s.LivenessCheck)

	// Root endpoint
	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "WACAST Core",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/docs",
		})
	})

	// Swagger UI - Serve OpenAPI documentation
	s.engine.GET("/api/docs", s.ServeSwaggerUI)
	s.engine.GET("/openapi.yaml", s.ServeOpenAPISpec)

	// API versioning group
	v1 := s.engine.Group("/api/v1")
	{
		// Register session routes (from handlers package)
		handlers.RegisterSessionRoutes(v1, s.sessionService)

		// Register message routes (from handlers package)
		handlers.RegisterMessageRoutes(v1, s.messageService)

		// Server info endpoints
		info := v1.Group("/info")
		{
			info.GET("/status", s.ServerStatus)
			info.GET("/stats", s.ServerStats)
		}
	}

	// WebSocket routes for real-time updates
	s.engine.GET("/ws/sessions/:device_id/qr", s.websocketHandler.ConnectQR)

	// QR code display page
	s.engine.GET("/qr/:device_id", s.websocketHandler.ServeQRPage)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	utils.Info("Starting HTTP server",
		zap.String("address", addr),
		zap.String("environment", s.config.Environment),
	)

	return s.engine.Run(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	utils.Info("Shutting down HTTP server")
	return nil
}

// ============================================================================
// Health Check Handlers
// ============================================================================

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Uptime    int64                  `json:"uptime_seconds"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthCheck returns overall health status
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

// ReadinessCheck returns readiness status (all dependencies ready)
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

// LivenessCheck returns liveness status (server is running)
func (s *Server) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ALIVE",
		Timestamp: time.Now().Unix(),
		Uptime:    int64(time.Since(s.startTime).Seconds()),
	})
}

// ============================================================================
// Server Info Handlers
// ============================================================================

// ServerStatus returns server status
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

// ServerStats returns server statistics
func (s *Server) ServerStats(c *gin.Context) {
	msgStats := s.messageService.GetQueueStats()
	activeSessions := len(s.sessionService.GetAllActiveSessions())

	c.JSON(http.StatusOK, gin.H{
		"sessions": gin.H{
			"active": activeSessions,
			"max":    25,
		},
		"messages":        msgStats,
		"uptime_seconds":  int64(time.Since(s.startTime).Seconds()),
		"timestamp":       time.Now().Unix(),
	})
}

// ServeSwaggerUI serves the Swagger UI for API documentation
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

// ServeOpenAPISpec serves the OpenAPI specification file
func (s *Server) ServeOpenAPISpec(c *gin.Context) {
	openapi := `openapi: 3.0.0
info:
  title: WACAST Core API
  description: Enterprise-grade WhatsApp Gateway
  version: 1.0.0
servers:
  - url: http://localhost:8080
    description: Local Development Server
  - url: https://api.wacast.io
    description: Production Server
paths:
  /health:
    get:
      summary: Overall health status
      tags: [Health]
      responses:
        '200':
          description: Healthy
  /health/ready:
    get:
      summary: Readiness probe
      tags: [Health]
      responses:
        '200':
          description: Ready
  /health/live:
    get:
      summary: Liveness probe
      tags: [Health]
      responses:
        '200':
          description: Live
  /api/v1/sessions:
    get:
      summary: List active sessions
      tags: [Sessions]
      responses:
        '200':
          description: List of sessions
  /api/v1/sessions/:device_id:
    get:
      summary: Get session status
      tags: [Sessions]
      parameters:
        - name: device_id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Session details
        '404':
          description: Session not found
  /api/v1/sessions/initiate:
    post:
      summary: Initiate new session
      tags: [Sessions]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
      responses:
        '202':
          description: Session initiated
  /api/v1/sessions/:device_id/stop:
    post:
      summary: Stop session
      tags: [Sessions]
      parameters:
        - name: device_id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Session stopped
  /api/v1/devices/:device_id/messages:
    post:
      summary: Send text message
      tags: [Messages]
      parameters:
        - name: device_id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
      responses:
        '202':
          description: Message queued
  /api/v1/messages/stats:
    get:
      summary: Get queue statistics
      tags: [Messages]
      responses:
        '200':
          description: Queue stats
  /api/v1/info/status:
    get:
      summary: Server status
      tags: [Server Info]
      responses:
        '200':
          description: Server status
  /api/v1/info/stats:
    get:
      summary: Server statistics
      tags: [Server Info]
      responses:
        '200':
          description: Server statistics
`

	c.Header("Content-Type", "text/yaml; charset=utf-8")
	c.String(http.StatusOK, openapi)
}
