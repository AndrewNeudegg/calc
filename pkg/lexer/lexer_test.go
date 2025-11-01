package lexer

import (
	"testing"
)

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"42", []TokenType{TokenNumber, TokenEOF}},
		{"3.14", []TokenType{TokenNumber, TokenEOF}},
		{"100.5", []TokenType{TokenNumber, TokenEOF}},
	}
	
	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		
		if len(tokens) != len(tt.expected) {
			t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expected), len(tokens))
			continue
		}
		
		for i, tok := range tokens {
			if tok.Type != tt.expected[i] {
				t.Errorf("input %q: token %d expected %s, got %s", tt.input, i, tt.expected[i], tok.Type)
			}
		}
	}
}

func TestLexerOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"3 + 4", []TokenType{TokenNumber, TokenPlus, TokenNumber, TokenEOF}},
		{"10 - 5", []TokenType{TokenNumber, TokenMinus, TokenNumber, TokenEOF}},
		{"6 * 7", []TokenType{TokenNumber, TokenMultiply, TokenNumber, TokenEOF}},
		{"8 / 2", []TokenType{TokenNumber, TokenDivide, TokenNumber, TokenEOF}},
		{"20%", []TokenType{TokenNumber, TokenPercent, TokenEOF}},
	}
	
	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		
		if len(tokens) != len(tt.expected) {
			t.Errorf("input %q: expected %d tokens, got %d", tt.input, len(tt.expected), len(tokens))
			continue
		}
		
		for i, tok := range tokens {
			if tok.Type != tt.expected[i] {
				t.Errorf("input %q: token %d expected %s, got %s", tt.input, i, tt.expected[i], tok.Type)
			}
		}
	}
}

func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"in", TokenIn},
		{"of", TokenOf},
		{"per", TokenPer},
		{"today", TokenToday},
		{"tomorrow", TokenTomorrow},
		{"yesterday", TokenYesterday},
		{"half", TokenHalf},
		{"double", TokenDouble},
	}
	
	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()
		
		if tok.Type != tt.expected {
			t.Errorf("input %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

func TestLexerUnits(t *testing.T) {
	tests := []struct {
		input string
		hasUnit bool
	}{
		{"10 m", true},
		{"5 km", true},
		{"100 kg", true},
		{"2 hours", true},
	}
	
	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		
		hasUnit := false
		for _, tok := range tokens {
			if tok.Type == TokenUnit {
				hasUnit = true
				break
			}
		}
		
		if hasUnit != tt.hasUnit {
			t.Errorf("input %q: expected hasUnit=%v, got %v", tt.input, tt.hasUnit, hasUnit)
		}
	}
}

func TestLexerCurrency(t *testing.T) {
	tests := []struct {
		input       string
		hasCurrency bool
	}{
		{"$100", true},
		{"£50", true},
		{"10 dollars", false},
	}
	
	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		
		hasCurrency := false
		for _, tok := range tokens {
			if tok.Type == TokenCurrency {
				hasCurrency = true
				break
			}
		}
		
		if hasCurrency != tt.hasCurrency {
			t.Errorf("input %q: expected hasCurrency=%v, got %v", tt.input, tt.hasCurrency, hasCurrency)
		}
	}
}
