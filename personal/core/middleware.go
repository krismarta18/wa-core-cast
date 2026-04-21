package main

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/utils"
)

// GinLogger returns a Gin middleware that logs requests using Zap
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Call the next handler
		c.Next()

		// Calculate request latency
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// Skip logging for health checks
		if path == "/health" || path == "/health/ready" || path == "/health/live" {
			return
		}

		// Determine log level based on status code
		var fields []zap.Field
		switch {
		case statusCode >= 400 && statusCode < 500:
			fields = append(fields,
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			)
			utils.Warn("HTTP request", fields...)
		case statusCode >= 500:
			if errorMsg, exists := c.Get("error_message"); exists {
				fields = append(fields, zap.String("error", errorMsg.(string)))
			}
			fields = append(fields,
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			)
			utils.Error("HTTP request", fields...)
		default:
			utils.Debug("HTTP request",
				zap.String("client_ip", clientIP),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
			)
		}
	}
}

// ErrorHandler returns a Gin middleware that handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				utils.Error("Panic recovered",
					zap.String("path", c.Request.URL.Path),
					zap.String("error", formatPanicError(err)),
				)

				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}

// NoOpLogger returns a Gin middleware that discards log output
func NoOpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Discard Gin's default logging
		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("REQUEST_ID")
		if requestID == "" {
			// Generate a new request ID if not already set
			requestID = generateRequestID()
			c.Set("REQUEST_ID", requestID)
		}

		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// CORS middleware for cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware applies rate limiting (placeholder for future implementation)
func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	// TODO: Implement rate limiting with token bucket or sliding window
	return func(c *gin.Context) {
		c.Next()
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// formatPanicError converts a panic error to a string
func formatPanicError(err interface{}) string {
	switch v := err.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "unknown panic"
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple implementation - can be enhanced with UUIDs
	return time.Now().Format("20060102150405") + "-" + generateRandomString(6)
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// LoggerConfig wraps Zap's logger to be compatible with Gin
type LoggerWriter struct {
	logger *zap.Logger
}

// Write implements io.Writer for Gin's logger
func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	lw.logger.Info(string(p))
	return len(p), nil
}

// DisinfectOutput returns an io.Writer that discards all output
func DisinfectOutput() io.Writer {
	return io.Discard
}
