package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/services/auth"
	"wacast/core/services/message"
	"wacast/core/utils"
)

const uploadDir = "./uploads"

// MessageHandler handles message-related endpoints
type MessageHandler struct {
	messageService *message.Service
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(svc *message.Service) *MessageHandler {
	return &MessageHandler{
		messageService: svc,
	}
}

// SendMessageRequest is the request body for sending a message
type SendMessageRequest struct {
	TargetJID string `json:"target_jid" binding:"required"`  // Recipient JID
	Content   string `json:"content" binding:"required"`     // Message content
	GroupID   *string `json:"group_id"`                      // Optional group ID
	Priority  int    `json:"priority"`                       // 1-5, default 3
}

// SendMessageWithMediaRequest is request for media message
type SendMessageWithMediaRequest struct {
	TargetJID   string `json:"target_jid" binding:"required"`  
	MediaURL    string `json:"media_url" binding:"required"`   
	ContentType string `json:"content_type" binding:"required"` // image, document, audio, video
	Caption     *string `json:"caption"`                       
	GroupID     *string `json:"group_id"`                      
}

// SendScheduledMessageRequest is request for scheduled message
type SendScheduledMessageRequest struct {
	TargetJID    string    `json:"target_jid"`
	Content      string    `json:"content"`
	ScheduledFor time.Time `json:"scheduled_for"`
	GroupID      *string   `json:"group_id"`
	MediaURL     string    `json:"media_url"`
	ContentType  string    `json:"content_type"`
	Caption      *string   `json:"caption"`
}

// MessageStatusResponse is response for message status
type MessageStatusResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"` // pending, sent, delivered, read, failed
	Timestamp int64  `json:"timestamp"`
}

// MessageSendResponse is response when sending a message
type MessageSendResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// QueueStatsResponse is response for queue statistics
type QueueStatsResponse struct {
	TotalSent    int64   `json:"total_sent"`
	TotalReceived int64  `json:"total_received"`
	TotalFailed  int64   `json:"total_failed"`
	Pending      int64   `json:"pending"`
	AvgLatency   float64 `json:"avg_latency_ms"`
}

// SendMessage sends a text message
// POST /devices/:device_id/messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	ctx := c.Request.Context()

	messageID, err := h.messageService.SendMessage(ctx, deviceID, req.TargetJID, req.Content, req.GroupID, nil)
	if err != nil {
		utils.Error("Failed to send message",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to send message: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, MessageSendResponse{
		MessageID: messageID,
		Status:    "pending",
		Timestamp: time.Now().Unix(),
	})
}

// SendMessageWithMedia sends a message with media attachment.
// Accepts multipart/form-data with fields: target_jid, content_type, caption (optional), file.
// POST /devices/:device_id/messages/media
func (h *MessageHandler) SendMessageWithMedia(c *gin.Context) {
	deviceID := c.Param("device_id")

	// Parse multipart form (max 50 MB)
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		// Fallback: try JSON (backward compat)
		var req SendMessageWithMediaRequest
		if jsonErr := c.ShouldBindJSON(&req); jsonErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: expected multipart/form-data or JSON"})
			return
		}
		ctx := c.Request.Context()
		messageID, err := h.messageService.SendMessageWithMedia(ctx, deviceID, req.TargetJID, req.MediaURL, req.ContentType, req.Caption, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, MessageSendResponse{MessageID: messageID, Status: "pending", Timestamp: time.Now().Unix()})
		return
	}

	targetJID := c.PostForm("target_jid")
	contentType := c.PostForm("content_type")
	captionStr := c.PostForm("caption")
	var caption *string
	if captionStr != "" {
		caption = &captionStr
	}

	if targetJID == "" || contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_jid and content_type are required"})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required in multipart form"})
		return
	}
	defer file.Close()

	// Save file to uploads dir
	mediaURL, err := saveUploadedFile(file, header.Filename)
	if err != nil {
		utils.Error("Failed to save uploaded file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	ctx := c.Request.Context()
	messageID, err := h.messageService.SendMessageWithMedia(ctx, deviceID, targetJID, mediaURL, contentType, caption, nil)
	if err != nil {
		utils.Error("Failed to send media message", zap.String("device_id", deviceID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageSendResponse{
		MessageID: messageID,
		Status:    "pending",
		Timestamp: time.Now().Unix(),
	})
}

// UploadMedia uploads a file and returns its public URL.
// POST /upload/media
func (h *MessageHandler) UploadMedia(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file field is required"})
		return
	}
	defer file.Close()

	url, err := saveUploadedFile(file, header.Filename)
	if err != nil {
		utils.Error("Failed to save file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      url,
		"filename": header.Filename,
		"size":     header.Size,
	})
}

// saveUploadedFile persists the multipart file to the uploads directory and returns its public URL.
func saveUploadedFile(file io.Reader, originalName string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return a URL that Go server can serve (see static file route)
	baseURL := "http://localhost:8080"
	return fmt.Sprintf("%s/uploads/%s", baseURL, filename), nil
}

// SendScheduledMessage sends a scheduled message
// POST /devices/:device_id/messages/scheduled
func (h *MessageHandler) SendScheduledMessage(c *gin.Context) {
	deviceID := c.Param("device_id")

	var targetJID string
	var content string
	var scheduledFor time.Time
	var contentType string
	var mediaURL *string
	var caption *string

	// Try to parse multipart form first (for file uploads)
	if err := c.Request.ParseMultipartForm(50 << 20); err == nil {
		targetJID = c.PostForm("target_jid")
		content = c.PostForm("content")
		captionStr := c.PostForm("caption")
		if captionStr != "" {
			caption = &captionStr
		}
		contentType = c.PostForm("content_type")
		scheduledAtStr := c.PostForm("scheduled_for")

		if scheduledAtStr != "" {
			var parseErr error
			scheduledFor, parseErr = time.Parse(time.RFC3339, scheduledAtStr)
			if parseErr != nil {
				// Retry with Date string format if RFC3339 fails (next.js often sends ISO strings)
				scheduledFor, parseErr = time.Parse("2006-01-02T15:04", scheduledAtStr)
				if parseErr != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled_for format (RFC3339 or ISO required)"})
					return
				}
			}
		}

		// Handle file if present
		file, header, fileErr := c.Request.FormFile("file")
		if fileErr == nil {
			defer file.Close()
			savedURL, saveErr := saveUploadedFile(file, header.Filename)
			if saveErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
			mediaURL = &savedURL
			if contentType == "" {
				contentType = "image" // Default to image if not specified
			}
		}
	} else {
		// Fallback to JSON
		var req SendScheduledMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: expected multipart/form-data or JSON"})
			return
		}
		targetJID = req.TargetJID
		content = req.Content
		scheduledFor = req.ScheduledFor
		if req.MediaURL != "" {
			mediaURL = &req.MediaURL
		}
		contentType = req.ContentType
		caption = req.Caption
	}

	if targetJID == "" || scheduledFor.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_jid and scheduled_for are required"})
		return
	}

	// Validate scheduled time is in future
	if scheduledFor.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheduled time must be in the future",
		})
		return
	}

	ctx := c.Request.Context()
	messageID, err := h.messageService.SendScheduledMessage(ctx, deviceID, targetJID, content, scheduledFor, mediaURL, contentType, caption, nil)
	if err != nil {
		utils.Error("Failed to schedule message", zap.String("device_id", deviceID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to schedule message: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, MessageSendResponse{
		MessageID: messageID,
		Status:    "scheduled",
		Timestamp: time.Now().Unix(),
	})
}

// GetMessageStatus retrieves the status of a message
// GET /messages/:message_id/status
func (h *MessageHandler) GetMessageStatus(c *gin.Context) {
	messageID := c.Param("message_id")

	status, err := h.messageService.GetMessageStatus(messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Message not found",
		})
		return
	}

	statusText := "unknown"
	switch status {
	case message.StatusPending:
		statusText = "pending"
	case message.StatusSent:
		statusText = "sent"
	case message.StatusDelivered:
		statusText = "delivered"
	case message.StatusRead:
		statusText = "read"
	case message.StatusFailed:
		statusText = "failed"
	}

	c.JSON(http.StatusOK, MessageStatusResponse{
		MessageID: messageID,
		Status:    statusText,
		Timestamp: time.Now().Unix(),
	})
}

// GetQueueStats retrieves message queue statistics
// GET /messages/stats
func (h *MessageHandler) GetQueueStats(c *gin.Context) {
	stats := h.messageService.GetQueueStats()

	c.JSON(http.StatusOK, QueueStatsResponse{
		TotalSent:    stats["total_sent"].(int64),
		TotalReceived: stats["total_received"].(int64),
		TotalFailed:  stats["total_failed"].(int64),
		Pending:      stats["pending"].(int64),
		AvgLatency:   stats["avg_latency_ms"].(float64),
	})
}

// GetFailedMessages retrieves failed messages for retry
// GET /messages/failed
func (h *MessageHandler) GetFailedMessages(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	failed, err := h.messageService.GetFailedMessages(limit)
	if err != nil {
		utils.Error("Failed to get failed messages", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve failed messages",
		})
		return
	}

	response := make([]gin.H, len(failed))
	for i, msg := range failed {
		response[i] = gin.H{
			"message_id": msg.ID,
			"device_id":  msg.DeviceID,
			"target_jid": msg.TargetJID,
			"content":    msg.Content,
			"retry_count": msg.RetryCount,
			"max_retries": msg.MaxRetries,
			"error":      msg.ErrorLog,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(response),
		"messages": response,
	})
}

// ProcessQueueManually triggers manual queue processing
// POST /messages/process
func (h *MessageHandler) ProcessQueueManually(c *gin.Context) {
	if err := h.messageService.ProcessQueue(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to process queue: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Queue processing triggered",
	})
}

// ListScheduledMessages returns any pending messages scheduled for the future
// GET /devices/:device_id/messages/scheduled
func (h *MessageHandler) ListScheduledMessages(c *gin.Context) {
	deviceID := c.Param("device_id")
	msgs, err := h.messageService.ListScheduledMessages(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

// ListMessageHistory returns sent or failed messages for a device
// GET /devices/:device_id/messages/history
func (h *MessageHandler) ListMessageHistory(c *gin.Context) {
	deviceID := c.Param("device_id")
	msgs, err := h.messageService.ListMessageHistory(deviceID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

// CancelScheduledMessage cancels a pending message
// DELETE /messages/:message_id
func (h *MessageHandler) CancelScheduledMessage(c *gin.Context) {
	messageID := c.Param("message_id")
	if err := h.messageService.CancelScheduledMessage(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scheduled message cancelled successfully"})
}

// RegisterMessageRoutes registers all message routes
func RegisterMessageRoutes(router interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}, messageService *message.Service, jwtSecret string, authService *auth.Service) {
	handler := NewMessageHandler(messageService)

	// Protected group
	group := router.Group("")
	group.Use(JWTAuthMiddleware(jwtSecret, authService))

	// Message sending endpoints
	messages := group.Group("/devices/:device_id/messages")
	{
		messages.POST("", handler.SendMessage)                    // Send text message
		messages.POST("/media", handler.SendMessageWithMedia)     // Send media message (multipart)
		messages.POST("/scheduled", handler.SendScheduledMessage) // Schedule message
		messages.GET("/scheduled", handler.ListScheduledMessages) // List pending scheduled
		messages.GET("/history", handler.ListMessageHistory)      // List history
	}

	// File upload
	group.POST("/upload/media", handler.UploadMedia) // Upload media file → returns URL

	// Queue management endpoints
	queue := group.Group("/messages")
	{
		queue.GET("/:message_id/status", handler.GetMessageStatus)  // Check message status
		queue.DELETE("/:message_id", handler.CancelScheduledMessage) // Cancel/Delete message
		queue.GET("/stats", handler.GetQueueStats)                 // Get queue stats
		queue.GET("/failed", handler.GetFailedMessages)            // List failed messages
		queue.POST("/process", handler.ProcessQueueManually)       // Manual process
	}
}
