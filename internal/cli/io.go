package cli

import (
	"fmt"
	"os"
)

// FileReader provides utilities for reading files.
type FileReader struct{}

// NewFileReader creates a new FileReader instance.
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadFile reads the contents of a file and returns it as a string.
func (fr *FileReader) ReadFile(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", filename, err)
	}

	return string(data), nil
}

// FileExists checks if a file exists and is readable.
func (fr *FileReader) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
