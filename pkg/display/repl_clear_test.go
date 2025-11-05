package display

import (
	"strings"
	"testing"
)

// Test that :clear resets the REPL session: history, variables, and prompt counter.
func TestREPL_ClearCommand_ResetsSession(t *testing.T) {
	// Isolate HOME so defaults are consistent
	t.Setenv("HOME", t.TempDir())
	r := NewREPL()

	// Seed some state
	v1 := r.EvaluateLine("x = 3")
	if v1.IsError() {
		t.Fatalf("unexpected error on seed assign: %s", v1.Error)
	}
	v2 := r.EvaluateLine("x + 1")
	if v2.IsError() || strings.TrimSpace(r.formatter.Format(v2)) != "4.00" {
		t.Fatalf("unexpected result before clear: %s (err=%v)", r.formatter.Format(v2), v2.IsError())
	}

	if got := len(r.ListLines()); got != 2 {
		t.Fatalf("expected 2 lines before clear, got %d", got)
	}
	if r.nextID != 3 {
		t.Fatalf("expected nextID=3 before clear, got %d", r.nextID)
	}

	// Issue :clear
	vClear := r.EvaluateLine(":clear")
	// Command returns sentinel error (no printed result)
	if !vClear.IsError() || vClear.Error != "" {
		t.Fatalf(":clear should return sentinel no-op error, got: %+v", vClear)
	}

	// Verify state reset
	if got := len(r.ListLines()); got != 0 {
		t.Fatalf("expected 0 lines after clear, got %d", got)
	}
	if r.nextID != 1 {
		t.Fatalf("expected nextID reset to 1 after clear, got %d", r.nextID)
	}

	// Variables should be cleared; using x should error now
	v3 := r.EvaluateLine("x + 1")
	if !v3.IsError() {
		t.Fatalf("expected undefined variable error after clear, got non-error: %s", r.formatter.Format(v3))
	}
	if !strings.Contains(v3.Error, "undefined variable") {
		t.Fatalf("expected undefined variable error message, got: %q", v3.Error)
	}
}
