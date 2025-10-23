package mcp

import (
	"fmt"
	"github.com/berrydev-ai/gqai/graphql"
	"github.com/berrydev-ai/gqai/tool"
)

// ToolsList handles the 'tools/list' MCP command.
func ToolsList(request JSONRPCRequest, config *graphql.GraphQLConfig) JSONRPCResponse {
	tools, err := tool.ToolsFromConfig(config)
	if err != nil {
		return errorResponse(request, InternalError, fmt.Sprintf("Error loading tools: %v", err))
	}
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  map[string]interface{}{"tools": tools},
	}
}
