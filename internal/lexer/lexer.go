package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
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
	case '[':
		tok = Token{Type: LEFT_BRACKET, Value: string(l.ch), Position: l.position}
		l.readChar()
	case ']':
		tok = Token{Type: RIGHT_BRACKET, Value: string(l.ch), Position: l.position}
		l.readChar()
	case ':':
		tok = Token{Type: COLON, Value: string(l.ch), Position: l.position}
		l.readChar()
	case ',':
		tok = Token{Type: COMMA, Value: string(l.ch), Position: l.position}
		l.readChar()
	case '"':
		return l.readString()
	case 0:
		tok = Token{Type: EOF, Value: "", Position: l.position}
	default:
		// Handle numbers, booleans, and null
		if l.ch == '-' || (l.ch >= '0' && l.ch <= '9') {
			return l.readNumber()
		} else if isAlpha(l.ch) {
			return l.readKeyword()
		} else {
			// Check if it's a valid JSON character that we don't support yet
			if unicode.IsPrint(rune(l.ch)) {
				return Token{Type: INVALID, Value: string(l.ch), Position: l.position},
					fmt.Errorf("unexpected character '%c' at %s", l.ch, l.position)
			} else {
				return Token{Type: INVALID, Value: fmt.Sprintf("\\x%02x", l.ch), Position: l.position},
					fmt.Errorf("unexpected character '\\x%02x' at %s", l.ch, l.position)
			}
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

// readString reads a JSON string token with escape sequence support.
func (l *lexer) readString() (Token, error) {
	position := l.position // Save the starting position
	var value []byte

	// Skip opening quote
	l.readChar()

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			if l.ch == 0 {
				return Token{Type: INVALID, Value: string(value), Position: position},
					fmt.Errorf("unterminated string at %s", position)
			}

			switch l.ch {
			case '"':
				value = append(value, '"')
			case '\\':
				value = append(value, '\\')
			case '/':
				value = append(value, '/')
			case 'b':
				value = append(value, '\b')
			case 'f':
				value = append(value, '\f')
			case 'n':
				value = append(value, '\n')
			case 'r':
				value = append(value, '\r')
			case 't':
				value = append(value, '\t')
			case 'u':
				// Handle Unicode escape sequence \uXXXX
				unicode, err := l.readUnicodeEscape()
				if err != nil {
					return Token{Type: INVALID, Value: string(value), Position: position}, err
				}
				value = append(value, unicode...)
			default:
				return Token{Type: INVALID, Value: string(value), Position: position},
					fmt.Errorf("invalid escape sequence '\\%c' at %s", l.ch, l.position)
			}
		} else {
			value = append(value, l.ch)
		}
		l.readChar()
	}

	if l.ch != '"' {
		return Token{Type: INVALID, Value: string(value), Position: position},
			fmt.Errorf("unterminated string at %s", position)
	}

	// Skip closing quote
	l.readChar()

	return Token{Type: STRING, Value: string(value), Position: position}, nil
}

// readUnicodeEscape reads a Unicode escape sequence \uXXXX and returns the UTF-8 bytes.
func (l *lexer) readUnicodeEscape() ([]byte, error) {
	l.readChar() // skip 'u'

	var hexDigits [4]byte
	for i := 0; i < 4; i++ {
		if l.ch == 0 {
			return nil, fmt.Errorf("incomplete Unicode escape sequence at %s", l.position)
		}
		if !isHexDigit(l.ch) {
			return nil, fmt.Errorf("invalid Unicode escape sequence '\\u%s' at %s", string(hexDigits[:i]), l.position)
		}
		hexDigits[i] = l.ch
		if i < 3 { // Don't advance past the last digit
			l.readChar()
		}
	}

	// Convert hex string to rune
	var codePoint rune
	for _, digit := range hexDigits {
		codePoint <<= 4
		switch {
		case digit >= '0' && digit <= '9':
			codePoint += rune(digit - '0')
		case digit >= 'A' && digit <= 'F':
			codePoint += rune(digit - 'A' + 10)
		case digit >= 'a' && digit <= 'f':
			codePoint += rune(digit - 'a' + 10)
		}
	}

	// Convert rune to UTF-8 bytes
	result := make([]byte, 4)
	n := utf8.EncodeRune(result, codePoint)
	return result[:n], nil
}

// isHexDigit returns true if the character is a valid hexadecimal digit.
func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F') || (ch >= 'a' && ch <= 'f')
}

// isAlpha returns true if the character is an alphabetic character.
func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit returns true if the character is a digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// readNumber reads a JSON number token with support for integers, floats, and scientific notation.
func (l *lexer) readNumber() (Token, error) {
	position := l.position // Save the starting position
	var value []byte

	// Handle optional minus sign
	if l.ch == '-' {
		value = append(value, l.ch)
		l.readChar()

		// After minus, we must have a digit
		if !isDigit(l.ch) {
			return Token{Type: INVALID, Value: string(value), Position: position},
				fmt.Errorf("invalid number format at %s", position)
		}
	}

	// Handle the integer part
	if l.ch == '0' {
		// If it starts with 0, it must be 0, 0.x, or 0ex (no leading zeros allowed)
		value = append(value, l.ch)
		l.readChar()

		// Check if there's an invalid leading zero (like 01, 02, etc.)
		if isDigit(l.ch) {
			return Token{Type: INVALID, Value: string(value), Position: position},
				fmt.Errorf("numbers cannot have leading zeros at %s", position)
		}
	} else {
		// Read all digits for the integer part
		for isDigit(l.ch) {
			value = append(value, l.ch)
			l.readChar()
		}
	}

	// Handle optional fractional part
	if l.ch == '.' {
		value = append(value, l.ch)
		l.readChar()

		// After decimal point, we must have at least one digit
		if !isDigit(l.ch) {
			return Token{Type: INVALID, Value: string(value), Position: position},
				fmt.Errorf("invalid number format: missing digits after decimal point at %s", position)
		}

		// Read all fractional digits
		for isDigit(l.ch) {
			value = append(value, l.ch)
			l.readChar()
		}
	}

	// Handle optional exponent part
	if l.ch == 'e' || l.ch == 'E' {
		value = append(value, l.ch)
		l.readChar()

		// Handle optional exponent sign
		if l.ch == '+' || l.ch == '-' {
			value = append(value, l.ch)
			l.readChar()
		}

		// After exponent marker (and optional sign), we must have at least one digit
		if !isDigit(l.ch) {
			return Token{Type: INVALID, Value: string(value), Position: position},
				fmt.Errorf("invalid number format: missing digits in exponent at %s", position)
		}

		// Read all exponent digits
		for isDigit(l.ch) {
			value = append(value, l.ch)
			l.readChar()
		}
	}

	return Token{Type: NUMBER, Value: string(value), Position: position}, nil
}

// readKeyword reads a JSON keyword (true, false, null).
func (l *lexer) readKeyword() (Token, error) {
	position := l.position // Save the starting position
	var value []byte

	// Read all alphabetic characters
	for isAlpha(l.ch) {
		value = append(value, l.ch)
		l.readChar()
	}

	keyword := string(value)

	// Validate the keyword
	switch keyword {
	case "true", "false":
		return Token{Type: BOOLEAN, Value: keyword, Position: position}, nil
	case "null":
		return Token{Type: NULL, Value: keyword, Position: position}, nil
	default:
		return Token{Type: INVALID, Value: keyword, Position: position},
			fmt.Errorf("invalid keyword '%s' at %s", keyword, position)
	}
}
