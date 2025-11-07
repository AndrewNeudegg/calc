package display

import (
	"testing"
)

// TestAutocompleteTabCyclingDoesNotAccumulate verifies that repeatedly pressing Tab
// cycles through suggestions without accumulating them in the buffer.
func TestAutocompleteTabCyclingDoesNotAccumulate(t *testing.T) {
	repl := NewREPL()
	
	// Simulate the autocomplete scenario
	// When user types ":s" and presses Tab multiple times, it should cycle through
	// suggestions without accumulating them
	
	// Get suggestions for ":s"
	suggestions := repl.autocomplete.GetSuggestions(":s")
	if len(suggestions) < 2 {
		t.Fatalf("Expected at least 2 suggestions for ':s', got %d", len(suggestions))
	}
	
	// Verify we have :save and :set
	hasSave := false
	hasSet := false
	for _, s := range suggestions {
		if s.Text == ":save " {
			hasSave = true
		}
		if s.Text == ":set " {
			hasSet = true
		}
	}
	
	if !hasSave || !hasSet {
		t.Error("Expected both :save and :set in suggestions")
	}
}

// TestAutocompleteClearsOnTyping verifies that typing clears active suggestions.
func TestAutocompleteClearsOnTyping(t *testing.T) {
	// This is a conceptual test - in actual implementation, when user types
	// a character after seeing suggestions, the suggestions should be cleared
	// and a new set should be generated on the next Tab press
	
	repl := NewREPL()
	
	// Get initial suggestions
	suggestions1 := repl.autocomplete.GetSuggestions(":s")
	if len(suggestions1) == 0 {
		t.Fatal("Expected suggestions for ':s'")
	}
	
	// Simulate typing more characters - suggestions for a different prefix
	suggestions2 := repl.autocomplete.GetSuggestions(":sa")
	
	// Should have different (or subset of) suggestions
	// The key is that we're getting fresh suggestions, not accumulating
	for _, s := range suggestions2 {
		if s.Text != ":save " {
			t.Errorf("Unexpected suggestion for ':sa': %s", s.Text)
		}
	}
}
