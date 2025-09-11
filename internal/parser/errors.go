package parser

import (
	"fmt"

	"github.com/VuNe/json-parser/internal/lexer"
)

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Position lexer.Position
	Token    lexer.Token
}

// Error implements the error interface.
func (e ParseError) Error() string {
	return fmt.Sprintf("parse error at %s: %s (token: %s)", e.Position, e.Message, e.Token.Type)
}

// NewParseError creates a new ParseError.
func NewParseError(message string, token lexer.Token) *ParseError {
	return &ParseError{
		Message:  message,
		Position: token.Position,
		Token:    token,
	}
}
