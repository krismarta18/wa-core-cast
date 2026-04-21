package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/models"
	"wacast/core/services/auth"
	"wacast/core/utils"
)

// ─────────────────────────────────────────────────────────────────────────────
// Handler
// ─────────────────────────────────────────────────────────────────────────────

// AuthHandler exposes HTTP endpoints for user authentication.
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{authService: svc}
}

// RegisterAuthRoutes attaches all auth routes to the supplied router group.
//
//	Public  : POST /auth/register, POST /auth/request-otp, POST /auth/verify-otp
//	Protected: GET /auth/me, POST /auth/logout  (requires JWTAuthMiddleware)
func RegisterAuthRoutes(v1 *gin.RouterGroup, svc *auth.Service, jwtSecret string) {
	h := NewAuthHandler(svc)

	public := v1.Group("/auth")
	{
		public.POST("/register", h.Register)
		public.POST("/request-otp", h.RequestOTP)
		public.POST("/verify-otp", h.VerifyOTP)
		public.POST("/refresh", h.RefreshToken)
	}

	protected := v1.Group("/auth")
	protected.Use(JWTAuthMiddleware(jwtSecret, svc))
	{
		protected.GET("/me", h.Me)
		protected.PUT("/profile", h.UpdateProfile)
		protected.POST("/logout", h.Logout)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Request / Response types
// ─────────────────────────────────────────────────────────────────────────────

type registerRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=8,max=20"`
	FullName    string `json:"full_name"    binding:"required,min=2,max=100"`
}

type requestOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=8,max=20"`
}

type verifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=8,max=20"`
	OTPCode     string `json:"otp_code"     binding:"required,len=6"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /auth/register
// ─────────────────────────────────────────────────────────────────────────────

// Register creates a new user account and sends an OTP to the given phone number.
//
//	Body: { "phone_number": "628xx", "full_name": "...", }
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	if err := h.authService.Register(c.Request.Context(), req.PhoneNumber, req.FullName); err != nil {
		switch {
		case errors.Is(err, auth.ErrPhoneAlreadyRegistered):
			conflict(c, "PHONE_ALREADY_REGISTERED", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("OTP sent — verify your number to complete registration", gin.H{
		"phone_number": req.PhoneNumber,
	}))
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /auth/request-otp
// ─────────────────────────────────────────────────────────────────────────────

// RequestOTP sends a new login OTP to the given registered phone number.
//
//	Body: { "phone_number": "628xx" }
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req requestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	if err := h.authService.RequestOTP(c.Request.Context(), req.PhoneNumber); err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			notFound(c, "USER_NOT_FOUND", err.Error())
		case errors.Is(err, auth.ErrUserBanned):
			forbidden(c, "USER_BANNED", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("OTP sent", gin.H{
		"phone_number": req.PhoneNumber,
	}))
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /auth/verify-otp
// ─────────────────────────────────────────────────────────────────────────────

// VerifyOTP validates the OTP and, on success, returns a signed JWT.
//
//	Body: { "phone_number": "628xx", "otp_code": "123456" }
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	result, err := h.authService.VerifyOTP(c.Request.Context(), req.PhoneNumber, req.OTPCode, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrOTPNotFound):
			notFound(c, "OTP_NOT_FOUND", err.Error())
		case errors.Is(err, auth.ErrOTPExpired):
			unprocessable(c, "OTP_EXPIRED", err.Error())
		case errors.Is(err, auth.ErrOTPInvalid):
			unprocessable(c, "OTP_INVALID", err.Error())
		case errors.Is(err, auth.ErrTooManyAttempts):
			tooManyRequests(c, "TOO_MANY_ATTEMPTS", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("Login successful", gin.H{
		"access_token": result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":   "Bearer",
		"expires_in":   result.ExpiresIn,
		"refresh_expires_in": result.RefreshExpiresIn,
		"user":         result.User.ToResponse(),
	}))
}

// RefreshToken rotates the current refresh token and issues a new access token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshSession(c.Request.Context(), req.RefreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrRefreshTokenInvalid):
			c.JSON(http.StatusUnauthorized, errorResponse("INVALID_REFRESH_TOKEN", err.Error()))
		case errors.Is(err, auth.ErrUserNotFound):
			notFound(c, "USER_NOT_FOUND", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("Token refreshed successfully", gin.H{
		"access_token": result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":   "Bearer",
		"expires_in":   result.ExpiresIn,
		"refresh_expires_in": result.RefreshExpiresIn,
		"user":         result.User.ToResponse(),
	}))
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /auth/me  (protected)
// ─────────────────────────────────────────────────────────────────────────────

// Me returns the currently authenticated user's profile.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString(ContextKeyUserID)

	user, err := h.authService.GetUser(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			notFound(c, "USER_NOT_FOUND", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("", gin.H{
		"user": user.ToResponse(),
	}))
}

// ─────────────────────────────────────────────────────────────────────────────
// PUT /auth/profile (protected)
// ─────────────────────────────────────────────────────────────────────────────

// UpdateProfile updates the currently authenticated user's profile.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString(ContextKeyUserID)

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			notFound(c, "USER_NOT_FOUND", err.Error())
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("Profile updated successfully", gin.H{
		"user": user.ToResponse(),
	}))
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /auth/logout  (protected)
// ─────────────────────────────────────────────────────────────────────────────

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken := c.GetString(ContextKeyAccessToken)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, errorResponse("MISSING_TOKEN", "Authorization token is required"))
		return
	}

	if err := h.authService.Logout(c.Request.Context(), accessToken); err != nil {
		switch {
		case errors.Is(err, auth.ErrSessionNotFound):
			c.JSON(http.StatusUnauthorized, errorResponse("SESSION_NOT_FOUND", err.Error()))
		default:
			internalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, successResponse("Logged out successfully", nil))
}

// ─────────────────────────────────────────────────────────────────────────────
// Response helpers
// ─────────────────────────────────────────────────────────────────────────────

func successResponse(message string, data gin.H) gin.H {
	r := gin.H{"success": true}
	if message != "" {
		r["message"] = message
	}
	if data != nil {
		for k, v := range data {
			r[k] = v
		}
	}
	return r
}

func errorResponse(code, message string) gin.H {
	return gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	}
}

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, errorResponse("VALIDATION_ERROR", msg))
}

func notFound(c *gin.Context, code, msg string) {
	c.JSON(http.StatusNotFound, errorResponse(code, msg))
}

func conflict(c *gin.Context, code, msg string) {
	c.JSON(http.StatusConflict, errorResponse(code, msg))
}

func forbidden(c *gin.Context, code, msg string) {
	c.JSON(http.StatusForbidden, errorResponse(code, msg))
}

func unprocessable(c *gin.Context, code, msg string) {
	c.JSON(http.StatusUnprocessableEntity, errorResponse(code, msg))
}

func tooManyRequests(c *gin.Context, code, msg string) {
	c.JSON(http.StatusTooManyRequests, errorResponse(code, msg))
}

func internalError(c *gin.Context, err error) {
	utils.Error("auth handler internal error",
		zap.String("path", c.FullPath()),
		zap.String("method", c.Request.Method),
		zap.Error(err),
	)

	resp := errorResponse("INTERNAL_ERROR", "an internal error occurred")
	if gin.Mode() != gin.ReleaseMode {
		resp["error"].(gin.H)["detail"] = err.Error()
	}

	c.JSON(http.StatusInternalServerError, resp)
	_ = c.Error(err)
}
