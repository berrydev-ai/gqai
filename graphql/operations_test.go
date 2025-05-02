package graphql

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOperations(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	operationsDir := filepath.Join(tempDir, "operations")
	err := os.MkdirAll(operationsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary operations directory: %v", err)
	}

	// Create a sample GraphQL operation file
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
	config := &GraphQLConfig{
		SingleProject: &GraphQLProject{
			Documents: []string{operationsDir},
		},
	}

	// Load operations
	operations, err := LoadOperations(config)
	if err != nil {
		t.Fatalf("LoadOperations returned an error: %v", err)
	}

	// Verify operations were loaded correctly
	if len(operations) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(operations))
	}

	// Check for the GetFilm operation
	operation, exists := operations["GetFilm"]
	if !exists {
		t.Fatal("GetFilm operation not found")
	}

	// Verify the operation properties
	if operation.Name != "GetFilm" {
		t.Fatalf("Expected operation name to be GetFilm, got %s", operation.Name)
	}

	if operation.OperationType != "query" {
		t.Fatalf("Expected operation type to be query, got %s", operation.OperationType)
	}
}
