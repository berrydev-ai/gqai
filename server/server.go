package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/berrydev-ai/gqai/graphql"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SetupMCPServer creates and configures the MCP server with tools from GraphQL config
func SetupMCPServer(config *graphql.GraphQLConfig) (*server.MCPServer, error) {
	serverOpts := []server.ServerOption{
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	}

	s := server.NewMCPServer("gqai", "0.0.4", serverOpts...)

	// Load GraphQL operations and create tools
	ops, err := graphql.LoadOperations(config)
	if err != nil {
		return nil, fmt.Errorf("failed to load GraphQL operations: %w", err)
	}

	// Create tools from GraphQL operations
	for _, op := range ops {
		tool := createToolFromOperation(config, op)
		handler := createToolHandler(config, op)

		s.AddTool(tool, handler)
	}

	return s, nil
}

// createToolFromOperation creates an MCP tool from a GraphQL operation
func createToolFromOperation(config *graphql.GraphQLConfig, op *graphql.Operation) mcp.Tool {
	// Create tool with basic description
	return mcp.NewTool(op.Name,
		mcp.WithDescription(fmt.Sprintf("Execute GraphQL %s operation: %s", op.OperationType, op.Name)),
		mcp.WithString("variables",
			mcp.Description("JSON string of GraphQL variables")),
	)
}

// createToolHandler creates a tool handler for a GraphQL operation
func createToolHandler(config *graphql.GraphQLConfig, op *graphql.Operation) server.ToolHandlerFunc {
	endpoint := config.SingleProject.Schema[0].URL
	headers := config.SingleProject.Schema[0].Headers

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get arguments from request
		args := request.GetArguments()
		if args == nil {
			args = make(map[string]interface{})
		}

		// Execute GraphQL operation
		result, err := graphql.Execute(endpoint, args, op, headers)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("GraphQL execution failed: %v", err)), nil
		}

		// Convert result to JSON string
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// StartServer starts the MCP server with the specified transport
func StartServer(s *server.MCPServer, transport string, address string) error {
	switch transport {
	case "stdio":
		fmt.Fprintf(os.Stderr, "Starting MCP server with stdio transport...\n")
		return server.ServeStdio(s)
	case "http":
		fmt.Fprintf(os.Stderr, "Starting MCP server with HTTP streaming transport on %s...\n", address)
		httpServer := server.NewStreamableHTTPServer(s)
		return httpServer.Start(address)
	case "sse":
		fmt.Fprintf(os.Stderr, "Starting MCP server with SSE transport on %s...\n", address)
		sseServer := server.NewSSEServer(s)
		return sseServer.Start(address)
	default:
		return fmt.Errorf("invalid transport type '%s'. Must be 'stdio', 'http', or 'sse'", transport)
	}
}
