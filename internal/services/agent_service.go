// Package services contains Vanguard's application business logic.
package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/prateekgupta3991/vanguard/internal/entities"
	"github.com/prateekgupta3991/vanguard/internal/repository"
)

// AgentService defines the agent-registry business operations used by handlers.
type AgentService interface {
	// Register creates or refreshes an agent registration.
	Register(ctx context.Context, registration entities.Registration) (registered entities.Agent, created bool, err error)
	// Get retrieves an agent by its Vanguard-assigned ID.
	Get(ctx context.Context, id string) (entities.Agent, error)
	// List retrieves all known agents.
	List(ctx context.Context) ([]entities.Agent, error)
}

// RegistryService implements AgentService using an injected repository.
type RegistryService struct {
	repository repository.AgentRepository
	now        func() time.Time
	newID      func() (string, error)
}

var _ AgentService = (*RegistryService)(nil)

// NewAgentService creates an agent registry service with the supplied repository.
func NewAgentService(agentRepository repository.AgentRepository) *RegistryService {
	return &RegistryService{
		repository: agentRepository,
		now:        func() time.Time { return time.Now().UTC() },
		newID:      newUUID,
	}
}

// Register creates or refreshes an agent registration.
func (s *RegistryService) Register(ctx context.Context, registration entities.Registration) (entities.Agent, bool, error) {
	id, err := s.newID()
	if err != nil {
		return entities.Agent{}, false, fmt.Errorf("generate agent ID: %w", err)
	}

	now := s.now()
	return s.repository.Register(ctx, entities.Agent{
		ID:          id,
		Identifier:  registration.Identifier,
		Name:        registration.Name,
		Framework:   registration.Framework,
		Metadata:    registration.Metadata,
		FirstSeenAt: now,
		LastSeenAt:  now,
	})
}

// Get retrieves an agent by its Vanguard-assigned ID.
func (s *RegistryService) Get(ctx context.Context, id string) (entities.Agent, error) {
	return s.repository.Get(ctx, id)
}

// List retrieves all known agents.
func (s *RegistryService) List(ctx context.Context) ([]entities.Agent, error) {
	return s.repository.List(ctx)
}

// newUUID generates an RFC 4122 version 4 UUID.
func newUUID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		value[0:4], value[4:6], value[6:8], value[8:10], value[10:16]), nil
}
