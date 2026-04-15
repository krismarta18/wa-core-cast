package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/services/session"
	"wacast/core/utils"
)

// SessionHandler handles WhatsApp session endpoints
type SessionHandler struct {
	sessionService *session.Service
	encryptionKey  string
	sessionTimeout int
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(svc *session.Service, encryptionKey string, sessionTimeout int) *SessionHandler {
	return &SessionHandler{
		sessionService: svc,
		encryptionKey:  encryptionKey,
		sessionTimeout: sessionTimeout,
	}
}

// InitiateSessionRequest is the request body for initiating a session
type InitiateSessionRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	DisplayName string `json:"display_name"`
}

type SessionStatusResponse struct {
	DeviceID    string `json:"device_id"`
	Status      int    `json:"status"` // 0=inactive, 1=active, 2=pending
	IsActive    bool   `json:"is_active"`
	Phone       string `json:"phone,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// GetSessionStatus retrieves the status of a session
// GET /sessions/:device_id
func (h *SessionHandler) GetSessionStatus(c *gin.Context) {
	deviceID := c.Param("device_id")

	status := h.sessionService.GetSessionStatus(deviceID)
	isActive := h.sessionService.IsSessionActive(deviceID)
	
	var phone, displayName string
	sessionData := h.sessionService.GetSession(deviceID)
	if sessionData != nil && sessionData.Config != nil {
		phone = sessionData.Config.Phone
		displayName = sessionData.Config.DisplayName
	}

	c.JSON(http.StatusOK, SessionStatusResponse{
		DeviceID:    deviceID,
		Status:      int(status),
		IsActive:    isActive,
		Phone:       phone,
		DisplayName: displayName,
	})
}

// GetAllActiveSessions retrieves all active sessions
// GET /sessions
func (h *SessionHandler) GetAllActiveSessions(c *gin.Context) {
	sessions := h.sessionService.GetAllActiveSessions()

	response := make([]SessionStatusResponse, len(sessions))
	for i, s := range sessions {
		var phone, displayName string
		if s.Config != nil {
			phone = s.Config.Phone
			displayName = s.Config.DisplayName
		}
		
		response[i] = SessionStatusResponse{
			DeviceID:    s.ID,
			Status:      int(s.Status),
			IsActive:    s.Status == session.SessionActive,
			Phone:       phone,
			DisplayName: displayName,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(sessions),
		"sessions": response,
	})
}

// GetQRCode retrieves the QR code for a pending session
// GET /sessions/:device_id/qr
// Query params:
//   - format=png (returns PNG image directly)
//   - format=json (default, returns JSON with base64 PNG)
func (h *SessionHandler) GetQRCode(c *gin.Context) {
	deviceID := c.Param("device_id")
	format := c.DefaultQuery("format", "json") // Default to JSON format

	// Get QR code string from session service
	qrCodeString := h.sessionService.GetQRCode(deviceID)
	if qrCodeString == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "QR code not available for this device",
		})
		return
	}

	// Generate PNG image from the QR string
	pngBytes, err := h.sessionService.GenerateQRCodeImage(qrCodeString)
	if err != nil {
		utils.Error("Failed to generate QR code image",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to generate QR code image: %v", err),
		})
		return
	}

	utils.Info("QR code generated",
		zap.String("device_id", deviceID),
		zap.String("qr_data", qrCodeString),
		zap.Int("png_size_bytes", len(pngBytes)),
	)

	// Return PNG image directly if requested
	if format == "png" {
		c.Header("Content-Disposition", "inline; filename=qrcode.png")
		c.Data(http.StatusOK, "image/png", pngBytes)
		return
	}

	// Return JSON response with base64-encoded PNG and alternative formats
	base64PNG := base64.StdEncoding.EncodeToString(pngBytes)
	
	// URL encode the QR string for third-party QR services
	encodedQR := url.QueryEscape(qrCodeString)
	
	// Use QR server for alternative URL-based access
	qrImageURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", encodedQR)

	// Return response with all formats available
	c.JSON(http.StatusOK, gin.H{
		"device_id":      deviceID,
		"qr_code_string": qrCodeString, // ✅ Raw string for WhatsApp linking
		"qr_code_image": gin.H{
			"base64_png": "data:image/png;base64," + base64PNG, // ✅ Base64 PNG for direct display in <img>
			"png_bytes":  len(pngBytes),                         // ✅ PNG size in bytes
			"url_format": qrImageURL,                            // ✅ Alternative URL for legacy systems
			"direct_url": fmt.Sprintf("/api/v1/sessions/%s/qr?format=png", deviceID), // ✅ Direct PNG endpoint
		},
		"status":  int(session.SessionPending),
		"message": "Scan this QR code with WhatsApp mobile app",
		"instructions": gin.H{
			"step1": "Open WhatsApp on your mobile device",
			"step2": "Go to Settings > Linked Devices",
			"step3": "Tap 'Link a Device'",
			"step4": "Scan the QR code displayed above",
			"step5": "Wait for connection confirmation (auto)",
		},
		"frontend_usage": gin.H{
			"method1":     "Use direct_url in <img src> tag (e.g., <img src='/api/v1/sessions/{device_id}/qr?format=png'>)",
			"method2":     "Display base64_png in <img src> tag for inline embedding",
			"method3":     "Use qr_code_string with qrcode.js library: new QRCode(element, qr_code_string)",
			"deprecated":  "url_format uses third-party QR API for fallback",
		},
	})
}

// InitiateSession starts a new WhatsApp session
// POST /sessions/initiate
func (h *SessionHandler) InitiateSession(c *gin.Context) {
	var req InitiateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	ctx := c.Request.Context()

	// Register QR code handler for this session
	h.sessionService.RegisterQRCodeCallback(req.DeviceID, func(event *session.QRCodeEvent) {
		utils.Info("QR code generated",
			zap.String("device_id", event.DeviceID),
		)
		// TODO: Send QR code to frontend via WebSocket or HTTP
	})

	// Register status callback
	h.sessionService.RegisterStatusCallback(req.DeviceID, func(event *session.ConnectionStatusEvent) {
		utils.Info("Connection status changed",
			zap.String("device_id", event.DeviceID),
			zap.Int("new_status", int(event.Status)),
		)
		// TODO: Notify frontend of status change
	})

	// Register message handler
	h.sessionService.RegisterMessageHandler(req.DeviceID, func(event *session.MessageReceivedEvent) {
		utils.Debug("Message received",
			zap.String("device_id", event.DeviceID),
			zap.String("from_jid", event.FromJID),
			zap.String("content", event.Content),
		)
		// TODO: Save message to database via message service
	})

	// Start the session
	cfg := &session.SessionConfig{
		DeviceID:       req.DeviceID,
		UserID:         req.UserID,
		Phone:          req.Phone,
		DisplayName:    req.DisplayName,
		EncryptionKey:  h.encryptionKey,
		SessionTimeout: h.sessionTimeout,
		ReconnectLimit: 5,
	}

	if err := h.sessionService.StartSession(ctx, cfg); err != nil {
		utils.Error("Failed to start session",
			zap.String("device_id", req.DeviceID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to start session: %v", err),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Session initiated. Scan QR code via GET /api/v1/sessions/{device_id}/qr",
		"device_id":   req.DeviceID,
		"status":      int(session.SessionPending),
		"next_step":   "GET /api/v1/sessions/" + req.DeviceID + "/qr to retrieve QR code",
		"poll_status": "GET /api/v1/sessions/" + req.DeviceID + " to check status",
	})
}

// StopSession stops an active WhatsApp session
// POST /sessions/:device_id/stop
func (h *SessionHandler) StopSession(c *gin.Context) {
	deviceID := c.Param("device_id")
	ctx := c.Request.Context()

	if err := h.sessionService.StopSession(ctx, deviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to stop session: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Session stopped",
		"device_id": deviceID,
	})
}

// Message sending is handled through message handlers at /api/v1/devices/:device_id/messages
// This maintains separation of concerns and allows proper queuing

// RegisterSessionRoutes registers all session routes
func RegisterSessionRoutes(router interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}, sessionService *session.Service, encryptionKey string, sessionTimeout int) {
	handler := NewSessionHandler(sessionService, encryptionKey, sessionTimeout)

	// Session endpoints
	sessions := router.Group("/sessions")
	{
		sessions.GET("", handler.GetAllActiveSessions)
		sessions.GET("/:device_id", handler.GetSessionStatus)
		sessions.GET("/:device_id/qr", handler.GetQRCode)
		sessions.POST("/initiate", handler.InitiateSession)
		sessions.POST("/:device_id/stop", handler.StopSession)
		// Message sending is handled through /api/v1/devices/:device_id/messages
	}
}
