package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/config"
	"github.com/prateekgupta3991/vanguard/internal/handler"
	"github.com/prateekgupta3991/vanguard/internal/repository"
	"github.com/prateekgupta3991/vanguard/internal/services"
)

// main composes Vanguard's dependencies and runs the HTTP server.
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	agentRepository := repository.NewMemoryRepository()
	agentService := services.NewAgentService(agentRepository)
	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           handler.NewRouter(agentService, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownSignal.Done()
		logger.Info("shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("HTTP server shutdown failed", "error", err)
		}
	}()

	logger.Info("starting Vanguard", "address", cfg.HTTPAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}
