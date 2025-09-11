package parser

import (
	"fmt"
	"testing"

	"github.com/VuNe/json-parser/internal/lexer"
)

// TestParser_EdgeCasesForCoverage tests edge cases to improve test coverage
func TestParser_EdgeCasesForCoverage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "deeply nested objects",
			input:       `{"a": {"b": {"c": {"d": {"e": "deep"}}}}}`,
			expectError: false,
		},
		{
			name:        "deeply nested arrays",
			input:       `[[[[["deep"]]]]]`,
			expectError: false,
		},
		{
			name:        "complex mixed nesting",
			input:       `{"users": [{"name": "John", "data": [1, 2, {"nested": true}]}]}`,
			expectError: false,
		},
		{
			name:        "array at root level",
			input:       `[1, 2, 3, "text", true, null, {"obj": "value"}]`,
			expectError: false,
		},
		{
			name:        "string at root level",
			input:       `"just a string"`,
			expectError: false,
		},
		{
			name:        "number at root level",
			input:       `42`,
			expectError: false,
		},
		{
			name:        "boolean at root level",
			input:       `true`,
			expectError: false,
		},
		{
			name:        "null at root level",
			input:       `null`,
			expectError: false,
		},
		{
			name:        "scientific notation edge cases",
			input:       `{"small": 1e-10, "large": 1E+100, "zero_exp": 123e0}`,
			expectError: false,
		},
		{
			name:        "unicode in keys and values",
			input:       `{"🔑": "🌍", "\u0041key": "value\u0042"}`,
			expectError: false,
		},
		{
			name:        "empty arrays and objects mixed",
			input:       `{"empty_obj": {}, "empty_arr": [], "nested": {"also_empty": []}}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewWithInput(l, tt.input) // Use NewWithInput for enhanced error reporting

			result, err := p.Parse()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Note: result can be nil for valid null values, which is expected
				if result == nil && tt.input != `null` && tt.input != `"null"` {
					t.Error("expected non-nil result for valid input")
				}
			}
		})
	}
}

// TestParser_ErrorRecoveryAndSuggestions tests the enhanced error reporting
func TestParser_ErrorRecoveryAndSuggestions(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
	}{
		{
			name:         "missing colon basic error",
			input:        `{"key" "value"}`,
			expectedType: "Syntax error",
		},
		{
			name:         "trailing comma basic error", 
			input:        `{"key": "value",}`,
			expectedType: "Syntax error",
		},
		{
			name:         "unterminated object basic error",
			input:        `{"key": "value"`,
			expectedType: "Syntax error",
		},
		{
			name:         "multiline error basic",
			input:        "{\n  \"key\":\n  \"value\",\n}",
			expectedType: "Syntax error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewWithInput(l, tt.input)

			_, err := p.Parse()

			if err == nil {
				t.Error("expected error but got none")
				return
			}

			errorStr := err.Error()
			
			// Check error type
			if !containsSubstring(errorStr, tt.expectedType) {
				t.Errorf("expected error type %q in error: %v", tt.expectedType, errorStr)
			}

			// Check that error has position information
			if !containsSubstring(errorStr, "line") && !containsSubstring(errorStr, "column") {
				t.Errorf("expected position information in error: %v", errorStr)
			}
		})
	}
}

// TestParser_MemoryAndPerformance tests performance characteristics
func TestParser_MemoryAndPerformance(t *testing.T) {
	// Test with a reasonably large but not excessive JSON
	largeJSON := generateLargeJSON(100) // 100 nested levels
	
	l := lexer.New(largeJSON)
	p := New(l)

	result, err := p.Parse()
	if err != nil {
		t.Errorf("unexpected error parsing large JSON: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result")
	}
}

// TestParser_ConcurrentParsing tests thread safety (even though our parser isn't explicitly thread-safe)
func TestParser_ConcurrentParsing(t *testing.T) {
	inputs := []string{
		`{"test": 1}`,
		`[1, 2, 3]`,
		`"simple string"`,
		`true`,
		`null`,
		`42.5`,
	}

	// Test multiple parsers can run without interfering
	for i, input := range inputs {
		t.Run(fmt.Sprintf("concurrent_%d", i), func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			
			l := lexer.New(input)
			p := New(l)

			result, err := p.Parse()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Note: null values result in nil, which is expected
			if result == nil && input != `null` {
				t.Error("expected non-nil result for non-null input")
			}
		})
	}
}

// TestParser_ParseValueOnly tests ParseValue method directly
func TestParser_ParseValueOnly(t *testing.T) {
	tests := []string{
		`{}`,
		`[]`,
		`"string"`,
		`123`,
		`true`,
		`null`,
	}

	for _, input := range tests {
		t.Run("parse_value_"+input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)

			result, err := p.ParseValue()
			if err != nil {
				t.Errorf("ParseValue() returned error: %v", err)
			}
			// Note: null values result in nil, which is expected
			if result == nil && input != `null` {
				t.Error("expected non-nil result for non-null input")
			}
		})
	}
}

// Helper function to generate large nested JSON for performance testing
func generateLargeJSON(depth int) string {
	if depth == 0 {
		return `"leaf"`
	}
	return `{"level": ` + generateLargeJSON(depth-1) + `}`
}
