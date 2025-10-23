package graphql

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExecuteWithHeaders tests that the headers are properly included in GraphQL requests
func TestExecuteWithHeaders(t *testing.T) {
	// Create a test server that checks for headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers were received
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header to be 'Bearer test-token', got '%s'", authHeader)
		}

		customHeader := r.Header.Get("Custom-Header")
		if customHeader != "test-value" {
			t.Errorf("Expected Custom-Header to be 'test-value', got '%s'", customHeader)
		}

		// Return a valid GraphQL response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"data": {"test": "success"}}`)
	}))
	defer server.Close()

	// Create a test operation
	op := &Operation{
		Name:          "TestQuery",
		Raw:           "query TestQuery { test }",
		OperationType: "query",
	}

	// Create headers
	headers := map[string]string{
		"Authorization": "Bearer test-token",
		"Custom-Header": "test-value",
	}

	// Execute the request with headers
	result, err := Execute(server.URL, nil, op, headers)
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}

	// Check that the result is correct
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be a map, got %T", result)
	}

	data, ok := resultMap["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result.data to be a map, got %T", resultMap["data"])
	}

	testValue, ok := data["test"].(string)
	if !ok || testValue != "success" {
		t.Fatalf("Expected result.data.test to be 'success', got %v", data["test"])
	}
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `invalid json`)
	}))
	defer server.Close()

	// Create a test operation
	op := &Operation{
		Name:          "TestQuery",
		Raw:           "query TestQuery { test }",
		OperationType: "query",
	}

	// Execute the request
	result, err := Execute(server.URL, nil, op, nil)
	if err == nil {
		t.Fatalf("Expected error for invalid JSON response, got nil")
	}
	if result != nil {
		t.Fatalf("Expected nil result for invalid JSON response, got %v", result)
	}
}

func TestExecuteWithNonOKStatus(t *testing.T) {
	// Create a test server that returns a 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "internal server error"}`)
	}))
	defer server.Close()

	// Create a test operation
	op := &Operation{
		Name:          "TestQuery",
		Raw:           "query TestQuery { test }",
		OperationType: "query",
	}

	// Execute the request
	result, err := Execute(server.URL, nil, op, nil)
	if err == nil {
		t.Fatalf("Expected error for non-OK status, got nil")
	}
	if result != nil {
		t.Fatalf("Expected nil result for non-OK status, got %v", result)
	}
}

func TestExecuteWithNetworkError(t *testing.T) {
	// Use an invalid URL to simulate network error
	op := &Operation{
		Name:          "TestQuery",
		Raw:           "query TestQuery { test }",
		OperationType: "query",
	}

	// Execute the request with invalid URL
	result, err := Execute("http://invalid.url.that.does.not.exist", nil, op, nil)
	if err == nil {
		t.Fatalf("Expected error for network failure, got nil")
	}
	if result != nil {
		t.Fatalf("Expected nil result for network failure, got %v", result)
	}
}
