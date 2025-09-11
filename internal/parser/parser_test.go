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
			name:        "non-empty object (not supported in Step 1)",
			input:       `{"key": "value"}`,
			expectError: true,
			errorMsg:    "non-empty objects are not supported in Step 1",
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
				} else if _, ok := value.(EmptyObject); !ok {
					t.Errorf("expected EmptyObject, got %T", value)
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

func TestNewEmptyObject(t *testing.T) {
	obj := NewEmptyObject()
	if obj == nil {
		t.Error("NewEmptyObject() returned nil")
	}

	if len(obj) != 0 {
		t.Errorf("NewEmptyObject() should be empty, got length %d", len(obj))
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
