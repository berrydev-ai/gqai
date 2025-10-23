package mcp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/berrydev-ai/gqai/graphql"
)

func TestSendResponse(t *testing.T) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      "test-id",
		Result:  map[string]interface{}{"status": "ok"},
	}

	sendResponse(encoder, response)

	// Check that something was written
	if buf.Len() == 0 {
		t.Error("Expected response to be written to buffer")
	}

	// Try to decode what was written
	var decoded JSONRPCResponse
	if err := json.NewDecoder(&buf).Decode(&decoded); err != nil {
		t.Errorf("Failed to decode written response: %v", err)
	}

	if decoded.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", decoded.JSONRPC)
	}
	if decoded.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", decoded.ID)
	}
}

func TestSendResponseWithError(t *testing.T) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	// Test with invalid data that can't be encoded
	invalidResponse := make(chan int) // channels can't be JSON encoded

	sendResponse(encoder, invalidResponse)

	// Should still write something (error will be logged but not returned)
	if buf.Len() == 0 {
		t.Error("Expected some output even for invalid response")
	}
}

func TestHandleInput(t *testing.T) {
	// This is a placeholder function, so we just test it doesn't panic
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "test-method",
	}

	// Should not panic
	handleInput(request)
}

// TestRunMCPStdIO is difficult to test directly since it reads from stdin in a loop
// We'll test the components it uses instead
func TestRunMCPStdIOComponents(t *testing.T) {
	config := &graphql.GraphQLConfig{
		SingleProject: &graphql.GraphQLProject{
			Schema: []graphql.SchemaPointer{
				{URL: "http://example.com/graphql"},
			},
		},
	}

	// Test that RouteMCPRequest works (used in the loop)
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2025-03-26",
		},
	}

	response := RouteMCPRequest(request, config)

	if response.Error != nil {
		t.Errorf("Expected no error from RouteMCPRequest, got %v", response.Error)
	}
}

// Test the JSON decoding part that happens in the loop
func TestJSONDecodingInLoop(t *testing.T) {
	// Create a JSON string that represents a valid request
	requestJSON := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"method": "initialize",
		"params": {
			"protocolVersion": "2025-03-26"
		}
	}`

	decoder := json.NewDecoder(strings.NewReader(requestJSON))

	var request JSONRPCRequest
	err := decoder.Decode(&request)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if request.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", request.JSONRPC)
	}
	if request.Method != "initialize" {
		t.Errorf("Expected method 'initialize', got '%s'", request.Method)
	}
}

func TestJSONRPCVersionCheck(t *testing.T) {
	// Test the version check logic from the loop
	request := JSONRPCRequest{
		JSONRPC: "1.0", // Invalid version
		ID:      "test-id",
		Method:  "test-method",
	}

	// This mimics the check in the loop
	if request.JSONRPC != "2.0" {
		// Should create an error response
		response := errorResponse(request, InvalidRequest, "Only JSON-RPC 2.0 is supported")

		if response.Error == nil {
			t.Error("Expected error for invalid JSON-RPC version")
		}
		if response.Error.Code != InvalidRequest {
			t.Errorf("Expected InvalidRequest error, got %d", response.Error.Code)
		}
	} else {
		t.Error("Expected version check to fail for '1.0'")
	}
}
