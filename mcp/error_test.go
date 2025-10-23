package mcp

import (
	"testing"
)

func TestErrorResponse(t *testing.T) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  "test-method",
	}

	response := errorResponse(request, InvalidRequest, "Test error message")

	// Check basic response structure
	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC to be '2.0', got '%s'", response.JSONRPC)
	}
	if response.ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got '%s'", response.ID)
	}
	if response.Result != nil {
		t.Error("Expected Result to be nil for error response")
	}

	// Check error structure
	if response.Error == nil {
		t.Fatal("Expected Error to be non-nil")
	}
	if response.Error.Code != InvalidRequest {
		t.Errorf("Expected error code to be %d, got %d", InvalidRequest, response.Error.Code)
	}
	if response.Error.Message != "Test error message" {
		t.Errorf("Expected error message to be 'Test error message', got '%s'", response.Error.Message)
	}
	if response.Error.Data != nil {
		t.Error("Expected error data to be nil")
	}
}

func TestErrorResponseWithDifferentCodes(t *testing.T) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      123,
		Method:  "test-method",
	}

	testCases := []struct {
		code    int
		message string
	}{
		{ParseError, "Parse error"},
		{InvalidRequest, "Invalid request"},
		{MethodNotFound, "Method not found"},
		{InvalidParams, "Invalid params"},
		{InternalError, "Internal error"},
	}

	for _, tc := range testCases {
		response := errorResponse(request, tc.code, tc.message)

		if response.Error.Code != tc.code {
			t.Errorf("Expected error code %d, got %d", tc.code, response.Error.Code)
		}
		if response.Error.Message != tc.message {
			t.Errorf("Expected error message '%s', got '%s'", tc.message, response.Error.Message)
		}
	}
}
