package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berrydev-ai/gqai/graphql"
)

func TestToolsCall(t *testing.T) {
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
    director
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a config that points to our temporary directory
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Test successful tool call
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/call",
		Params: map[string]any{
			"name":      "GetFilm",
			"arguments": map[string]any{"id": "123"},
		},
	}

	response := ToolsCall(request, config)

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	if response.Result == nil {
		t.Error("Expected result to be non-nil")
	}

	result, ok := response.Result.(CallToolResult)
	if !ok {
		t.Fatal("Expected result to be CallToolResult")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	if result.Content[0].Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", result.Content[0].Type)
	}
}

func TestToolsCallMissingParams(t *testing.T) {
	config := &graphql.GraphQLConfig{}

	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/call",
		Params:  nil,
	}

	response := mcp.ToolsCall(request, config)

	if response.Error == nil {
		t.Error("Expected error for missing params")
	}
	if response.Error.Code != InvalidParams {
		t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
	}
}

func TestToolsCallMissingToolName(t *testing.T) {
	config := &graphql.GraphQLConfig{}

	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/call",
		Params: map[string]any{
			"arguments": map[string]any{"id": "123"},
		},
	}

	response := ToolsCall(request, config)

	if response.Error == nil {
		t.Error("Expected error for missing tool name")
	}
	if response.Error.Code != InvalidParams {
		t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
	}
}

func TestToolsCallToolNotFound(t *testing.T) {
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Documents: []string{"/nonexistent"},
		},
	}

	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/call",
		Params: map[string]any{
			"name":      "NonExistentTool",
			"arguments": map[string]any{},
		},
	}

	response := ToolsCall(request, config)

	if response.Error == nil {
		t.Error("Expected error for tool not found")
	}
	if response.Error.Code != InternalError {
		t.Errorf("Expected InternalError, got %d", response.Error.Code)
	}
}

func TestToolsCallWithNilResponse(t *testing.T) {
	// TODO: This test requires mocking the LoadTool function to return a tool
	// that returns nil. Currently skipped due to import cycle issues when
	// trying to test internal MCP functions from within the MCP package.
	// Consider refactoring to allow better testability, such as dependency injection
	// or moving test utilities to a separate package.
	t.Skip("Skipping test that requires function mocking - import cycle issue")
}

func TestToolsCallWithInvalidResponseType(t *testing.T) {
	// Skip this test as it requires function mocking
	t.Skip("Skipping test that requires function mocking")
}
