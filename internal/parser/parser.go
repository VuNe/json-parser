package parser

import (
	"fmt"

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

// ParseValue parses a JSON value (for Step 1, only empty objects).
func (p *parser) ParseValue() (JSONValue, error) {
	switch p.currentToken.Type {
	case lexer.LEFT_BRACE:
		return p.parseObject()
	case lexer.EOF:
		return nil, NewParseError("unexpected end of input", p.currentToken)
	default:
		return nil, NewParseError("expected JSON value", p.currentToken)
	}
}

// parseObject parses a JSON object (for Step 1, only empty objects).
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

	// For Step 1, we only support empty objects
	// Check if the next token is a closing brace
	if p.currentToken.Type == lexer.RIGHT_BRACE {
		// This is an empty object {}
		p.nextToken() // consume the closing brace
		return NewEmptyObject(), nil
	}

	// If we get here, it means there's content inside the braces
	// For Step 1, this is not supported
	return nil, NewParseError("non-empty objects are not supported in Step 1", p.currentToken)
}

// expectToken checks if the current token matches the expected type and advances.
func (p *parser) expectToken(expected lexer.TokenType) error {
	if p.currentToken.Type != expected {
		return NewParseError(
			fmt.Sprintf("expected %s, got %s", expected, p.currentToken.Type),
			p.currentToken,
		)
	}
	p.nextToken()
	return nil
}
