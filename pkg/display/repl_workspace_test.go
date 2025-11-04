package display

import (
	"os"
	"strings"
	"testing"
)

func TestWorkspaceSaveAndOpenRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Use a temp working dir so files don't leak
	wd := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(wd)
	defer os.Chdir(oldWd)

	r := NewREPL()
	_ = r.EvaluateLine("x=2")
	_ = r.EvaluateLine("y=3")

	// Save to file
	_ = r.EvaluateLine(":save work.calc")

	// Verify contents
	b, err := os.ReadFile("work.calc")
	if err != nil {
		t.Fatalf("reading work.calc: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, "x=2") || !strings.Contains(s, "y=3") {
		t.Fatalf("workspace file missing lines: %q", s)
	}

	// Open in a fresh REPL
	r2 := NewREPL()
	_ = r2.EvaluateLine(":open work.calc")

	lines := r2.ListLines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after open, got %d", len(lines))
	}
	if strings.TrimSpace(lines[0].Input) != "x=2" || strings.TrimSpace(lines[1].Input) != "y=3" {
		t.Fatalf("unexpected lines after open: %#v", []string{lines[0].Input, lines[1].Input})
	}

	// Sanity: environment restored so expressions work
	v := r2.EvaluateLine("x + y")
	if v.IsError() || int(v.Number) != 5 {
		t.Fatalf("expected x+y == 5, got %+v", v)
	}
}
