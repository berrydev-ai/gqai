package graphql

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGraphQLConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	// Create a mock config that matches the expected structure for the GraphQLConfig
	configContent := `
schema:
  - http://localhost:4000/graphql:
      headers:
        Authorization: Token
documents:
  - operations/**/*.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	// Test using the actual LoadGraphQLConfig function
	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	// Verify the config was loaded correctly
	if config == nil {
		t.Fatal("LoadGraphQLConfig returned nil config")
	}

	if config.SingleProject == nil {
		t.Fatal("SingleProject is nil")
	}

	if len(config.SingleProject.Schema) != 1 {
		t.Fatalf("Expected 1 schema, got %d", len(config.SingleProject.Schema))
	}

	if config.SingleProject.Schema[0].URL != "http://localhost:4000/graphql" {
		t.Fatalf("Expected schema URL to be http://localhost:4000/graphql, got %s",
			config.SingleProject.Schema[0].URL)
	}

	headers := config.SingleProject.Schema[0].Headers
	if headers == nil {
		t.Fatal("Headers is nil")
	}

	if headers["Authorization"] != "Token" {
		t.Fatalf("Expected Authorization header to be 'Token', got '%s'",
			headers["Authorization"])
	}

	if len(config.SingleProject.Documents) != 1 {
		t.Fatalf("Expected 1 document path, got %d", len(config.SingleProject.Documents))
	}

	if config.SingleProject.Documents[0] != "operations/**/*.graphql" {
		t.Fatalf("Expected document path to be operations/**/*.graphql, got %s",
			config.SingleProject.Documents[0])
	}
}

// TestSchemaURLAsKey tests parsing schema configuration where the URL is a key
func TestSchemaURLAsKey(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	// Create a config with URL as key format
	configContent := `
schema:
  - http://localhost:4000/graphql:
      headers:
        Authorization: Bearer token
        Custom-Header: test-value
documents:
  - operations/**/*.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	// Load the config file
	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	// Verify the config is not nil
	if config == nil {
		t.Fatal("LoadGraphQLConfig returned nil config")
	}

	// Check if SingleProject was populated
	if config.SingleProject == nil {
		t.Fatal("SingleProject is nil")
	}

	// Check schema points
	if len(config.SingleProject.Schema) == 0 {
		t.Fatal("No schema pointers found")
	}

	// Check the URL
	if config.SingleProject.Schema[0].URL != "http://localhost:4000/graphql" {
		t.Fatalf("Expected URL to be http://localhost:4000/graphql, got %s",
			config.SingleProject.Schema[0].URL)
	}

	// Check headers
	headers := config.SingleProject.Schema[0].Headers
	if headers == nil {
		t.Fatal("Headers map is nil")
	}

	// Check specific headers
	if headers["Authorization"] != "Bearer token" {
		t.Fatalf("Expected Authorization header to be 'Bearer token', got '%s'",
			headers["Authorization"])
	}

	if headers["Custom-Header"] != "test-value" {
		t.Fatalf("Expected Custom-Header to be 'test-value', got '%s'",
			headers["Custom-Header"])
	}
}

func TestLoadGraphQLConfigWithURLAsKey(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	// Create a config with URL as key format
	configContent := `
schema:
  - http://localhost:4000/graphql:
      headers:
        Authorization: Token
        Custom-Header: test-value
documents:
  - ./documents/foo.graphql
  - ./documents/bar.graphql
  - ./documents/baz.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	// Load the config file
	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	// Verify the config was loaded correctly
	if config == nil {
		t.Fatal("LoadGraphQLConfig returned nil config")
	}

	if config.SingleProject == nil {
		t.Fatal("SingleProject is nil")
	}

	if len(config.SingleProject.Schema) != 1 {
		t.Fatalf("Expected 1 schema, got %d", len(config.SingleProject.Schema))
	}

	// Check the URL
	schemaPtr := config.SingleProject.Schema[0]
	if schemaPtr.URL != "http://localhost:4000/graphql" {
		t.Fatalf("Expected URL to be http://localhost:4000/graphql, got %s", schemaPtr.URL)
	}

	// Check headers
	if schemaPtr.Headers == nil {
		t.Fatal("Headers is nil")
	}

	if schemaPtr.Headers["Authorization"] != "Token" {
		t.Fatalf("Expected Authorization header to be 'Token', got '%s'",
			schemaPtr.Headers["Authorization"])
	}

	if schemaPtr.Headers["Custom-Header"] != "test-value" {
		t.Fatalf("Expected Custom-Header to be 'test-value', got '%s'",
			schemaPtr.Headers["Custom-Header"])
	}

	// Check documents
	if len(config.SingleProject.Documents) != 3 {
		t.Fatalf("Expected 3 documents, got %d", len(config.SingleProject.Documents))
	}

	expectedDocs := []string{
		"./documents/foo.graphql",
		"./documents/bar.graphql",
		"./documents/baz.graphql",
	}

	for i, doc := range config.SingleProject.Documents {
		if doc != expectedDocs[i] {
			t.Fatalf("Expected document %d to be %s, got %s", i, expectedDocs[i], doc)
		}
	}
}

func TestLoadGraphQLConfigWithSimpleURL(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	// Create a config with simple URL format
	configContent := `
schema: http://localhost:4000/graphql
documents: ./documents/*.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	// Load the config file
	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	// Verify the config was loaded correctly
	if config == nil {
		t.Fatal("LoadGraphQLConfig returned nil config")
	}

	if config.SingleProject == nil {
		t.Fatal("SingleProject is nil")
	}

	if len(config.SingleProject.Schema) != 1 {
		t.Fatalf("Expected 1 schema, got %d", len(config.SingleProject.Schema))
	}

	// Check the URL
	schemaPtr := config.SingleProject.Schema[0]
	if schemaPtr.URL != "http://localhost:4000/graphql" {
		t.Fatalf("Expected URL to be http://localhost:4000/graphql, got %s", schemaPtr.URL)
	}

	// Check that there are no headers
	if schemaPtr.Headers != nil && len(schemaPtr.Headers) > 0 {
		t.Fatalf("Expected no headers, got %v", schemaPtr.Headers)
	}

	// Check documents
	if len(config.SingleProject.Documents) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(config.SingleProject.Documents))
	}

	if config.SingleProject.Documents[0] != "./documents/*.graphql" {
		t.Fatalf("Expected document to be ./documents/*.graphql, got %s",
			config.SingleProject.Documents[0])
	}
}

func TestLoadGraphQLConfigWithEnvVarHeader(t *testing.T) {
	// Set an environment variable for the test
	os.Setenv("TEST_AUTH_TOKEN", "env-token-value")
	t.Cleanup(func() { os.Unsetenv("TEST_AUTH_TOKEN") })

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	configContent := `
schema:
  - http://localhost:4000/graphql:
      headers:
        Authorization: Bearer ${TEST_AUTH_TOKEN}
documents:
  - operations/**/*.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	headers := config.SingleProject.Schema[0].Headers
	if headers["Authorization"] != "Bearer env-token-value" {
		t.Fatalf("Expected Authorization header to be 'Bearer env-token-value', got '%s'", headers["Authorization"])
	}
}

func TestLoadGraphQLConfigWithEnvVarHeaderDefault(t *testing.T) {
	// Unset the env var to test default value
	os.Unsetenv("TEST_AUTH_TOKEN_DEFAULT")
	t.Cleanup(func() { os.Unsetenv("TEST_AUTH_TOKEN_DEFAULT") })

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	configContent := `
schema:
  - http://localhost:4000/graphql:
      headers:
        Authorization: Bearer ${TEST_AUTH_TOKEN_DEFAULT:-default-token}
documents:
  - operations/**/*.graphql
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	headers := config.SingleProject.Schema[0].Headers
	if headers["Authorization"] != "Bearer default-token" {
		t.Fatalf("Expected Authorization header to be 'Bearer default-token', got '%s'", headers["Authorization"])
	}
}

func TestLoadGraphQLConfigWithEnvVarsEverywhere(t *testing.T) {
	os.Setenv("TEST_SCHEMA_URL", "http://localhost:4000/graphql")
	os.Setenv("TEST_DOC_PATH", "operations/**/*.graphql")
	os.Setenv("TEST_INCLUDE", "operations/include.graphql")
	os.Setenv("TEST_EXCLUDE", "operations/exclude.graphql")
	os.Setenv("TEST_HEADER", "header-value")
	t.Cleanup(func() {
		os.Unsetenv("TEST_SCHEMA_URL")
		os.Unsetenv("TEST_DOC_PATH")
		os.Unsetenv("TEST_INCLUDE")
		os.Unsetenv("TEST_EXCLUDE")
		os.Unsetenv("TEST_HEADER")
	})

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "graphqlconfig.yml")

	configContent := `
schema:
  - ${TEST_SCHEMA_URL}:
      headers:
        X-Test-Header: ${TEST_HEADER}
documents:
  - ${TEST_DOC_PATH}
include: ${TEST_INCLUDE}
exclude: ${TEST_EXCLUDE}
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}

	config, err := LoadGraphQLConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGraphQLConfig returned an error: %v", err)
	}

	if config.SingleProject.Schema[0].URL != "http://localhost:4000/graphql" {
		t.Fatalf("Expected schema URL to be expanded, got %s", config.SingleProject.Schema[0].URL)
	}
	if config.SingleProject.Documents[0] != "operations/**/*.graphql" {
		t.Fatalf("Expected document path to be expanded, got %s", config.SingleProject.Documents[0])
	}
	if config.SingleProject.Include[0] != "operations/include.graphql" {
		t.Fatalf("Expected include to be expanded, got %s", config.SingleProject.Include[0])
	}
	if config.SingleProject.Exclude[0] != "operations/exclude.graphql" {
		t.Fatalf("Expected exclude to be expanded, got %s", config.SingleProject.Exclude[0])
	}
	headers := config.SingleProject.Schema[0].Headers
	if headers["X-Test-Header"] != "header-value" {
		t.Fatalf("Expected header to be expanded, got %s", headers["X-Test-Header"])
	}
}
