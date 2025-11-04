package display

import (
	"bufio"
	"bytes"
	"testing"
)

// helper to run an editor session with given input bytes and history
func runEditor(t *testing.T, history []string, input []byte) (string, bool, bool) {
	t.Helper()
	ed := NewEditor("> ", history)
	r := bufio.NewReader(bytes.NewReader(input))
	var out bytes.Buffer
	line, aborted, eof := ed.ReadLine(r, &out)
	return line, aborted, eof
}

func TestEditor_Insert_Backspace(t *testing.T) {
	// a b c <backspace> d <enter>
	input := []byte{'a', 'b', 'c', 0x7f, 'd', '\n'}
	line, aborted, eof := runEditor(t, nil, input)
	if aborted || eof {
		t.Fatalf("unexpected aborted=%v eof=%v", aborted, eof)
	}
	if line != "abd" {
		t.Fatalf("expected 'abd', got %q", line)
	}
}

func TestEditor_Move_Insert(t *testing.T) {
	// a b c <Left> <Left> X <enter>
	input := []byte{'a', 'b', 'c', 0x1b, '[', 'D', 0x1b, '[', 'D', 'X', '\n'}
	line, aborted, eof := runEditor(t, nil, input)
	if aborted || eof {
		t.Fatalf("unexpected aborted=%v eof=%v", aborted, eof)
	}
	if line != "aXbc" {
		t.Fatalf("expected 'aXbc', got %q", line)
	}
}

func TestEditor_WordNav_DeletePrevWord(t *testing.T) {
	// "hello world" ESC b (word left) Ctrl-W (delete prev word) Enter
	input := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', 0x1b, 'b', 0x17, '\n'}
	line, aborted, eof := runEditor(t, nil, input)
	if aborted || eof {
		t.Fatalf("unexpected aborted=%v eof=%v", aborted, eof)
	}
	if line != "world" {
		t.Fatalf("expected 'world', got %q", line)
	}
}

func TestEditor_HistoryRecall(t *testing.T) {
	history := []string{"one", "two", "three"}
	// Up, Enter
	input := []byte{0x1b, '[', 'A', '\n'}
	line, aborted, eof := runEditor(t, history, input)
	if aborted || eof {
		t.Fatalf("unexpected aborted=%v eof=%v", aborted, eof)
	}
	if line != "three" {
		t.Fatalf("expected 'three', got %q", line)
	}
}

func TestEditor_CtrlA_CtrlE(t *testing.T) {
	// a b c Ctrl-A X Ctrl-E Y Enter => XabcY
	input := []byte{'a', 'b', 'c', 0x01, 'X', 0x05, 'Y', '\n'}
	line, aborted, eof := runEditor(t, nil, input)
	if aborted || eof {
		t.Fatalf("unexpected aborted=%v eof=%v", aborted, eof)
	}
	if line != "XabcY" {
		t.Fatalf("expected 'XabcY', got %q", line)
	}
}
