package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParser_Prev(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedType   interface{}
		expectedOffset int
		expectedAbs    bool
	}{
		{
			name:           "prev keyword",
			input:          "prev",
			expectedType:   &PrevExpr{},
			expectedOffset: 0,
			expectedAbs:    false,
		},
		{
			name:           "prev~",
			input:          "prev~",
			expectedType:   &PrevExpr{},
			expectedOffset: 1,
			expectedAbs:    false,
		},
		{
			name:           "prev~1",
			input:          "prev~1",
			expectedType:   &PrevExpr{},
			expectedOffset: 1,
			expectedAbs:    false,
		},
		{
			name:           "prev~5",
			input:          "prev~5",
			expectedType:   &PrevExpr{},
			expectedOffset: 5,
			expectedAbs:    false,
		},
		{
			name:           "prev~10",
			input:          "prev~10",
			expectedType:   &PrevExpr{},
			expectedOffset: 10,
			expectedAbs:    false,
		},
		{
			name:           "prev#1 (absolute)",
			input:          "prev#1",
			expectedType:   &PrevExpr{},
			expectedOffset: 1,
			expectedAbs:    true,
		},
		{
			name:           "prev#15 (absolute)",
			input:          "prev#15",
			expectedType:   &PrevExpr{},
			expectedOffset: 15,
			expectedAbs:    true,
		},
		{
			name:           "prev#100 (absolute)",
			input:          "prev#100",
			expectedType:   &PrevExpr{},
			expectedOffset: 100,
			expectedAbs:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			// Remove EOF token
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			p := New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			prevExpr, ok := expr.(*PrevExpr)
			if !ok {
				t.Fatalf("Expected *PrevExpr, got %T", expr)
			}

			if prevExpr.Offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, prevExpr.Offset)
			}

			if prevExpr.Absolute != tt.expectedAbs {
				t.Errorf("Expected absolute %v, got %v", tt.expectedAbs, prevExpr.Absolute)
			}
		})
	}
}

func TestParser_PrevInExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "addition with prev",
			input: "10 + prev",
		},
		{
			name:  "multiplication with prev~1",
			input: "prev~1 * 2",
		},
		{
			name:  "prev + prev~1",
			input: "prev + prev~1",
		},
		{
			name:  "assignment with prev",
			input: "x = prev",
		},
		{
			name:  "complex expression",
			input: "(prev + 10) * prev~1",
		},
		{
			name:  "addition with prev#15",
			input: "10 + prev#15",
		},
		{
			name:  "multiplication with prev#10",
			input: "prev#10 * 42",
		},
		{
			name:  "mixing relative and absolute",
			input: "prev#15 + prev~2",
		},
		{
			name:  "assignment with absolute prev",
			input: "x = prev#5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			// Remove EOF token
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			p := New(tokens)
			_, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
		})
	}
}

func TestParser_PrevValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "prev~0 is valid",
			input:       "prev~0",
			expectError: false,
		},
		{
			name:        "prev~100 is valid",
			input:       "prev~100",
			expectError: false,
		},
		{
			name:        "prev#1 is valid",
			input:       "prev#1",
			expectError: false,
		},
		{
			name:        "prev#100 is valid",
			input:       "prev#100",
			expectError: false,
		},
		{
			name:        "prev# without number is invalid",
			input:       "prev#",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			// Remove EOF token
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			p := New(tokens)
			_, err := p.Parse()
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
