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

	// Return the result wrapped in MCP response format
	if resp == nil {
		resp = map[string]any{}
	}

	// Check if the response is a valid JSON object
	if _, ok := resp.(map[string]any); !ok {
		return errorResponse(request, InternalError, fmt.Sprintf("Tool %v returned an invalid response", toolName))
	}

	// Convert the json response to a json-escaped string
	resp_text, err := JSONEscapedString(resp)
	if err != nil {
		return errorResponse(request, InternalError, fmt.Sprintf("Error converting response to JSON string: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: resp_text,
			},
		},
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
		Error:   nil,
	}
}
