package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berrydev-ai/gqai/graphql"
)

func TestNewSSEServer(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewSSEServer(config)

	if server == nil {
		t.Error("Expected non-nil server")
	}
	if server.config != config {
		t.Error("Expected server config to match input config")
	}
	if server.clients == nil {
		t.Error("Expected clients map to be initialized")
	}
}

func TestSSEServerHandleMessageWrongMethod(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewSSEServer(config)

	req := httptest.NewRequest("GET", "/message?sessionId=test", nil)
	w := httptest.NewRecorder()

	server.HandleMessage(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestSSEServerHandleMessageMissingSessionID(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewSSEServer(config)

	req := httptest.NewRequest("POST", "/message", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	server.HandleMessage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Missing sessionId parameter"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestSSEServerHandleMessageInvalidSession(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewSSEServer(config)

	req := httptest.NewRequest("POST", "/message?sessionId=invalid", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	server.HandleMessage(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedBody := "Invalid session"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestSSEServerHandleMessageInvalidJSON(t *testing.T) {
	config := &graphql.GraphQLConfig{}
	server := NewSSEServer(config)

	// Add a mock client
	server.clientsMux.Lock()
	server.clients["test-session"] = &SSEClient{
		sessionID: "test-session",
		writer:    &mockResponseWriter{},
		done:      make(chan struct{}),
	}
	server.clientsMux.Unlock()

	req := httptest.NewRequest("POST", "/message?sessionId=test-session", strings.NewReader("invalid json"))
	w := httptest.NewRecorder()

	server.HandleMessage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Invalid JSON"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestSSEServerHandleMessageValidRequest(t *testing.T) {
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
		},
	}
	server := NewSSEServer(config)

	// Create a mock writer that captures what gets written
	var capturedData string
	mockWriter := &mockResponseWriter{
		writeFunc: func(data []byte) (int, error) {
			capturedData += string(data)
			return len(data), nil
		},
	}

	// Add a mock client
	server.clientsMux.Lock()
	server.clients["test-session"] = &SSEClient{
		sessionID: "test-session",
		writer:    mockWriter,
		done:      make(chan struct{}),
	}
	server.clientsMux.Unlock()

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

	req := httptest.NewRequest("POST", "/message?sessionId=test-session", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	server.HandleMessage(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, w.Code)
	}

	// Check that response was sent via SSE
	if !strings.Contains(capturedData, "event: message") {
		t.Error("Expected SSE message event to be sent")
	}
	if !strings.Contains(capturedData, `"jsonrpc":"2.0"`) {
		t.Error("Expected JSON-RPC response to be sent")
	}
}

// mockResponseWriter implements http.ResponseWriter for testing
type mockResponseWriter struct {
	header     http.Header
	writeFunc  func([]byte) (int, error)
	statusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	if m.writeFunc != nil {
		return m.writeFunc(data)
	}
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func (m *mockResponseWriter) Flush() {
	// Mock flush - do nothing
}
