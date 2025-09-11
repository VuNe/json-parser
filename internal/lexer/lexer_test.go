package lexer

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"simple braces", "{}"},
		{"with whitespace", " { } "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input)
			if lexer == nil {
				t.Fatal("New() returned nil")
			}
		})
	}
}

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []Token
	}{
		{
			name:  "empty object",
			input: "{}",
			expectedTokens: []Token{
				{Type: LEFT_BRACE, Value: "{", Position: Position{Line: 1, Column: 1, Offset: 0}},
				{Type: RIGHT_BRACE, Value: "}", Position: Position{Line: 1, Column: 2, Offset: 1}},
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 3, Offset: 2}},
			},
		},
		{
			name:  "empty object with whitespace",
			input: " { } ",
			expectedTokens: []Token{
				{Type: LEFT_BRACE, Value: "{", Position: Position{Line: 1, Column: 2, Offset: 1}},
				{Type: RIGHT_BRACE, Value: "}", Position: Position{Line: 1, Column: 4, Offset: 3}},
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 6, Offset: 5}},
			},
		},
		{
			name:  "empty object with newlines",
			input: "{\n}",
			expectedTokens: []Token{
				{Type: LEFT_BRACE, Value: "{", Position: Position{Line: 1, Column: 1, Offset: 0}},
				{Type: RIGHT_BRACE, Value: "}", Position: Position{Line: 2, Column: 1, Offset: 2}},
				{Type: EOF, Value: "", Position: Position{Line: 2, Column: 2, Offset: 3}},
			},
		},
		{
			name:  "empty string",
			input: "",
			expectedTokens: []Token{
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 1, Offset: 0}},
			},
		},
		{
			name:  "whitespace only",
			input: "   \t\n  ",
			expectedTokens: []Token{
				{Type: EOF, Value: "", Position: Position{Line: 2, Column: 3, Offset: 7}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.expectedTokens {
				token, err := l.NextToken()
				if err != nil {
					t.Fatalf("NextToken() error = %v", err)
				}

				if token.Type != expected.Type {
					t.Errorf("token %d: expected type %v, got %v", i, expected.Type, token.Type)
				}
				if token.Value != expected.Value {
					t.Errorf("token %d: expected value %q, got %q", i, expected.Value, token.Value)
				}
				if token.Position.Line != expected.Position.Line {
					t.Errorf("token %d: expected line %d, got %d", i, expected.Position.Line, token.Position.Line)
				}
				if token.Position.Column != expected.Position.Column {
					t.Errorf("token %d: expected column %d, got %d", i, expected.Position.Column, token.Position.Column)
				}
				if token.Position.Offset != expected.Position.Offset {
					t.Errorf("token %d: expected offset %d, got %d", i, expected.Position.Offset, token.Position.Offset)
				}
			}
		})
	}
}

func TestLexer_NextToken_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError bool
		expectedToken TokenType
	}{
		{
			name:          "invalid character",
			input:         "a",
			expectedError: true,
			expectedToken: INVALID,
		},
		{
			name:          "invalid character in braces",
			input:         "{a}",
			expectedError: false, // First token should be LEFT_BRACE
			expectedToken: LEFT_BRACE,
		},
		{
			name:          "control character",
			input:         "\x01",
			expectedError: true,
			expectedToken: INVALID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			token, err := l.NextToken()

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if token.Type != tt.expectedToken {
				t.Errorf("expected token type %v, got %v", tt.expectedToken, token.Type)
			}
		})
	}
}

func TestLexer_HasMore(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool // HasMore() result after each NextToken() call
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []bool{false},
		},
		{
			name:     "single token",
			input:    "{",
			expected: []bool{true, false},
		},
		{
			name:     "multiple tokens",
			input:    "{}",
			expected: []bool{true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.expected {
				hasMore := l.HasMore()
				if hasMore != expected {
					t.Errorf("HasMore() call %d: expected %v, got %v", i, expected, hasMore)
				}

				// Advance to next token for next iteration
				if hasMore || i == len(tt.expected)-1 {
					_, _ = l.NextToken()
				}
			}
		})
	}
}

func TestLexer_StringTokenization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Token
	}{
		{
			name:  "empty string",
			input: `""`,
			expected: Token{
				Type:     STRING,
				Value:    "",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "simple string",
			input: `"hello"`,
			expected: Token{
				Type:     STRING,
				Value:    "hello",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with escaped quote",
			input: `"hello \"world\""`,
			expected: Token{
				Type:     STRING,
				Value:    `hello "world"`,
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with escaped backslash",
			input: `"hello\\world"`,
			expected: Token{
				Type:     STRING,
				Value:    `hello\world`,
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with newline escape",
			input: `"hello\nworld"`,
			expected: Token{
				Type:     STRING,
				Value:    "hello\nworld",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with tab escape",
			input: `"hello\tworld"`,
			expected: Token{
				Type:     STRING,
				Value:    "hello\tworld",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with all basic escapes",
			input: `"test\"\\\b\f\n\r\t"`,
			expected: Token{
				Type:     STRING,
				Value:    "test\"\\\b\f\n\r\t",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with unicode escape",
			input: `"hello\u0041world"`,
			expected: Token{
				Type:     STRING,
				Value:    "helloAworld",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
		{
			name:  "string with forward slash escape",
			input: `"hello\/world"`,
			expected: Token{
				Type:     STRING,
				Value:    "hello/world",
				Position: Position{Line: 1, Column: 1, Offset: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			token, err := l.NextToken()

			if err != nil {
				t.Fatalf("NextToken() error = %v", err)
			}

			if token.Type != tt.expected.Type {
				t.Errorf("expected type %v, got %v", tt.expected.Type, token.Type)
			}
			if token.Value != tt.expected.Value {
				t.Errorf("expected value %q, got %q", tt.expected.Value, token.Value)
			}
			if token.Position.Line != tt.expected.Position.Line {
				t.Errorf("expected line %d, got %d", tt.expected.Position.Line, token.Position.Line)
			}
			if token.Position.Column != tt.expected.Position.Column {
				t.Errorf("expected column %d, got %d", tt.expected.Position.Column, token.Position.Column)
			}
		})
	}
}

func TestLexer_StringErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "unterminated string",
			input:         `"hello`,
			expectedError: "unterminated string",
		},
		{
			name:          "unterminated string with escape",
			input:         `"hello\`,
			expectedError: "unterminated string",
		},
		{
			name:          "invalid escape sequence",
			input:         `"hello\x"`,
			expectedError: "invalid escape sequence",
		},
		{
			name:          "incomplete unicode escape",
			input:         `"hello\u123"`,
			expectedError: "invalid Unicode escape sequence",
		},
		{
			name:          "invalid unicode escape",
			input:         `"hello\uGGGG"`,
			expectedError: "invalid Unicode escape sequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			token, err := l.NextToken()

			if err == nil {
				t.Error("expected error but got none")
			}
			if token.Type != INVALID {
				t.Errorf("expected INVALID token, got %v", token.Type)
			}
			if !containsSubstring(err.Error(), tt.expectedError) {
				t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestLexer_ColonCommaTokens(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []Token
	}{
		{
			name:  "colon token",
			input: ":",
			expectedTokens: []Token{
				{Type: COLON, Value: ":", Position: Position{Line: 1, Column: 1, Offset: 0}},
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 2, Offset: 1}},
			},
		},
		{
			name:  "comma token",
			input: ",",
			expectedTokens: []Token{
				{Type: COMMA, Value: ",", Position: Position{Line: 1, Column: 1, Offset: 0}},
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 2, Offset: 1}},
			},
		},
		{
			name:  "key-value structure tokens",
			input: `"key":"value"`,
			expectedTokens: []Token{
				{Type: STRING, Value: "key", Position: Position{Line: 1, Column: 1, Offset: 0}},
				{Type: COLON, Value: ":", Position: Position{Line: 1, Column: 6, Offset: 5}},
				{Type: STRING, Value: "value", Position: Position{Line: 1, Column: 7, Offset: 6}},
				{Type: EOF, Value: "", Position: Position{Line: 1, Column: 14, Offset: 13}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.expectedTokens {
				token, err := l.NextToken()
				if err != nil {
					t.Fatalf("NextToken() error = %v", err)
				}

				if token.Type != expected.Type {
					t.Errorf("token %d: expected type %v, got %v", i, expected.Type, token.Type)
				}
				if token.Value != expected.Value {
					t.Errorf("token %d: expected value %q, got %q", i, expected.Value, token.Value)
				}
				if token.Position.Line != expected.Position.Line {
					t.Errorf("token %d: expected line %d, got %d", i, expected.Position.Line, token.Position.Line)
				}
				if token.Position.Column != expected.Position.Column {
					t.Errorf("token %d: expected column %d, got %d", i, expected.Position.Column, token.Position.Column)
				}
			}
		})
	}
}

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		expected string
	}{
		{
			name:     "start position",
			position: Position{Line: 1, Column: 1, Offset: 0},
			expected: "line 1, column 1",
		},
		{
			name:     "middle position",
			position: Position{Line: 5, Column: 10, Offset: 42},
			expected: "line 5, column 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.position.String()
			if result != tt.expected {
				t.Errorf("Position.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// Helper function to check if a string contains a substring (already exists in parser_test.go)
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
