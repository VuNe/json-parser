package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/VuNe/json-parser/internal/lexer"
)

// BenchmarkParser_SimpleObject benchmarks parsing of simple objects
func BenchmarkParser_SimpleObject(b *testing.B) {
	input := `{"name": "John", "age": 30, "active": true}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		_, err := p.Parse()
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// BenchmarkParser_SimpleObjectVsStdLib compares our parser vs Go's standard library
func BenchmarkParser_SimpleObjectVsStdLib(b *testing.B) {
	input := `{"name": "John", "age": 30, "active": true, "salary": 75000.50}`

	b.Run("OurParser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := lexer.New(input)
			p := New(l)
			_, err := p.Parse()
			if err != nil {
				b.Fatalf("Parse failed: %v", err)
			}
		}
	})

	b.Run("StdLib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result interface{}
			err := json.Unmarshal([]byte(input), &result)
			if err != nil {
				b.Fatalf("Unmarshal failed: %v", err)
			}
		}
	})
}

// BenchmarkParser_Arrays benchmarks array parsing
func BenchmarkParser_Arrays(b *testing.B) {
	// Small array
	smallArray := `[1, 2, 3, "four", true, null]`

	// Medium array
	mediumArray := `[` + strings.Repeat(`{"id": 1, "name": "item"}, `, 100)[:len(strings.Repeat(`{"id": 1, "name": "item"}, `, 100))-2] + `]`

	// Large array
	largeArray := `[` + strings.Repeat(`"item", `, 1000)[:len(strings.Repeat(`"item", `, 1000))-2] + `]`

	b.Run("SmallArray", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := lexer.New(smallArray)
			p := New(l)
			_, err := p.Parse()
			if err != nil {
				b.Fatalf("Parse failed: %v", err)
			}
		}
	})

	b.Run("MediumArray", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := lexer.New(mediumArray)
			p := New(l)
			_, err := p.Parse()
			if err != nil {
				b.Fatalf("Parse failed: %v", err)
			}
		}
	})

	b.Run("LargeArray", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := lexer.New(largeArray)
			p := New(l)
			_, err := p.Parse()
			if err != nil {
				b.Fatalf("Parse failed: %v", err)
			}
		}
	})
}

// BenchmarkParser_DeepNesting benchmarks deeply nested structures
func BenchmarkParser_DeepNesting(b *testing.B) {
	depths := []int{10, 50, 100, 200}

	for _, depth := range depths {
		input := generateNestedJSON(depth)

		b.Run(fmt.Sprintf("Depth%d", depth), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.New(input)
				p := New(l)
				_, err := p.Parse()
				if err != nil {
					b.Fatalf("Parse failed at depth %d: %v", depth, err)
				}
			}
		})
	}
}

// BenchmarkParser_StringProcessing benchmarks string parsing with various escape scenarios
func BenchmarkParser_StringProcessing(b *testing.B) {
	tests := map[string]string{
		"SimpleString":  `{"text": "Hello World"}`,
		"EscapedString": `{"text": "Hello \"World\" with \n newlines and \t tabs"}`,
		"UnicodeString": `{"text": "Hello \u0041\u0042\u0043 Unicode \u1234\u5678"}`,
		"LongString":    `{"text": "` + strings.Repeat("A very long string with lots of content. ", 100) + `"}`,
		"ManyStrings":   `{"a": "text", "b": "more text", "c": "even more text", "d": "lots of text", "e": "much text"}`,
	}

	for name, input := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.New(input)
				p := New(l)
				_, err := p.Parse()
				if err != nil {
					b.Fatalf("Parse failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkParser_NumberProcessing benchmarks number parsing
func BenchmarkParser_NumberProcessing(b *testing.B) {
	tests := map[string]string{
		"Integers":     `{"a": 1, "b": 42, "c": -123, "d": 0}`,
		"Floats":       `{"a": 1.5, "b": -3.14159, "c": 0.123456789}`,
		"Scientific":   `{"a": 1e10, "b": 6.022e23, "c": 1.23e-4, "d": -1.5E+10}`,
		"Mixed":        `{"int": 42, "float": 3.14, "sci": 1e6, "neg": -123.45}`,
		"LargeNumbers": `{"big": 1.7976931348623157e+308, "small": 1e-308, "zero": 0}`,
	}

	for name, input := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.New(input)
				p := New(l)
				_, err := p.Parse()
				if err != nil {
					b.Fatalf("Parse failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkParser_ErrorHandling benchmarks error cases (should fail fast)
func BenchmarkParser_ErrorHandling(b *testing.B) {
	tests := map[string]string{
		"TrailingComma":      `{"key": "value",}`,
		"MissingColon":       `{"key" "value"}`,
		"UnterminatedString": `{"key": "value`,
		"InvalidNumber":      `{"key": 123.}`,
		"MismatchedBraces":   `{"key": "value"]`,
	}

	for name, input := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.New(input)
				p := New(l)
				_, err := p.Parse()
				if err == nil {
					b.Fatalf("Expected error but got none for %s", name)
				}
				// Expected to fail - measure how quickly we detect errors
			}
		})
	}
}

// BenchmarkParser_MemoryAllocation benchmarks memory allocation patterns
func BenchmarkParser_MemoryAllocation(b *testing.B) {
	input := `{"users": [{"id": 1, "name": "Alice", "tags": ["admin", "user"]}, {"id": 2, "name": "Bob", "tags": ["user"]}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		_, err := p.Parse()
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// Helper function to generate nested JSON for benchmarking
func generateNestedJSON(depth int) string {
	if depth <= 0 {
		return `"leaf"`
	}
	return `{"level": ` + generateNestedJSON(depth-1) + `}`
}
