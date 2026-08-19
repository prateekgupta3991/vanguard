// Package repository contains persistence contracts and implementations.
package repository

import (
	"context"
	"errors"

	"github.com/prateekgupta3991/vanguard/internal/entities"
)

// ErrNotFound indicates that an agent does not exist in the registry.
var ErrNotFound = errors.New("agent not found")

// AgentRepository defines persistence operations required by the registry.
type AgentRepository interface {
	// Register creates an agent or refreshes an existing registration.
	Register(ctx context.Context, candidate entities.Agent) (registered entities.Agent, created bool, err error)
	// Get retrieves an agent by its Vanguard-assigned ID.
	Get(ctx context.Context, id string) (entities.Agent, error)
	// List retrieves all registered agents.
	List(ctx context.Context) ([]entities.Agent, error)
}
