package handler

import (
	"errors"
	"log/slog"
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/entities"
	"github.com/prateekgupta3991/vanguard/internal/handler/dto"
	"github.com/prateekgupta3991/vanguard/internal/repository"
	"github.com/prateekgupta3991/vanguard/internal/services"
)

// AgentHandler defines the HTTP operations exposed by the agent registry.
type AgentHandler interface {
	// Register handles agent registration and refresh requests.
	Register(c *gin.Context)
	// Get handles retrieval of one agent.
	Get(c *gin.Context)
	// List handles retrieval of all agents.
	List(c *gin.Context)
}

// agentHandler implements AgentHandler with injected business logic and logging.
type agentHandler struct {
	service services.AgentService
	logger  *slog.Logger
}

var _ AgentHandler = (*agentHandler)(nil)

// NewAgentHandler creates an agent HTTP handler with injected dependencies.
func NewAgentHandler(service services.AgentService, logger *slog.Logger) AgentHandler {
	return &agentHandler{service: service, logger: logger}
}

// Register validates and handles an agent registration request.
func (h *agentHandler) Register(c *gin.Context) {
	var request dto.RegisterAgentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Info("agent registration rejected", "reason", err.Error())
		respondError(c, stdhttp.StatusBadRequest, "INVALID_REQUEST", "Agent registration request is invalid", "VALIDATION_ERROR", err.Error())
		return
	}

	registered, created, err := h.service.Register(c.Request.Context(), entities.Registration{
		Identifier: request.Identifier,
		Name:       request.Name,
		Framework:  request.Framework,
		Metadata:   request.Metadata,
	})
	if err != nil {
		h.logger.Error("agent registration failed", "identifier", request.Identifier, "error", err)
		respondError(c, stdhttp.StatusInternalServerError, "AGENT_REGISTRATION_FAILED", "Agent registration failed", "INTERNAL_ERROR", err.Error())
		return
	}

	status := stdhttp.StatusOK
	code := "AGENT_UPDATED"
	description := "Agent registration updated"
	if created {
		status = stdhttp.StatusCreated
		code = "AGENT_REGISTERED"
		description = "Agent registered successfully"
	}
	h.logger.Info("agent registration processed", "agent_id", registered.ID, "identifier", registered.Identifier, "created", created)
	respond(c, status, registered, code, description)
}

// Get handles retrieval of an agent by ID.
func (h *agentHandler) Get(c *gin.Context) {
	registered, err := h.service.Get(c.Request.Context(), c.Param("agentId"))
	if errors.Is(err, repository.ErrNotFound) {
		h.logger.Info("agent lookup returned no result", "agent_id", c.Param("agentId"))
		respondError(c, stdhttp.StatusNotFound, "AGENT_NOT_FOUND", "Agent was not found", "AGENT_NOT_FOUND", "No agent exists with the supplied ID")
		return
	}
	if err != nil {
		h.logger.Error("agent lookup failed", "agent_id", c.Param("agentId"), "error", err)
		respondError(c, stdhttp.StatusInternalServerError, "AGENT_LOOKUP_FAILED", "Agent lookup failed", "INTERNAL_ERROR", err.Error())
		return
	}
	h.logger.Info("agent retrieved", "agent_id", registered.ID)
	respond(c, stdhttp.StatusOK, registered, "AGENT_FOUND", "Agent retrieved successfully")
}

// List handles retrieval of all registered agents.
func (h *agentHandler) List(c *gin.Context) {
	agents, err := h.service.List(c.Request.Context())
	if err != nil {
		h.logger.Error("agent listing failed", "error", err)
		respondError(c, stdhttp.StatusInternalServerError, "AGENT_LIST_FAILED", "Agent listing failed", "INTERNAL_ERROR", err.Error())
		return
	}
	h.logger.Info("agents listed", "count", len(agents))
	respond(c, stdhttp.StatusOK, agents, "AGENTS_LISTED", "Agents retrieved successfully")
}
