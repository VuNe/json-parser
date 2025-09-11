package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	fr := NewFileReader()
	if fr == nil {
		t.Fatal("NewFileReader() returned nil")
	}
}

func TestFileReader_ReadFile(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	// Create a test file with known content
	testContent := `{}`
	testFile := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create an empty file
	emptyFile := filepath.Join(tempDir, "empty.json")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create empty test file: %v", err)
	}

	nonExistentFile := filepath.Join(tempDir, "nonexistent.json")

	tests := []struct {
		name            string
		filename        string
		expectError     bool
		expectedContent string
	}{
		{
			name:            "valid file",
			filename:        testFile,
			expectError:     false,
			expectedContent: testContent,
		},
		{
			name:            "empty file",
			filename:        emptyFile,
			expectError:     false,
			expectedContent: "",
		},
		{
			name:        "non-existent file",
			filename:    nonExistentFile,
			expectError: true,
		},
		{
			name:        "empty filename",
			filename:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := NewFileReader()

			content, err := fr.ReadFile(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if content != tt.expectedContent {
					t.Errorf("expected content %q, got %q", tt.expectedContent, content)
				}
			}
		})
	}
}

func TestFileReader_FileExists(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "exists.json")
	if err := os.WriteFile(testFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	nonExistentFile := filepath.Join(tempDir, "nonexistent.json")

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "existing file",
			filename: testFile,
			expected: true,
		},
		{
			name:     "non-existent file",
			filename: nonExistentFile,
			expected: false,
		},
		{
			name:     "empty filename",
			filename: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := NewFileReader()

			exists := fr.FileExists(tt.filename)

			if exists != tt.expected {
				t.Errorf("expected FileExists(%q) = %v, got %v", tt.filename, tt.expected, exists)
			}
		})
	}
}
