package display

import (
	"testing"
)

// TestEditorTabCycling tests that Tab cycling doesn't accumulate suggestions in the buffer
func TestEditorTabCycling(t *testing.T) {
	ed := NewEditor("test> ", []string{})
	
	// Set up a mock autocomplete function that returns multiple suggestions
	suggestions := []Suggestion{
		{Text: ":save ", Display: ":save <file>", Category: "command"},
		{Text: ":set ", Display: ":set <key> <val>", Category: "command"},
		{Text: ":show ", Display: ":show", Category: "command"},
	}
	
	ed.SetAutocompleteFn(func(input string) []Suggestion {
		// Return suggestions if input is ":s"
		if input == ":s" {
			return suggestions
		}
		return nil
	})
	
	// Simulate typing ":s"
	ed.buf = []rune(":s")
	ed.cur = 2
	
	// First Tab - should apply first suggestion
	ed.handleTab()
	result1 := string(ed.buf)
	if result1 != ":save " {
		t.Errorf("First Tab: expected ':save ', got '%s'", result1)
	}
	
	// Second Tab - should cycle to second suggestion, not append
	ed.handleTab()
	result2 := string(ed.buf)
	if result2 != ":set " {
		t.Errorf("Second Tab: expected ':set ', got '%s'", result2)
	}
	
	// Third Tab - should cycle to third suggestion
	ed.handleTab()
	result3 := string(ed.buf)
	if result3 != ":show " {
		t.Errorf("Third Tab: expected ':show ', got '%s'", result3)
	}
	
	// Fourth Tab - should wrap back to first suggestion
	ed.handleTab()
	result4 := string(ed.buf)
	if result4 != ":save " {
		t.Errorf("Fourth Tab (wrap): expected ':save ', got '%s'", result4)
	}
	
	// Verify the buffer length is reasonable (not accumulating)
	if len(ed.buf) > 10 {
		t.Errorf("Buffer too long: %d characters, expected <= 10. Content: '%s'", len(ed.buf), string(ed.buf))
	}
}

// TestEditorTypingClearsSuggestions tests that typing clears active suggestions
func TestEditorTypingClearsSuggestions(t *testing.T) {
	ed := NewEditor("test> ", []string{})
	
	suggestions := []Suggestion{
		{Text: "myVar", Display: "myVar", Category: "variable"},
		{Text: "myValue", Display: "myValue", Category: "variable"},
	}
	
	ed.SetAutocompleteFn(func(input string) []Suggestion {
		if input == "my" {
			return suggestions
		}
		return nil
	})
	
	// Set up initial state with suggestions
	ed.buf = []rune("my")
	ed.cur = 2
	ed.suggestions = suggestions
	ed.suggestIndex = 0
	ed.originalBuf = []rune("my")
	
	// Verify suggestions are active
	if ed.suggestIndex != 0 {
		t.Error("Suggestions should be active before typing")
	}
	
	// Type a character
	ed.insertRune('V')
	
	// Verify suggestions are cleared
	if ed.suggestIndex != -1 {
		t.Error("Suggestions should be cleared after typing")
	}
	if ed.originalBuf != nil {
		t.Error("originalBuf should be cleared after typing")
	}
	if ed.suggestions != nil {
		t.Error("suggestions should be cleared after typing")
	}
}
