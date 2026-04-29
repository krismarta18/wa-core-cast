package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wacast/core/services/warming"
)

type WarmingHandler struct {
	warmingService *warming.Service
}

func NewWarmingHandler(warmingService *warming.Service) *WarmingHandler {
	return &WarmingHandler{
		warmingService: warmingService,
	}
}

type StartWarmingRequest struct {
	DeviceIDs       []string `json:"device_ids" binding:"required"`
	DurationMinutes int      `json:"duration_minutes" binding:"required,min=1"`
}

func (h *WarmingHandler) StartWarming(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	var req StartWarmingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	deviceUUIDs := []uuid.UUID{}
	for _, idStr := range req.DeviceIDs {
		dID, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID: " + idStr, "success": false})
			return
		}
		deviceUUIDs = append(deviceUUIDs, dID)
	}

	if err := h.warmingService.StartWarming(c.Request.Context(), userID, deviceUUIDs, req.DurationMinutes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sesi warming berhasil dimulai.",
	})
}

func (h *WarmingHandler) StopWarming(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	if err := h.warmingService.StopWarming(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sesi warming telah dihentikan.",
	})
}

func (h *WarmingHandler) GetStatus(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	status := h.warmingService.GetStatus(userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}
