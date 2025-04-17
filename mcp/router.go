package mcp

import (
	"fmt"
	"github.com/fotoetienne/gqai/graphql"
	"log"
)

func RouteMCPRequest(request JSONRPCRequest, config *graphql.Config) JSONRPCResponse {
	switch request.Method {

	case initializeMethod:
		return mcpInitialize(request)

	case "notifications/initialized", "initialized":
		log.Printf("Server initialized successfully")
		return JSONRPCResponse{}

	case "tools/list":
		return ToolsList(request, config)

	case "tools/call":
		return ToolsCall(request, config)

	default:
		return errorResponse(request, MethodNotFound, fmt.Sprintf("Method '%s' not found", request.Method))
	}
}
