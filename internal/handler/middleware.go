package handler

import (
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// requestLogger records each completed HTTP request at info level.
func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()
		logger.Info("HTTP request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	}
}

// recoveryLogger recovers panics and records them as server errors.
func recoveryLogger(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered while handling HTTP request", "panic", recovered)
		respondError(c, stdhttp.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", "INTERNAL_ERROR", "Internal server error")
	})
}
