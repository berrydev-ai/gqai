package tool

import (
	"fmt"
	"github.com/fotoetienne/gqai/graphql"
)

func ToolsFromConfig(config *graphql.GraphQLConfig) ([]*MCPTool, error) {
	ops, err := graphql.LoadOperations(config)
	if err != nil {
		return nil, err
	}

	var tools []*MCPTool
	for _, op := range ops {
		tools = append(tools, toolFromOperation(config, op))
	}
	return tools, nil
}

func LoadTool(config *graphql.GraphQLConfig, name string) (*MCPTool, error) {
	ops, err := graphql.LoadOperations(config)
	if err != nil {
		return nil, err
	}

	op := ops[name]
	if op != nil {
		return toolFromOperation(config, op), nil
	}

	return nil, fmt.Errorf("tool %s not found", name)
}

func toolFromOperation(config *graphql.GraphQLConfig, op *graphql.Operation) *MCPTool {
	inputSchema, _ := ExtractInputSchema(op.Raw)
	endpoint := config.SingleProject.Schema[0].URL
	headers := config.SingleProject.Schema[0].Headers
	
	return &MCPTool{
		Name:        op.Name,
		Description: "", // TODO: maybe use docstring/comments?
		InputSchema: inputSchema,
		Execute: func(input map[string]any) (any, error) {
			return graphql.ExecuteWithHeaders(endpoint, input, op, headers)
		},
		Annotations: struct {
			Title           string `json:"title,omitempty"`
			ReadOnlyHint    bool   `json:"readOnlyHint"`
			DestructiveHint bool   `json:"destructiveHint"`
			IdempotentHint  bool   `json:"idempotentHint"`
			OpenWorldHint   bool   `json:"openWorldHint"`
		}{
			Title:           op.Name,
			ReadOnlyHint:    op.OperationType == "query",
			DestructiveHint: op.OperationType == "mutation",
			IdempotentHint:  op.OperationType == "query",
			OpenWorldHint:   true,
		},
	}
}
