package lexer

import (
	"fmt"
	"unicode"
)

// Lexer interface defines the contract for tokenizing JSON input.
type Lexer interface {
	NextToken() (Token, error)
	HasMore() bool
	Position() Position
}

// lexer is the concrete implementation of the Lexer interface.
type lexer struct {
	input    string
	position Position
	current  int  // current position in input (points to current char)
	ch       byte // current char under examination
}

// New creates a new lexer instance for the given input string.
func New(input string) Lexer {
	l := &lexer{
		input: input,
		position: Position{
			Line:   1,
			Column: 1,
			Offset: 0,
		},
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input.
func (l *lexer) readChar() {
	if l.current >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.current]
	}

	// Update position tracking
	if l.current > 0 && l.input[l.current-1] == '\n' {
		l.position.Line++
		l.position.Column = 1
	} else if l.current > 0 {
		l.position.Column++
	}

	l.position.Offset = l.current
	l.current++
}

// skipWhitespace skips whitespace characters (space, tab, newline, carriage return).
func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken scans the input and returns the next token.
func (l *lexer) NextToken() (Token, error) {
	var tok Token

	l.skipWhitespace()

	// Capture the current position for the token
	tok.Position = l.position

	switch l.ch {
	case '{':
		tok = Token{Type: LEFT_BRACE, Value: string(l.ch), Position: l.position}
		l.readChar()
	case '}':
		tok = Token{Type: RIGHT_BRACE, Value: string(l.ch), Position: l.position}
		l.readChar()
	case 0:
		tok = Token{Type: EOF, Value: "", Position: l.position}
	default:
		// Check if it's a valid JSON character that we don't support yet
		if unicode.IsPrint(rune(l.ch)) {
			return Token{Type: INVALID, Value: string(l.ch), Position: l.position},
				fmt.Errorf("unexpected character '%c' at %s", l.ch, l.position)
		} else {
			return Token{Type: INVALID, Value: fmt.Sprintf("\\x%02x", l.ch), Position: l.position},
				fmt.Errorf("unexpected character '\\x%02x' at %s", l.ch, l.position)
		}
	}

	return tok, nil
}

// HasMore returns true if there are more tokens to process.
func (l *lexer) HasMore() bool {
	return l.ch != 0
}

// Position returns the current position in the input.
func (l *lexer) Position() Position {
	return l.position
}
