# Nested Objects and Arrays Support (Step 4)

## Why

Adding support for nested objects and arrays completes the core JSON parsing functionality, enabling our parser to handle complex, real-world JSON structures. This step introduces recursive parsing, which is fundamental to processing hierarchical data. It transforms our parser from a simple flat object processor into a full-featured JSON parser capable of handling any valid JSON structure, including deeply nested combinations of objects and arrays.

## What

Implement comprehensive support for nested JSON structures by adding array parsing capabilities and recursive object/array processing. This includes handling mixed nested structures like arrays of objects, objects containing arrays, and arbitrary nesting depth.

**Key Features to Implement:**
1. **Array Support**: Parse JSON arrays with proper bracket tokenization and element processing
2. **Recursive Parsing**: Implement recursive descent parsing for nested objects and arrays
3. **Mixed Structures**: Handle objects containing arrays and arrays containing objects
4. **Enhanced Lexer**: Add array bracket tokenization (`[`, `]`)
5. **Memory Management**: Efficient handling of nested data structures

## Acceptance Criteria

**Array Processing Requirements:**
- [ ] Parser handles empty arrays: `{"items": []}`
- [ ] Parser handles arrays with primitive values: `{"numbers": [1, 2, 3]}`
- [ ] Parser handles arrays with mixed types: `{"mixed": [1, "text", true, null]}`
- [ ] Parser handles arrays with string values: `{"names": ["Alice", "Bob", "Charlie"]}`
- [ ] Parser validates proper array syntax (brackets, comma separation)
- [ ] Parser rejects malformed arrays (missing brackets, trailing commas)

**Nested Object Requirements:**
- [ ] Parser handles objects within objects: `{"person": {"name": "John", "age": 30}}`
- [ ] Parser handles multiple levels of nesting: `{"a": {"b": {"c": "deep"}}}`
- [ ] Parser handles objects within arrays: `{"users": [{"name": "Alice"}, {"name": "Bob"}]}`
- [ ] Parser handles arrays within objects within arrays: `{"data": [{"items": [1, 2, 3]}]}`

**Complex Structure Requirements:**
- [ ] Parser handles arrays of objects: `[{"id": 1}, {"id": 2}]`
- [ ] Parser handles arrays of arrays: `[[1, 2], [3, 4]]`
- [ ] Parser handles deeply nested combinations with arbitrary depth
- [ ] Parser handles empty nested structures: `{"obj": {}, "arr": []}`
- [ ] Parser maintains type information throughout nested structures

**Lexer Enhancements:**
- [ ] New token types: `LEFT_BRACKET`, `RIGHT_BRACKET`
- [ ] Array bracket tokenization and validation
- [ ] Proper token stream handling for nested structures
- [ ] Position tracking through nested contexts
- [ ] Error reporting with nested context information

**Parser Enhancements:**
- [ ] Recursive parsing function for objects and arrays
- [ ] Parse context management for nested structures
- [ ] Array element parsing with type detection
- [ ] Proper comma handling in arrays and objects
- [ ] Stack management for parsing state

**Data Structure Mapping:**
- [ ] JSON Arrays → Go `[]interface{}`
- [ ] Nested objects maintain `map[string]interface{}` structure
- [ ] Mixed arrays preserve element types
- [ ] Proper interface{} usage for dynamic typing
- [ ] Memory-efficient nested structure allocation

**Recursive Parsing Logic:**
- [ ] `ParseObject()` method for object parsing
- [ ] `ParseArray()` method for array parsing  
- [ ] `ParseValue()` method that dispatches to appropriate parser
- [ ] Proper recursion depth handling
- [ ] State management across recursive calls

**Error Handling Requirements:**
- [ ] Context-aware error messages showing nesting level
- [ ] Position tracking through nested structures
- [ ] Stack overflow protection for deeply nested structures
- [ ] Clear error messages for mismatched brackets
- [ ] Proper error propagation from recursive calls

**Testing Requirements:**
- [ ] Unit tests for array tokenization and parsing
- [ ] Unit tests for nested object parsing at various depths
- [ ] Unit tests for mixed structure combinations
- [ ] Unit tests for malformed nested structures
- [ ] Integration tests using provided Step 4 test files
- [ ] Stress tests for deeply nested structures
- [ ] Memory usage tests for large nested structures
- [ ] All previous step tests continue to pass

**Performance Requirements:**
- [ ] Efficient memory allocation for nested structures
- [ ] Minimal stack usage in recursive parsing
- [ ] No memory leaks in complex nested structures
- [ ] Reasonable performance for deeply nested JSON (100+ levels)

**Quality Requirements:**
- [ ] Clean separation between object and array parsing logic
- [ ] Consistent error handling across all parsing contexts
- [ ] Maintainable recursive parsing implementation
- [ ] Proper resource cleanup in error conditions
