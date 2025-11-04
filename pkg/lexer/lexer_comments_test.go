package lexer

import "testing"

func TestLexerSkipsDoubleSlashComments(t *testing.T) {
	l := New("1 + 2 // comment here\n+ 3")
	toks := l.AllTokens()
	// Expect token sequence: 1, +, 2, +, 3, EOF
	want := []TokenType{TokenNumber, TokenPlus, TokenNumber, TokenPlus, TokenNumber, TokenEOF}
	if len(toks) != len(want) {
		t.Fatalf("unexpected token length: %d vs %d", len(toks), len(want))
	}
	for i, tt := range want {
		if toks[i].Type != tt {
			t.Fatalf("token %d: got %v, want %v (literal %q)", i, toks[i].Type, tt, toks[i].Literal)
		}
	}
}
