// Package entities contains Vanguard's core domain models and business logic.
package entities

import "time"

// Agent represents an agent known to the Vanguard control plane.
type Agent struct {
	ID          string         `json:"id"`
	Identifier  string         `json:"identifier"`
	Name        string         `json:"name"`
	Framework   string         `json:"framework"`
	Metadata    map[string]any `json:"metadata"`
	FirstSeenAt time.Time      `json:"firstSeenAt"`
	LastSeenAt  time.Time      `json:"lastSeenAt"`
}

// Registration contains the agent identity observed by an enforcement point.
type Registration struct {
	Identifier string
	Name       string
	Framework  string
	Metadata   map[string]any
}
