package config

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultHTTPAddress     = ":8081"
	defaultShutdownTimeout = 10 * time.Second
)

// Config contains runtime configuration for the Vanguard server.
type Config struct {
	HTTPAddress     string
	ShutdownTimeout time.Duration
}

// Load reads and validates configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddress:     valueOrDefault("VANGUARD_HTTP_ADDRESS", defaultHTTPAddress),
		ShutdownTimeout: defaultShutdownTimeout,
	}

	if raw := os.Getenv("VANGUARD_SHUTDOWN_TIMEOUT"); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse VANGUARD_SHUTDOWN_TIMEOUT: %w", err)
		}
		if timeout <= 0 {
			return Config{}, fmt.Errorf("VANGUARD_SHUTDOWN_TIMEOUT must be positive")
		}
		cfg.ShutdownTimeout = timeout
	}

	return cfg, nil
}

// valueOrDefault returns an environment value or its fallback.
func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
