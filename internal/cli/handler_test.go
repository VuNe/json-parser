package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	handler := New()
	if handler == nil {
		t.Fatal("New() returned nil")
	}

	// Initial exit code should be 0
	if handler.ExitCode() != 0 {
		t.Errorf("initial exit code should be 0, got %d", handler.ExitCode())
	}
}

func TestHandler_ParseString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectError  bool
		expectedExit int
	}{
		{
			name:         "valid empty object",
			input:        "{}",
			expectError:  false,
			expectedExit: 0,
		},
		{
			name:         "valid empty object with whitespace",
			input:        " { } ",
			expectError:  false,
			expectedExit: 0,
		},
		{
			name:         "empty input",
			input:        "",
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "invalid JSON",
			input:        "{",
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "non-empty object",
			input:        `{"key": "value"}`,
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "invalid character",
			input:        "invalid",
			expectError:  true,
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New()

			err := handler.ParseString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if handler.ExitCode() != tt.expectedExit {
				t.Errorf("expected exit code %d, got %d", tt.expectedExit, handler.ExitCode())
			}
		})
	}
}

func TestHandler_ParseFile(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	// Create test files
	validFile := filepath.Join(tempDir, "valid.json")
	if err := os.WriteFile(validFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	invalidFile := filepath.Join(tempDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte("{"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	emptyFile := filepath.Join(tempDir, "empty.json")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	nonExistentFile := filepath.Join(tempDir, "nonexistent.json")

	tests := []struct {
		name         string
		filename     string
		expectError  bool
		expectedExit int
	}{
		{
			name:         "valid JSON file",
			filename:     validFile,
			expectError:  false,
			expectedExit: 0,
		},
		{
			name:         "invalid JSON file",
			filename:     invalidFile,
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "empty JSON file",
			filename:     emptyFile,
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "non-existent file",
			filename:     nonExistentFile,
			expectError:  true,
			expectedExit: 1,
		},
		{
			name:         "empty filename",
			filename:     "",
			expectError:  true,
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New()

			err := handler.ParseFile(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if handler.ExitCode() != tt.expectedExit {
				t.Errorf("expected exit code %d, got %d", tt.expectedExit, handler.ExitCode())
			}
		})
	}
}

func TestHandler_ExitCode(t *testing.T) {
	tests := []struct {
		name         string
		parseInput   string
		expectError  bool
		expectedExit int
	}{
		{
			name:         "success case",
			parseInput:   "{}",
			expectError:  false,
			expectedExit: 0,
		},
		{
			name:         "failure case",
			parseInput:   "invalid",
			expectError:  true,
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New()

			// Initially should be 0
			if handler.ExitCode() != 0 {
				t.Errorf("initial exit code should be 0, got %d", handler.ExitCode())
			}

			// Parse something to potentially change exit code
			_ = handler.ParseString(tt.parseInput)

			if handler.ExitCode() != tt.expectedExit {
				t.Errorf("expected exit code %d, got %d", tt.expectedExit, handler.ExitCode())
			}
		})
	}
}
