# JSON Parser

A high-performance, production-ready JSON parser written in Go with comprehensive error reporting and validation.

## Features

- **Fast JSON Parsing**: Competitive performance with Go's standard library (within 17% for typical workloads)
- **Enhanced Error Reporting**: Detailed error messages with line/column positions, suggestions, and context
- **Comprehensive Validation**: Passes 20/20 official JSON test cases, correctly rejects 20/23 invalid cases
- **Full JSON Support**: Objects, arrays, strings, numbers, booleans, null, nested structures
- **Unicode Support**: Full UTF-8 and escape sequence handling including `\uXXXX`
- **Production Ready**: Extensive test coverage, benchmarked, and stress-tested

## Installation

```bash
go build -o json-parser ./cmd/json-parser
```

## Usage

### Command Line

```bash
# Parse and validate a JSON file
./json-parser example.json

# Exit codes:
# 0 = Valid JSON
# 1 = Invalid JSON or file error
```

### As a Library

```go
import (
    "github.com/VuNe/json-parser/internal/lexer"
    "github.com/VuNe/json-parser/internal/parser"
)

// Basic parsing
input := `{"name": "John", "age": 30}`
l := lexer.New(input)
p := parser.New(l)
result, err := p.Parse()

// Enhanced error reporting  
p := parser.NewWithInput(l, input)
result, err := p.Parse()
if err != nil {
    fmt.Printf("Parse error: %v\n", err) // Includes line/column info and suggestions
}
```

## Architecture

The parser follows a clean 3-layer architecture:

1. **Lexer**: Tokenizes raw JSON input with position tracking
2. **Parser**: Recursive descent parser that builds structured data
3. **CLI**: Command-line interface with file I/O and exit code management

## Performance

Benchmarked against Go's standard library:
- **Speed**: ~3961 ns/op (vs stdlib 3372 ns/op) - 17% slower
- **Memory**: 768 B/op (vs stdlib 640 B/op) - 20% more memory  
- **Error Detection**: Very fast at 1100-2300 ns/op

The slight overhead provides significantly better error messages and diagnostic information.

## Error Handling

Enhanced error messages include:

```
Syntax error at line 2, column 15: missing colon after object key
Expected ':', but found '"'
Near: 2| "key" "value"
                ^
Suggestion: Add a ':' after the object key
```

## Supported JSON Features

- ✅ Objects with string keys
- ✅ Arrays with mixed types
- ✅ Strings with escape sequences (`\"`, `\\`, `\n`, `\t`, `\r`, `\b`, `\f`, `\/`)
- ✅ Unicode escapes (`\uXXXX`)
- ✅ Numbers (integers, floats, scientific notation)
- ✅ Booleans (`true`, `false`)
- ✅ Null values
- ✅ Nested structures (objects and arrays)
- ✅ Whitespace handling
- ❌ Comments (not part of JSON spec)
- ❌ Trailing commas (not part of JSON spec)
- ❌ Single quotes (not part of JSON spec)

## Testing

Comprehensive test suite includes:
- Unit tests for all components (75%+ coverage)
- Official JSON test suite (43 test cases)
- Performance benchmarks  
- Integration tests
- Edge case and stress testing

```bash
# Run all tests
make test

# Run benchmarks
go test ./internal/parser -bench=. -benchmem

# Run with coverage
go test -cover ./...
```

## Development

### Building
```bash
make help          # Show available commands
make test          # Run test suite
make clean         # Clean build artifacts
```

### Project Structure
```
├── cmd/json-parser/       # CLI application
├── internal/
│   ├── lexer/            # Tokenization
│   ├── parser/           # JSON grammar parsing  
│   └── cli/              # CLI interface
├── test/                 # Test files and data
└── docs/                 # Documentation
```

## License

This project is part of a JSON parser implementation challenge.