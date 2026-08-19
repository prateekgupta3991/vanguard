package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prateekgupta3991/vanguard/internal/handler/dto"
	"github.com/prateekgupta3991/vanguard/internal/repository"
	"github.com/prateekgupta3991/vanguard/internal/services"
)

func TestAgentRegistrationLifecycle(t *testing.T) {
	handler := newTestHandler()

	created := performRequest(t, handler, http.MethodPost, "/api/v1/agents", `{
		"identifier":"codex:developer-laptop",
		"name":"Codex",
		"framework":"codex",
		"metadata":{"environment":"development"}
	}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", created.Code, http.StatusCreated, created.Body.String())
	}
	createdAgent := responseData(t, created)
	createdID, ok := createdAgent["id"].(string)
	if !ok || createdID == "" {
		t.Fatalf("created agent ID = %#v, want a non-empty string", createdAgent["id"])
	}

	updated := performRequest(t, handler, http.MethodPost, "/api/v1/agents", `{
		"identifier":"codex:developer-laptop",
		"name":"Codex CLI",
		"framework":"codex"
	}`)
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d; body = %s", updated.Code, http.StatusOK, updated.Body.String())
	}
	updatedAgent := responseData(t, updated)
	if updatedAgent["id"] != createdID {
		t.Fatalf("updated agent ID = %#v, want %q", updatedAgent["id"], createdID)
	}
	if updatedAgent["name"] != "Codex CLI" {
		t.Fatalf("updated agent name = %#v, want Codex CLI", updatedAgent["name"])
	}

	fetched := performRequest(t, handler, http.MethodGet, "/api/v1/agents/"+createdID, "")
	if fetched.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d; body = %s", fetched.Code, http.StatusOK, fetched.Body.String())
	}

	listed := performRequest(t, handler, http.MethodGet, "/api/v1/agents", "")
	if listed.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d; body = %s", listed.Code, http.StatusOK, listed.Body.String())
	}
	var listResponse struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(listed.Body.Bytes(), &listResponse); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listResponse.Data) != 1 {
		t.Fatalf("listed agents = %d, want 1", len(listResponse.Data))
	}
}

func TestRegisterAgentRejectsInvalidRequest(t *testing.T) {
	response := performRequest(t, newTestHandler(), http.MethodPost, "/api/v1/agents", `{"name":"Codex"}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}

	var envelope dto.VanguardResponse
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if envelope.Error == nil || envelope.Error.ErrorCode != "VALIDATION_ERROR" {
		t.Fatalf("error = %#v, want VALIDATION_ERROR", envelope.Error)
	}
}

func TestGetAgentReturnsNotFound(t *testing.T) {
	response := performRequest(t, newTestHandler(), http.MethodGet, "/api/v1/agents/unknown", "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func newTestHandler() http.Handler {
	service := services.NewAgentService(repository.NewMemoryRepository())
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return NewHandler(service, logger)
}

func performRequest(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func responseData(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return envelope.Data
}
