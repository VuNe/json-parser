package parser

import (
	"fmt"
	"strings"

	"github.com/VuNe/json-parser/internal/lexer"
)

// ErrorType represents the type of parsing error.
type ErrorType int

const (
	LexicalError  ErrorType = iota // Invalid characters, malformed strings, numbers
	SyntaxError                    // Unexpected tokens, missing brackets, colons
	SemanticError                  // Duplicate keys, invalid values in context
)

// String returns a human-readable representation of the error type.
func (et ErrorType) String() string {
	switch et {
	case LexicalError:
		return "Lexical"
	case SyntaxError:
		return "Syntax"
	case SemanticError:
		return "Semantic"
	default:
		return "Unknown"
	}
}

// ParseError represents an enhanced error that occurred during parsing.
type ParseError struct {
	Type        ErrorType
	Message     string
	Position    lexer.Position
	Token       lexer.Token
	Expected    []string // What was expected
	Found       string   // What was actually found
	JSONSnippet string   // Snippet of JSON around the error
	Suggestion  string   // Recovery suggestion
	SourceInput string   // Original input for context
}

// Error implements the error interface with enhanced formatting.
func (e ParseError) Error() string {
	var parts []string

	// Start with error type and basic message
	parts = append(parts, fmt.Sprintf("%s error at %s: %s", e.Type, e.Position, e.Message))

	// Add expected vs found context
	if len(e.Expected) > 0 && e.Found != "" {
		expectedStr := strings.Join(e.Expected, " or ")
		parts = append(parts, fmt.Sprintf("Expected %s, but found %s", expectedStr, e.Found))
	}

	// Add JSON snippet with position marker
	if e.JSONSnippet != "" {
		parts = append(parts, fmt.Sprintf("Near: %s", e.JSONSnippet))
	}

	// Add recovery suggestion
	if e.Suggestion != "" {
		parts = append(parts, fmt.Sprintf("Suggestion: %s", e.Suggestion))
	}

	return strings.Join(parts, "\n")
}

// NewParseError creates a basic ParseError (backward compatibility).
func NewParseError(message string, token lexer.Token) *ParseError {
	return &ParseError{
		Type:     SyntaxError,
		Message:  message,
		Position: token.Position,
		Token:    token,
		Found:    token.Type.String(),
	}
}

// NewDetailedParseError creates an enhanced ParseError with full context.
func NewDetailedParseError(errorType ErrorType, message string, token lexer.Token, expected []string, suggestion string, sourceInput string) *ParseError {
	pe := &ParseError{
		Type:        errorType,
		Message:     message,
		Position:    token.Position,
		Token:       token,
		Expected:    expected,
		Found:       fmt.Sprintf("'%s' (%s)", token.Value, token.Type),
		Suggestion:  suggestion,
		SourceInput: sourceInput,
	}

	// Generate JSON snippet with position marker
	pe.JSONSnippet = pe.generateJSONSnippet()

	return pe
}

// NewLexicalError creates a lexical error (for lexer errors).
func NewLexicalError(message string, token lexer.Token, suggestion string, sourceInput string) *ParseError {
	return NewDetailedParseError(LexicalError, message, token, nil, suggestion, sourceInput)
}

// NewSyntaxError creates a syntax error with expected tokens.
func NewSyntaxError(message string, token lexer.Token, expected []string, suggestion string, sourceInput string) *ParseError {
	return NewDetailedParseError(SyntaxError, message, token, expected, suggestion, sourceInput)
}

// NewSemanticError creates a semantic error.
func NewSemanticError(message string, token lexer.Token, suggestion string, sourceInput string) *ParseError {
	return NewDetailedParseError(SemanticError, message, token, nil, suggestion, sourceInput)
}

// generateJSONSnippet creates a snippet of JSON around the error position with a position marker.
func (e *ParseError) generateJSONSnippet() string {
	if e.SourceInput == "" {
		return ""
	}

	lines := strings.Split(e.SourceInput, "\n")
	if e.Position.Line < 1 || e.Position.Line > len(lines) {
		return ""
	}

	lineIdx := e.Position.Line - 1
	line := lines[lineIdx]

	// Create a snippet showing the problematic line with a pointer
	var snippet strings.Builder

	// Add line number and content
	snippet.WriteString(fmt.Sprintf("%d| %s\n", e.Position.Line, line))

	// Add pointer line showing where the error occurred
	pointer := strings.Repeat(" ", len(fmt.Sprintf("%d| ", e.Position.Line)))
	if e.Position.Column > 0 && e.Position.Column <= len(line) {
		pointer += strings.Repeat(" ", e.Position.Column-1) + "^"
	}
	snippet.WriteString(pointer)

	return snippet.String()
}

// Common error suggestions
const (
	SuggestionMissingColon        = "Add a ':' after the object key"
	SuggestionMissingComma        = "Add a ',' between object properties or array elements"
	SuggestionRemoveTrailingComma = "Remove the trailing comma before closing '}' or ']'"
	SuggestionCloseString         = "Add a closing quote '\"' to terminate the string"
	SuggestionEscapeCharacter     = "Escape special characters with a backslash '\\'"
	SuggestionValidNumber         = "Ensure numbers follow JSON format (no leading zeros, proper decimal/exponent notation)"
	SuggestionCloseObject         = "Add a closing '}' to complete the object"
	SuggestionCloseArray          = "Add a closing ']' to complete the array"
	SuggestionStringKey           = "Object keys must be strings enclosed in double quotes"
	SuggestionValidKeyword        = "Use lowercase for JSON keywords: 'true', 'false', 'null'"
)
