package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/services/message"
	"wacast/core/utils"
)

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
	TargetJID   string    `json:"target_jid" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	ScheduledFor time.Time `json:"scheduled_for" binding:"required"`
	GroupID     *string   `json:"group_id"`
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

	messageID, err := h.messageService.SendMessage(ctx, deviceID, req.TargetJID, req.Content, req.GroupID)
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

// SendMessageWithMedia sends a message with media attachment
// POST /devices/:device_id/messages/media
func (h *MessageHandler) SendMessageWithMedia(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req SendMessageWithMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	ctx := c.Request.Context()

	messageID, err := h.messageService.SendMessageWithMedia(ctx, deviceID, req.TargetJID, req.MediaURL, req.ContentType, req.Caption)
	if err != nil {
		utils.Error("Failed to send media message",
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

// SendScheduledMessage sends a scheduled message
// POST /devices/:device_id/messages/scheduled
func (h *MessageHandler) SendScheduledMessage(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req SendScheduledMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Validate scheduled time is in future
	if req.ScheduledFor.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheduled time must be in the future",
		})
		return
	}

	ctx := c.Request.Context()

	messageID, err := h.messageService.SendScheduledMessage(ctx, deviceID, req.TargetJID, req.Content, req.ScheduledFor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to schedule message: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id":   messageID,
		"status":       "scheduled",
		"scheduled_for": req.ScheduledFor,
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

// RegisterMessageRoutes registers all message routes
func RegisterMessageRoutes(router interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}, messageService *message.Service) {
	handler := NewMessageHandler(messageService)

	// Message sending endpoints
	messages := router.Group("/devices/:device_id/messages")
	{
		messages.POST("", handler.SendMessage)                    // Send text message
		messages.POST("/media", handler.SendMessageWithMedia)     // Send media message
		messages.POST("/scheduled", handler.SendScheduledMessage) // Schedule message
	}

	// Queue management endpoints
	queue := router.Group("/messages")
	{
		queue.GET("/:message_id/status", handler.GetMessageStatus)  // Check message status
		queue.GET("/stats", handler.GetQueueStats)                 // Get queue stats
		queue.GET("/failed", handler.GetFailedMessages)            // List failed messages
		queue.POST("/process", handler.ProcessQueueManually)       // Manual process
	}
}
