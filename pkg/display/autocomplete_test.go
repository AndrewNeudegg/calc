package display

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/currency"
	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/settings"
	"github.com/andrewneudegg/calc/pkg/units"
)

func TestAutocompleteSuggestsCommands(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()

	ac := NewAutocompleteEngine(env, u, c, s)

	// Test command suggestions
	suggestions := ac.GetSuggestions(":h")
	if len(suggestions) == 0 {
		t.Fatal("Expected command suggestions for ':h'")
	}

	found := false
	for _, sugg := range suggestions {
		if sugg.Text == ":help" && sugg.Category == "command" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ':help' command in suggestions")
	}
}

func TestAutocompleteSuggestsVariables(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()

	// Add some variables
	env.SetVariable("myVar", evaluator.NewNumber(42))
	env.SetVariable("myOtherVar", evaluator.NewNumber(100))

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("my")
	if len(suggestions) < 2 {
		t.Fatalf("Expected at least 2 variable suggestions, got %d", len(suggestions))
	}

	// Check that both variables are suggested
	foundVar1 := false
	foundVar2 := false
	for _, sugg := range suggestions {
		if sugg.Text == "myVar" && sugg.Category == "variable" {
			foundVar1 = true
		}
		if sugg.Text == "myOtherVar" && sugg.Category == "variable" {
			foundVar2 = true
		}
	}

	if !foundVar1 || !foundVar2 {
		t.Error("Expected both 'myVar' and 'myOtherVar' in suggestions")
	}
}

func TestAutocompleteSuggestsFunctions(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("su")
	found := false
	for _, sugg := range suggestions {
		if sugg.Text == "sum(" && sugg.Category == "function" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'sum' function in suggestions for 'su'")
	}
}

func TestAutocompleteSuggestsKeywords(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()
	s.FuzzyMode = true

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("hal")
	found := false
	for _, sugg := range suggestions {
		if sugg.Text == "half of " && sugg.Category == "keyword" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'half of' keyword in suggestions for 'hal'")
	}
}

func TestAutocompleteRespectsSettings(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()
	s.Autocomplete = false

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions(":h")
	if len(suggestions) != 0 {
		t.Error("Expected no suggestions when autocomplete is disabled")
	}
}

func TestAutocompleteFuzzyKeywordsRespectSettings(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()
	s.FuzzyMode = false

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("hal")
	for _, sugg := range suggestions {
		if sugg.Text == "half of " {
			t.Error("Fuzzy keyword 'half of' should not appear when fuzzy mode is off")
		}
	}

	// Date keywords should still appear
	suggestions = ac.GetSuggestions("tod")
	found := false
	for _, sugg := range suggestions {
		if sugg.Text == "today" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Date keyword 'today' should appear even when fuzzy mode is off")
	}
}

func TestGetLastWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello world", "world"},
		{"10 + myVar", "myVar"},
		{"sum(a, b", "b"},
		{"value * ", ""},
		{"", ""},
		{"x+y", "y"},
	}

	for _, tt := range tests {
		result := getLastWord(tt.input)
		if result != tt.expected {
			t.Errorf("getLastWord(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestAutocompleteSuggestsCurrencies(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("us")
	found := false
	for _, sugg := range suggestions {
		if sugg.Text == "usd" && sugg.Category == "currency" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'usd' currency in suggestions for 'us'")
	}
}

func TestAutocompleteSuggestsUnits(t *testing.T) {
	env := evaluator.NewEnvironment()
	u := units.NewSystem()
	c := currency.NewSystem()
	s := settings.Default()

	ac := NewAutocompleteEngine(env, u, c, s)

	suggestions := ac.GetSuggestions("k")
	foundKm := false
	foundKg := false
	for _, sugg := range suggestions {
		if sugg.Text == "km" && sugg.Category == "unit" {
			foundKm = true
		}
		if sugg.Text == "kg" && sugg.Category == "unit" {
			foundKg = true
		}
	}

	if !foundKm || !foundKg {
		t.Error("Expected 'km' and 'kg' units in suggestions for 'k'")
	}
}
