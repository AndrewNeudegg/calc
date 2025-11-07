package display

import (
	"testing"
)

// Integration tests for autocomplete in REPL context

func TestAutocompleteIntegrationWithVariables(t *testing.T) {
	// Create a REPL and set up some variables
	repl := NewREPL()
	
	// Evaluate some lines to create variables
	repl.EvaluateLine("myRate = 42")
	repl.EvaluateLine("myValue = 100")
	repl.EvaluateLine("total = 500")
	
	// Test that autocomplete suggests these variables
	suggestions := repl.autocomplete.GetSuggestions("my")
	
	foundRate := false
	foundValue := false
	for _, s := range suggestions {
		if s.Text == "myRate" {
			foundRate = true
		}
		if s.Text == "myValue" {
			foundValue = true
		}
	}
	
	if !foundRate || !foundValue {
		t.Error("Expected both myRate and myValue in suggestions")
	}
}

func TestAutocompleteIntegrationWithSettings(t *testing.T) {
	repl := NewREPL()
	
	// Disable autocomplete
	repl.settings.Autocomplete = false
	repl.autocomplete.settings = repl.settings
	
	suggestions := repl.autocomplete.GetSuggestions(":h")
	if len(suggestions) != 0 {
		t.Error("Expected no suggestions when autocomplete is disabled")
	}
	
	// Re-enable autocomplete
	repl.settings.Autocomplete = true
	suggestions = repl.autocomplete.GetSuggestions(":h")
	if len(suggestions) == 0 {
		t.Error("Expected suggestions when autocomplete is enabled")
	}
}

func TestAutocompleteIntegrationWithFuzzyMode(t *testing.T) {
	repl := NewREPL()
	
	// Enable fuzzy mode
	repl.settings.FuzzyMode = true
	repl.autocomplete = NewAutocompleteEngine(repl.env, repl.env.Units(), repl.env.Currency(), repl.settings)
	
	suggestions := repl.autocomplete.GetSuggestions("hal")
	foundHalf := false
	for _, s := range suggestions {
		if s.Text == "half of " {
			foundHalf = true
			break
		}
	}
	if !foundHalf {
		t.Error("Expected 'half of' keyword when fuzzy mode is on")
	}
	
	// Disable fuzzy mode
	repl.settings.FuzzyMode = false
	repl.autocomplete = NewAutocompleteEngine(repl.env, repl.env.Units(), repl.env.Currency(), repl.settings)
	
	suggestions = repl.autocomplete.GetSuggestions("hal")
	for _, s := range suggestions {
		if s.Text == "half of " {
			t.Error("'half of' keyword should not appear when fuzzy mode is off")
		}
	}
}

func TestAutocompleteIntegrationAllCategories(t *testing.T) {
	repl := NewREPL()
	
	// Test commands
	cmdSuggestions := repl.autocomplete.GetSuggestions(":s")
	if len(cmdSuggestions) == 0 {
		t.Error("Expected command suggestions for ':s'")
	}
	
	// Test variables after defining them
	repl.EvaluateLine("testVar = 10")
	varSuggestions := repl.autocomplete.GetSuggestions("test")
	foundVar := false
	for _, s := range varSuggestions {
		if s.Text == "testVar" && s.Category == "variable" {
			foundVar = true
			break
		}
	}
	if !foundVar {
		t.Error("Expected variable suggestion for 'test'")
	}
	
	// Test functions
	fnSuggestions := repl.autocomplete.GetSuggestions("av")
	foundFn := false
	for _, s := range fnSuggestions {
		if s.Text == "average(" && s.Category == "function" {
			foundFn = true
			break
		}
	}
	if !foundFn {
		t.Error("Expected function suggestion for 'av'")
	}
	
	// Test units
	unitSuggestions := repl.autocomplete.GetSuggestions("km")
	foundUnit := false
	for _, s := range unitSuggestions {
		if s.Text == "km" && s.Category == "unit" {
			foundUnit = true
			break
		}
	}
	if !foundUnit {
		t.Error("Expected unit suggestion for 'km'")
	}
	
	// Test currencies
	currSuggestions := repl.autocomplete.GetSuggestions("gb")
	foundCurr := false
	for _, s := range currSuggestions {
		if s.Text == "gbp" && s.Category == "currency" {
			foundCurr = true
			break
		}
	}
	if !foundCurr {
		t.Error("Expected currency suggestion for 'gb'")
	}
	
	// Test keywords
	kwSuggestions := repl.autocomplete.GetSuggestions("tod")
	foundKw := false
	for _, s := range kwSuggestions {
		if s.Text == "today" && s.Category == "keyword" {
			foundKw = true
			break
		}
	}
	if !foundKw {
		t.Error("Expected keyword suggestion for 'tod'")
	}
}
