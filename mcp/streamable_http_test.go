package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/berrydev-ai/gqai/graphql"
)

func TestNewStreamableHTTPServer(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	if server == nil {
		t.Error("Expected non-nil server")
	}
	if server.config != config {
		t.Error("Expected server config to match input config")
	}
	if server.sessions == nil {
		t.Error("Expected sessions map to be initialized")
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPWrongMethod(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("PUT", "/mcp", nil)
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPGet(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("GET", "/mcp", nil)
	w := httptest.NewRecorder()

	// This will block because of the infinite loop, so we need to close the connection
	go func() {
		// Simulate client disconnect after a short time
		time.Sleep(10 * time.Millisecond)
		// The test will timeout and the defer will clean up
	}()

	server.HandleStreamableHTTP(w, req)

	// Check that session was created and initial event was sent
	body := w.Body.String()
	if !strings.Contains(body, "event: session") {
		t.Error("Expected session event to be sent")
	}
	if !strings.Contains(body, "http_") {
		t.Error("Expected session ID to be sent")
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPPostMissingSessionID(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Missing Mcp-Session-Id header"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPPostInvalidSession(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader("{}"))
	req.Header.Set("Mcp-Session-Id", "invalid-session")
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedBody := "Invalid session"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPPostInvalidJSON(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	// Add a mock session
	server.sessionsMux.Lock()
	server.sessions["test-session"] = &StreamableHTTPSession{
		sessionID: "test-session",
		responses: make(chan JSONRPCResponse, 10),
		done:      make(chan struct{}),
	}
	server.sessionsMux.Unlock()

	req := httptest.NewRequest("POST", "/mcp", strings.NewReader("invalid json"))
	req.Header.Set("Mcp-Session-Id", "test-session")
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Invalid JSON"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPPostValidRequest(t *testing.T) {
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
		},
	}
	server := NewStreamableHTTPServer(config)

	// Add a mock session
	session := &StreamableHTTPSession{
		sessionID: "test-session",
		responses: make(chan JSONRPCResponse, 10),
		done:      make(chan struct{}),
	}
	server.sessionsMux.Lock()
	server.sessions["test-session"] = session
	server.sessionsMux.Unlock()

	// Create a valid initialize request
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2025-03-26",
		},
	}
	requestBody, _ := json.Marshal(request)

	req := httptest.NewRequest("POST", "/mcp", bytes.NewReader(requestBody))
	req.Header.Set("Mcp-Session-Id", "test-session")
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, w.Code)
	}

	// Check that response was sent to session channel
	select {
	case response := <-session.responses:
		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}
	default:
		t.Error("Expected response to be sent to session channel")
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPDeleteMissingSessionID(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("DELETE", "/mcp", nil)
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Missing Mcp-Session-Id header"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPDeleteValidSession(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	// Add a mock session
	session := &StreamableHTTPSession{
		sessionID: "test-session",
		responses: make(chan JSONRPCResponse, 10),
		done:      make(chan struct{}),
	}
	server.sessionsMux.Lock()
	server.sessions["test-session"] = session
	server.sessionsMux.Unlock()

	req := httptest.NewRequest("DELETE", "/mcp", nil)
	req.Header.Set("Mcp-Session-Id", "test-session")
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Check that session was removed
	server.sessionsMux.RLock()
	_, exists := server.sessions["test-session"]
	server.sessionsMux.RUnlock()

	if exists {
		t.Error("Expected session to be removed")
	}
}

func TestStreamableHTTPServerHandleStreamableHTTPDeleteInvalidSession(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewStreamableHTTPServer(config)

	req := httptest.NewRequest("DELETE", "/mcp", nil)
	req.Header.Set("Mcp-Session-Id", "nonexistent-session")
	w := httptest.NewRecorder()

	server.HandleStreamableHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}
