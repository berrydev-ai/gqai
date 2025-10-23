package mcp

import (
	"testing"

	"github.com/berrydev-ai/gqai/graphql"
)

func TestRouteMCPRequest(t *testing.T) {
	// Create a minimal configuration for testing
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
			Documents: []string{"./operations"},
		},
	}

	// Test the initialize method
	initRequest := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "initialize",
		Params: map[string]interface{}{
			"rootUri":         "file:///workspace",
			"protocolVersion": "0.2.0", // Add the missing protocol version
		},
	}

	initResponse := RouteMCPRequest(initRequest, config)
	if initResponse.ID != "1" {
		t.Errorf("Expected response ID to be 1, got %s", initResponse.ID)
	}
	if initResponse.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC to be 2.0, got %s", initResponse.JSONRPC)
	}
	if initResponse.Error != nil {
		t.Errorf("Expected no error, got %v", initResponse.Error)
	}

	// Test an unsupported method
	unknownRequest := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "2",
		Method:  "unknown_method",
	}

	unknownResponse := RouteMCPRequest(unknownRequest, config)
	if unknownResponse.Error == nil {
		t.Error("Expected error for unknown method, got nil")
	}
	if unknownResponse.Error.Code != MethodNotFound {
		t.Errorf("Expected error code to be %d, got %d", MethodNotFound, unknownResponse.Error.Code)
	}
}
