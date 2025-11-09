package display

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

// TestEditorCtrlC tests that Ctrl-C cancels the current line
func TestEditorCtrlC(t *testing.T) {
	ed := NewEditor("> ", nil)
	
	// Simulate typing some text
	ed.buf = []rune("test input")
	ed.cur = len(ed.buf)
	
	// Create a reader with Ctrl-C (0x03)
	input := bytes.NewBufferString("\x03")
	reader := bufio.NewReader(input)
	output := &bytes.Buffer{}
	
	// Read a line - should abort
	line, aborted, eof := ed.ReadLine(reader, output)
	
	if !aborted {
		t.Error("Expected Ctrl-C to abort the line")
	}
	
	if eof {
		t.Error("Ctrl-C should not signal EOF")
	}
	
	if line != "" {
		t.Errorf("Expected empty line after Ctrl-C, got: %s", line)
	}
	
	// Buffer should be cleared
	if len(ed.buf) != 0 {
		t.Errorf("Expected buffer to be cleared after Ctrl-C, got length: %d", len(ed.buf))
	}
}

// TestEditorCtrlDEOF tests that Ctrl-D on empty line signals EOF
func TestEditorCtrlDEOF(t *testing.T) {
	ed := NewEditor("> ", nil)
	
	// Create a reader with Ctrl-D (0x04) on empty buffer
	input := bytes.NewBufferString("\x04")
	reader := bufio.NewReader(input)
	output := &bytes.Buffer{}
	
	// Read a line - should signal EOF
	line, aborted, eof := ed.ReadLine(reader, output)
	
	if aborted {
		t.Error("Ctrl-D on empty line should not signal abort")
	}
	
	if !eof {
		t.Error("Ctrl-D on empty line should signal EOF")
	}
	
	if line != "" {
		t.Errorf("Expected empty line, got: %s", line)
	}
}

// TestEditorCtrlDDelete tests that Ctrl-D deletes character when buffer is not empty
func TestEditorCtrlDDelete(t *testing.T) {
	ed := NewEditor("> ", nil)
	
	// Set up buffer with text and cursor in the middle
	ed.buf = []rune("test")
	ed.cur = 1 // After 't'
	
	// Create a reader with Ctrl-D (0x04) then Enter
	input := bytes.NewBufferString("\x04\r")
	reader := bufio.NewReader(input)
	output := &bytes.Buffer{}
	
	// Read a line
	line, aborted, eof := ed.ReadLine(reader, output)
	
	if aborted || eof {
		t.Error("Unexpected abort or EOF")
	}
	
	// Should have deleted the 'e' character
	if line != "tst" {
		t.Errorf("Expected 'tst' after Ctrl-D delete, got: %s", line)
	}
}

// TestEditorEnterSubmitsLine tests that Enter submits the line
func TestEditorEnterSubmitsLine(t *testing.T) {
	ed := NewEditor("> ", nil)
	
	// Create input with text and Enter
	input := bytes.NewBufferString("hello\r")
	reader := bufio.NewReader(input)
	output := &bytes.Buffer{}
	
	// Read a line
	line, aborted, eof := ed.ReadLine(reader, output)
	
	if aborted || eof {
		t.Error("Unexpected abort or EOF")
	}
	
	if line != "hello" {
		t.Errorf("Expected 'hello', got: %s", line)
	}
	
	// Check that output contains \r\n
	if !strings.Contains(output.String(), "\r\n") {
		t.Error("Expected output to contain \\r\\n after Enter")
	}
}

// TestEditorBackspace tests that backspace deletes character
func TestEditorBackspace(t *testing.T) {
	ed := NewEditor("> ", nil)
	
	// Create input with text, backspace, and Enter
	input := bytes.NewBufferString("abc\x7f\r") // 0x7f is backspace
	reader := bufio.NewReader(input)
	output := &bytes.Buffer{}
	
	// Read a line
	line, aborted, eof := ed.ReadLine(reader, output)
	
	if aborted || eof {
		t.Error("Unexpected abort or EOF")
	}
	
	if line != "ab" {
		t.Errorf("Expected 'ab' after backspace, got: %s", line)
	}
}
