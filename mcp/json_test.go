package mcp

import (
	"strings"
	"testing"
)

func TestPrettyJSON(t *testing.T) {
	testData := map[string]interface{}{
		"name":    "test",
		"value":   42,
		"nested":  map[string]interface{}{"key": "value"},
		"array":   []int{1, 2, 3},
	}

	result := PrettyJSON(testData)

	// Check that result is not empty
	if result == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Check that result contains expected content
	if !strings.Contains(result, `"name": "test"`) {
		t.Error("Expected result to contain name field")
	}
	if !strings.Contains(result, `"value": 42`) {
		t.Error("Expected result to contain value field")
	}
	if !strings.Contains(result, `"nested"`) {
		t.Error("Expected result to contain nested object")
	}
	if !strings.Contains(result, `"array"`) {
		t.Error("Expected result to contain array")
	}

	// Check that result has proper indentation (contains newlines and spaces)
	if !strings.Contains(result, "\n") {
		t.Error("Expected result to have newlines for formatting")
	}
	if !strings.Contains(result, "    ") {
		t.Error("Expected result to have indentation spaces")
	}
}

func TestPrettyJSONWithInvalidData(t *testing.T) {
	// Test with data that can't be marshaled (function)
	invalidData := func() {}

	result := PrettyJSON(invalidData)

	// Should return empty string for invalid data
	if result != "" {
		t.Errorf("Expected empty string for invalid data, got: %s", result)
	}
}

func TestJSONEscapedString(t *testing.T) {
	testData := map[string]interface{}{
		"html":    "<script>alert('xss')</script>",
		"quotes":  `"double" 'single'`,
		"normal":  "safe text",
		"unicode": "café",
	}

	result, err := JSONEscapedString(testData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that result is not empty
	if result == "" {
		t.Error("Expected non-empty escaped JSON string")
	}

	// Check that HTML characters are escaped
	if !strings.Contains(result, `\u003cscript\u003e`) {
		t.Error("Expected HTML tags to be escaped")
	}
	if !strings.Contains(result, `\u0026`) || strings.Contains(result, "&") {
		// json.HTMLEscape escapes & to \u0026
		t.Error("Expected ampersands to be escaped")
	}
}

func TestJSONEscapedStringWithInvalidData(t *testing.T) {
	// Test with data that can't be marshaled (function)
	invalidData := func() {}

	result, err := JSONEscapedString(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data")
	}
	if result != "" {
		t.Errorf("Expected empty string for invalid data, got: %s", result)
	}
}

func TestJSONEscapedStringWithSpecialCharacters(t *testing.T) {
	specialData := map[string]interface{}{
		"lessThan":    "<",
		"greaterThan": ">",
		"ampersand":   "&",
		"quote":       "\"",
		"apostrophe":  "'",
	}

	result, err := JSONEscapedString(specialData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that special characters are properly escaped
	if !strings.Contains(result, `\u003c`) {
		t.Error("Expected < to be escaped")
	}
	if !strings.Contains(result, `\u003e`) {
		t.Error("Expected > to be escaped")
	}
	if !strings.Contains(result, `\u0026`) {
		t.Error("Expected & to be escaped")
	}
	if !strings.Contains(result, `\"`) {
		t.Error("Expected \" to be escaped")
	}
}
