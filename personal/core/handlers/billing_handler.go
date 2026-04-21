package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wacast/core/services/auth"
	"wacast/core/services/billing"
)

type BillingHandler struct {
	billingService *billing.Service
}

func NewBillingHandler(svc *billing.Service) *BillingHandler {
	return &BillingHandler{billingService: svc}
}

func RegisterBillingRoutes(v1 *gin.RouterGroup, svc *billing.Service, jwtSecret string, authService *auth.Service) {
	h := NewBillingHandler(svc)

	protected := v1.Group("/billing")
	protected.Use(JWTAuthMiddleware(jwtSecret, authService))
	{
		protected.GET("/overview", h.GetOverview)
	}
}

func (h *BillingHandler) GetOverview(c *gin.Context) {
	userID := c.GetString(ContextKeyUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errorResponse("UNAUTHORIZED", "Missing user context"))
		return
	}

	overview, err := h.billingService.GetOverview(c.Request.Context(), userID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, successResponse("", gin.H{
		"billing": overview,
	}))
}

