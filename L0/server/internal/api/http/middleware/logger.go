package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	const mark = "http.middleware.Logger"

	return func(c *gin.Context) {
		logger := slog.With(slog.String("mark", mark))
		start := time.Now()

		// Log request start
		logger.Debug("Request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		// Process request
		c.Next()

		// Gather metrics
		duration := time.Since(start).String()
		status := c.Writer.Status()

		// Form basic log attributes
		logArgs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"duration", duration,
			"status", status,
			"ip", c.ClientIP(),
		}

		// Check for errors
		errs := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errs) > 0 {
			// Log all errors
			for _, err := range errs {
				logger.Error("Request failed", append(logArgs, "error", err.Error())...)
			}
			return
		}

		// Log successful request completion based on status level
		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("Server error", logArgs...)
		case status >= http.StatusBadRequest:
			logger.Warn("Client error", logArgs...)
		default:
			logger.Info("Request completed", logArgs...)
		}
	}
}
