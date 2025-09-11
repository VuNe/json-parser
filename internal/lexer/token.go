package lexer

import "fmt"

// TokenType represents the type of a token.
type TokenType int

const (
	// Special tokens
	INVALID TokenType = iota
	EOF

	// JSON structural tokens
	LEFT_BRACE    // {
	RIGHT_BRACE   // }
	LEFT_BRACKET  // [
	RIGHT_BRACKET // ]
	COLON         // :
	COMMA         // ,

	// JSON value tokens (for future use)
	STRING  // "string"
	NUMBER  // 123, 123.45
	BOOLEAN // true, false
	NULL    // null
)

// Token represents a token with its type, value, and position.
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

// String returns a string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case INVALID:
		return "INVALID"
	case EOF:
		return "EOF"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case NULL:
		return "NULL"
	default:
		return fmt.Sprintf("TokenType(%d)", int(t))
	}
}

// String returns a string representation of the token.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %s", t.Type, t.Value, t.Position)
}
