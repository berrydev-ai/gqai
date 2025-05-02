package tool

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fotoetienne/gqai/graphql"
)

func TestToolsFromConfig(t *testing.T) {
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
    releaseDate
  }
}
`
	queryPath := filepath.Join(operationsDir, "get_film.graphql")
	err = os.WriteFile(queryPath, []byte(queryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample GraphQL file: %v", err)
	}

	// Create a sample GraphQL mutation file
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

	// Create a config that points to our temporary directory
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{operationsDir},
		},
	}

	// Get tools from config
	tools, err := ToolsFromConfig(config)
	if err != nil {
		t.Fatalf("ToolsFromConfig returned an error: %v", err)
	}

	// Verify tools were created correctly
	if len(tools) != 2 {
		t.Fatalf("Expected 2 tools, got %d", len(tools))
	}

	// Check for tools by name
	var getFilmTool, addFilmTool *MCPTool
	for _, tool := range tools {
		if tool.Name == "GetFilm" {
			getFilmTool = tool
		} else if tool.Name == "AddFilm" {
			addFilmTool = tool
		}
	}

	// Verify GetFilm tool
	if getFilmTool == nil {
		t.Fatal("GetFilm tool not found")
	}
	if !getFilmTool.Annotations.ReadOnlyHint {
		t.Error("Expected GetFilm tool to be read-only")
	}
	if getFilmTool.Annotations.DestructiveHint {
		t.Error("Expected GetFilm tool not to be destructive")
	}

	// Verify AddFilm tool
	if addFilmTool == nil {
		t.Fatal("AddFilm tool not found")
	}
	if !addFilmTool.Annotations.DestructiveHint {
		t.Error("Expected AddFilm tool to be destructive")
	}
	if addFilmTool.Annotations.ReadOnlyHint {
		t.Error("Expected AddFilm tool not to be read-only")
	}
}

func TestLoadTool(t *testing.T) {
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
    releaseDate
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

	// Load a specific tool
	tool, err := LoadTool(config, "GetFilm")
	if err != nil {
		t.Fatalf("LoadTool returned an error: %v", err)
	}

	// Verify tool was loaded correctly
	if tool == nil {
		t.Fatal("LoadTool returned nil tool")
	}
	if tool.Name != "GetFilm" {
		t.Fatalf("Expected tool name to be GetFilm, got %s", tool.Name)
	}

	// Test loading a non-existent tool
	_, err = LoadTool(config, "NonExistentTool")
	if err == nil {
		t.Error("Expected error when loading non-existent tool, got nil")
	}
}
