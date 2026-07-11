package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware returns a Gin middleware that logs each HTTP request using structured slog.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process the request
		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Grab userID from context if set by AuthMiddleware
		userID, _ := c.Get("userID")

		// Build log attrs
		attrs := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
		}

		if userID != nil {
			attrs = append(attrs, slog.Any("user_id", userID))
		}

		// Log errors attached to the context
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				attrs = append(attrs, slog.String("error", e.Error()))
			}
			slog.Error("Request selesai dengan error", attrs...)
			return
		}

		// Log level based on status code
		if statusCode >= 500 {
			slog.Error("Request selesai", attrs...)
		} else if statusCode >= 400 {
			slog.Warn("Request selesai", attrs...)
		} else {
			slog.Info("Request selesai", attrs...)
		}
	}
}
