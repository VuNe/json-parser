# JSON Parser Architecture

## Overview

The JSON parser follows a clean 3-layer architecture designed for incremental development through 5 progressive complexity steps. The design emphasizes separation of concerns, extensibility, and maintainability.

## Architecture Diagram

```mermaid
graph TD
    A[CLI Interface] --> B[Parser]
    B --> C[Lexer]
    C --> D[Input Stream]
    
    B --> E[Token Stream]
    C --> E
    
    F[Error Handler] --> A
    B --> F
    C --> F
    
    G[JSON Value Types] --> B
    H[Token Types] --> C
    H --> B
    
    subgraph "Core Components"
        B
        C
        F
    end
    
    subgraph "Data Types"
        G
        H
    end
    
    subgraph "I/O Layer"
        A
        D
    end
```

## Component Details

### 1. Lexer (Tokenizer) Layer

**Responsibility**: Convert raw JSON string input into a stream of meaningful tokens.

**Key Interfaces**:
```go
type Lexer interface {
    NextToken() (Token, error)
    HasMore() bool
    Position() Position
}

type Token struct {
    Type     TokenType
    Value    string
    Position Position
}
```

**Token Types**:
- `LEFT_BRACE` (`{`)
- `RIGHT_BRACE` (`}`)
- `LEFT_BRACKET` (`[`)
- `RIGHT_BRACKET` (`]`)
- `COLON` (`:`)
- `COMMA` (`,`)
- `STRING` (quoted strings)
- `NUMBER` (integers/floats)
- `BOOLEAN` (`true`/`false`)
- `NULL` (`null`)
- `EOF` (end of input)
- `INVALID` (error token)

**Features**:
- Position tracking for error reporting
- String escape sequence handling
- Number format validation
- Whitespace skipping

### 2. Parser Layer

**Responsibility**: Process token stream according to JSON grammar rules and build structured data.

**Key Interfaces**:
```go
type Parser interface {
    Parse() (interface{}, error)
    ParseValue() (interface{}, error)
}

type JSONValue interface{}
```

**Parsing Strategy**:
- Recursive descent parser
- Token lookahead for decision making
- Context-aware error messages
- Progressive complexity through 5 steps

**Data Structure Mapping**:
- JSON Objects → Go `map[string]interface{}`
- JSON Arrays → Go `[]interface{}`
- JSON Strings → Go `string`
- JSON Numbers → Go `float64` or `int64`
- JSON Booleans → Go `bool`
- JSON Null → Go `nil`

### 3. CLI Interface Layer

**Responsibility**: Handle file I/O, coordinate parsing, and provide appropriate exit codes.

**Key Interface**:
```go
type CLIHandler interface {
    ParseFile(filename string) error
    ParseString(input string) error
    ExitCode() int
}
```

**Features**:
- File input handling
- Standard input support
- Exit code management (0=valid, 1=invalid)
- Error message formatting

### 4. Error Handling System

**Responsibility**: Provide detailed, position-aware error information.

**Key Types**:
```go
type Position struct {
    Line   int
    Column int
    Offset int
}

type ParseError struct {
    Message  string
    Position Position
    Token    Token
}
```

**Error Categories**:
- Lexical errors (invalid characters, malformed strings)
- Syntax errors (unexpected tokens, missing brackets)
- Semantic errors (duplicate keys, invalid values)

## Progressive Implementation Strategy

### Step 1: Empty Object Foundation
- **Lexer**: `{`, `}`, whitespace handling
- **Parser**: Empty object validation
- **CLI**: Basic file I/O and exit codes

### Step 2: String Key-Value Pairs
- **Lexer**: String tokenization, `:`, `,`
- **Parser**: Key-value pair parsing
- **Error Handling**: String validation

### Step 3: Multiple Data Types
- **Lexer**: Numbers, booleans, null
- **Parser**: Type-specific parsing logic
- **Value System**: All JSON primitive types

### Step 4: Nested Structures
- **Lexer**: Array brackets `[`, `]`
- **Parser**: Recursive parsing for objects and arrays
- **Memory Management**: Efficient nested structure handling

### Step 5: Comprehensive Validation
- **Error System**: Enhanced error reporting
- **Performance**: Memory optimization
- **Testing**: Complete test coverage

## Design Principles

### 1. Separation of Concerns
- Lexer handles character-level processing
- Parser handles grammar and structure
- CLI handles I/O and user interaction

### 2. Interface-Based Design
- Components interact through well-defined interfaces
- Easy to test and mock individual layers
- Extensible for future enhancements

### 3. Progressive Complexity
- Each step builds upon previous functionality
- Minimal changes required between steps
- Clear validation criteria for each step

### 4. Error-First Design
- Position tracking throughout parsing
- Meaningful error messages for debugging
- Graceful error recovery where possible

### 5. Performance Considerations
- Minimal memory allocations in hot paths
- Single-pass parsing where possible
- Efficient string handling

## Package Structure

```
internal/
├── lexer/
│   ├── lexer.go          # Core lexer implementation
│   ├── token.go          # Token types and utilities
│   └── position.go       # Position tracking
├── parser/
│   ├── parser.go         # Main parser logic
│   ├── values.go         # JSON value handling
│   └── errors.go         # Parser-specific errors
└── cli/
    ├── handler.go        # CLI interface implementation
    └── io.go            # File I/O utilities

cmd/
└── json-parser/
    └── main.go          # Application entry point
```

## Validation Against Implementation Plan

✅ **3-Layer Architecture**: Lexer → Parser → CLI Interface  
✅ **Progressive Steps**: Supports all 5 complexity levels  
✅ **Test-Driven**: Interface design enables comprehensive testing  
✅ **Error Handling**: Position-aware error reporting system  
✅ **Performance**: Minimal allocation design in hot paths  
✅ **Maintainability**: Clean separation and interface-based design  
✅ **Extensibility**: Easy to add new token types and parsing rules  

This architecture provides a solid foundation for implementing the JSON parser according to the specified plan while maintaining clean, testable, and maintainable code.
