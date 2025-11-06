package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParser_Prev(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType interface{}
		expectedVal  interface{}
	}{
		{
			name:         "prev keyword",
			input:        "prev",
			expectedType: &PrevExpr{},
			expectedVal:  0,
		},
		{
			name:         "prev~",
			input:        "prev~",
			expectedType: &PrevExpr{},
			expectedVal:  1,
		},
		{
			name:         "prev~1",
			input:        "prev~1",
			expectedType: &PrevExpr{},
			expectedVal:  1,
		},
		{
			name:         "prev~5",
			input:        "prev~5",
			expectedType: &PrevExpr{},
			expectedVal:  5,
		},
		{
			name:         "prev~10",
			input:        "prev~10",
			expectedType: &PrevExpr{},
			expectedVal:  10,
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

			if prevExpr.Offset != tt.expectedVal {
				t.Errorf("Expected offset %d, got %d", tt.expectedVal, prevExpr.Offset)
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
