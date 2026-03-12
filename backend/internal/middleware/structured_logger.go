package middleware

import (
	"time"

	"github.com/devsherkhane/trello-clone/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StructuredLogger logs details of HTTP requests using zap
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Create fields for logger
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		// Log errors if any occurred during the request cycle
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
			logger.Log.Error("HTTP request failed", fields...)
		} else {
			// Skip logging healthy health checks to avoid spanning the logs
			if path == "/health" && c.Writer.Status() == 200 {
				return
			}
			logger.Log.Info("HTTP request", fields...)
		}
	}
}
