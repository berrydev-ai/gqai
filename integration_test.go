package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fotoetienne/gqai/graphql"
	"github.com/fotoetienne/gqai/mcp"
	"github.com/fotoetienne/gqai/tool"
)

func TestFullIntegration(t *testing.T) {
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
    director
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	mutationContent := `
mutation AddFilm($film: FilmInput!) {
  addFilm(film: $film) {
    id
    title
  }
}
`
	mutationPath := filepath.Join(operationsDir, "add_film.graphql")
	err = os.WriteFile(mutationPath, []byte(mutationContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a GraphQL config
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Test 1: Load operations
	operations, err := graphql.LoadOperations(config)
	if err != nil {
		t.Fatalf("Failed to load operations: %v", err)
	}
	if len(operations) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(operations))
	}

	// Test 2: Create tools from config
	tools, err := tool.ToolsFromConfig(config)
	if err != nil {
		t.Fatalf("Failed to create tools from config: %v", err)
	}
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	// Test 3: MCP tools/list
	request := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-1",
		Method:  "tools/list",
	}

	response := mcp.RouteMCPRequest(request, config)
	if response.Error != nil {
		t.Errorf("tools/list failed: %v", response.Error)
	}

	// Test 4: MCP tools/call
	callRequest := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-2",
		Method:  "tools/call",
		Params: map[string]any{
			"name":      "GetFilm",
			"arguments": map[string]any{"id": "123"},
		},
	}

	callResponse := mcp.RouteMCPRequest(callRequest, config)
	if callResponse.Error != nil {
		t.Errorf("tools/call failed: %v", callResponse.Error)
	}

	// Test 5: HTTP server integration
	// Create a test server using the streamable HTTP handler
	server := &mcp.StreamableHTTPServer{
		Config:   config,
		Sessions: make(map[string]*mcp.StreamableHTTPSession),
	}

	// Test GET (establish session)
	getReq := httptest.NewRequest("GET", "/mcp", nil)
	getW := httptest.NewRecorder()
	server.HandleStreamableHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected GET status 200, got %d", getW.Code)
	}

	// Extract session ID from response
	body := getW.Body.String()
	sessionID := ""
	if strings.Contains(body, "data: ") {
		// Extract session ID from SSE format
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "data: ") {
				sessionID = strings.TrimPrefix(line, "data: ")
				break
			}
		}
	}

	if sessionID == "" {
		t.Error("Failed to extract session ID from GET response")
	}

	// Test POST (send message) with the session ID
	postPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "test-3",
		"method":  "tools/list",
	}
	payloadBytes, _ := json.Marshal(postPayload)

	postReq := httptest.NewRequest("POST", "/mcp", bytes.NewReader(payloadBytes))
	postReq.Header.Set("Mcp-Session-Id", sessionID)
	postW := httptest.NewRecorder()

	server.HandleStreamableHTTP(postW, postReq)

	if postW.Code != http.StatusAccepted {
		t.Errorf("Expected POST status 202, got %d", postW.Code)
	}

	// Test 6: Config loading
	configPath := filepath.Join(tempDir, ".graphqlrc.yml")
	configContent := `
schema: http://example.com/graphql
documents: ./operations/*.graphql
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	loadedConfig, err := graphql.LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.SingleProject.Schema[0].URL != "http://example.com/graphql" {
		t.Errorf("Expected schema URL 'http://example.com/graphql', got '%s'", loadedConfig.SingleProject.Schema[0].URL)
	}

	if len(loadedConfig.SingleProject.Documents) != 1 {
		t.Errorf("Expected 1 document pattern, got %d", len(loadedConfig.SingleProject.Documents))
	}
}

func TestErrorHandlingIntegration(t *testing.T) {
	// Test error handling throughout the system
	config := &graphql.GraphQLConfig{}

	// Test MCP with invalid method
	request := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-error",
		Method:  "invalid_method",
	}

	response := mcp.RouteMCPRequest(request, config)
	if response.Error == nil {
		t.Error("Expected error for invalid method")
	}
	if response.Error.Code != mcp.MethodNotFound {
		t.Errorf("Expected MethodNotFound error, got %d", response.Error.Code)
	}

	// Test tools/call with missing tool
	callRequest := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-error-2",
		Method:  "tools/call",
		Params: map[string]any{
			"name": "NonExistentTool",
		},
	}

	callResponse := mcp.RouteMCPRequest(callRequest, config)
	if callResponse.Error == nil {
		t.Error("Expected error for non-existent tool")
	}
}
