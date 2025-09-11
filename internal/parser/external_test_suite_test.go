package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/VuNe/json-parser/internal/lexer"
)

// TestOfficialJSONTestSuite runs our parser against the official JSON test cases
func TestOfficialJSONTestSuite(t *testing.T) {
	// Get test data directory
	testDataDir := filepath.Join("..", "..", "test", "external", "test", "external", "json_org")

	// Check if test directory exists
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		t.Skip("External test data not found, skipping official test suite")
		return
	}

	// Read all test files
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("Error reading test directory: %v", err)
	}

	var validTests, invalidTests []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			fileName := file.Name()
			if strings.HasPrefix(fileName, "valid_") {
				validTests = append(validTests, fileName)
			} else if strings.HasPrefix(fileName, "invalid_") {
				invalidTests = append(invalidTests, fileName)
			}
		}
	}

	t.Logf("Found %d valid and %d invalid test cases", len(validTests), len(invalidTests))

	// Test valid cases
	t.Run("ValidJSONTests", func(t *testing.T) {
		for _, testFile := range validTests {
			t.Run(testFile, func(t *testing.T) {
				filePath := filepath.Join(testDataDir, testFile)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Error reading test file %s: %v", testFile, err)
				}

				l := lexer.New(string(content))
				p := NewWithInput(l, string(content))

				result, err := p.Parse()
				if err != nil {
					t.Errorf("Valid JSON test %s failed: %v", testFile, err)
					t.Logf("Content: %s", string(content))
				} else {
					// Verify we got some result (except for null which is nil)
					if result == nil && !strings.Contains(string(content), "null") && string(content) != "null" {
						t.Errorf("Valid JSON test %s returned nil result", testFile)
					}
				}
			})
		}
	})

	// Test invalid cases
	t.Run("InvalidJSONTests", func(t *testing.T) {
		for _, testFile := range invalidTests {
			t.Run(testFile, func(t *testing.T) {
				filePath := filepath.Join(testDataDir, testFile)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Error reading test file %s: %v", testFile, err)
				}

				l := lexer.New(string(content))
				p := NewWithInput(l, string(content))

				result, err := p.Parse()
				if err == nil {
					t.Errorf("Invalid JSON test %s should have failed but succeeded", testFile)
					t.Logf("Content: %s", string(content))
					t.Logf("Result: %v", result)
				} else {
					// Log the error for inspection (this is expected)
					t.Logf("Invalid JSON test %s correctly failed with: %v", testFile, err)
				}
			})
		}
	})
}

// TestOfficialJSONTestSuitePerformance runs performance tests on larger JSON files
func TestOfficialJSONTestSuitePerformance(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test", "external", "test", "external", "json_org")

	performanceTests := []string{
		"valid_deep_nesting.json",
		"valid_long_string.json",
		"valid_mixed_nesting.json",
	}

	for _, testFile := range performanceTests {
		t.Run("perf_"+testFile, func(t *testing.T) {
			filePath := filepath.Join(testDataDir, testFile)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Skipf("Performance test file %s not found: %v", testFile, err)
				return
			}

			// Run parsing multiple times to check for consistency
			for i := 0; i < 10; i++ {
				l := lexer.New(string(content))
				p := New(l)

				result, err := p.Parse()
				if err != nil {
					t.Errorf("Performance test %s iteration %d failed: %v", testFile, i, err)
					break
				}

				// Verify result is consistent (basic check)
				if result == nil && !strings.Contains(string(content), "null") {
					t.Errorf("Performance test %s iteration %d returned unexpected nil", testFile, i)
					break
				}
			}
		})
	}
}

// TestJSONSpecificationCompliance tests compliance with JSON specification edge cases
func TestJSONSpecificationCompliance(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		shouldPass  bool
		description string
	}{
		{
			name:        "empty_string",
			json:        "",
			shouldPass:  false,
			description: "Empty input should be rejected",
		},
		{
			name:        "whitespace_only",
			json:        "   \t\n\r   ",
			shouldPass:  false,
			description: "Whitespace-only input should be rejected",
		},
		{
			name:        "multiple_values",
			json:        `{} {}`,
			shouldPass:  false,
			description: "Multiple top-level values should be rejected",
		},
		{
			name:        "leading_zeros",
			json:        `{"num": 007}`,
			shouldPass:  false,
			description: "Leading zeros in numbers should be rejected",
		},
		{
			name:        "hex_numbers",
			json:        `{"num": 0xFF}`,
			shouldPass:  false,
			description: "Hexadecimal numbers should be rejected",
		},
		{
			name:        "comments",
			json:        `{"key": "value" /* comment */}`,
			shouldPass:  false,
			description: "Comments should be rejected in strict JSON",
		},
		{
			name:        "single_quotes",
			json:        `{'key': 'value'}`,
			shouldPass:  false,
			description: "Single quotes should be rejected",
		},
		{
			name:        "unquoted_keys",
			json:        `{key: "value"}`,
			shouldPass:  false,
			description: "Unquoted keys should be rejected",
		},
		{
			name:        "trailing_commas",
			json:        `{"key": "value",}`,
			shouldPass:  false,
			description: "Trailing commas should be rejected",
		},
		{
			name:        "unicode_bom",
			json:        "\uFEFF{}",
			shouldPass:  true, // Some parsers accept BOM, others don't
			description: "Unicode BOM handling",
		},
		{
			name:        "large_numbers",
			json:        `{"big": 1.7976931348623157e+308}`,
			shouldPass:  true,
			description: "Large numbers within float64 range should be accepted",
		},
		{
			name:        "very_deep_nesting",
			json:        generateVeryDeepNesting(1000),
			shouldPass:  true,
			description: "Very deep nesting should be handled gracefully",
		},
		{
			name:        "null_bytes",
			json:        "{\"key\": \"value\x00\"}",
			shouldPass:  false,
			description: "Null bytes in strings should be rejected",
		},
		{
			name:        "control_characters",
			json:        "{\"key\": \"value\t\"}",
			shouldPass:  false,
			description: "Unescaped control characters should be rejected",
		},
		{
			name:        "escaped_control_characters",
			json:        `{"key": "value\t"}`,
			shouldPass:  true,
			description: "Properly escaped control characters should be accepted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.json)
			p := NewWithInput(l, tt.json)

			result, err := p.Parse()

			if tt.shouldPass {
				if err != nil {
					t.Errorf("Test %s (%s) should pass but failed: %v", tt.name, tt.description, err)
				}
			} else {
				if err == nil {
					t.Errorf("Test %s (%s) should fail but passed. Result: %v", tt.name, tt.description, result)
				}
			}
		})
	}
}

// Helper function to generate very deep nesting for stress testing
func generateVeryDeepNesting(depth int) string {
	if depth == 0 {
		return `"end"`
	}
	return `{"level": ` + generateVeryDeepNesting(depth-1) + `}`
}

// TestExternalTestSuiteCoverage ensures we test a comprehensive range of cases
func TestExternalTestSuiteCoverage(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test", "external", "test", "external", "json_org")

	// Count test files by category
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Skipf("Test data directory not accessible: %v", err)
		return
	}

	categories := map[string]int{
		"valid":   0,
		"invalid": 0,
		"total":   0,
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			categories["total"]++
			if strings.HasPrefix(file.Name(), "valid_") {
				categories["valid"]++
			} else if strings.HasPrefix(file.Name(), "invalid_") {
				categories["invalid"]++
			}
		}
	}

	t.Logf("Test suite coverage:")
	t.Logf("  Total test files: %d", categories["total"])
	t.Logf("  Valid cases: %d", categories["valid"])
	t.Logf("  Invalid cases: %d", categories["invalid"])

	// Ensure we have a good balance of test cases
	if categories["total"] < 20 {
		t.Errorf("Expected at least 20 test cases, found %d", categories["total"])
	}

	if categories["valid"] < 10 {
		t.Errorf("Expected at least 10 valid test cases, found %d", categories["valid"])
	}

	if categories["invalid"] < 10 {
		t.Errorf("Expected at least 10 invalid test cases, found %d", categories["invalid"])
	}

	// Test case balance should be reasonable
	ratio := float64(categories["valid"]) / float64(categories["invalid"])
	if ratio < 0.5 || ratio > 2.0 {
		t.Logf("Warning: Test case ratio (valid/invalid = %.2f) might be unbalanced", ratio)
	}
}
