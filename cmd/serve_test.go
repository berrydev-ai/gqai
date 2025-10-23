package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/berrydev-ai/gqai/graphql"
	"github.com/gorilla/mux"
)

func TestListToolsHandler(t *testing.T) {
	// Create a temporary directory for operations
	tempDir := t.TempDir()
	operationsDir := filepath.Join(tempDir, "operations")
	err := os.MkdirAll(operationsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary operations directory: %v", err)
	}

	// Create sample GraphQL files
	queryContent := `
query GetFilm($id: ID!) {
  film(id: $id) {
    title
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a config that points to our temporary directory
	config = &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Test successful list tools
	req := httptest.NewRequest("GET", "/tools", nil)
	w := httptest.NewRecorder()

	listToolsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Try to decode the response
	var tools []interface{}
	if err := json.NewDecoder(w.Body).Decode(&tools); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
}

func TestListToolsHandlerError(t *testing.T) {
	// Set config to nil to cause an error
	config = nil

	req := httptest.NewRequest("GET", "/tools", nil)
	w := httptest.NewRecorder()

	listToolsHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedBody := "Error loading tools"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestCallToolHandler(t *testing.T) {
	// Create a temporary directory for operations
	tempDir := t.TempDir()
	operationsDir := filepath.Join(tempDir, "operations")
	err := os.MkdirAll(operationsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary operations directory: %v", err)
	}

	// Create a sample GraphQL query file
	queryContent := `
query GetFilm($id: ID!) {
  film(id: $id) {
    title
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a config that points to our temporary directory
	config = &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Test successful tool call
	payload := map[string]interface{}{
		"toolName": "GetFilm",
		"input":    map[string]interface{}{"id": "123"},
	}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/tools/call", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	callToolHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Try to decode the response
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if _, exists := response["output"]; !exists {
		t.Error("Expected response to contain 'output' field")
	}
}

func TestCallToolHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/tools/call", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	callToolHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Invalid JSON"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestCallToolHandlerToolNotFound(t *testing.T) {
	// Create a config with no operations
	config = &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Documents: []string{"/nonexistent"},
		},
	}

	payload := map[string]interface{}{
		"toolName": "NonExistentTool",
		"input":    map[string]interface{}{},
	}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/tools/call", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	callToolHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedBody := "Error loading tool"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestServeHandler(t *testing.T) {
	// Create a temporary directory for operations
	tempDir := t.TempDir()
	operationsDir := filepath.Join(tempDir, "operations")
	err := os.MkdirAll(operationsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary operations directory: %v", err)
	}

	// Create a sample GraphQL query file
	queryContent := `
query GetFilm($id: ID!) {
  film(id: $id) {
    title
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a config that points to our temporary directory
	config = &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Test successful serve handler
	payload := map[string]interface{}{
		"input": map[string]interface{}{"id": "123"},
	}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/tools/GetFilm", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	// Set up mux vars for the tool name
	vars := map[string]string{
		"name": "GetFilm",
	}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()

	serveHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Try to decode the response
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if _, exists := response["output"]; !exists {
		t.Error("Expected response to contain 'output' field")
	}
}

func TestServeHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/tools/TestTool", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	serveHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "Invalid JSON"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected response to contain '%s', got '%s'", expectedBody, w.Body.String())
	}
}
