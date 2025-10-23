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

 func TestLoadOperationsWithInvalidGraphQL(t *testing.T) {
  // Create a temporary directory
  tempDir := t.TempDir()
  operationsDir := filepath.Join(tempDir, "operations")
  err := os.MkdirAll(operationsDir, 0755)
  if err != nil {
  	t.Fatalf("Failed to create temporary operations directory: %v", err)
  }

  // Create a file with invalid GraphQL
  invalidQueryContent := `
 invalid graphql syntax here
 `
  invalidPath := filepath.Join(operationsDir, "invalid.graphql")
  err = os.WriteFile(invalidPath, []byte(invalidQueryContent), 0644)
  if err != nil {
  	t.Fatalf("Failed to create invalid GraphQL file: %v", err)
  }

  // Create a config that points to our temporary directory
  config := &GraphQLConfig{
  	SingleProject: &GraphQLProject{
  		Documents: []string{operationsDir},
  	},
  }

  // Load operations - should return error for invalid GraphQL
  operations, err := LoadOperations(config)
  if err == nil {
  	t.Fatalf("Expected error for invalid GraphQL, got nil")
  }
  if operations != nil {
  	t.Fatalf("Expected nil operations for invalid GraphQL, got %v", operations)
  }
 }

 func TestLoadOperationsWithEmptyDirectory(t *testing.T) {
  // Create a temporary empty directory
  tempDir := t.TempDir()
  operationsDir := filepath.Join(tempDir, "operations")
  err := os.MkdirAll(operationsDir, 0755)
  if err != nil {
  	t.Fatalf("Failed to create temporary operations directory: %v", err)
  }

  // Create a config that points to our empty directory
  config := &GraphQLConfig{
  	SingleProject: &GraphQLProject{
  		Documents: []string{operationsDir},
  	},
  }

  // Load operations - should succeed with empty map
  operations, err := LoadOperations(config)
  if err != nil {
  	t.Fatalf("LoadOperations returned an error for empty directory: %v", err)
  }
  if len(operations) != 0 {
  	t.Fatalf("Expected empty operations map, got %d operations", len(operations))
  }
 }

 func TestLoadOperationsWithNonExistentDirectory(t *testing.T) {
  // Create a config that points to a non-existent directory
  config := &GraphQLConfig{
  	SingleProject: &GraphQLProject{
  		Documents: []string{"/non/existent/directory"},
  	},
  }

  // Load operations - should return error
  operations, err := LoadOperations(config)
  if err == nil {
  	t.Fatalf("Expected error for non-existent directory, got nil")
  }
  if operations != nil {
  	t.Fatalf("Expected nil operations for non-existent directory, got %v", operations)
  }
 }

 func TestLoadOperationsWithMultipleOperationsInOneFile(t *testing.T) {
  // Create a temporary directory
  tempDir := t.TempDir()
  operationsDir := filepath.Join(tempDir, "operations")
  err := os.MkdirAll(operationsDir, 0755)
  if err != nil {
  	t.Fatalf("Failed to create temporary operations directory: %v", err)
  }

  // Create a file with multiple operations
  multiOpContent := `
 query GetFilm($id: ID!) {
   film(id: $id) {
     title
     director
   }
 }

 mutation AddFilm($film: FilmInput!) {
   addFilm(film: $film) {
     id
     title
   }
 }
 `
  multiPath := filepath.Join(operationsDir, "multi.graphql")
  err = os.WriteFile(multiPath, []byte(multiOpContent), 0644)
  if err != nil {
  	t.Fatalf("Failed to create multi-operation GraphQL file: %v", err)
  }

  // Create a config that points to our temporary directory
  config := &graphql.GraphQLConfig{
  	SingleProject: &graphql.GraphQLProject{
  		Documents: []string{operationsDir},
  	},
  }

  // Load operations
  operations, err := LoadOperations(config)
  if err != nil {
  	t.Fatalf("LoadOperations returned an error: %v", err)
  }

  // Verify both operations were loaded
  if len(operations) != 2 {
  	t.Fatalf("Expected 2 operations, got %d", len(operations))
  }

  getFilmOp, exists := operations["GetFilm"]
  if !exists {
  	t.Fatal("GetFilm operation not found")
  }
  if getFilmOp.OperationType != "query" {
  	t.Fatalf("Expected GetFilm operation type to be query, got %s", getFilmOp.OperationType)
  }

  addFilmOp, exists := operations["AddFilm"]
  if !exists {
  	t.Fatal("AddFilm operation not found")
  }
  if addFilmOp.OperationType != "mutation" {
  	t.Fatalf("Expected AddFilm operation type to be mutation, got %s", addFilmOp.OperationType)
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
