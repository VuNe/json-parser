package lexer

import "fmt"

// Position represents a position in the source text with line and column numbers.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // 0-based byte offset
}

// String returns a human-readable representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}
