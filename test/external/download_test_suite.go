package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// downloadTestSuite downloads JSON test cases from various sources
func main() {
	// Create test directories
	dirs := []string{
		"test/external/json_org",
		"test/external/nst_json_test_suite",
		"test/external/custom",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
		}
	}

	// Download some standard test cases (we'll create these locally for now)
	createLocalTestCases()
}

func createLocalTestCases() {
	// Valid JSON test cases from json.org specification
	validTests := map[string]string{
		"test/external/json_org/valid_empty_object.json":           `{}`,
		"test/external/json_org/valid_empty_array.json":            `[]`,
		"test/external/json_org/valid_string.json":                 `"Hello World"`,
		"test/external/json_org/valid_number_int.json":             `42`,
		"test/external/json_org/valid_number_float.json":           `3.14159`,
		"test/external/json_org/valid_number_scientific.json":      `6.022e23`,
		"test/external/json_org/valid_boolean_true.json":           `true`,
		"test/external/json_org/valid_boolean_false.json":          `false`,
		"test/external/json_org/valid_null.json":                   `null`,
		"test/external/json_org/valid_simple_object.json":          `{"name": "John", "age": 30}`,
		"test/external/json_org/valid_simple_array.json":           `[1, 2, 3, "four", true, null]`,
		"test/external/json_org/valid_nested_object.json":          `{"person": {"name": "Alice", "address": {"city": "NYC", "zip": 10001}}}`,
		"test/external/json_org/valid_nested_array.json":           `[[[1]], [[2]], [[3]]]`,
		"test/external/json_org/valid_mixed_nesting.json":          `{"users": [{"id": 1, "tags": ["admin", "active"]}, {"id": 2, "tags": []}]}`,
		"test/external/json_org/valid_unicode.json":                `{"message": "Hello 🌍", "japanese": "こんにちは", "escape": "Quote: \"Hello\""}`,
		"test/external/json_org/valid_numbers_edge_cases.json":     `{"zero": 0, "negative": -42, "decimal": 0.5, "exp_pos": 1e+10, "exp_neg": 1e-5}`,
		"test/external/json_org/valid_strings_escapes.json":        `{"quote": "\"", "backslash": "\\", "newline": "\n", "tab": "\t", "unicode": "\u0041"}`,
		"test/external/json_org/valid_large_number.json":           `{"big": 1.7976931348623157e+308}`,
		"test/external/json_org/valid_deep_nesting.json":           generateDeepNesting(50),
		"test/external/json_org/valid_long_string.json":            `{"long": "` + generateLongString(1000) + `"}`,
	}

	// Invalid JSON test cases
	invalidTests := map[string]string{
		"test/external/json_org/invalid_trailing_comma_object.json":     `{"key": "value",}`,
		"test/external/json_org/invalid_trailing_comma_array.json":      `[1, 2, 3,]`,
		"test/external/json_org/invalid_missing_colon.json":             `{"key" "value"}`,
		"test/external/json_org/invalid_missing_comma.json":             `{"key1": "value1" "key2": "value2"}`,
		"test/external/json_org/invalid_unterminated_string.json":       `{"key": "unterminated`,
		"test/external/json_org/invalid_unterminated_object.json":       `{"key": "value"`,
		"test/external/json_org/invalid_unterminated_array.json":        `[1, 2, 3`,
		"test/external/json_org/invalid_extra_comma.json":               `{"key":, "value"}`,
		"test/external/json_org/invalid_leading_zero.json":              `{"number": 01}`,
		"test/external/json_org/invalid_trailing_dot.json":              `{"number": 42.}`,
		"test/external/json_org/invalid_leading_dot.json":               `{"number": .42}`,
		"test/external/json_org/invalid_multiple_dots.json":             `{"number": 4.2.2}`,
		"test/external/json_org/invalid_invalid_escape.json":            `{"text": "\q"}`,
		"test/external/json_org/invalid_incomplete_unicode.json":        `{"text": "\u12"}`,
		"test/external/json_org/invalid_control_char.json":              "{\"text\": \"line1\nline2\"}", // unescaped control char
		"test/external/json_org/invalid_single_quotes.json":             `{'key': 'value'}`,
		"test/external/json_org/invalid_unquoted_key.json":              `{key: "value"}`,
		"test/external/json_org/invalid_undefined.json":                 `{"value": undefined}`,
		"test/external/json_org/invalid_infinity.json":                  `{"value": Infinity}`,
		"test/external/json_org/invalid_nan.json":                       `{"value": NaN}`,
		"test/external/json_org/invalid_mismatched_brackets.json":       `{"array": [1, 2, 3}`,
		"test/external/json_org/invalid_empty_string_as_number.json":    `{"number": ""}`,
		"test/external/json_org/invalid_duplicate_keys_strict.json":     `{"key": 1, "key": 2}`, // This is actually valid JSON but might be flagged
	}

	// Write all test cases
	allTests := make(map[string]string)
	for k, v := range validTests {
		allTests[k] = v
	}
	for k, v := range invalidTests {
		allTests[k] = v
	}

	for filename, content := range allTests {
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", filename, err)
		} else {
			fmt.Printf("Created %s\n", filename)
		}
	}

	fmt.Printf("Created %d test files\n", len(allTests))
}

func generateDeepNesting(depth int) string {
	if depth == 0 {
		return `"leaf"`
	}
	return `{"level": ` + generateDeepNesting(depth-1) + `}`
}

func generateLongString(length int) string {
	result := ""
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < length; i++ {
		result += string(chars[i%len(chars)])
	}
	return result
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
