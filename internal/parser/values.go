package parser

// JSONValue represents a JSON value of any type.
type JSONValue any

// EmptyObject represents an empty JSON object {}.
type EmptyObject map[string]any

// NewEmptyObject creates a new empty JSON object.
func NewEmptyObject() EmptyObject {
	return make(EmptyObject)
}
