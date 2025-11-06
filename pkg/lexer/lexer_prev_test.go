package lexer

import (
	"testing"
)

func TestLexer_Prev(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "prev keyword",
			input: "prev",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "prev with tilde",
			input: "prev~",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev~"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "prev~1",
			input: "prev~1",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev~1"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "prev~5",
			input: "prev~5",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev~5"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "prev~10",
			input: "prev~10",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev~10"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "expression with prev",
			input: "10 + prev",
			expected: []Token{
				{Type: TokenNumber, Literal: "10"},
				{Type: TokenPlus, Literal: "+"},
				{Type: TokenPrev, Literal: "prev"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "expression with prev~1",
			input: "prev~1 * 2",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev~1"},
				{Type: TokenMultiply, Literal: "*"},
				{Type: TokenNumber, Literal: "2"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			name:  "multiple prev references",
			input: "prev + prev~1",
			expected: []Token{
				{Type: TokenPrev, Literal: "prev"},
				{Type: TokenPlus, Literal: "+"},
				{Type: TokenPrev, Literal: "prev~1"},
				{Type: TokenEOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tokens := l.AllTokens()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, expected := range tt.expected {
				if tokens[i].Type != expected.Type {
					t.Errorf("Token %d: expected type %v, got %v", i, expected.Type, tokens[i].Type)
				}
				if tokens[i].Literal != expected.Literal {
					t.Errorf("Token %d: expected literal %q, got %q", i, expected.Literal, tokens[i].Literal)
				}
			}
		})
	}
}
