package graphql

import (
	"fmt"
	"net/textproto"

	"github.com/spf13/viper"
)

type SchemaPointer struct {
	URL     string            `yaml:"-"` // URL will be set programmatically
	Headers map[string]string `yaml:"headers,omitempty"`
}

type GraphQLProject struct {
	Schema     []SchemaPointer `yaml:"schema"`
	Documents  []string        `yaml:"documents"`
	Extensions map[string]any  `yaml:"extensions"`
	Include    []string        `yaml:"include"`
	Exclude    []string        `yaml:"exclude"`
}

type GraphQLProjects struct {
	Projects map[string]GraphQLProject `yaml:"projects"`
}

type GraphQLConfig struct {
	SingleProject *GraphQLProject  `yaml:",inline"`
	MultiProjects *GraphQLProjects `yaml:",inline"`
}

// normalizeHeader returns the canonical format of the header key
// This uses the standard library's textproto.CanonicalMIMEHeaderKey
// which handles proper capitalization of HTTP headers
func normalizeHeader(key string) string {
	return textproto.CanonicalMIMEHeaderKey(key)
}

func LoadGraphQLConfig(path string) (*GraphQLConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	// Create a new config object
	config := &GraphQLConfig{
		SingleProject: &GraphQLProject{},
	}

	// Parse schema configuration
	if err := parseSchema(config); err != nil {
		return nil, err
	}

	// Parse documents configuration
	if err := parseDocuments(config); err != nil {
		return nil, err
	}

	return config, nil
}

// parseSchema handles the schema configuration which can be either a string or an array
func parseSchema(config *GraphQLConfig) error {
	schema := viper.Get("schema")
	if schema == nil {
		return nil
	}

	// Case 1: schema is a simple string URL
	if urlStr, ok := schema.(string); ok {
		config.SingleProject.Schema = append(config.SingleProject.Schema, SchemaPointer{
			URL: urlStr,
		})
		return nil
	}

	// Case 2: schema is an array
	schemaArr, ok := schema.([]interface{})
	if !ok {
		return fmt.Errorf("schema must be a string URL or an array")
	}

	for _, item := range schemaArr {
		// Case 2.1: array item is a simple string URL
		if urlStr, ok := item.(string); ok {
			config.SingleProject.Schema = append(config.SingleProject.Schema, SchemaPointer{
				URL: urlStr,
			})
			continue
		}

		// Case 2.2: array item is a map with URL as key
		urlMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip invalid items
		}

		for url, value := range urlMap {
			schemaPtr := SchemaPointer{URL: url}

			// Extract headers if they exist
			if valueMap, ok := value.(map[string]interface{}); ok {
				if headersMap, ok := valueMap["headers"].(map[string]interface{}); ok {
					schemaPtr.Headers = make(map[string]string)
					for key, val := range headersMap {
						if valStr, ok := val.(string); ok {
							// Use standard header normalization
							normalKey := normalizeHeader(key)
							schemaPtr.Headers[normalKey] = valStr
						}
					}
				}
			}

			config.SingleProject.Schema = append(config.SingleProject.Schema, schemaPtr)
			break // Only process the first URL in the map
		}
	}

	return nil
}

// parseDocuments handles the documents configuration which can be either a string or an array
func parseDocuments(config *GraphQLConfig) error {
	documents := viper.Get("documents")
	if documents == nil {
		return nil
	}

	// Case 1: documents is a simple string
	if docStr, ok := documents.(string); ok {
		config.SingleProject.Documents = append(config.SingleProject.Documents, docStr)
		return nil
	}

	// Case 2: documents is an array
	docArr, ok := documents.([]interface{})
	if !ok {
		return fmt.Errorf("documents must be a string path or an array")
	}

	for _, item := range docArr {
		if docStr, ok := item.(string); ok {
			config.SingleProject.Documents = append(config.SingleProject.Documents, docStr)
		}
	}

	return nil
}
