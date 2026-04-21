package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wacast/core/services/analytics"
	"wacast/core/services/auth"
)

type AnalyticHandler struct {
	analyticService *analytics.Service
}

func NewAnalyticHandler(svc *analytics.Service) *AnalyticHandler {
	return &AnalyticHandler{analyticService: svc}
}

func RegisterAnalyticRoutes(v1 *gin.RouterGroup, svc *analytics.Service, jwtSecret string, authService *auth.Service) {
	h := NewAnalyticHandler(svc)

	protected := v1.Group("/analytics")
	protected.Use(JWTAuthMiddleware(jwtSecret, authService))
	{
		protected.GET("/usage", h.GetUsageStats)
		protected.GET("/failures", h.GetFailureAnalytics)
	}
}

// GetUsageStats handles usage statistics requests
// GET /api/v1/analytics/usage
func (h *AnalyticHandler) GetUsageStats(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.analyticService.GetUsageStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetFailureAnalytics handles failure rate analytic requests
// GET /api/v1/analytics/failures
func (h *AnalyticHandler) GetFailureAnalytics(c *gin.Context) {
	userIDStr := c.GetString(ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.analyticService.GetFailureAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
