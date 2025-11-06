package lexer

import (
	"testing"
)

func TestLexerNumbersWithCommas(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		tokenType TokenType
	}{
		// US format (comma as thousand separator, period as decimal)
		{"1,234", "1,234", TokenNumber},
		{"1,234.56", "1,234.56", TokenNumber},
		{"55,101.10", "55,101.10", TokenNumber},
		{"31,432", "31,432", TokenNumber},
		{"1,500,000", "1,500,000", TokenNumber},
		{"12,345,678.90", "12,345,678.90", TokenNumber},
		
		// European format (period as thousand separator, comma as decimal)
		{"1.234", "1.234", TokenNumber},
		{"1.234,56", "1.234,56", TokenNumber},
		{"65.342,10", "65.342,10", TokenNumber},
		{"12.345.678,90", "12.345.678,90", TokenNumber},
		
		// Edge cases
		{"100,000", "100,000", TokenNumber},
		{"100.000", "100.000", TokenNumber},
		{"50,000.50", "50,000.50", TokenNumber},
		{"50.000,50", "50.000,50", TokenNumber},
		
		// Simple numbers (no separators)
		{"42", "42", TokenNumber},
		{"3.14", "3.14", TokenNumber},
		{"100.5", "100.5", TokenNumber},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.tokenType {
			t.Errorf("input %q: expected token type %s, got %s", tt.input, tt.tokenType, tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("input %q: expected literal %q, got %q", tt.input, tt.expected, tok.Literal)
		}
	}
}

func TestLexerCurrencyWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []TokenType
		expectedLiterals []string
	}{
		{
			input: "$31,432",
			expectedTokens: []TokenType{TokenCurrency, TokenNumber, TokenEOF},
			expectedLiterals: []string{"$", "31,432", ""},
		},
		{
			input: "$55,101.10",
			expectedTokens: []TokenType{TokenCurrency, TokenNumber, TokenEOF},
			expectedLiterals: []string{"$", "55,101.10", ""},
		},
		{
			input: "€65.342,10",
			expectedTokens: []TokenType{TokenCurrency, TokenNumber, TokenEOF},
			expectedLiterals: []string{"€", "65.342,10", ""},
		},
		{
			input: "£1,234,567.89",
			expectedTokens: []TokenType{TokenCurrency, TokenNumber, TokenEOF},
			expectedLiterals: []string{"£", "1,234,567.89", ""},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()

		if len(tokens) != len(tt.expectedTokens) {
			t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expectedTokens), len(tokens))
			continue
		}

		for i, tok := range tokens {
			if tok.Type != tt.expectedTokens[i] {
				t.Errorf("input %q: token %d expected type %s, got %s", tt.input, i, tt.expectedTokens[i], tok.Type)
			}
			if tok.Literal != tt.expectedLiterals[i] {
				t.Errorf("input %q: token %d expected literal %q, got %q", tt.input, i, tt.expectedLiterals[i], tok.Literal)
			}
		}
	}
}

func TestLexerExpressionWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []TokenType
	}{
		{
			input: "1,234 + 5,678",
			expectedTokens: []TokenType{TokenNumber, TokenPlus, TokenNumber, TokenEOF},
		},
		{
			input: "10,000 * 2.5",
			expectedTokens: []TokenType{TokenNumber, TokenMultiply, TokenNumber, TokenEOF},
		},
		{
			input: "1.234,56 + 2.345,67",
			expectedTokens: []TokenType{TokenNumber, TokenPlus, TokenNumber, TokenEOF},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()

		if len(tokens) != len(tt.expectedTokens) {
			t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expectedTokens), len(tokens))
			continue
		}

		for i, tok := range tokens {
			if tok.Type != tt.expectedTokens[i] {
				t.Errorf("input %q: token %d expected %s, got %s", tt.input, i, tt.expectedTokens[i], tok.Type)
			}
		}
	}
}
