package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIIntegration tests the complete CLI interface end-to-end
func TestCLIIntegration(t *testing.T) {
	// Build the binary first
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	t.Run("ValidJSONFiles", func(t *testing.T) {
		validFiles := map[string]string{
			"simple.json":  `{"name": "test", "value": 42}`,
			"array.json":   `[1, 2, 3, "four", true, null]`,
			"nested.json":  `{"user": {"id": 1, "profile": {"name": "John"}}}`,
			"empty.json":   `{}`,
			"complex.json": `{"users": [{"id": 1, "tags": ["admin"]}, {"id": 2, "tags": []}]}`,
		}

		for filename, content := range validFiles {
			t.Run(filename, func(t *testing.T) {
				tempFile := createTempFile(t, filename, content)
				defer os.Remove(tempFile)

				cmd := exec.Command(binaryPath, tempFile)
				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					t.Errorf("Command failed for valid JSON %s: %v", filename, err)
					t.Logf("Stderr: %s", stderr.String())
				}

				// Check exit code (should be 0 for valid JSON)
				if cmd.ProcessState.ExitCode() != 0 {
					t.Errorf("Expected exit code 0, got %d for %s", cmd.ProcessState.ExitCode(), filename)
				}

				// Should have no output for valid JSON (unless we add verbose mode)
				if stderr.Len() > 0 {
					t.Logf("Stderr output for valid JSON %s: %s", filename, stderr.String())
				}
			})
		}
	})

	t.Run("InvalidJSONFiles", func(t *testing.T) {
		invalidFiles := map[string]string{
			"trailing_comma.json":    `{"key": "value",}`,
			"missing_colon.json":     `{"key" "value"}`,
			"unterminated.json":      `{"key": "value"`,
			"invalid_number.json":    `{"num": 123.}`,
			"mismatched_brackets.json": `{"array": [1, 2, 3}`,
		}

		for filename, content := range invalidFiles {
			t.Run(filename, func(t *testing.T) {
				tempFile := createTempFile(t, filename, content)
				defer os.Remove(tempFile)

				cmd := exec.Command(binaryPath, tempFile)
				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				
				// Should exit with error for invalid JSON
				if err == nil {
					t.Errorf("Command should have failed for invalid JSON %s", filename)
				}

				// Check exit code (should be 1 for invalid JSON)
				if cmd.ProcessState.ExitCode() != 1 {
					t.Errorf("Expected exit code 1, got %d for %s", cmd.ProcessState.ExitCode(), filename)
				}

				// Should have error message
				if stderr.Len() == 0 {
					t.Errorf("Expected error message for invalid JSON %s", filename)
				}

				errorMsg := stderr.String()
				// Error messages should contain position information
				if !strings.Contains(errorMsg, "line") || !strings.Contains(errorMsg, "column") {
					t.Errorf("Error message should contain position info for %s. Got: %s", filename, errorMsg)
				}
			})
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "non_existent_file.json")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Error("Command should have failed for non-existent file")
		}

		if cmd.ProcessState.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", cmd.ProcessState.ExitCode())
		}

		if !strings.Contains(stderr.String(), "does not exist") {
			t.Errorf("Expected 'does not exist' message, got: %s", stderr.String())
		}
	})

	t.Run("NoArguments", func(t *testing.T) {
		cmd := exec.Command(binaryPath)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Error("Command should have failed with no arguments")
		}

		if cmd.ProcessState.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", cmd.ProcessState.ExitCode())
		}

		if !strings.Contains(stderr.String(), "Usage") {
			t.Errorf("Expected usage message, got: %s", stderr.String())
		}
	})
}

// TestCLIWithTestDataFiles tests CLI with the actual test data files
func TestCLIWithTestDataFiles(t *testing.T) {
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	// Test with step test files
	testDataDirs := []string{
		filepath.Join("..", "test", "testdata"),
	}

	for _, testDir := range testDataDirs {
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			continue // Skip if directory doesn't exist
		}

		files, err := os.ReadDir(testDir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			t.Run(file.Name(), func(t *testing.T) {
				filePath := filepath.Join(testDir, file.Name())
				
				cmd := exec.Command(binaryPath, filePath)
				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				
				// Determine expected result based on filename
				shouldPass := strings.Contains(file.Name(), "valid_")
				
				if shouldPass {
					if err != nil {
						t.Errorf("Valid file %s should have passed: %v", file.Name(), err)
						t.Logf("Stderr: %s", stderr.String())
					}
					if cmd.ProcessState.ExitCode() != 0 {
						t.Errorf("Valid file %s should have exit code 0, got %d", file.Name(), cmd.ProcessState.ExitCode())
					}
				} else if strings.Contains(file.Name(), "invalid_") {
					if err == nil {
						t.Errorf("Invalid file %s should have failed", file.Name())
					}
					if cmd.ProcessState.ExitCode() != 1 {
						t.Errorf("Invalid file %s should have exit code 1, got %d", file.Name(), cmd.ProcessState.ExitCode())
					}
				}
			})
		}
	}
}

// TestConcurrentParsing tests that multiple parser instances can run concurrently
func TestConcurrentParsing(t *testing.T) {
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	// Create multiple test files
	testFiles := make([]string, 10)
	for i := 0; i < 10; i++ {
		content := fmt.Sprintf(`{"test": %d, "data": [1, 2, %d]}`, i, i*10)
		testFiles[i] = createTempFile(t, fmt.Sprintf("concurrent_%d.json", i), content)
		defer os.Remove(testFiles[i])
	}

	// Run multiple commands concurrently
	done := make(chan bool, len(testFiles))
	errors := make(chan error, len(testFiles))

	for _, file := range testFiles {
		go func(filename string) {
			cmd := exec.Command(binaryPath, filename)
			err := cmd.Run()
			if err != nil {
				errors <- fmt.Errorf("file %s failed: %v", filename, err)
			} else if cmd.ProcessState.ExitCode() != 0 {
				errors <- fmt.Errorf("file %s had exit code %d", filename, cmd.ProcessState.ExitCode())
			}
			done <- true
		}(file)
	}

	// Wait for all to complete
	for i := 0; i < len(testFiles); i++ {
		<-done
	}

	// Check for errors
	close(errors)
	for err := range errors {
		t.Error(err)
	}
}

// TestStressTest runs stress testing with larger JSON files
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	// Create large JSON file
	var largeJSON strings.Builder
	largeJSON.WriteString(`{"items": [`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			largeJSON.WriteString(`, `)
		}
		largeJSON.WriteString(fmt.Sprintf(`{"id": %d, "name": "item_%d", "data": [%d, %d, %d]}`, 
			i, i, i*1, i*2, i*3))
	}
	largeJSON.WriteString(`]}`)

	tempFile := createTempFile(t, "large.json", largeJSON.String())
	defer os.Remove(tempFile)

	cmd := exec.Command(binaryPath, tempFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Errorf("Stress test failed: %v", err)
		t.Logf("Stderr: %s", stderr.String())
	}

	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("Stress test should have exit code 0, got %d", cmd.ProcessState.ExitCode())
	}
}

// Helper functions

func buildBinary(t *testing.T) string {
	binaryPath := filepath.Join(t.TempDir(), "json-parser-test")
	
	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/json-parser")
	cmd.Dir = ".."
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v\nStderr: %s", err, stderr.String())
	}
	
	return binaryPath
}

func createTempFile(t *testing.T, name, content string) string {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, name)
	
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file %s: %v", name, err)
	}
	
	return filePath
}
