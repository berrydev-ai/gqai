package mcp

import (
	"fmt"
	"github.com/fotoetienne/gqai/graphql"
	"github.com/fotoetienne/gqai/tool"
)

// ToolsCall handles the 'tools/call' MCP command.
func ToolsCall(request JSONRPCRequest, config *graphql.Config) JSONRPCResponse {
	// Check if the tool name is provided in the request
	if request.Params == nil {
		return errorResponse(request, InvalidParams, "Params must include tool name")
	}
	toolName, ok := request.Params.(map[string]any)["name"].(string)
	if !ok {
		return errorResponse(request, InvalidParams, "Tool name is required")
	}

	// Load tool
	tool, err := tool.LoadTool(config, toolName)
	if err != nil {
		return errorResponse(request, InternalError, err.Error())
	}

	// Execute the tool with the provided input
	input, ok := request.Params.(map[string]any)["arguments"].(map[string]any)
	resp, err := tool.Execute(input)
	if err != nil {
		return errorResponse(request, InternalError, fmt.Sprintf("Error executing tool %v: %v", toolName, err))
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  resp,
		Error:   nil,
	}
}
