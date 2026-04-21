package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"wacast/core/services/message"
	"wacast/core/services/settings"
	"wacast/core/utils"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type SettingHandler struct {
	settingsService *settings.Service
	messageService  *message.Service
}

func NewSettingHandler(settingsService *settings.Service, messageService *message.Service) *SettingHandler {
	return &SettingHandler{
		settingsService: settingsService,
		messageService:  messageService,
	}
}

func (h *SettingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	settings := rg.Group("/settings")
	{
		settings.GET("", h.GetAllSettings)
		settings.PUT("/:key", h.UpdateSetting)
	}
}

func (h *SettingHandler) GetAllSettings(c *gin.Context) {
	data, err := h.settingsService.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare as a list of SystemSetting-like objects for the UI
	type entry struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	list := make([]entry, 0, len(data))
	for k, v := range data {
		list = append(list, entry{Key: k, Value: v})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"settings": list,
	})
}

func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var input struct {
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// 1. Update in Database
	if err := h.settingsService.UpdateSetting(c.Request.Context(), key, input.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting in database"})
		return
	}

	// 2. Sync to Message Service if it's an Anti-Bot setting
	if strings.HasPrefix(key, "anti_bot_") {
		h.syncAntiBotSettings(key, input.Value)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Setting %s updated successfully", key),
	})
}

func (h *SettingHandler) syncAntiBotSettings(key, value string) {
	config := h.messageService.GetConfig()

	switch key {
	case "anti_bot_enabled":
		config.AntiBotEnabled = (strings.ToLower(value) == "true")
	case "anti_bot_suffix_length":
		var length int
		if _, err := fmt.Sscanf(value, "%d", &length); err == nil {
			config.RandomSuffixLength = length
		}
	}

	h.messageService.UpdateConfig(&config)
	utils.Info("Anti-Bot configuration synced in-memory",
		zap.String("key", key),
		zap.String("value", value),
	)
}
