package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wacast/core/services/auth"
	"wacast/core/services/autoresponse"
)

type AutoResponseHandler struct {
	svc *autoresponse.Service
}

func NewAutoResponseHandler(svc *autoresponse.Service) *AutoResponseHandler {
	return &AutoResponseHandler{svc: svc}
}

func RegisterAutoResponseRoutes(router *gin.RouterGroup, svc *autoresponse.Service, jwtSecret string, authSvc *auth.Service) {
	h := NewAutoResponseHandler(svc)
	
	arGroup := router.Group("/auto-response")
	arGroup.Use(JWTAuthMiddleware(jwtSecret, authSvc))
	{
		arGroup.GET("/keywords", h.GetKeywords)
		arGroup.POST("/keywords", h.CreateKeyword)
		arGroup.PUT("/keywords/:id", h.UpdateKeyword)
		arGroup.DELETE("/keywords/:id", h.DeleteKeyword)
		arGroup.PATCH("/keywords/:id/toggle", h.ToggleKeyword)

		arGroup.GET("/templates", h.GetTemplates)
		arGroup.POST("/templates", h.CreateTemplate)
		arGroup.PUT("/templates/:id", h.UpdateTemplate)
		arGroup.DELETE("/templates/:id", h.DeleteTemplate)
	}
}

// GetKeywords godoc
func (h *AutoResponseHandler) GetKeywords(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	keywords, err := h.svc.GetKeywords(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keywords: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keywords": keywords})
}

// CreateKeyword godoc
func (h *AutoResponseHandler) CreateKeyword(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	var req autoresponse.CreateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	kw, err := h.svc.CreateKeyword(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keyword: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, kw)
}

// UpdateKeyword godoc
func (h *AutoResponseHandler) UpdateKeyword(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID: " + err.Error()})
		return
	}

	var req autoresponse.UpdateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	kw, err := h.svc.UpdateKeyword(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update keyword: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, kw)
}

// DeleteKeyword godoc
func (h *AutoResponseHandler) DeleteKeyword(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID: " + err.Error()})
		return
	}

	if err := h.svc.DeleteKeyword(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete keyword: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Keyword deleted successfully"})
}

// ToggleKeyword godoc
func (h *AutoResponseHandler) ToggleKeyword(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID: " + err.Error()})
		return
	}

	kw, err := h.svc.ToggleKeyword(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle keyword: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, kw)
}

// -- Templates --

// GetTemplates godoc
func (h *AutoResponseHandler) GetTemplates(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	templates, err := h.svc.GetTemplates(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// CreateTemplate godoc
func (h *AutoResponseHandler) CreateTemplate(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	var req autoresponse.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	tmpl, err := h.svc.CreateTemplate(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tmpl)
}

// UpdateTemplate godoc
func (h *AutoResponseHandler) UpdateTemplate(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID: " + err.Error()})
		return
	}

	var req autoresponse.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	tmpl, err := h.svc.UpdateTemplate(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tmpl)
}

// DeleteTemplate godoc
func (h *AutoResponseHandler) DeleteTemplate(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID: " + err.Error()})
		return
	}

	if err := h.svc.DeleteTemplate(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}
