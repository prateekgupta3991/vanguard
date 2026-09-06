package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("VANGUARD_HTTP_ADDRESS", "")
	t.Setenv("VANGUARD_SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.HTTPAddress != ":8081" {
		t.Fatalf("HTTPAddress = %q, want %q", cfg.HTTPAddress, ":8081")
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %v, want %v", cfg.ShutdownTimeout, 10*time.Second)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("VANGUARD_HTTP_ADDRESS", "127.0.0.1:9090")
	t.Setenv("VANGUARD_SHUTDOWN_TIMEOUT", "3s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.HTTPAddress != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddress = %q, want %q", cfg.HTTPAddress, "127.0.0.1:9090")
	}
	if cfg.ShutdownTimeout != 3*time.Second {
		t.Fatalf("ShutdownTimeout = %v, want %v", cfg.ShutdownTimeout, 3*time.Second)
	}
}

func TestLoadRejectsInvalidShutdownTimeout(t *testing.T) {
	t.Setenv("VANGUARD_SHUTDOWN_TIMEOUT", "never")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
}
