# Error Handling Guide

This guide explains how the JSON parser handles errors and how to interpret error messages.

## Error Types

The parser categorizes errors into three types:

### 1. Lexical Errors
Invalid characters or malformed tokens at the character level.

**Examples:**
- Invalid escape sequences: `{"text": "\q"}`
- Incomplete Unicode escapes: `{"text": "\u12"}`
- Unterminated strings: `{"text": "hello`
- Invalid number formats: `{"num": 123.}`

### 2. Syntax Errors  
Valid tokens in invalid arrangements according to JSON grammar.

**Examples:**
- Missing colons: `{"key" "value"}`
- Missing commas: `{"a": 1 "b": 2}`
- Trailing commas: `{"key": "value",}`
- Mismatched brackets: `{"array": [1, 2, 3}}`

### 3. Semantic Errors
Contextually inappropriate values (currently unused but available for future features).

## Enhanced Error Messages

When using `parser.NewWithInput()`, you get enhanced error messages with:

### Position Information
Every error includes precise line and column numbers:
```
Syntax error at line 2, column 15: missing colon after object key
```

### Expected vs Found Context
Clear indication of what was expected:
```
Expected ':', but found '"'
```

### JSON Snippet with Position Marker
Visual indication of the error location:
```
Near: 2| "key" "value"
               ^
```

### Recovery Suggestions
Helpful suggestions for common mistakes:
```
Suggestion: Add a ':' after the object key
```

## Common Error Patterns

### Missing Colon
```json
{"key" "value"}
```
**Error:** `Syntax error at line 1, column 8: missing colon after object key`
**Fix:** Add `:` between key and value: `{"key": "value"}`

### Trailing Comma
```json  
{"key": "value",}
```
**Error:** `Syntax error at line 1, column 17: trailing comma before object close`
**Fix:** Remove trailing comma: `{"key": "value"}`

### Unterminated String
```json
{"key": "unterminated
```
**Error:** `Syntax error at line 1, column 22: expected JSON value`  
**Fix:** Add closing quote: `{"key": "unterminated"}`

### Invalid Number
```json
{"number": 01}
```
**Error:** `Syntax error at line 1, column 13: expected JSON value`
**Fix:** Remove leading zero: `{"number": 1}` 

### Mismatched Brackets
```json
{"array": [1, 2, 3}
```
**Error:** `Syntax error at line 1, column 19: expected ',' or ']'`
**Fix:** Use correct bracket: `{"array": [1, 2, 3]}`

## Error Handling in Code

### Basic Error Handling
```go
l := lexer.New(input)  
p := parser.New(l)
result, err := p.Parse()
if err != nil {
    fmt.Printf("Parse failed: %v\n", err)
    return
}
```

### Enhanced Error Reporting
```go
l := lexer.New(input)
p := parser.NewWithInput(l, input)  // Enhanced version
result, err := p.Parse()
if err != nil {
    // Get detailed error with position info and suggestions
    fmt.Printf("JSON Error:\n%v\n", err)
    return
}
```

### Error Type Checking
```go
if parseErr, ok := err.(*parser.ParseError); ok {
    switch parseErr.Type {
    case parser.LexicalError:
        fmt.Println("Character or token level error")
    case parser.SyntaxError: 
        fmt.Println("Grammar or structure error")
    case parser.SemanticError:
        fmt.Println("Contextual or meaning error")
    }
}
```

## CLI Error Codes

The command-line interface uses standard exit codes:

- **0**: Success - JSON is valid
- **1**: Error - JSON is invalid or file cannot be read

## Best Practices

1. **Always check errors**: JSON parsing can fail for many reasons
2. **Use enhanced parsing**: `NewWithInput()` provides better diagnostics  
3. **Display full error messages**: They include helpful context and suggestions
4. **Handle file errors**: Check for file access issues separately from parse errors
5. **Log position information**: Line/column info helps with debugging large JSON files

## Performance Considerations

- Error detection is fast (1100-2300 ns/op)
- Enhanced error reporting adds minimal overhead
- Failed parsing stops immediately when first error is encountered
- Error messages are constructed only when errors occur

## Error Message Localization

Currently, all error messages are in English. The error structure supports localization:

```go
type ParseError struct {
    Type        ErrorType
    Message     string      // Localizable message
    Position    Position    // Universal position info
    Expected    []string    // Localizable expected tokens
    Suggestion  string      // Localizable suggestion
    // ...
}
```

This allows for future internationalization while preserving the structured error information.
