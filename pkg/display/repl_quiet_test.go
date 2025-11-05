package display

import (
	"testing"
)

// Test that :quiet suppresses assignment result printing by returning a sentinel no-op value.
func TestQuietSuppressesAssignmentOutput(t *testing.T) {
	r := NewREPL()

	// Baseline: assignment should return a non-error value when quiet is off
	v1 := r.EvaluateLine("x = 42")
	if v1.IsError() && v1.Error == "" {
		t.Fatalf("unexpected sentinel for assignment with quiet off")
	}

	// Turn quiet on via command
	vq := r.EvaluateLine(":quiet on")
	if !vq.IsError() || vq.Error != "" {
		// command returns sentinel
	}
	if !r.IsQuiet() {
		t.Fatalf("quiet mode should be enabled after :quiet on")
	}

	// Now assignment should be suppressed (sentinel)
	v2 := r.EvaluateLine("y = 7")
	if !v2.IsError() || v2.Error != "" {
		t.Fatalf("expected sentinel no-op error for assignment in quiet mode, got: %+v", v2)
	}

	// Non-assignment (e.g., a print) should still produce a value
	v3 := r.EvaluateLine("print(\"val {x}\")")
	if v3.IsError() && v3.Error == "" {
		t.Fatalf("print should not be suppressed in quiet mode")
	}

	// Turn quiet off
	_ = r.EvaluateLine(":quiet off")
	if r.IsQuiet() {
		t.Fatalf("quiet mode should be disabled after :quiet off")
	}

	// Assignment should again return a value (not sentinel)
	v4 := r.EvaluateLine("z = 1")
	if v4.IsError() && v4.Error == "" {
		t.Fatalf("assignment should not be suppressed when quiet is off")
	}
}
