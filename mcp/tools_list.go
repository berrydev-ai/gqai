package mcp

import (
	"fmt"
	"github.com/fotoetienne/gqai/graphql"
	"github.com/fotoetienne/gqai/tool"
)

// ToolsList handles the 'tools/list' MCP command.
func ToolsList(request JSONRPCRequest, config *graphql.Config) JSONRPCResponse {
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
