package appserver

import (
	"fmt"
	"net"
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
	"wacast/core/services/integration"
	"wacast/core/services/license"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/services/settings"
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
	settingsService  *settings.Service
	integrationService *integration.Service
	port             int
	host             string
	startTime        time.Time
	websocketHandler *handlers.WebSocketHandler
	httpInstance     *http.Server
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
		settingsService:  settings.NewService(db),
		integrationService: integration.NewService(db),
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

	s.engine.GET("/setup", s.ServeSetupPage)
	s.registerStaticRoutes()

	s.engine.GET("/api/docs", s.ServeSwaggerUI)
	s.engine.GET("/openapi.yaml", s.ServeOpenAPISpec)

	s.engine.Static("/uploads", "./uploads")

	v1 := s.engine.Group("/api/v1")
	{
		handlers.RegisterAuthRoutes(v1, s.authService, s.config.JWTSecret)
		handlers.RegisterBillingRoutes(v1, s.billingService, s.config.JWTSecret, s.authService)
		handlers.RegisterSessionRoutes(v1, s.sessionService, s.config.EncryptionKey, s.config.SessionTimeout, s.config.JWTSecret, s.authService)
		handlers.RegisterMessageRoutes(v1, s.messageService, s.config.JWTSecret, s.authService, s.integrationService)
		handlers.RegisterContactRoutes(v1, s.contactService, s.config.JWTSecret, s.authService)
		handlers.RegisterAnalyticRoutes(v1, s.analyticService, s.config.JWTSecret, s.authService)
		handlers.RegisterBroadcastRoutes(v1, s.broadcastService, s.config.JWTSecret, s.authService)
		handlers.RegisterAutoResponseRoutes(v1, s.autoresponseService, s.config.JWTSecret, s.authService)

		settingHandler := handlers.NewSettingHandler(s.settingsService, s.messageService)
		settingHandler.RegisterRoutes(v1)

		configHandler := handlers.NewConfigHandler(s.db, s.config)
		configHandler.RegisterRoutes(v1)

		licenseHandler := handlers.NewLicenseHandler(s.licenseService)
		licenseHandler.RegisterRoutes(v1)

		infoHandler := handlers.NewInfoHandler(s.settingsService)
		infoHandler.RegisterRoutes(v1)

		info := v1.Group("/info")
		{
			info.GET("/status", s.ServerStatus)
			info.GET("/stats", s.ServerStats)
		}

		integrationHandler := handlers.NewIntegrationHandler(s.integrationService, s.config.JWTSecret, s.authService)
		integrationHandler.RegisterRoutes(v1)
	}

	s.engine.GET("/ws/sessions/:device_id/qr", s.websocketHandler.ConnectQR)
	s.engine.GET("/qr/:device_id", s.websocketHandler.ServeQRPage)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Retry logic for Windows port lingering using manual Listen
	var listener net.Listener
	var err error
	for i := 0; i < 15; i++ {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			utils.Warn("Port 8080 is busy (Windows socket lingering), retrying...", zap.Int("attempt", i+1), zap.Error(err))
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return fmt.Errorf("failed to bind to port %d after multiple attempts: %v", s.port, err)
	}

	s.httpInstance = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	utils.Info("Successfully bound to port, starting HTTP server", zap.String("address", addr))
	
	if err := s.httpInstance.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	utils.Info("Shutting down HTTP server")
	if s.httpInstance != nil {
		return s.httpInstance.Close() // Force close to release port immediately
	}
	return nil
}

func (s *Server) licenseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if !strings.HasPrefix(path, "/api/") || 
		   strings.Contains(path, "/api/v1/license") || 
		   strings.Contains(path, "/api/v1/config") ||
		   strings.Contains(path, "/health") ||
		   strings.Contains(path, "/uploads") {
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

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Uptime    int64                  `json:"uptime_seconds"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func (s *Server) ServeSetupPage(c *gin.Context) {
	hwid, _ := utils.GetHWID()
	if hwid == "" {
		hwid = "UNKNOWN-HWID"
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WACAST PRO - Activation</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #0b0e14;
            --card: #151921;
            --accent: #2563eb;
            --accent-hover: #1d4ed8;
            --text: #f8fafc;
            --text-dim: #94a3b8;
            --success: #10b981;
            --danger: #ef4444;
            --border: #262f40;
        }
        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg);
            color: var(--text);
            margin: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background-image: radial-gradient(circle at 50% -20%, #1e293b 0%, #0b0e14 100%);
        }
        .container {
            width: 100%;
            max-width: 440px;
            padding: 20px;
            animation: fadeIn 0.5s ease-out;
        }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        .card {
            background: var(--card);
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
            border: 1px solid var(--border);
            position: relative;
            overflow: hidden;
        }
        .card::before {
            content: '';
            position: absolute;
            top: 0; left: 0; right: 0; height: 4px;
            background: linear-gradient(90deg, var(--accent), #8b5cf6);
        }
        h1 {
            margin: 0 0 10px 0;
            font-size: 28px;
            font-weight: 700;
            text-align: center;
            background: linear-gradient(to bottom, #fff, #94a3b8);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        p.subtitle {
            color: var(--text-dim);
            text-align: center;
            font-size: 15px;
            margin-bottom: 32px;
            line-height: 1.5;
        }
        .section-label {
            display: block;
            margin-bottom: 12px;
            font-size: 12px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 1.5px;
            color: var(--accent);
        }
        .hwid-box {
            display: flex;
            gap: 12px;
            background: #090c10;
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 12px 16px;
            align-items: center;
            margin-bottom: 28px;
            transition: border-color 0.2s;
        }
        .hwid-box:hover { border-color: #334155; }
        .hwid-value {
            flex: 1;
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            color: #cbd5e1;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .btn-copy {
            background: #1e293b;
            color: #fff;
            border: 1px solid #334155;
            border-radius: 6px;
            padding: 6px 12px;
            font-size: 11px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
        }
        .btn-copy:hover { background: #334155; border-color: #475569; }

        input {
            width: 100%;
            padding: 14px 18px;
            background: #090c10;
            border: 1px solid var(--border);
            border-radius: 12px;
            color: white;
            font-size: 15px;
            outline: none;
            box-sizing: border-box;
            transition: all 0.2s;
            margin-bottom: 20px;
        }
        input:focus { border-color: var(--accent); box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.1); }
        
        .btn-primary {
            width: 100%;
            padding: 16px;
            background: var(--accent);
            color: white;
            border: none;
            border-radius: 12px;
            font-size: 16px;
            font-weight: 700;
            cursor: pointer;
            transition: all 0.2s;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        }
        .btn-primary:hover { background: var(--accent-hover); transform: translateY(-1px); }
        .btn-primary:active { transform: translateY(0); }
        .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }

        .contact-card {
            margin-top: 32px;
            padding: 20px;
            background: rgba(16, 185, 129, 0.03);
            border: 1px solid rgba(16, 185, 129, 0.1);
            border-radius: 16px;
            text-align: center;
        }
        .contact-text { font-size: 13px; color: var(--text-dim); margin-bottom: 12px; }
        .wa-link {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: #059669;
            color: white;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 30px;
            font-size: 14px;
            font-weight: 600;
            transition: all 0.2s;
        }
        .wa-link:hover { background: #047857; transform: scale(1.05); }

        #message {
            margin-top: 20px;
            padding: 14px;
            border-radius: 12px;
            font-size: 14px;
            display: none;
            text-align: center;
            font-weight: 500;
        }
        .msg-success { background: rgba(16, 185, 129, 0.1); color: var(--success); border: 1px solid rgba(16, 185, 129, 0.2); }
        .msg-error { background: rgba(239, 68, 68, 0.1); color: var(--danger); border: 1px solid rgba(239, 68, 68, 0.2); }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>Activation</h1>
            <p class="subtitle">Enter your license key to unlock WACAST PRO features.</p>
            
            <span class="section-label">Device Hardware ID</span>
            <div class="hwid-box">
                <div class="hwid-value" id="hwidDisplay">` + hwid + `</div>
                <button class="btn-copy" onclick="copyHWID()">Copy ID</button>
            </div>

            <span class="section-label">License Key</span>
            <input type="text" id="licenseKey" placeholder="XXXX-XXXX-XXXX-XXXX" autocomplete="off">
            
            <button id="btnActivate" class="btn-primary" onclick="activateLicense()">Activate System</button>
            
            <div id="message"></div>

            <div class="contact-card">
                <div class="contact-text">Don't have a license key yet?</div>
                <a href="#" id="waLink" class="wa-link">
                    <svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M.057 24l1.687-6.163c-1.041-1.804-1.588-3.849-1.587-5.946.003-6.556 5.338-11.891 11.893-11.891 3.181.001 6.167 1.24 8.413 3.488 2.245 2.248 3.481 5.236 3.48 8.414-.003 6.557-5.338 11.892-11.893 11.892-1.99-.001-3.951-.5-5.688-1.448l-6.305 1.654zm6.597-3.807c1.676.995 3.276 1.591 5.392 1.592 5.448 0 9.886-4.438 9.889-9.885.002-5.462-4.415-9.89-9.881-9.892-5.452 0-9.887 4.434-9.889 9.884-.001 2.225.651 3.891 1.746 5.634l-.999 3.648 3.742-.981zm11.387-5.464c-.074-.124-.272-.198-.57-.347-.297-.149-1.758-.868-2.031-.967-.272-.099-.47-.149-.669.149-.198.297-.768.967-.941 1.165-.173.198-.347.223-.644.074-.297-.149-1.255-.462-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.695.248-1.29.173-1.414z"/></svg>
                    Contact via WhatsApp
                </a>
            </div>
        </div>
    </div>

    <script>
        const hwid = '` + hwid + `';
        const waNumber = '6285887373722';
        
        // Setup WhatsApp link with HWID
        document.getElementById('waLink').href = 'https://wa.me/' + waNumber + '?text=' + 
            encodeURIComponent('Halo Admin, saya ingin aktivasi WACAST PRO.\n\nHWID saya: ' + hwid);

        function copyHWID() {
            navigator.clipboard.writeText(hwid).then(() => {
                const btn = document.querySelector('.btn-copy');
                const originalText = btn.innerText;
                btn.innerText = 'Copied!';
                btn.style.borderColor = '#10b981';
                btn.style.color = '#10b981';
                setTimeout(() => {
                    btn.innerText = originalText;
                    btn.style.borderColor = '#334155';
                    btn.style.color = '#fff';
                }, 2000);
            });
        }

        async function activateLicense() {
            const key = document.getElementById('licenseKey').value.trim();
            const btn = document.getElementById('btnActivate');
            const msg = document.getElementById('message');
            
            if (!key) {
                showMessage('License key is required', 'error');
                return;
            }
            
            btn.disabled = true;
            btn.innerText = 'Activating...';
            msg.style.display = 'none';
            
            try {
                const response = await fetch('/api/v1/license/activate', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ key: key })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    showMessage('Success! License activated. Redirecting...', 'success');
                    setTimeout(() => {
                        window.location.href = '/';
                    }, 2000);
                } else {
                    showMessage(data.message || 'Activation failed', 'error');
                    btn.disabled = false;
                    btn.innerText = 'Activate License';
                }
            } catch (e) {
                showMessage('Connection error. Please try again.', 'error');
                btn.disabled = false;
                btn.innerText = 'Activate License';
            }
        }
        
        function showMessage(text, type) {
            const msg = document.getElementById('message');
            msg.innerText = text;
            msg.className = type === 'success' ? 'msg-success' : 'msg-error';
            msg.style.display = 'block';
        }
    </script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}


