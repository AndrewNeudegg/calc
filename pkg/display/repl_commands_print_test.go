package display

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test that command outputs do not produce a calculation result line like "= 0.00".
func TestCommandPrintsNoResult_SaveOpenHelp(t *testing.T) {
	// Isolate settings path by pointing HOME to a temp dir
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	// Work in a temp CWD for workspace files
	tmpCwd := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpCwd)
	defer os.Chdir(oldWd)

	r := NewREPL()

	// Capture stdout
	oldStdout := os.Stdout
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = pw

	// 1) :help
	_ = r.EvaluateLine(":help")
	// Add a line so save has content
	_ = r.EvaluateLine("t=55")
	// 2) :save session.log (ensure settings write goes under tmp HOME and workspace in CWD)
	_ = r.EvaluateLine(":save session.log")
	// 3) :open test.calc (no-op unless file exists)
	_ = r.EvaluateLine(":open test.calc")

	// Restore stdout and read captured
	pw.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, pr); err != nil {
		t.Fatalf("copy: %v", err)
	}
	pr.Close()

	out := buf.String()

	// Must not contain a result line after commands
	if strings.Contains(out, "\n   = ") {
		t.Fatalf("command outputs should not print result lines; got:\n%s", out)
	}

	// Sanity: messages should be present
	if !strings.Contains(out, "Available commands:") || !strings.Contains(out, "saved to session.log") {
		t.Fatalf("expected help/open/save messages missing; got:\n%s", out)
	}

	// Settings file should exist under HOME/.config/calc/settings.json
	cfg := filepath.Join(tmpHome, ".config", "calc", "settings.json")
	if _, err := os.Stat(cfg); err != nil {
		t.Fatalf("expected settings saved at %s: %v", cfg, err)
	}

	// Workspace file should be created in CWD
	if _, err := os.Stat("session.log"); err != nil {
		t.Fatalf("expected session.log in CWD: %v", err)
	}
}
