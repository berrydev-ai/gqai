package tool

import (
	"github.com/fotoetienne/gqai/graphql"
)

func ToolsFromConfig(config *graphql.Config) ([]*MCPTool, error) {
	ops, err := graphql.LoadOperations(config)
	if err != nil {
		return nil, err
	}

	tools := []*MCPTool{}
	for name, op := range ops {
		// Avoid closure capture bug
		toolName := name
		toolOp := op

		inputSchema, _ := ExtractInputSchema(op.Raw)

		tools = append(tools, &MCPTool{
			Name:        toolName,
			Description: "", // TODO: maybe use docstring/comments?
			InputSchema: inputSchema,
			Execute: func(input map[string]any) (any, error) {
				return Execute(config, name, input)
			},
			Annotations: struct {
				Title           string `json:"title,omitempty"`
				ReadOnlyHint    bool   `json:"readOnlyHint"`
				DestructiveHint bool   `json:"destructiveHint"`
				IdempotentHint  bool   `json:"idempotentHint"`
				OpenWorldHint   bool   `json:"openWorldHint"`
			}{
				Title:           toolName,
				ReadOnlyHint:    toolOp.OperationType == "query",
				DestructiveHint: toolOp.OperationType == "mutation",
				IdempotentHint:  toolOp.OperationType == "query",
				OpenWorldHint:   true,
			},
		})
	}
	return tools, nil
}
