package display

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

// Editor provides a minimal line editor with history and control key support.
type Editor struct {
	prompt         string
	buf            []rune
	cur            int
	hist           []string
	hIndex         int // -1 means new entry (not on history)
	hlFn           func(string) string
	autocompleteFn func(string) []Suggestion
	suggestions    []Suggestion
	suggestIndex   int // Current suggestion index (-1 means no active suggestion)
}

// NewEditor creates a new editor instance for a single line entry.
func NewEditor(prompt string, history []string) *Editor {
	return &Editor{
		prompt:       prompt,
		buf:          []rune{},
		cur:          0,
		hist:         append([]string{}, history...),
		hIndex:       -1,
		suggestIndex: -1,
	}
}

// SetHighlighter sets an optional function to colorize the buffer when rendering.
func (e *Editor) SetHighlighter(fn func(string) string) { e.hlFn = fn }

// SetAutocompleteFn sets the autocomplete function.
func (e *Editor) SetAutocompleteFn(fn func(string) []Suggestion) { e.autocompleteFn = fn }

// ReadLine reads a line using raw key processing. It returns the line, whether it was aborted (Ctrl-C), and whether EOF (Ctrl-D on empty).
func (e *Editor) ReadLine(r *bufio.Reader, w io.Writer) (string, bool, bool) {
	e.render(w)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", false, true
		}
		switch b {
		case '\r', '\n':
			// Submit line
			fmt.Fprint(w, "\r\n")
			return string(e.buf), false, false
		case 0x01: // Ctrl-A
			e.cur = 0
		case 0x05: // Ctrl-E
			e.cur = len(e.buf)
		case 0x02: // Ctrl-B
			if e.cur > 0 {
				e.cur--
			}
		case 0x06: // Ctrl-F
			if e.cur < len(e.buf) {
				e.cur++
			}
		case 0x0b: // Ctrl-K
			e.buf = e.buf[:e.cur]
		case 0x15: // Ctrl-U
			e.buf = e.buf[e.cur:]
			e.cur = 0
		case 0x17: // Ctrl-W delete word before
			start := e.wordLeft()
			e.buf = append(e.buf[:start], e.buf[e.cur:]...)
			e.cur = start
		case 0x7f, 0x08: // Backspace / Ctrl-H
			if e.cur > 0 {
				e.buf = append(e.buf[:e.cur-1], e.buf[e.cur:]...)
				e.cur--
			}
		case 0x04: // Ctrl-D
			if len(e.buf) == 0 {
				return "", false, true
			}
			if e.cur < len(e.buf) {
				e.buf = append(e.buf[:e.cur], e.buf[e.cur+1:]...)
			}
		case 0x03: // Ctrl-C abort line
			e.buf = e.buf[:0]
			e.cur = 0
			return "", true, false
		case 0x09: // Tab - trigger autocomplete
			e.handleTab()
		case 0x1b: // ESC sequence
			e.handleEscape(r)
		default:
			// Insert printable rune (assuming UTF-8 single-byte for ASCII; extend minimal multi-byte support)
			if b < 0x80 && (b == ' ' || b >= 0x21) {
				e.insertRune(rune(b))
			} else if b&0xE0 == 0xC0 { // 2-byte UTF-8 start
				rest := make([]byte, 1)
				io.ReadFull(r, rest)
				rr := utf8Decode([]byte{b, rest[0]})
				e.insertRune(rr)
			} else if b&0xF0 == 0xE0 { // 3-byte
				rest := make([]byte, 2)
				io.ReadFull(r, rest)
				rr := utf8Decode(append([]byte{b}, rest...))
				e.insertRune(rr)
			} else if b&0xF8 == 0xF0 { // 4-byte
				rest := make([]byte, 3)
				io.ReadFull(r, rest)
				rr := utf8Decode(append([]byte{b}, rest...))
				e.insertRune(rr)
			}
		}
		e.render(w)
	}
}

func (e *Editor) insertRune(rn rune) {
	if e.cur == len(e.buf) {
		e.buf = append(e.buf, rn)
	} else {
		e.buf = append(e.buf[:e.cur], append([]rune{rn}, e.buf[e.cur:]...)...)
	}
	e.cur++
}

func (e *Editor) handleEscape(r *bufio.Reader) {
	// Look ahead for known sequences
	// Could be: ESC [ A/B/C/D, ESC [1;5C/D, ESC b/f, ESC [3~ (Delete), ESC [H/F, ESC [Z (Shift+Tab)
	seq, _ := r.Peek(1)
	if len(seq) == 0 {
		// Bare ESC - clear suggestions
		e.clearSuggestions()
		return
	}
	if seq[0] == '[' { // CSI
		r.ReadByte()
		// Read the rest up to a letter
		var buf bytes.Buffer
		for {
			b, _ := r.ReadByte()
			if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' || b == '~' {
				// command byte
				cmd := b
				param := buf.String()
				e.handleCSI(cmd, param)
				return
			}
			buf.WriteByte(b)
			if buf.Len() > 8 {
				return
			}
		}
	} else {
		// Possibly ESC b / ESC f
		b, _ := r.ReadByte()
		switch b {
		case 'b':
			e.cur = e.wordLeft()
		case 'f':
			e.cur = e.wordRight()
		}
	}
}

func (e *Editor) handleCSI(cmd byte, param string) {
	switch cmd {
	case 'A': // Up
		e.clearSuggestions()
		e.historyPrev()
	case 'B': // Down
		e.clearSuggestions()
		e.historyNext()
	case 'C': // Right
		e.clearSuggestions()
		if stringsHasSuffix(param, "1;5") { // Ctrl-Right
			e.cur = e.wordRight()
		} else if e.cur < len(e.buf) {
			e.cur++
		}
	case 'D': // Left
		e.clearSuggestions()
		if stringsHasSuffix(param, "1;5") { // Ctrl-Left
			e.cur = e.wordLeft()
		} else if e.cur > 0 {
			e.cur--
		}
	case 'H': // Home
		e.clearSuggestions()
		e.cur = 0
	case 'F': // End
		e.clearSuggestions()
		e.cur = len(e.buf)
	case 'Z': // Shift+Tab
		e.handleShiftTab()
	case '~': // Delete key: ESC [3~
		e.clearSuggestions()
		if param == "3" {
			if e.cur < len(e.buf) {
				e.buf = append(e.buf[:e.cur], e.buf[e.cur+1:]...)
			}
		}
	}
}

func (e *Editor) historyPrev() {
	if len(e.hist) == 0 {
		return
	}
	if e.hIndex == -1 {
		e.hIndex = len(e.hist) - 1
	} else if e.hIndex > 0 {
		e.hIndex--
	}
	e.buf = []rune(e.hist[e.hIndex])
	e.cur = len(e.buf)
}

func (e *Editor) historyNext() {
	if len(e.hist) == 0 {
		return
	}
	if e.hIndex == -1 {
		return
	}
	if e.hIndex < len(e.hist)-1 {
		e.hIndex++
		e.buf = []rune(e.hist[e.hIndex])
		e.cur = len(e.buf)
		return
	}
	// Past the newest -> empty new line
	e.hIndex = -1
	e.buf = e.buf[:0]
	e.cur = 0
}

func (e *Editor) wordLeft() int {
	i := e.cur
	// Skip any spaces before
	for i > 0 && unicode.IsSpace(e.buf[i-1]) {
		i--
	}
	// Skip word
	for i > 0 && isWordRune(e.buf[i-1]) {
		i--
	}
	return i
}

func (e *Editor) wordRight() int {
	i := e.cur
	// Skip spaces
	for i < len(e.buf) && unicode.IsSpace(e.buf[i]) {
		i++
	}
	// Skip word
	for i < len(e.buf) && isWordRune(e.buf[i]) {
		i++
	}
	return i
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func (e *Editor) render(w io.Writer) {
	// Move to line start, clear line, print prompt and buffer, then move cursor back if needed
	fmt.Fprint(w, "\r\x1b[2K")
	fmt.Fprint(w, e.prompt)
	content := string(e.buf)
	if e.hlFn != nil {
		content = e.hlFn(content)
	}
	fmt.Fprint(w, content)
	
	// Save cursor position
	cursorOffset := len(e.buf) - e.cur
	
	// Show current suggestion hint if available
	if len(e.suggestions) > 0 && e.suggestIndex >= 0 && e.suggestIndex < len(e.suggestions) {
		sugg := e.suggestions[e.suggestIndex]
		// Show suggestion hint in gray after the buffer
		fmt.Fprintf(w, " \x1b[90m[%s (%d/%d)]\x1b[0m", sugg.Display, e.suggestIndex+1, len(e.suggestions))
	}
	
	// Move cursor to correct position (back from end of buffer)
	if cursorOffset > 0 {
		// Need to account for the suggestion hint we just printed
		// Move cursor back to the correct position in the buffer
		fmt.Fprintf(w, "\r")
		fmt.Fprint(w, e.prompt)
		// Print buffer content up to cursor position
		contentUpToCursor := string(e.buf[:e.cur])
		if e.hlFn != nil {
			// Can't use highlighter for partial content, just print raw
			fmt.Fprint(w, contentUpToCursor)
		} else {
			fmt.Fprint(w, contentUpToCursor)
		}
	}
}

// Helper: minimal UTF-8 to rune decoder for known-sized buffer
func utf8Decode(b []byte) rune {
	// Very small decoder for 2-4 byte sequences
	switch len(b) {
	case 2:
		return rune(b[0]&0x1F)<<6 | rune(b[1]&0x3F)
	case 3:
		return rune(b[0]&0x0F)<<12 | rune(b[1]&0x3F)<<6 | rune(b[2]&0x3F)
	case 4:
		return rune(b[0]&0x07)<<18 | rune(b[1]&0x3F)<<12 | rune(b[2]&0x3F)<<6 | rune(b[3]&0x3F)
	default:
		if len(b) == 1 {
			return rune(b[0])
		}
		return 0
	}
}

func stringsHasSuffix(s, suf string) bool {
	if len(suf) > len(s) {
		return false
	}
	return s[len(s)-len(suf):] == suf
}

// handleTab handles Tab key press for autocomplete.
func (e *Editor) handleTab() {
	if e.autocompleteFn == nil {
		return
	}

	// If we don't have active suggestions, get them
	if e.suggestIndex == -1 {
		e.suggestions = e.autocompleteFn(string(e.buf))
		if len(e.suggestions) == 0 {
			return
		}
		e.suggestIndex = 0
	} else {
		// Cycle to next suggestion
		e.suggestIndex = (e.suggestIndex + 1) % len(e.suggestions)
	}

	// Apply the suggestion
	e.applySuggestion()
}

// handleShiftTab handles Shift+Tab for reverse cycling (called from handleCSI).
func (e *Editor) handleShiftTab() {
	if e.autocompleteFn == nil || len(e.suggestions) == 0 {
		return
	}

	// Cycle to previous suggestion
	e.suggestIndex--
	if e.suggestIndex < 0 {
		e.suggestIndex = len(e.suggestions) - 1
	}

	// Apply the suggestion
	e.applySuggestion()
}

// applySuggestion replaces the last word with the current suggestion.
func (e *Editor) applySuggestion() {
	if e.suggestIndex < 0 || e.suggestIndex >= len(e.suggestions) {
		return
	}

	suggestion := e.suggestions[e.suggestIndex]
	lastWord := getLastWord(string(e.buf))
	
	if lastWord == "" {
		// Just append the suggestion
		for _, r := range suggestion.Text {
			e.buf = append(e.buf, r)
		}
		e.cur = len(e.buf)
		return
	}

	// Find and replace the last word
	bufStr := string(e.buf)
	lastIdx := -1
	for i := len(bufStr) - len(lastWord); i >= 0; i-- {
		if i+len(lastWord) <= len(bufStr) && bufStr[i:i+len(lastWord)] == lastWord {
			lastIdx = i
			break
		}
	}

	if lastIdx >= 0 {
		// Replace the last word
		newBuf := []rune(bufStr[:lastIdx])
		for _, r := range suggestion.Text {
			newBuf = append(newBuf, r)
		}
		e.buf = newBuf
		e.cur = len(e.buf)
	}
}

// clearSuggestions clears the active suggestions.
func (e *Editor) clearSuggestions() {
	e.suggestions = nil
	e.suggestIndex = -1
}

