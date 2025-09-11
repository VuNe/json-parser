package cli

import (
	"fmt"
	"os"

	"github.com/VuNe/json-parser/internal/lexer"
	"github.com/VuNe/json-parser/internal/parser"
)

// CLIHandler interface defines the contract for handling CLI operations.
type CLIHandler interface {
	ParseFile(filename string) error
	ParseString(input string) error
	ExitCode() int
}

// handler is the concrete implementation of CLIHandler.
type handler struct {
	fileReader *FileReader
	exitCode   int
}

// New creates a new CLI handler instance.
func New() CLIHandler {
	return &handler{
		fileReader: NewFileReader(),
		exitCode:   0, // Default to success
	}
}

// ParseFile reads a file and parses its JSON content.
func (h *handler) ParseFile(filename string) error {
	// Check if file exists first
	if !h.fileReader.FileExists(filename) {
		h.exitCode = 1
		return fmt.Errorf("file '%s' does not exist or is not readable", filename)
	}

	// Read the file content
	content, err := h.fileReader.ReadFile(filename)
	if err != nil {
		h.exitCode = 1
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse the content
	return h.ParseString(content)
}

// ParseString parses the given JSON string.
func (h *handler) ParseString(input string) error {
	// Create lexer and parser with enhanced error reporting
	lex := lexer.New(input)
	p := parser.NewWithInput(lex, input)

	// Parse the JSON
	_, err := p.Parse()
	if err != nil {
		h.exitCode = 1
		return fmt.Errorf("JSON parsing failed: %w", err)
	}

	// If we reach here, parsing was successful
	h.exitCode = 0
	return nil
}

// ExitCode returns the current exit code.
func (h *handler) ExitCode() int {
	return h.exitCode
}

// Run is a convenience method that handles command line arguments and exits.
func Run() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	handler := New()

	err := handler.ParseFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(handler.ExitCode())
}
