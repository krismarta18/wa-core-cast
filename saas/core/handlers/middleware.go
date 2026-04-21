package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wacast/core/services/auth"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "SESSION_REVOKED",
					"message": "Session is no longer active",
				},
			})
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
