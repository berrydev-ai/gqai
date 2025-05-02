package graphql

import (
	"fmt"
	_ "os"
	"strings"

	"github.com/spf13/viper"
)

type SchemaPointer struct {
	URL     string            `yaml:"-"` // URL will be set programmatically
	Headers map[string]string `yaml:"headers,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (s *SchemaPointer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to unmarshal as string directly (simple URL with no headers)
	var stringURL string
	if err := unmarshal(&stringURL); err == nil {
		s.URL = stringURL
		return nil
	}

	// Try to unmarshal as map with URL as key and headers as nested value
	var schemaMap map[string]interface{}
	if err := unmarshal(&schemaMap); err != nil {
		return fmt.Errorf("schema must be a string URL or a map with URL as key")
	}

	// There should be only one key which is the URL
	for url, value := range schemaMap {
		s.URL = url
		
		// Check if the value is a map that contains headers
		if valueMap, ok := value.(map[interface{}]interface{}); ok {
			if headersMap, found := valueMap["headers"]; found {
				if headers, ok := headersMap.(map[interface{}]interface{}); ok {
					s.Headers = make(map[string]string)
					for k, v := range headers {
						if keyStr, ok := k.(string); ok {
							if valStr, ok := v.(string); ok {
								s.Headers[keyStr] = valStr
							}
						}
					}
				}
				break // Only process the first URL
			}
		}
	}
	
	return nil
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

	// Handle the schema property (can be a string or an array)
	schema := viper.Get("schema")
	if schema != nil {
		// Case 1: schema is a simple string URL
		if urlStr, ok := schema.(string); ok {
			config.SingleProject.Schema = append(config.SingleProject.Schema, SchemaPointer{
				URL: urlStr,
			})
		}
		
		// Case 2: schema is an array
		if schemaArr, ok := schema.([]interface{}); ok {
			for _, item := range schemaArr {
				// Case 2.1: array item is a simple string URL
				if urlStr, ok := item.(string); ok {
					config.SingleProject.Schema = append(config.SingleProject.Schema, SchemaPointer{
						URL: urlStr,
					})
					continue
				}
				
				// Case 2.2: array item is a map with URL as key (map[string]interface{})
				if urlMap, ok := item.(map[string]interface{}); ok {
					for url, value := range urlMap {
						schemaPtr := SchemaPointer{URL: url}
						
						// Extract headers if they exist
						if valueMap, ok := value.(map[string]interface{}); ok {
							if headersMap, ok := valueMap["headers"].(map[string]interface{}); ok {
								schemaPtr.Headers = make(map[string]string)
								for key, val := range headersMap {
									if valStr, ok := val.(string); ok {
											// Restore original case for common header names
										// This is necessary because Viper normalizes keys to lowercase
										headerKey := key
										switch strings.ToLower(key) {
										case "authorization":
											headerKey = "Authorization"
										case "content-type":
											headerKey = "Content-Type"
										case "accept":
											headerKey = "Accept"
										case "user-agent":
											headerKey = "User-Agent"
										case "custom-header":
											headerKey = "Custom-Header"
										}
										
										// Store header with restored case
										schemaPtr.Headers[headerKey] = valStr
									}
								}
							}
						}
						
						config.SingleProject.Schema = append(config.SingleProject.Schema, schemaPtr)
						break // Only process the first URL in the map
					}
				}
			}
		}
	}

	// Handle documents (similar approach)
	documents := viper.Get("documents")
	if documents != nil {
		// Case 1: documents is a simple string
		if docStr, ok := documents.(string); ok {
			config.SingleProject.Documents = append(config.SingleProject.Documents, docStr)
		}
		
		// Case 2: documents is an array
		if docArr, ok := documents.([]interface{}); ok {
			for _, item := range docArr {
				if docStr, ok := item.(string); ok {
					config.SingleProject.Documents = append(config.SingleProject.Documents, docStr)
				}
			}
		}
	}

	return config, nil
}
