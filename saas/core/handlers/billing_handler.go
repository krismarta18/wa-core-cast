package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"wacast/core/models"
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
		protected.POST("/checkout", h.CheckoutDummy)
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

func (h *BillingHandler) CheckoutDummy(c *gin.Context) {
	userID := c.GetString(ContextKeyUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errorResponse("UNAUTHORIZED", "Missing user context"))
		return
	}

	var req models.BillingCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	result, err := h.billingService.CheckoutDummy(c.Request.Context(), userID, req.PlanID)
	if err != nil {
		switch {
		case errors.Is(err, billing.ErrBillingPlanNotFound):
			notFound(c, "BILLING_PLAN_NOT_FOUND", err.Error())
		case errors.Is(err, billing.ErrBillingPlanInactive):
			unprocessable(c, "BILLING_PLAN_INACTIVE", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("Dummy payment completed and subscription activated", gin.H{
		"checkout": result,
	}))
}
