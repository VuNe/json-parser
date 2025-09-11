package parser

// JSONValue represents a JSON value of any type.
type JSONValue any

// JSONObject represents a JSON object with string keys.
type JSONObject map[string]any

// NewJSONObject creates a new JSON object.
func NewJSONObject() JSONObject {
	return make(JSONObject)
}

// EmptyObject represents an empty JSON object {}.
// This is kept for backward compatibility with Step 1.
type EmptyObject map[string]any

// NewEmptyObject creates a new empty JSON object.
func NewEmptyObject() EmptyObject {
	return make(EmptyObject)
}
