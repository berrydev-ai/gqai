package tool

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestExtractInputSchema(t *testing.T) {
	// Test with a valid query with variables
	query := `
query GetFilm($id: ID!, $title: String) {
  film(id: $id, title: $title) {
    id
    title
  }
}
`
	schema, err := ExtractInputSchema(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	// Check basic structure
	if schema["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check id parameter (required)
	idProp, ok := props["id"]
	if !ok {
		t.Fatal("Expected 'id' property")
	}
	idMap, ok := idProp.(map[string]any)
	if !ok {
		t.Fatal("Expected id property to be a map")
	}
	if idMap["type"] != "string" {
		t.Errorf("Expected id type 'string', got %v", idMap["type"])
	}

	// Check title parameter (optional)
	titleProp, ok := props["title"]
	if !ok {
		t.Fatal("Expected 'title' property")
	}
	titleMap, ok := titleProp.(map[string]any)
	if !ok {
		t.Fatal("Expected title property to be a map")
	}
	if titleMap["type"] != "string" {
		t.Errorf("Expected title type 'string', got %v", titleMap["type"])
	}

	// Check required array
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Expected required to be a string slice")
	}
	if len(required) != 1 || required[0] != "id" {
		t.Errorf("Expected required ['id'], got %v", required)
	}
}

func TestExtractInputSchemaNoVariables(t *testing.T) {
	// Test with a query that has no variables
	query := `
query GetFilms {
  films {
    id
    title
  }
}
`
	schema, err := ExtractInputSchema(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	if len(props) != 0 {
		t.Errorf("Expected no properties, got %d", len(props))
	}

	// Should not have required field
	if _, hasRequired := schema["required"]; hasRequired {
		t.Error("Expected no required field for query with no variables")
	}
}

func TestExtractInputSchemaInvalidGraphQL(t *testing.T) {
	// Test with invalid GraphQL
	query := "invalid graphql syntax"
	schema, err := ExtractInputSchema(query)
	if err == nil {
		t.Error("Expected error for invalid GraphQL")
	}
	if schema != nil {
		t.Errorf("Expected nil schema for invalid GraphQL, got %v", schema)
	}
}

func TestExtractInputSchemaNoOperations(t *testing.T) {
	// Test with GraphQL that has no operations
	query := "fragment Test on Film { id }"
	schema, err := ExtractInputSchema(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if schema != nil {
		t.Errorf("Expected nil schema for GraphQL with no operations, got %v", schema)
	}
}

func TestExtractInputSchemaWithDifferentTypes(t *testing.T) {
	// Test with various GraphQL types
	query := `
mutation CreateFilm($title: String!, $year: Int, $rating: Float, $available: Boolean, $tags: [String]) {
  createFilm(input: { title: $title, year: $year, rating: $rating, available: $available, tags: $tags }) {
    id
  }
}
`
	schema, err := ExtractInputSchema(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check String type
	if props["title"].(map[string]any)["type"] != "string" {
		t.Error("Expected title to be string type")
	}

	// Check Int type
	if props["year"].(map[string]any)["type"] != "integer" {
		t.Error("Expected year to be integer type")
	}

	// Check Float type
	if props["rating"].(map[string]any)["type"] != "number" {
		t.Error("Expected rating to be number type")
	}

	// Check Boolean type
	if props["available"].(map[string]any)["type"] != "boolean" {
		t.Error("Expected available to be boolean type")
	}

	// Check array type
	if props["tags"].(map[string]any)["type"] != "array" {
		t.Error("Expected tags to be array type")
	}

	// Check required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Expected required to be a string slice")
	}
	if len(required) != 1 || required[0] != "title" {
		t.Errorf("Expected required ['title'], got %v", required)
	}
}

func TestGraphqlTypeToJSONSchemaType(t *testing.T) {
	tests := []struct {
		graphqlType string
		isArray     bool
		expected    string
	}{
		{"String", false, "string"},
		{"ID", false, "string"},
		{"Int", false, "integer"},
		{"Float", false, "number"},
		{"Boolean", false, "boolean"},
		{"CustomType", false, "string"}, // fallback
		{"String", true, "array"},       // array case
	}

	for _, test := range tests {
		astType := &ast.Type{
			NamedType: test.graphqlType,
		}
		if test.isArray {
			astType.Elem = &ast.Type{} // simulate array
		}

		result := graphqlTypeToJSONSchemaType(astType)
		if result != test.expected {
			t.Errorf("Expected %s for %s (array: %v), got %s", test.expected, test.graphqlType, test.isArray, result)
		}
	}
}
