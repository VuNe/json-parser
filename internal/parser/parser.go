package parser

import (
	"strconv"

	"github.com/VuNe/json-parser/internal/lexer"
)

// Parser interface defines the contract for parsing JSON tokens.
type Parser interface {
	Parse() (JSONValue, error)
	ParseValue() (JSONValue, error)
}

// parser is the concrete implementation of the Parser interface.
type parser struct {
	lexer        lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
}

// New creates a new parser instance with the given lexer.
func New(l lexer.Lexer) Parser {
	p := &parser{lexer: l}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances both currentToken and peekToken.
func (p *parser) nextToken() {
	p.currentToken = p.peekToken
	var err error
	p.peekToken, err = p.lexer.NextToken()
	if err != nil {
		// For now, create an invalid token on lexer error
		p.peekToken = lexer.Token{
			Type:     lexer.INVALID,
			Value:    err.Error(),
			Position: p.lexer.Position(),
		}
	}
}

// Parse parses the complete JSON input and returns the parsed value.
func (p *parser) Parse() (JSONValue, error) {
	value, err := p.ParseValue()
	if err != nil {
		return nil, err
	}

	// Ensure we're at the end of input after parsing a valid value
	if p.currentToken.Type != lexer.EOF {
		return nil, NewParseError("expected EOF after JSON value", p.currentToken)
	}

	return value, nil
}

// ParseValue parses a JSON value (supports objects, arrays, and all primitive types).
func (p *parser) ParseValue() (JSONValue, error) {
	return p.parseValue()
}

// parseObject parses a JSON object with string key-value pairs.
func (p *parser) parseObject() (JSONValue, error) {
	if p.currentToken.Type != lexer.LEFT_BRACE {
		return nil, NewParseError("expected '{'", p.currentToken)
	}

	// Move past the opening brace
	p.nextToken()

	// Check if we hit EOF before finding the closing brace
	if p.currentToken.Type == lexer.EOF {
		return nil, NewParseError("expected '}'", p.currentToken)
	}

	obj := NewJSONObject()

	// Check if it's an empty object
	if p.currentToken.Type == lexer.RIGHT_BRACE {
		p.nextToken() // consume the closing brace
		return obj, nil
	}

	// Parse key-value pairs
	for {
		// Expect string key
		if p.currentToken.Type != lexer.STRING {
			return nil, NewParseError("expected string key", p.currentToken)
		}

		key := p.currentToken.Value
		p.nextToken()

		// Expect colon
		if p.currentToken.Type != lexer.COLON {
			return nil, NewParseError("expected ':'", p.currentToken)
		}
		p.nextToken()

		// Parse value (supports all JSON types)
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		// Check for comma or closing brace
		if p.currentToken.Type == lexer.RIGHT_BRACE {
			p.nextToken() // consume the closing brace
			break
		} else if p.currentToken.Type == lexer.COMMA {
			p.nextToken() // consume the comma

			// After comma, we must have another key-value pair or it's an error
			if p.currentToken.Type == lexer.RIGHT_BRACE {
				return nil, NewParseError("trailing comma not allowed", p.currentToken)
			}
		} else {
			return nil, NewParseError("expected ',' or '}'", p.currentToken)
		}
	}

	return obj, nil
}

// parseArray parses a JSON array with comma-separated values.
func (p *parser) parseArray() (JSONValue, error) {
	if p.currentToken.Type != lexer.LEFT_BRACKET {
		return nil, NewParseError("expected '['", p.currentToken)
	}

	// Move past the opening bracket
	p.nextToken()

	// Check if we hit EOF before finding the closing bracket
	if p.currentToken.Type == lexer.EOF {
		return nil, NewParseError("expected ']'", p.currentToken)
	}

	var arr []any

	// Check if it's an empty array
	if p.currentToken.Type == lexer.RIGHT_BRACKET {
		p.nextToken() // consume the closing bracket
		return arr, nil
	}

	// Parse array elements
	for {
		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		// Check for comma or closing bracket
		if p.currentToken.Type == lexer.RIGHT_BRACKET {
			p.nextToken() // consume the closing bracket
			break
		} else if p.currentToken.Type == lexer.COMMA {
			p.nextToken() // consume the comma

			// After comma, we must have another value or it's an error
			if p.currentToken.Type == lexer.RIGHT_BRACKET {
				return nil, NewParseError("trailing comma not allowed", p.currentToken)
			}
		} else {
			return nil, NewParseError("expected ',' or ']'", p.currentToken)
		}
	}

	return arr, nil
}

// parseValue parses a JSON value (supports objects, arrays, strings, numbers, booleans, and null).
func (p *parser) parseValue() (JSONValue, error) {
	switch p.currentToken.Type {
	case lexer.LEFT_BRACE:
		return p.parseObject()
	case lexer.LEFT_BRACKET:
		return p.parseArray()
	case lexer.STRING:
		value := p.currentToken.Value
		p.nextToken()
		return value, nil
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.BOOLEAN:
		return p.parseBoolean()
	case lexer.NULL:
		return p.parseNull()
	case lexer.EOF:
		return nil, NewParseError("unexpected end of input", p.currentToken)
	case lexer.INVALID, lexer.RIGHT_BRACE, lexer.RIGHT_BRACKET, lexer.COLON, lexer.COMMA:
		return nil, NewParseError("expected JSON value", p.currentToken)
	default:
		return nil, NewParseError("expected JSON value", p.currentToken)
	}
}

// parseNumber parses a JSON number token and returns the appropriate Go type.
func (p *parser) parseNumber() (JSONValue, error) {
	value := p.currentToken.Value
	p.nextToken()

	// Try to parse as integer first
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intVal, nil
	}

	// If integer parsing fails, try float64
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal, nil
	}

	// If both fail, return error
	return nil, NewParseError("invalid number format", p.currentToken)
}

// parseBoolean parses a JSON boolean token.
func (p *parser) parseBoolean() (JSONValue, error) {
	value := p.currentToken.Value
	p.nextToken()

	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return nil, NewParseError("invalid boolean value", p.currentToken)
	}
}

// parseNull parses a JSON null token.
func (p *parser) parseNull() (JSONValue, error) {
	value := p.currentToken.Value
	p.nextToken()

	if value == "null" {
		return nil, nil
	}

	return nil, NewParseError("invalid null value", p.currentToken)
}
