package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fotoetienne/gqai/graphql"
)

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func Execute(config *graphql.Config, toolName string, input map[string]any) (any, error) {
	// Load operations from the config
	ops, err := graphql.LoadOperations(config)
	if err != nil {
		return nil, err
	}

	// Find the specific operation for the tool
	op, ok := ops[toolName]
	if !ok {
		// if there are no operations, return an error
		if len(ops) == 0 {
			return nil, fmt.Errorf("no operations found in %s", config.Documents)
		}
		// debug: print available operations
		fmt.Println("Available operations: ")
		for name := range ops {
			fmt.Println("    ", name)
		}
		return nil, fmt.Errorf("operation %s not found", toolName)
	}

	// POST to schema URL
	reqBody := gqlRequest{
		Query:     op.Raw,
		Variables: input,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(config.Schema, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("GraphQL request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	return result, nil
}
