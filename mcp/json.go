package mcp

import (
	"bytes"
	"encoding/json"
	"log"
)

// PrettyJSON takes any value and returns a formatted JSON string representation.
// If the input cannot be marshaled to JSON, it returns an error.
func PrettyJSON(v interface{}) string {
	// First marshal the object to JSON
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to marshal to JSON: %v", err)
		return ""
	}

	// Create a buffer for pretty printing
	var prettyJSON bytes.Buffer

	// Use json.Indent to format the JSON with standard indentation
	err = json.Indent(&prettyJSON, jsonBytes, "", "    ")
	if err != nil {
		log.Printf("failed to indent JSON: %v", err)
		return ""
	}

	return prettyJSON.String()
}

// JSONEscapedString takes any value and returns a JSON-escaped string representation.
func JSONEscapedString(v interface{}) (string, error) {
	// First marshal the object to JSON
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// Create a buffer for JSON escaping
	var escapedJSON bytes.Buffer
	// Use json.HTMLEscape to escape the JSON string
	json.HTMLEscape(&escapedJSON, jsonBytes)

	return escapedJSON.String(), nil
}
