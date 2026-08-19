// Package handler contains Vanguard's Gin HTTP transport layer.
package handler

import (
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/handler/dto"
	"github.com/prateekgupta3991/vanguard/internal/services"
)

// NewHandler creates the Vanguard HTTP router with injected dependencies.
func NewHandler(agentService services.AgentService, logger *slog.Logger) stdhttp.Handler {
	router := gin.New()
	router.Use(requestLogger(logger), recoveryLogger(logger))
	router.HandleMethodNotAllowed = true

	agents := NewAgentHandler(agentService, logger)

	router.GET("/healthz", statusHandler)
	router.GET("/readyz", statusHandler)
	router.POST("/api/v1/agents", agents.Register)
	router.GET("/api/v1/agents", agents.List)
	router.GET("/api/v1/agents/:agentId", agents.Get)
	return router
}

// statusHandler reports process health and service readiness.
func statusHandler(c *gin.Context) {
	respond(c, stdhttp.StatusOK, gin.H{"status": "ok"}, "VANGUARD_HEALTHY", "Vanguard is healthy")
}

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

// respond writes a successful Vanguard response envelope.
func respond(c *gin.Context, status int, data any, code, description string) {
	c.JSON(status, dto.VanguardResponse{Data: data, Code: code, Description: description})
}

// respondError writes an unsuccessful Vanguard response envelope.
func respondError(c *gin.Context, status int, code, description, errorCode, errorDescription string) {
	c.JSON(status, dto.VanguardResponse{
		Code:        code,
		Description: description,
		Error: &dto.VanguardError{
			ErrorCode: errorCode,
			ErrorDesc: errorDescription,
		},
	})
}
