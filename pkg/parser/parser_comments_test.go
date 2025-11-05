package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

// Trailing comments should be ignored by the lexer so the parser sees a normal statement.
func TestParse_TrailingComment_Assignment(t *testing.T) {
	input := "x = 3 // set x"
	l := lexer.New(input)
	tokens := l.AllTokens()
	p := New(tokens)

	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	assign, ok := expr.(*AssignExpr)
	if !ok {
		t.Fatalf("expected AssignExpr, got %T", expr)
	}
	if assign.Name != "x" {
		t.Fatalf("expected variable name 'x', got %q", assign.Name)
	}
}

func TestParse_TrailingComment_Expression(t *testing.T) {
	input := "3 + 4 // sum"
	l := lexer.New(input)
	tokens := l.AllTokens()
	p := New(tokens)

	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if _, ok := expr.(*BinaryExpr); !ok {
		t.Fatalf("expected BinaryExpr, got %T", expr)
	}
}

// Comment-only input results in no tokens except EOF; the parser should report an error.
// Comment handling (treating as no-op) is the REPL's responsibility; this test asserts
// the parser behavior remains strict when given only EOF.
func TestParse_CommentOnly_IsError(t *testing.T) {
	input := "// just a comment"
	l := lexer.New(input)
	tokens := l.AllTokens()
	p := New(tokens)

	if _, err := p.Parse(); err == nil {
		t.Fatal("expected parse error for comment-only input, got nil")
	}
}
