package display

import (
	"strings"
	"testing"
)

// Test that comment-only lines are treated as no-op and do not print results
func TestEvaluateLine_CommentOnly(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	r := NewREPL()

	v := r.EvaluateLine("// this is a comment-only line")
	if !v.IsError() || v.Error != "" {
		t.Fatalf("expected sentinel no-op error, got: %+v", v)
	}

	// Ensure that subsequent valid input still works
	v2 := r.EvaluateLine("3 + 4")
	if v2.IsError() {
		t.Fatalf("unexpected error on valid input: %s", v2.Error)
	}
	if strings.TrimSpace(r.formatter.Format(v2)) != "7.00" {
		t.Fatalf("unexpected result: %s", r.formatter.Format(v2))
	}
}

// Trailing comments after an expression should be ignored and not affect evaluation
func TestEvaluateLine_TrailingComment(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	r := NewREPL()
	v := r.EvaluateLine("x = 3 // set x")
	if v.IsError() {
		t.Fatalf("unexpected error: %s", v.Error)
	}
	if strings.TrimSpace(r.formatter.Format(v)) != "3.00" {
		t.Fatalf("unexpected result: %s", r.formatter.Format(v))
	}

	v2 := r.EvaluateLine("x + 4 // use x")
	if v2.IsError() {
		t.Fatalf("unexpected error: %s", v2.Error)
	}
	if strings.TrimSpace(r.formatter.Format(v2)) != "7.00" {
		t.Fatalf("unexpected result: %s", r.formatter.Format(v2))
	}
}

// Multiple consecutive comment-only lines should be treated as no-ops
func TestEvaluateLine_MultipleCommentOnlyLines(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	r := NewREPL()
	v1 := r.EvaluateLine("// header")
	if !v1.IsError() || v1.Error != "" {
		t.Fatalf("expected sentinel no-op for first comment, got: %+v", v1)
	}
	v2 := r.EvaluateLine("   // another comment with leading spaces")
	if !v2.IsError() || v2.Error != "" {
		t.Fatalf("expected sentinel no-op for second comment, got: %+v", v2)
	}
	v3 := r.EvaluateLine("1 + 2 // after comments")
	if v3.IsError() {
		t.Fatalf("unexpected error after comments: %s", v3.Error)
	}
	if strings.TrimSpace(r.formatter.Format(v3)) != "3.00" {
		t.Fatalf("unexpected result: %s", r.formatter.Format(v3))
	}
}
