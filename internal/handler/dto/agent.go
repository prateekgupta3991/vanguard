package dto

// RegisterAgentRequest contains the client-provided agent identity and metadata.
type RegisterAgentRequest struct {
	Identifier string         `json:"identifier" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Framework  string         `json:"framework"`
	Metadata   map[string]any `json:"metadata"`
}
