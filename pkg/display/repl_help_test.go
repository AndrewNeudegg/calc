package display

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestHelpOutputRendering verifies that invoking :help prints with CRLF line breaks
// so each line starts at column 0 in raw-mode terminals.
func TestHelpOutputRendering(t *testing.T) {
	r := NewREPL()

	// Capture stdout
	oldStdout := os.Stdout
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = pw

	// Evaluate :help (printed by EvaluateLine)
	_ = r.EvaluateLine(":help")

	// Restore stdout and read captured output
	pw.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, pr); err != nil {
		t.Fatalf("copy: %v", err)
	}
	pr.Close()

	out := buf.String()

	// All line separators should be CRLF to prevent ragged left margins
	if strings.Contains(out, "\n") && !strings.Contains(out, "\r\n") {
		t.Fatalf("expected CRLF line endings, got output with bare LF: %q", out)
	}

	// Ensure specific key lines appear and start at expected indentation (two spaces or none)
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	// Trim any trailing empty line from ensured CRLF
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		t.Fatalf("no output captured")
	}

	// First header line should have no leading spaces
	if strings.HasPrefix(lines[0], " ") {
		t.Fatalf("first line should not start with spaces: %q", lines[0])
	}

	// Find a couple of known lines to check indentation/content
	wantPairs := []string{
		"  :save <file>       Save current workspace",
		"  :open <file>       Open a workspace file",
		"  :help              Show this help",
		"Available settings:",
		"  precision <n>         Number of decimal places (default: 2)",
	}
	for _, want := range wantPairs {
		found := false
		for _, ln := range lines {
			if ln == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected to find line: %q\nGot:\n%q", want, strings.Join(lines, "\n"))
		}
	}
}
