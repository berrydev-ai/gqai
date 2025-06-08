package graphql

import (
	"fmt"
	"net/textproto"
	"os"
	"regexp"

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

// expandEnvVars replaces ${VAR}, $VAR, or ${VAR:-default} in a string with the value from the environment or a default
func expandEnvVars(s string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(:-([^}]*))?\}|\$([A-Za-z0-9_]+)`) // matches ${VAR}, ${VAR:-default}, or $VAR
	return re.ReplaceAllStringFunc(s, func(m string) string {
		if m[0] == '$' && len(m) > 1 && m[1] == '{' {
			// ${VAR} or ${VAR:-default}
			inner := m[2 : len(m)-1]
			if idx := regexp.MustCompile(`:-`).FindStringIndex(inner); idx != nil {
				varName := inner[:idx[0]]
				defaultVal := inner[idx[1]:]
				if val, ok := os.LookupEnv(varName); ok {
					return val
				}
				return defaultVal
			} else {
				varName := inner
				if val, ok := os.LookupEnv(varName); ok {
					return val
				}
				return m // leave as is if not found
			}
		} else if m[0] == '$' {
			// $VAR
			varName := m[1:]
			if val, ok := os.LookupEnv(varName); ok {
				return val
			}
			return m // leave as is if not found
		}
		return m
	})
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

	// Parse include configuration
	if err := parseInclude(config); err != nil {
		return nil, err
	}

	// Parse exclude configuration
	if err := parseExclude(config); err != nil {
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
			URL: expandEnvVars(urlStr),
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
				URL: expandEnvVars(urlStr),
			})
			continue
		}

		// Case 2.2: array item is a map with URL as key
		urlMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip invalid items
		}

		for url, value := range urlMap {
			schemaPtr := SchemaPointer{URL: expandEnvVars(url)}

			// Extract headers if they exist
			if valueMap, ok := value.(map[string]interface{}); ok {
				if headersMap, ok := valueMap["headers"].(map[string]interface{}); ok {
					schemaPtr.Headers = make(map[string]string)
					for key, val := range headersMap {
						if valStr, ok := val.(string); ok {
							// Use standard header normalization
							normalKey := normalizeHeader(key)
							// Expand env vars in header values
							schemaPtr.Headers[normalKey] = expandEnvVars(valStr)
						}
					}
				}
				// Expand env vars in other string fields in valueMap if needed
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
		config.SingleProject.Documents = append(config.SingleProject.Documents, expandEnvVars(docStr))
		return nil
	}

	// Case 2: documents is an array
	docArr, ok := documents.([]interface{})
	if !ok {
		return fmt.Errorf("documents must be a string path or an array")
	}

	for _, item := range docArr {
		if docStr, ok := item.(string); ok {
			config.SingleProject.Documents = append(config.SingleProject.Documents, expandEnvVars(docStr))
		}
	}

	return nil
}

// parseInclude handles the include configuration which can be either a string or an array
func parseInclude(config *GraphQLConfig) error {
	include := viper.Get("include")
	if include == nil {
		return nil
	}

	if config.SingleProject == nil {
		config.SingleProject = &GraphQLProject{}
	}

	// Case 1: include is a simple string
	if incStr, ok := include.(string); ok {
		config.SingleProject.Include = append(config.SingleProject.Include, expandEnvVars(incStr))
		return nil
	}

	// Case 2: include is an array
	incArr, ok := include.([]interface{})
	if !ok {
		return fmt.Errorf("include must be a string path or an array")
	}

	for _, item := range incArr {
		if incStr, ok := item.(string); ok {
			config.SingleProject.Include = append(config.SingleProject.Include, expandEnvVars(incStr))
		}
	}

	return nil
}

// parseExclude handles the exclude configuration which can be either a string or an array
func parseExclude(config *GraphQLConfig) error {
	exclude := viper.Get("exclude")
	if exclude == nil {
		return nil
	}

	if config.SingleProject == nil {
		config.SingleProject = &GraphQLProject{}
	}

	// Case 1: exclude is a simple string
	if excStr, ok := exclude.(string); ok {
		config.SingleProject.Exclude = append(config.SingleProject.Exclude, expandEnvVars(excStr))
		return nil
	}

	// Case 2: exclude is an array
	excArr, ok := exclude.([]interface{})
	if !ok {
		return fmt.Errorf("exclude must be a string path or an array")
	}

	for _, item := range excArr {
		if excStr, ok := item.(string); ok {
			config.SingleProject.Exclude = append(config.SingleProject.Exclude, expandEnvVars(excStr))
		}
	}

	return nil
}
