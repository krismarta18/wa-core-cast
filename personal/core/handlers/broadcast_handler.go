package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"wacast/core/models"
	"wacast/core/services/auth"
	"wacast/core/services/billing"
	"wacast/core/services/broadcast"
	"wacast/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BroadcastHandler struct {
	service *broadcast.Service
}

func NewBroadcastHandler(service *broadcast.Service) *BroadcastHandler {
	return &BroadcastHandler{service: service}
}

func (h *BroadcastHandler) CreateCampaign(c *gin.Context) {
	userIDStr, _ := c.Get(ContextKeyUserID) // Set by JWTAuthMiddleware
	userID, _ := uuid.Parse(fmt.Sprintf("%v", userIDStr))

	var req struct {
		models.CreateBroadcastCampaignRequest
		Recipients []string `json:"recipients" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := h.service.CreateCampaign(userID, req.CreateBroadcastCampaignRequest, req.Recipients)
	if err != nil {
		utils.Error("Failed to create campaign", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign.ToResponse())
}

func (h *BroadcastHandler) StartCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	if err := h.service.StartCampaign(c.Request.Context(), campaignID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == billing.ErrMessageLimitReached.Error() || strings.Contains(err.Error(), billing.ErrMessageLimitReached.Error()) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign started successfully"})
}

func (h *BroadcastHandler) ListCampaigns(c *gin.Context) {
	userIDStr, _ := c.Get(ContextKeyUserID)
	userID, _ := uuid.Parse(fmt.Sprintf("%v", userIDStr))

	campaigns, err := h.service.ListCampaigns(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*models.BroadcastCampaignResponse, len(campaigns))
	for i, cp := range campaigns {
		response[i] = cp.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"broadcasts": response})
}

func (h *BroadcastHandler) GetCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := h.service.GetCampaign(campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func RegisterBroadcastRoutes(router interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}, service *broadcast.Service, jwtSecret string, authService *auth.Service) {
	handler := NewBroadcastHandler(service)

	g := router.Group("/broadcasts")
	g.Use(JWTAuthMiddleware(jwtSecret, authService))
	{
		g.POST("", handler.CreateCampaign)
		g.POST("/:id/start", handler.StartCampaign)
		g.GET("", handler.ListCampaigns)
		g.GET("/:id", handler.GetCampaign)
	}
}
