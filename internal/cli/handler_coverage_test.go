package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHandler_EdgeCases tests edge cases to improve coverage
func TestHandler_EdgeCases(t *testing.T) {
	h := New()

	// Test with empty filename
	err := h.ParseFile("")
	if err == nil {
		t.Error("expected error for empty filename")
	}
	if h.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", h.ExitCode())
	}

	// Test with directory instead of file
	tempDir := t.TempDir()
	err = h.ParseFile(tempDir)
	if err == nil {
		t.Error("expected error when parsing directory")
	}
	if h.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", h.ExitCode())
	}
}

// TestHandler_LargeFile tests parsing of larger JSON files
func TestHandler_LargeFile(t *testing.T) {
	h := New()

	// Create a temporary large JSON file
	tempFile := filepath.Join(t.TempDir(), "large.json")
	f, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	// Write a large JSON object
	f.WriteString("{\n")
	for i := 0; i < 1000; i++ {
		f.WriteString(`  "key` + string(rune('0'+i%10)) + `": "value` + string(rune('0'+i%10)) + `"`)
		if i < 999 {
			f.WriteString(",\n")
		} else {
			f.WriteString("\n")
		}
	}
	f.WriteString("}")

	err = h.ParseFile(tempFile)
	if err != nil {
		t.Errorf("unexpected error parsing large file: %v", err)
	}
	if h.ExitCode() != 0 {
		t.Errorf("expected exit code 0, got %d", h.ExitCode())
	}
}

// TestHandler_MalformedFiles tests various malformed JSON files
func TestHandler_MalformedFiles(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"incomplete_string", `{"key": "value`},
		{"invalid_escape", `{"key": "value\q"}`},
		{"trailing_comma", `{"key": "value",}`},
		{"missing_colon", `{"key" "value"}`},
		{"invalid_number", `{"key": 123.}`},
		{"invalid_keyword", `{"key": True}`},
		{"nested_incomplete", `{"outer": {"inner": }`},
		{"array_trailing_comma", `{"arr": [1, 2, 3,]}`},
		{"mixed_brackets", `{"arr": [1, 2, 3}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			
			// Create temp file with malformed content
			tempFile := filepath.Join(t.TempDir(), "malformed.json")
			err := os.WriteFile(tempFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			err = h.ParseFile(tempFile)
			if err == nil {
				t.Error("expected error for malformed JSON")
			}
			if h.ExitCode() != 1 {
				t.Errorf("expected exit code 1, got %d", h.ExitCode())
			}
		})
	}
}

// TestHandler_UnicodeFile tests parsing files with Unicode content
func TestHandler_UnicodeFile(t *testing.T) {
	h := New()

	// Create temp file with Unicode content
	tempFile := filepath.Join(t.TempDir(), "unicode.json")
	unicodeContent := `{
		"english": "Hello World",
		"japanese": "こんにちは世界",
		"emoji": "🌍🚀✨",
		"unicode_escape": "\u0041\u0042\u0043",
		"mixed": "Hello 世界 🌍"
	}`
	
	err := os.WriteFile(tempFile, []byte(unicodeContent), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	err = h.ParseFile(tempFile)
	if err != nil {
		t.Errorf("unexpected error parsing Unicode file: %v", err)
	}
	if h.ExitCode() != 0 {
		t.Errorf("expected exit code 0, got %d", h.ExitCode())
	}
}

// TestHandler_EmptyAndWhitespaceFiles tests edge cases with empty/whitespace files
func TestHandler_EmptyAndWhitespaceFiles(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{"completely_empty", "", true},
		{"only_whitespace", "   \n\t\r  ", true},
		{"only_newlines", "\n\n\n", true},
		{"valid_with_whitespace", "  {}  ", false},
		{"valid_multiline", "{\n  \n}", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			
			tempFile := filepath.Join(t.TempDir(), "test.json")
			err := os.WriteFile(tempFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			err = h.ParseFile(tempFile)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if h.ExitCode() != 1 {
					t.Errorf("expected exit code 1, got %d", h.ExitCode())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if h.ExitCode() != 0 {
					t.Errorf("expected exit code 0, got %d", h.ExitCode())
				}
			}
		})
	}
}
