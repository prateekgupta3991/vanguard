package repository

import (
	"context"
	"sort"
	"sync"

	"github.com/prateekgupta3991/vanguard/internal/entities"
)

// MemoryRepository stores agents in memory and is safe for concurrent use.
type MemoryRepository struct {
	mu           sync.RWMutex
	byID         map[string]entities.Agent
	byIdentifier map[string]string
}

var _ AgentRepository = (*MemoryRepository)(nil)

// NewMemoryRepository creates an empty in-memory agent repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		byID:         make(map[string]entities.Agent),
		byIdentifier: make(map[string]string),
	}
}

// Register creates an agent or refreshes the agent with the same identifier.
func (r *MemoryRepository) Register(_ context.Context, candidate entities.Agent) (entities.Agent, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id, exists := r.byIdentifier[candidate.Identifier]; exists {
		registered := r.byID[id]
		registered.Name = candidate.Name
		registered.Framework = candidate.Framework
		registered.Metadata = cloneMetadata(candidate.Metadata)
		registered.LastSeenAt = candidate.LastSeenAt
		r.byID[id] = registered
		return cloneAgent(registered), false, nil
	}

	candidate.Metadata = cloneMetadata(candidate.Metadata)
	r.byID[candidate.ID] = candidate
	r.byIdentifier[candidate.Identifier] = candidate.ID
	return cloneAgent(candidate), true, nil
}

// Get returns an agent by its Vanguard-assigned ID.
func (r *MemoryRepository) Get(_ context.Context, id string) (entities.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	registered, exists := r.byID[id]
	if !exists {
		return entities.Agent{}, ErrNotFound
	}
	return cloneAgent(registered), nil
}

// List returns all registered agents ordered by their first-seen time.
func (r *MemoryRepository) List(_ context.Context) ([]entities.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]entities.Agent, 0, len(r.byID))
	for _, registered := range r.byID {
		agents = append(agents, cloneAgent(registered))
	}
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].FirstSeenAt.Before(agents[j].FirstSeenAt)
	})
	return agents, nil
}

// cloneAgent copies an agent so callers cannot mutate repository state.
func cloneAgent(source entities.Agent) entities.Agent {
	source.Metadata = cloneMetadata(source.Metadata)
	return source
}

// cloneMetadata copies the top-level metadata map.
func cloneMetadata(source map[string]any) map[string]any {
	if source == nil {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}
