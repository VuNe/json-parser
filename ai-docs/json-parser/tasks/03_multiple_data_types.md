# Multiple Data Types Support (Step 3)

## Why

Supporting all JSON primitive data types (numbers, booleans, null) makes our parser capable of handling real-world JSON data. This step introduces type-specific parsing logic and validation, transforming our parser from a string-only tool into a comprehensive JSON processor. It establishes the foundation for proper JSON value representation and type safety that will be essential for nested structures in later steps.

## What

Extend the parser to support all JSON primitive data types as values: strings, numbers (integers and floats), booleans (`true`/`false`), and `null`. This includes implementing robust number parsing, boolean keyword recognition, and proper type representation in Go.

**Key Features to Implement:**
1. **Enhanced Lexer**: Tokenize numbers, boolean keywords (`true`, `false`), and `null`
2. **Number Parser**: Handle integers, floats, scientific notation, and edge cases
3. **Type System**: Proper Go type mapping for JSON values
4. **Value Parser**: Unified parsing logic for all JSON value types
5. **Enhanced Validation**: Type-specific error handling and validation

## Acceptance Criteria

**Functional Requirements:**
- [ ] Parser handles string values: `{"name": "John"}`
- [ ] Parser handles integer values: `{"age": 30}`
- [ ] Parser handles float values: `{"price": 19.99}`
- [ ] Parser handles boolean values: `{"active": true, "deleted": false}`
- [ ] Parser handles null values: `{"optional": null}`
- [ ] Parser handles mixed types: `{"name": "John", "age": 30, "active": true, "data": null}`
- [ ] All Step 1 and Step 2 test cases continue to pass

**Number Processing Requirements:**
- [ ] Support for positive and negative integers: `-42`, `123`
- [ ] Support for decimal numbers: `3.14`, `-0.5`
- [ ] Support for scientific notation: `1.23e-4`, `1E+10`
- [ ] Support for zero: `0`, `0.0`
- [ ] Rejection of invalid numbers: `01`, `3.`, `.5`, `infinity`
- [ ] Proper Go type mapping: integers as `int64`, floats as `float64`

**Boolean and Null Processing:**
- [ ] Recognition of `true` and `false` keywords
- [ ] Recognition of `null` keyword
- [ ] Case sensitivity enforcement (reject `True`, `FALSE`, `NULL`)
- [ ] Proper Go type mapping: `bool` for booleans, `nil` for null

**Lexer Enhancements:**
- [ ] New token types: `NUMBER`, `BOOLEAN`, `NULL`
- [ ] Number tokenization with format validation
- [ ] Keyword recognition logic
- [ ] Improved error reporting for invalid tokens
- [ ] Efficient token value extraction

**Parser Enhancements:**
- [ ] Value type detection and parsing dispatch
- [ ] Number parsing with proper error handling
- [ ] Boolean and null value processing
- [ ] Type-safe value storage in Go data structures
- [ ] Comprehensive value validation

**Data Structure Mapping:**
- [ ] JSON Objects → Go `map[string]interface{}`
- [ ] JSON Strings → Go `string`
- [ ] JSON Numbers → Go `int64` or `float64` (context-appropriate)
- [ ] JSON Booleans → Go `bool`
- [ ] JSON Null → Go `nil`

**Testing Requirements:**
- [ ] Unit tests for number tokenization and parsing (integers, floats, scientific)
- [ ] Unit tests for boolean and null parsing
- [ ] Unit tests for mixed-type objects
- [ ] Unit tests for number edge cases and error conditions
- [ ] Integration tests using provided Step 3 test files
- [ ] Type assertion tests for proper Go type mapping
- [ ] Performance tests for number parsing efficiency

**Error Handling Requirements:**
- [ ] Clear error messages for invalid number formats
- [ ] Position information for all parsing errors
- [ ] Distinction between lexical and semantic errors
- [ ] Graceful handling of unexpected tokens

**Quality Requirements:**
- [ ] Efficient number parsing without excessive string allocations
- [ ] Proper IEEE 754 floating-point handling
- [ ] No precision loss in number representation
- [ ] Maintained backward compatibility with previous steps
