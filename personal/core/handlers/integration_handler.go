package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"wacast/core/models"
	"wacast/core/services/auth"
	"wacast/core/services/integration"
)

type IntegrationHandler struct {
	integrationService *integration.Service
	jwtSecret          string
	authService        *auth.Service
}

func NewIntegrationHandler(svc *integration.Service, jwtSecret string, authService *auth.Service) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService: svc,
		jwtSecret:          jwtSecret,
		authService:        authService,
	}
}

func (h *IntegrationHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	integration := v1.Group("/integration")
	integration.Use(JWTAuthMiddleware(h.jwtSecret, h.authService))
	{
		integration.GET("/keys", h.ListAPIKeys)
		integration.POST("/keys", h.CreateAPIKey)
		integration.DELETE("/keys/:id", h.DeleteAPIKey)
		integration.GET("/webhook", h.GetWebhookSettings)
		integration.PUT("/webhook", h.UpdateWebhookSettings)
	}
}

func (h *IntegrationHandler) ListAPIKeys(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	keys, err := h.integrationService.ListAPIKeys(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar API Key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func (h *IntegrationHandler) CreateAPIKey(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama Key wajib diisi"})
		return
	}

	response, err := h.integrationService.CreateAPIKey(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat API Key"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *IntegrationHandler) DeleteAPIKey(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)
	idStr := c.Param("id")
	id, _ := uuid.Parse(idStr)

	if err := h.integrationService.DeleteAPIKey(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus API Key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API Key berhasil dihapus"})
}

func (h *IntegrationHandler) GetWebhookSettings(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	settings, err := h.integrationService.GetWebhookSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pengaturan Webhook"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *IntegrationHandler) UpdateWebhookSettings(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, _ := uuid.Parse(userIDStr)

	var settings models.WebhookSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data pengaturan tidak valid"})
		return
	}

	if err := h.integrationService.UpdateWebhookSettings(userID, &settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan Webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan Webhook berhasil diperbarui"})
}
