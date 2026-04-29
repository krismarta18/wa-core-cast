package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wacast/core/services/auth"
	"wacast/core/services/integration"
	"wacast/core/utils"
)

// Context keys injected by JWTAuthMiddleware.
const (
	ContextKeyUserID      = "userID"
	ContextKeyPhoneNumber = "phoneNumber"
	ContextKeyFullName    = "fullName"
	ContextKeyAccessToken = "accessToken"
)

// JWTAuthMiddleware validates the Bearer token in the Authorization header.
// On success it injects user claims into the gin context.
// On failure it aborts with 401 Unauthorized.
func JWTAuthMiddleware(jwtSecret string, authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header must be: Bearer <token>",
				},
			})
			return
		}

		claims, err := utils.ValidateJWT(parts[1], jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Token is invalid or expired",
				},
			})
			return
		}

		if err := authService.ValidateSession(c.Request.Context(), parts[1]); err != nil {
			if errors.Is(err, auth.ErrSessionNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "SESSION_REVOKED",
						"message": "Session is no longer active",
					},
				})
			} else {
				// Internal error (e.g. database locked). Return 500 so frontend doesn't clear session.
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error during session validation",
					},
				})
			}
			return
		}

		// Inject claims into context for downstream handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyPhoneNumber, claims.PhoneNumber)
		c.Set(ContextKeyFullName, claims.FullName)
		c.Set(ContextKeyAccessToken, parts[1])

		c.Next()
	}
}

// IntegrationAuthMiddleware allows either JWT or API Key authentication.
// Priority: JWT (Authorization Bearer) first, then API Key (X-API-Key).
func IntegrationAuthMiddleware(jwtSecret string, authService *auth.Service, integrationService *integration.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Try JWT first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				claims, err := utils.ValidateJWT(parts[1], jwtSecret)
				if err == nil {
					if err := authService.ValidateSession(c.Request.Context(), parts[1]); err == nil {
						c.Set(ContextKeyUserID, claims.UserID)
						c.Set(ContextKeyAccessToken, parts[1])
						c.Next()
						return
					}
				}
			}
		}

		// 2. Try API Key
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Also check Bearer in case user uses it for API Key (common)
			if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				apiKey = strings.SplitN(authHeader, " ", 2)[1]
			}
		}

		if apiKey != "" {
			userID, err := integrationService.ValidateAPIKey(apiKey)
			if err == nil {
				c.Set(ContextKeyUserID, userID.String())
				c.Next()
				return
			}
		}

		// 3. Fallback to 401
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Valid JWT or API Key (X-API-Key) is required",
			},
		})
	}
}
