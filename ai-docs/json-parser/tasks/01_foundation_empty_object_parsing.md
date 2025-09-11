# Foundation + Empty Object Parsing (Step 1)

## Why

This task establishes the foundational architecture and implements the most basic JSON parsing capability. By starting with empty object parsing (`{}`), we create a solid base with the core 3-layer architecture (Lexer → Parser → CLI) while delivering immediate value. This approach allows us to validate our architectural decisions early and provides a working parser that can correctly identify valid vs invalid JSON, which is essential for all subsequent features.

## What

Implement the foundational JSON parser architecture with basic lexer, parser, and CLI components that can:

- Parse valid empty JSON objects (`{}`)
- Detect and reject invalid JSON input
- Return appropriate exit codes (0 for valid JSON, 1 for invalid)
- Provide basic error reporting

**Core Components to Implement:**
1. **Basic Lexer**: Tokenize `{`, `}`, whitespace, and detect invalid characters
2. **Basic Parser**: Validate empty object structure according to JSON grammar
3. **CLI Interface**: Handle file input/output and manage exit codes
4. **Error Handling**: Basic position tracking and error reporting
5. **Project Structure**: Set up internal packages (lexer, parser, cli)

## Acceptance Criteria

**Functional Requirements:**
- [ ] Parser successfully validates `{}` as valid JSON and returns exit code 0
- [ ] Parser correctly identifies invalid JSON (malformed brackets, extra characters) and returns exit code 1
- [ ] CLI accepts filename as command-line argument and reads file content
- [ ] CLI handles missing files gracefully with appropriate error messages
- [ ] Basic lexer can tokenize left brace, right brace, and EOF tokens
- [ ] Parser validates that tokens follow empty object grammar rules

**Technical Requirements:**
- [ ] Code organized into internal packages: `internal/lexer/`, `internal/parser/`, `internal/cli/`
- [ ] All core interfaces defined: `Lexer`, `Parser`, `CLIHandler`
- [ ] Basic error types implemented: `ParseError` with position information
- [ ] Token types defined: `LEFT_BRACE`, `RIGHT_BRACE`, `EOF`, `INVALID`

**Testing Requirements:**
- [ ] Unit tests for lexer tokenization of empty objects
- [ ] Unit tests for parser validation logic
- [ ] Unit tests for CLI exit code handling
- [ ] Integration tests using provided Step 1 test files
- [ ] Test coverage for error scenarios (malformed input)

**Quality Requirements:**
- [ ] Code follows Go best practices and project coding standards
- [ ] All functions have appropriate documentation
- [ ] No linter errors in golangci-lint
- [ ] Memory allocations minimized in core parsing paths
