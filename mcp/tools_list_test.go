package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fotoetienne/gqai/graphql"
)

func TestToolsList(t *testing.T) {
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

	mutationContent := `
mutation AddFilm($film: FilmInput!) {
  addFilm(film: $film) {
    id
  }
}
`
	mutationPath := filepath.Join(operationsDir, "add_film.graphql")
	err = os.WriteFile(mutationPath, []byte(mutationContent), 0644)
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

	// Test successful tools list
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/list",
	}

	response := ToolsList(request, config)

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	if response.Result == nil {
		t.Error("Expected result to be non-nil")
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	tools, ok := result["tools"]
	if !ok {
		t.Error("Expected result to contain 'tools' key")
	}

	toolsSlice, ok := tools.([]*interface{})
	if !ok {
		t.Error("Expected tools to be a slice")
	}

	if len(toolsSlice) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(toolsSlice))
	}
}

func TestToolsListWithEmptyConfig(t *testing.T) {
	// Test with config that has no operations
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Documents: []string{"/nonexistent"},
		},
	}

	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/list",
	}

	response := ToolsList(request, config)

	if response.Error != nil {
		t.Errorf("Expected no error for empty config, got %v", response.Error)
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	tools, ok := result["tools"]
	if !ok {
		t.Error("Expected result to contain 'tools' key")
	}

	toolsSlice, ok := tools.([]*interface{})
	if !ok {
		t.Error("Expected tools to be a slice")
	}

	if len(toolsSlice) != 0 {
		t.Errorf("Expected 0 tools for empty config, got %d", len(toolsSlice))
	}
}

func TestToolsListWithInvalidConfig(t *testing.T) {
	// Test with config that causes an error
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Documents: []string{}, // Empty documents array
		},
	}

	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "tools/list",
	}

	response := ToolsList(request, config)

	if response.Error == nil {
		t.Error("Expected error for invalid config")
	}
	if response.Error.Code != InternalError {
		t.Errorf("Expected InternalError, got %d", response.Error.Code)
	}
}
