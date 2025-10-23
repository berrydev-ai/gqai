package mcp

import (
	"testing"
)

func TestMcpInitialize(t *testing.T) {
	// Test successful initialization with supported protocol version
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params: map[string]interface{}{
			"protocolVersion": "2025-03-26",
		},
	}

	response := mcpInitialize(request)

	// Check response structure
	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC to be '2.0', got '%s'", response.JSONRPC)
	}
	if response.ID != "test-id" {
		t.Errorf("Expected ID to be 'test-id', got '%s'", response.ID)
	}
	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// Check result
	result, ok := response.Result.(InitializeResult)
	if !ok {
		t.Fatal("Expected result to be InitializeResult")
	}

	if result.ProtocolVersion != "2025-03-26" {
		t.Errorf("Expected protocol version '2025-03-26', got '%s'", result.ProtocolVersion)
	}

	if result.ServerInfo.Name != "gqai" {
		t.Errorf("Expected server name 'gqai', got '%s'", result.ServerInfo.Name)
	}

	if result.ServerInfo.Version != "0.0.4" {
		t.Errorf("Expected server version '0.0.4', got '%s'", result.ServerInfo.Version)
	}

	if result.Capabilities.Tools == nil {
		t.Error("Expected tools capabilities to be non-nil")
	}
}

func TestMcpInitializeWithUnsupportedVersion(t *testing.T) {
	// Test initialization with unsupported protocol version
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params: map[string]interface{}{
			"protocolVersion": "unsupported-version",
		},
	}

	response := mcpInitialize(request)

	// Should fallback to latest supported version
	result, ok := response.Result.(InitializeResult)
	if !ok {
		t.Fatal("Expected result to be InitializeResult")
	}

	if result.ProtocolVersion != "2025-03-26" {
		t.Errorf("Expected fallback to latest version '2025-03-26', got '%s'", result.ProtocolVersion)
	}
}

func TestMcpInitializeWithOlderSupportedVersion(t *testing.T) {
	// Test initialization with older supported protocol version
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
		},
	}

	response := mcpInitialize(request)

	result, ok := response.Result.(InitializeResult)
	if !ok {
		t.Fatal("Expected result to be InitializeResult")
	}

	if result.ProtocolVersion != "2024-11-05" {
		t.Errorf("Expected protocol version '2024-11-05', got '%s'", result.ProtocolVersion)
	}
}

func TestMcpInitializeMissingParams(t *testing.T) {
	// Test initialization with missing parameters
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params:  nil,
	}

	response := mcpInitialize(request)

	if response.Error == nil {
		t.Error("Expected error for missing params")
	}
	if response.Error.Code != InvalidParams {
		t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
	}
}

func TestMcpInitializeInvalidParamsFormat(t *testing.T) {
	// Test initialization with invalid parameters format
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params:  "invalid-params",
	}

	response := mcpInitialize(request)

	if response.Error == nil {
		t.Error("Expected error for invalid params format")
	}
	if response.Error.Code != InvalidParams {
		t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
	}
}

func TestMcpInitializeMissingProtocolVersion(t *testing.T) {
	// Test initialization with missing protocol version
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      "test-id",
		Method:  initializeMethod,
		Params: map[string]interface{}{
			"someOtherField": "value",
		},
	}

	response := mcpInitialize(request)

	if response.Error == nil {
		t.Error("Expected error for missing protocol version")
	}
	if response.Error.Code != InvalidParams {
		t.Errorf("Expected InvalidParams error, got %d", response.Error.Code)
	}
}
