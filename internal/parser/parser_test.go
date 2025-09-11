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
