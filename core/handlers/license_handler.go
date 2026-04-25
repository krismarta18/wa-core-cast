package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wacast/core/services/license"
)

type LicenseHandler struct {
	service *license.Service
}

func NewLicenseHandler(s *license.Service) *LicenseHandler {
	return &LicenseHandler{service: s}
}

func (h *LicenseHandler) RegisterRoutes(r *gin.RouterGroup) {
	licenseGroup := r.Group("/license")
	{
		licenseGroup.GET("/status", h.GetStatus)
		licenseGroup.POST("/activate", h.Activate)
	}
}

func (h *LicenseHandler) GetStatus(c *gin.Context) {
	status, err := h.service.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    status,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

func (h *LicenseHandler) Activate(c *gin.Context) {
	var req struct {
		Key string `json:"key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "License key is required"})
		return
	}

	if err := h.service.Activate(req.Key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "License activated successfully",
	})
}
