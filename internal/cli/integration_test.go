package cli

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestStep2Integration(t *testing.T) {
	validTests := []struct {
		filename string
		name     string
	}{
		{"step2_valid_simple_pair.json", "simple key-value pair"},
		{"step2_valid_multiple_pairs.json", "multiple key-value pairs"},
		{"step2_valid_empty_strings.json", "empty string values"},
		{"step2_valid_with_escapes.json", "strings with escape sequences"},
		{"step2_valid_unicode.json", "strings with unicode escapes"},
		{"step2_valid_with_whitespace.json", "object with whitespace"},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the test file
			filePath := filepath.Join("../../test/testdata", tt.filename)
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tt.filename, err)
			}

			// Test with CLI handler
			handler := New()
			err = handler.ParseString(string(content))

			if err != nil {
				t.Errorf("Expected valid JSON in %s, got error: %v", tt.filename, err)
			}

			if handler.ExitCode() != 0 {
				t.Errorf("Expected exit code 0 for valid JSON in %s, got %d", tt.filename, handler.ExitCode())
			}
		})
	}
}

func TestStep2IntegrationInvalid(t *testing.T) {
	invalidTests := []struct {
		filename string
		name     string
	}{
		{"step2_invalid_missing_colon.json", "missing colon"},
		{"step2_invalid_missing_comma.json", "missing comma"},
		{"step2_invalid_trailing_comma.json", "trailing comma"},
		{"step2_invalid_unterminated_string.json", "unterminated string"},
		{"step2_invalid_bad_escape.json", "invalid escape sequence"},
		{"step2_invalid_incomplete_unicode.json", "incomplete unicode escape"},
		{"step2_invalid_non_string_key.json", "non-string key"},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the test file
			filePath := filepath.Join("../../test/testdata", tt.filename)
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tt.filename, err)
			}

			// Test with CLI handler
			handler := New()
			err = handler.ParseString(string(content))

			if err == nil {
				t.Errorf("Expected error for invalid JSON in %s, but got none", tt.filename)
			}

			if handler.ExitCode() != 1 {
				t.Errorf("Expected exit code 1 for invalid JSON in %s, got %d", tt.filename, handler.ExitCode())
			}
		})
	}
}

func TestStep1BackwardCompatibility(t *testing.T) {
	step1Tests := []struct {
		filename string
		name     string
		valid    bool
	}{
		{"step1_valid_empty.json", "empty object", true},
		{"step1_valid_empty_with_whitespace.json", "empty object with whitespace", true},
		{"step1_valid_empty_multiline.json", "empty object multiline", true},
		{"step1_invalid_empty.json", "invalid empty file", false},
		{"step1_invalid_extra_content.json", "extra content", false},
		{"step1_invalid_missing_close.json", "missing close brace", false},
		{"step1_invalid_missing_open.json", "missing open brace", false},
		{"step1_invalid_non_empty.json", "non-empty object (now valid in Step 2)", true}, // This should now be valid!
	}

	for _, tt := range step1Tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the test file
			filePath := filepath.Join("../../test/testdata", tt.filename)
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tt.filename, err)
			}

			// Test with CLI handler
			handler := New()
			err = handler.ParseString(string(content))

			if tt.valid {
				if err != nil {
					t.Errorf("Expected valid JSON in %s, got error: %v", tt.filename, err)
				}
				if handler.ExitCode() != 0 {
					t.Errorf("Expected exit code 0 for valid JSON in %s, got %d", tt.filename, handler.ExitCode())
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid JSON in %s, but got none", tt.filename)
				}
				if handler.ExitCode() != 1 {
					t.Errorf("Expected exit code 1 for invalid JSON in %s, got %d", tt.filename, handler.ExitCode())
				}
			}
		})
	}
}
