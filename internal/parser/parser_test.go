package parser

import (
	"testing"

	"github.com/VuNe/json-parser/internal/lexer"
)

func TestNew(t *testing.T) {
	l := lexer.New("{}")
	p := New(l)
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid empty object",
			input:       "{}",
			expectError: false,
		},
		{
			name:        "valid empty object with whitespace",
			input:       " { } ",
			expectError: false,
		},
		{
			name:        "valid empty object with newlines",
			input:       "{\n}",
			expectError: false,
		},
		{
			name:        "empty input",
			input:       "",
			expectError: true,
			errorMsg:    "unexpected end of input",
		},
		{
			name:        "only whitespace",
			input:       "   ",
			expectError: true,
			errorMsg:    "unexpected end of input",
		},
		{
			name:        "missing closing brace",
			input:       "{",
			expectError: true,
			errorMsg:    "expected '}'",
		},
		{
			name:        "missing opening brace",
			input:       "}",
			expectError: true,
			errorMsg:    "expected JSON value",
		},
		{
			name:        "extra content after valid JSON",
			input:       "{}extra",
			expectError: true,
			errorMsg:    "expected EOF after JSON value",
		},
		{
			name:        "simple key-value pair",
			input:       `{"key": "value"}`,
			expectError: false,
		},
		{
			name:        "invalid character",
			input:       "a",
			expectError: true,
			errorMsg:    "expected JSON value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			value, err := p.Parse()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != "" {
					// Check if the error message contains the expected substring
					found := false
					if tt.errorMsg != "" {
						found = containsSubstring(err.Error(), tt.errorMsg)
					}
					if !found {
						t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// For valid cases, check the value type
				if value == nil {
					t.Error("expected non-nil value for successful parse")
				} else if _, ok := value.(JSONObject); !ok {
					t.Errorf("expected JSONObject, got %T", value)
				}
			}
		})
	}
}

func TestParser_ParseValue(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid empty object",
			input:       "{}",
			expectError: false,
		},
		{
			name:        "EOF",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid token",
			input:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			value, err := p.ParseValue()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if value == nil {
					t.Error("expected non-nil value")
				}
			}
		})
	}
}

func TestParseError_Error(t *testing.T) {
	token := lexer.Token{
		Type:     lexer.INVALID,
		Value:    "x",
		Position: lexer.Position{Line: 1, Column: 1, Offset: 0},
	}

	err := NewParseError("test error message", token)

	errorStr := err.Error()
	if errorStr == "" {
		t.Error("Error() returned empty string")
	}

	// Check that the error contains expected components
	expectedComponents := []string{"test error message", "line 1, column 1", "INVALID"}
	for _, component := range expectedComponents {
		if !containsSubstring(errorStr, component) {
			t.Errorf("error string %q should contain %q", errorStr, component)
		}
	}
}

func TestParser_KeyValuePairs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    map[string]any
	}{
		{
			name:        "single key-value pair",
			input:       `{"key": "value"}`,
			expectError: false,
			expected:    map[string]any{"key": "value"},
		},
		{
			name:        "multiple key-value pairs",
			input:       `{"key1": "value1", "key2": "value2"}`,
			expectError: false,
			expected:    map[string]any{"key1": "value1", "key2": "value2"},
		},
		{
			name:        "empty string values",
			input:       `{"key1": "", "key2": "value"}`,
			expectError: false,
			expected:    map[string]any{"key1": "", "key2": "value"},
		},
		{
			name:        "string with escape sequences",
			input:       `{"key": "hello\nworld"}`,
			expectError: false,
			expected:    map[string]any{"key": "hello\nworld"},
		},
		{
			name:        "key with escape sequences",
			input:       `{"key\twith\ttabs": "value"}`,
			expectError: false,
			expected:    map[string]any{"key\twith\ttabs": "value"},
		},
		{
			name:        "nested object",
			input:       `{"outer": {"inner": "value"}}`,
			expectError: false,
			expected: map[string]any{
				"outer": map[string]any{"inner": "value"},
			},
		},
		{
			name:        "object with whitespace",
			input:       `{ "key" : "value" }`,
			expectError: false,
			expected:    map[string]any{"key": "value"},
		},
		{
			name:        "object with newlines",
			input:       "{\n  \"key\": \"value\"\n}",
			expectError: false,
			expected:    map[string]any{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			value, err := p.Parse()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if value == nil {
					t.Error("expected non-nil value")
					return
				}

				obj, ok := value.(JSONObject)
				if !ok {
					t.Errorf("expected JSONObject, got %T", value)
					return
				}

				if !deepEqual(obj, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, obj)
				}
			}
		})
	}
}

func TestParser_KeyValuePairErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMsg string
		expectError bool
	}{
		{
			name:        "missing colon",
			input:       `{"key" "value"}`,
			expectedMsg: "expected ':'",
			expectError: true,
		},
		{
			name:        "missing comma",
			input:       `{"key1": "value1" "key2": "value2"}`,
			expectedMsg: "expected ',' or '}'",
			expectError: true,
		},
		{
			name:        "trailing comma",
			input:       `{"key": "value",}`,
			expectedMsg: "trailing comma not allowed",
			expectError: true,
		},
		{
			name:        "non-string key",
			input:       `{123: "value"}`,
			expectedMsg: "expected string key",
			expectError: true,
		},
		{
			name:        "missing value",
			input:       `{"key":}`,
			expectedMsg: "expected JSON value",
			expectError: true,
		},
		{
			name:        "missing closing brace",
			input:       `{"key": "value"`,
			expectedMsg: "expected ',' or '}'",
			expectError: true,
		},
		{
			name:        "empty key",
			input:       `{"": "value"}`,
			expectedMsg: "",
			expectError: false, // Empty keys are valid in JSON
		},
		{
			name:        "duplicate keys",
			input:       `{"key": "value1", "key": "value2"}`,
			expectedMsg: "",
			expectError: false, // Duplicate keys are handled by overwriting
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				// This is a valid case, skip error checking
				return
			}

			l := lexer.New(tt.input)
			p := New(l)

			_, err := p.Parse()

			if err == nil {
				t.Error("expected error but got none")
			} else if tt.expectedMsg != "" && !containsSubstring(err.Error(), tt.expectedMsg) {
				t.Errorf("expected error containing %q, got %q", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestParser_StringValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    `{"key": "hello"}`,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    `{"key": ""}`,
			expected: "",
		},
		{
			name:     "string with quotes",
			input:    `{"key": "hello \"world\""}`,
			expected: `hello "world"`,
		},
		{
			name:     "string with backslashes",
			input:    `{"key": "hello\\world"}`,
			expected: `hello\world`,
		},
		{
			name:     "string with newlines",
			input:    `{"key": "hello\nworld"}`,
			expected: "hello\nworld",
		},
		{
			name:     "string with unicode",
			input:    `{"key": "hello\u0041world"}`,
			expected: "helloAworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			value, err := p.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			obj, ok := value.(JSONObject)
			if !ok {
				t.Fatalf("expected JSONObject, got %T", value)
			}

			actual := obj["key"]
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestNewJSONObject(t *testing.T) {
	obj := NewJSONObject()
	if obj == nil {
		t.Error("NewJSONObject() returned nil")
	}

	if len(obj) != 0 {
		t.Errorf("NewJSONObject() should be empty, got length %d", len(obj))
	}
}

func TestNewEmptyObject(t *testing.T) {
	obj := NewEmptyObject()
	if obj == nil {
		t.Error("NewEmptyObject() returned nil")
	}

	if len(obj) != 0 {
		t.Errorf("NewEmptyObject() should be empty, got length %d", len(obj))
	}
}

// Helper function to deeply compare two map structures
func deepEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v1 := range a {
		v2, ok := b[k]
		if !ok {
			return false
		}

		if !deepEqualValue(v1, v2) {
			return false
		}
	}

	return true
}

// Helper function to compare any values recursively
func deepEqualValue(a, b any) bool {
	switch val1 := a.(type) {
	case string:
		val2, ok := b.(string)
		return ok && val1 == val2
	case map[string]any:
		val2, ok := b.(map[string]any)
		if ok {
			return deepEqual(val1, val2)
		}
		// Also check for JSONObject type
		if jsonObj, ok := b.(JSONObject); ok {
			return deepEqual(val1, map[string]any(jsonObj))
		}
		return false
	case JSONObject:
		// Convert JSONObject to map[string]any for comparison
		return deepEqualValue(map[string]any(val1), b)
	default:
		return a == b
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestParser_NumberParsing tests the parser's ability to parse number values correctly.
func TestParser_NumberParsing(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue any
		expectedType  string
	}{
		{
			name:          "positive integer",
			input:         `{"value": 123}`,
			expectedValue: int64(123),
			expectedType:  "int64",
		},
		{
			name:          "negative integer",
			input:         `{"value": -456}`,
			expectedValue: int64(-456),
			expectedType:  "int64",
		},
		{
			name:          "zero",
			input:         `{"value": 0}`,
			expectedValue: int64(0),
			expectedType:  "int64",
		},
		{
			name:          "positive float",
			input:         `{"value": 123.45}`,
			expectedValue: 123.45,
			expectedType:  "float64",
		},
		{
			name:          "negative float",
			input:         `{"value": -67.89}`,
			expectedValue: -67.89,
			expectedType:  "float64",
		},
		{
			name:          "float starting with zero",
			input:         `{"value": 0.123}`,
			expectedValue: 0.123,
			expectedType:  "float64",
		},
		{
			name:          "scientific notation positive exponent",
			input:         `{"value": 1.23e+10}`,
			expectedValue: 1.23e+10,
			expectedType:  "float64",
		},
		{
			name:          "scientific notation negative exponent",
			input:         `{"value": 1.23e-4}`,
			expectedValue: 1.23e-4,
			expectedType:  "float64",
		},
		{
			name:          "scientific notation uppercase E",
			input:         `{"value": 6.022E23}`,
			expectedValue: 6.022e23,
			expectedType:  "float64",
		},
		{
			name:          "integer scientific notation",
			input:         `{"value": 1E+10}`,
			expectedValue: 1e+10,
			expectedType:  "float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			result, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}

			obj, ok := result.(JSONObject)
			if !ok {
				t.Fatalf("Expected JSONObject, got %T", result)
			}

			value, exists := obj["value"]
			if !exists {
				t.Fatalf("Expected 'value' key not found")
			}

			// Check type and value
			switch tt.expectedType {
			case "int64":
				if intVal, ok := value.(int64); !ok {
					t.Errorf("Expected int64, got %T", value)
				} else if intVal != tt.expectedValue {
					t.Errorf("Expected %v, got %v", tt.expectedValue, intVal)
				}
			case "float64":
				if floatVal, ok := value.(float64); !ok {
					t.Errorf("Expected float64, got %T", value)
				} else if floatVal != tt.expectedValue {
					t.Errorf("Expected %v, got %v", tt.expectedValue, floatVal)
				}
			}
		})
	}
}

// TestParser_BooleanParsing tests the parser's ability to parse boolean values correctly.
func TestParser_BooleanParsing(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue bool
	}{
		{
			name:          "true value",
			input:         `{"active": true}`,
			expectedValue: true,
		},
		{
			name:          "false value",
			input:         `{"active": false}`,
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			result, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}

			obj, ok := result.(JSONObject)
			if !ok {
				t.Fatalf("Expected JSONObject, got %T", result)
			}

			value, exists := obj["active"]
			if !exists {
				t.Fatalf("Expected 'active' key not found")
			}

			boolVal, ok := value.(bool)
			if !ok {
				t.Errorf("Expected bool, got %T", value)
			} else if boolVal != tt.expectedValue {
				t.Errorf("Expected %v, got %v", tt.expectedValue, boolVal)
			}
		})
	}
}

// TestParser_NullParsing tests the parser's ability to parse null values correctly.
func TestParser_NullParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "null value",
			input: `{"optional": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			result, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}

			obj, ok := result.(JSONObject)
			if !ok {
				t.Fatalf("Expected JSONObject, got %T", result)
			}

			value, exists := obj["optional"]
			if !exists {
				t.Fatalf("Expected 'optional' key not found")
			}

			if value != nil {
				t.Errorf("Expected nil, got %v", value)
			}
		})
	}
}

// TestParser_MixedTypes tests the parser's ability to parse objects with multiple data types.
func TestParser_MixedTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name:  "mixed primitive types",
			input: `{"name": "John", "age": 30, "active": true, "data": null, "balance": -123.45}`,
			expected: map[string]any{
				"name":    "John",
				"age":     int64(30),
				"active":  true,
				"data":    nil,
				"balance": -123.45,
			},
		},
		{
			name:  "all number formats",
			input: `{"int": 42, "float": 3.14, "scientific": 1.23e-4, "zero": 0, "negative": -789}`,
			expected: map[string]any{
				"int":        int64(42),
				"float":      3.14,
				"scientific": 1.23e-4,
				"zero":       int64(0),
				"negative":   int64(-789),
			},
		},
		{
			name:  "all boolean combinations",
			input: `{"enabled": true, "disabled": false}`,
			expected: map[string]any{
				"enabled":  true,
				"disabled": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			result, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}

			obj, ok := result.(JSONObject)
			if !ok {
				t.Fatalf("Expected JSONObject, got %T", result)
			}

			// Check each expected key-value pair
			for key, expectedValue := range tt.expected {
				actualValue, exists := obj[key]
				if !exists {
					t.Errorf("Expected key '%s' not found", key)
					continue
				}

				if expectedValue == nil {
					if actualValue != nil {
						t.Errorf("Key '%s': expected nil, got %v", key, actualValue)
					}
				} else if actualValue != expectedValue {
					t.Errorf("Key '%s': expected %v (%T), got %v (%T)",
						key, expectedValue, expectedValue, actualValue, actualValue)
				}
			}

			// Check that we have the right number of keys
			if len(obj) != len(tt.expected) {
				t.Errorf("Expected %d keys, got %d", len(tt.expected), len(obj))
			}
		})
	}
}

// TestParser_TypeAssertions tests that we can properly assert types from parsed values.
func TestParser_TypeAssertions(t *testing.T) {
	input := `{"str": "hello", "num": 42, "float": 3.14, "bool": true, "null": null}`

	l := lexer.New(input)
	p := New(l)

	result, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	obj, ok := result.(JSONObject)
	if !ok {
		t.Fatalf("Expected JSONObject, got %T", result)
	}

	// Test string assertion
	if strVal, ok := obj["str"].(string); !ok || strVal != "hello" {
		t.Errorf("String assertion failed: got %v (%T)", obj["str"], obj["str"])
	}

	// Test int64 assertion
	if intVal, ok := obj["num"].(int64); !ok || intVal != 42 {
		t.Errorf("Int64 assertion failed: got %v (%T)", obj["num"], obj["num"])
	}

	// Test float64 assertion
	if floatVal, ok := obj["float"].(float64); !ok || floatVal != 3.14 {
		t.Errorf("Float64 assertion failed: got %v (%T)", obj["float"], obj["float"])
	}

	// Test bool assertion
	if boolVal, ok := obj["bool"].(bool); !ok || boolVal != true {
		t.Errorf("Bool assertion failed: got %v (%T)", obj["bool"], obj["bool"])
	}

	// Test nil assertion
	if obj["null"] != nil {
		t.Errorf("Nil assertion failed: got %v (%T)", obj["null"], obj["null"])
	}
}
