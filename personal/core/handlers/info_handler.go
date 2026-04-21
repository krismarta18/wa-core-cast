package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wacast/core/services/settings"
)

type InfoHandler struct {
	settingsService *settings.Service
}

func NewInfoHandler(settingsService *settings.Service) *InfoHandler {
	return &InfoHandler{settingsService: settingsService}
}

func (h *InfoHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	info := v1.Group("/info")
	{
		info.GET("/config", h.GetPublicConfig)
	}
}

func (h *InfoHandler) GetPublicConfig(c *gin.Context) {
	settings, err := h.settingsService.GetSettings(c.Request.Context())
	if err != nil {
		// Fallback to empty settings instead of 500 error
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"config":  map[string]string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"config": settings,
	})
}
