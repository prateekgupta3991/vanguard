// Package handler contains Vanguard's Gin HTTP transport layer.
package handler

import (
	"log/slog"
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/services"
)

// NewRouter creates the Vanguard HTTP router and registers every API route.
func NewRouter(agentService services.AgentService, logger *slog.Logger) stdhttp.Handler {
	router := gin.New()
	router.Use(requestLogger(logger), recoveryLogger(logger))
	router.HandleMethodNotAllowed = true

	agents := NewAgentHandler(agentService, logger)

	router.GET("/healthz", statusHandler)
	router.GET("/readyz", statusHandler)

	apiV1 := router.Group("/api/v1")
	apiV1.POST("/agents", agents.Register)
	apiV1.GET("/agents", agents.List)
	apiV1.GET("/agents/:agentId", agents.Get)

	return router
}
