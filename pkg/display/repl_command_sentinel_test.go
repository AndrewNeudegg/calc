package display

import "testing"

// Ensure commands return a sentinel value that causes callers to skip printing results.
func TestEvaluateLineCommandSentinel(t *testing.T) {
	r := NewREPL()
	v := r.EvaluateLine(":help")
	if !v.IsError() || v.Error != "" {
		t.Fatalf("commands should return sentinel error value (IsError=true, Error=''); got IsError=%v, Error=%q", v.IsError(), v.Error)
	}
}
