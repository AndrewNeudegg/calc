package display

import (
	"sort"
	"strings"

	"github.com/andrewneudegg/calc/pkg/currency"
	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/settings"
	"github.com/andrewneudegg/calc/pkg/units"
)

// Suggestion represents a single autocomplete suggestion.
type Suggestion struct {
	Text        string // The text to insert
	Display     string // The text to display (may include description)
	Category    string // Category label (command, variable, function, unit, currency, keyword)
	Description string // Optional short description
}

// AutocompleteEngine generates suggestions based on the current input.
type AutocompleteEngine struct {
	commands  []Suggestion
	functions []Suggestion
	keywords  []Suggestion
	env       *evaluator.Environment
	units     *units.System
	currency  *currency.System
	settings  *settings.Settings
}

// NewAutocompleteEngine creates a new autocomplete engine.
func NewAutocompleteEngine(env *evaluator.Environment, u *units.System, c *currency.System, s *settings.Settings) *AutocompleteEngine {
	ac := &AutocompleteEngine{
		env:      env,
		units:    u,
		currency: c,
		settings: s,
	}
	ac.initCommands()
	ac.initFunctions()
	ac.initKeywords()
	return ac
}

func (ac *AutocompleteEngine) initCommands() {
	ac.commands = []Suggestion{
		{Text: ":help", Display: ":help", Category: "command", Description: "Show available commands"},
		{Text: ":save ", Display: ":save <file>", Category: "command", Description: "Save current workspace"},
		{Text: ":open ", Display: ":open <file>", Category: "command", Description: "Open a workspace file"},
		{Text: ":set ", Display: ":set <key> <val>", Category: "command", Description: "Set a preference"},
		{Text: ":clear", Display: ":clear", Category: "command", Description: "Clear screen and reset session"},
		{Text: ":quiet ", Display: ":quiet [on|off]", Category: "command", Description: "Toggle quiet mode"},
		{Text: ":tz ", Display: ":tz list", Category: "command", Description: "List timezones"},
		{Text: ":quit", Display: ":quit", Category: "command", Description: "Exit the program"},
		{Text: ":exit", Display: ":exit", Category: "command", Description: "Exit the program"},
		{Text: ":q", Display: ":q", Category: "command", Description: "Exit the program"},
	}
}

func (ac *AutocompleteEngine) initFunctions() {
	ac.functions = []Suggestion{
		{Text: "sum(", Display: "sum(...)", Category: "function", Description: "Sum of all arguments"},
		{Text: "average(", Display: "average(...)", Category: "function", Description: "Average of arguments"},
		{Text: "mean(", Display: "mean(...)", Category: "function", Description: "Mean of arguments"},
		{Text: "min(", Display: "min(...)", Category: "function", Description: "Minimum of arguments"},
		{Text: "max(", Display: "max(...)", Category: "function", Description: "Maximum of arguments"},
		{Text: "print(\"", Display: "print(\"...\")", Category: "function", Description: "Print with variable interpolation"},
	}
}

func (ac *AutocompleteEngine) initKeywords() {
	// Only include keywords when fuzzy mode is enabled
	fuzzyKeywords := []Suggestion{
		{Text: "half of ", Display: "half of", Category: "keyword", Description: "Half of value"},
		{Text: "double ", Display: "double", Category: "keyword", Description: "Double value"},
		{Text: "twice ", Display: "twice", Category: "keyword", Description: "Twice value"},
		{Text: "three quarters of ", Display: "three quarters of", Category: "keyword", Description: "Three quarters of value"},
		{Text: "increase ", Display: "increase X by Y%", Category: "keyword", Description: "Increase by percentage"},
		{Text: "decrease ", Display: "decrease X by Y%", Category: "keyword", Description: "Decrease by percentage"},
	}

	// Date/time keywords (always available)
	dateKeywords := []Suggestion{
		{Text: "now", Display: "now", Category: "keyword", Description: "Current date/time"},
		{Text: "today", Display: "today", Category: "keyword", Description: "Current date"},
		{Text: "tomorrow", Display: "tomorrow", Category: "keyword", Description: "Tomorrow's date"},
		{Text: "yesterday", Display: "yesterday", Category: "keyword", Description: "Yesterday's date"},
		{Text: "next week", Display: "next week", Category: "keyword", Description: "Date one week ahead"},
		{Text: "last week", Display: "last week", Category: "keyword", Description: "Date one week ago"},
		{Text: "next month", Display: "next month", Category: "keyword", Description: "Date one month ahead"},
	}

	// Previous result keywords (always available)
	prevKeywords := []Suggestion{
		{Text: "prev", Display: "prev", Category: "keyword", Description: "Most recent result"},
		{Text: "prev~", Display: "prev~N", Category: "keyword", Description: "Result N steps back"},
		{Text: "prev#", Display: "prev#N", Category: "keyword", Description: "Result at line N"},
	}

	ac.keywords = append(dateKeywords, prevKeywords...)
	if ac.settings.FuzzyMode {
		ac.keywords = append(ac.keywords, fuzzyKeywords...)
	}
}

// GetSuggestions returns suggestions based on the current input.
func (ac *AutocompleteEngine) GetSuggestions(input string) []Suggestion {
	if !ac.settings.Autocomplete {
		return nil
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var suggestions []Suggestion

	// Command suggestions (input starts with :)
	if strings.HasPrefix(input, ":") {
		for _, cmd := range ac.commands {
			if strings.HasPrefix(cmd.Text, input) {
				suggestions = append(suggestions, cmd)
			}
		}
		return suggestions
	}

	// Get the last word or partial word being typed
	lastWord := getLastWord(input)
	if lastWord == "" {
		return nil
	}

	// Variable suggestions
	variables := ac.getVariables()
	for _, v := range variables {
		if strings.HasPrefix(strings.ToLower(v.Text), strings.ToLower(lastWord)) {
			suggestions = append(suggestions, v)
		}
	}

	// Function suggestions
	for _, fn := range ac.functions {
		if strings.HasPrefix(strings.ToLower(fn.Text), strings.ToLower(lastWord)) {
			suggestions = append(suggestions, fn)
		}
	}

	// Keyword suggestions
	for _, kw := range ac.keywords {
		if strings.HasPrefix(strings.ToLower(kw.Text), strings.ToLower(lastWord)) {
			suggestions = append(suggestions, kw)
		}
	}

	// Unit suggestions
	unitSuggestions := ac.getUnits(lastWord)
	suggestions = append(suggestions, unitSuggestions...)

	// Currency suggestions
	currencySuggestions := ac.getCurrencies(lastWord)
	suggestions = append(suggestions, currencySuggestions...)

	// Sort by relevance (exact match first, then alphabetically)
	sort.SliceStable(suggestions, func(i, j int) bool {
		// Prioritize exact prefix matches
		iExact := strings.HasPrefix(strings.ToLower(suggestions[i].Text), strings.ToLower(lastWord))
		jExact := strings.HasPrefix(strings.ToLower(suggestions[j].Text), strings.ToLower(lastWord))
		if iExact != jExact {
			return iExact
		}
		return suggestions[i].Text < suggestions[j].Text
	})

	return suggestions
}

func (ac *AutocompleteEngine) getVariables() []Suggestion {
	var suggestions []Suggestion
	// Access environment variables
	varNames := ac.env.GetVariableNames()
	for _, name := range varNames {
		suggestions = append(suggestions, Suggestion{
			Text:     name,
			Display:  name,
			Category: "variable",
		})
	}
	return suggestions
}

func (ac *AutocompleteEngine) getUnits(prefix string) []Suggestion {
	var suggestions []Suggestion
	// Common units to suggest
	commonUnits := []string{
		"m", "cm", "mm", "km", "ft", "in", "mi",
		"kg", "g", "mg", "lb", "oz",
		"s", "min", "h", "day", "week", "month", "year",
		"l", "ml", "gal", "pt",
		"°c", "°f", "k",
	}

	for _, u := range commonUnits {
		if strings.HasPrefix(strings.ToLower(u), strings.ToLower(prefix)) && ac.units.IsUnit(u) {
			suggestions = append(suggestions, Suggestion{
				Text:     u,
				Display:  u,
				Category: "unit",
			})
		}
	}

	return suggestions
}

func (ac *AutocompleteEngine) getCurrencies(prefix string) []Suggestion {
	var suggestions []Suggestion
	// Common currencies
	currencies := []struct {
		code   string
		symbol string
		name   string
	}{
		{"usd", "$", "US Dollar"},
		{"gbp", "£", "British Pound"},
		{"eur", "€", "Euro"},
		{"jpy", "¥", "Japanese Yen"},
		{"aud", "AUD", "Australian Dollar"},
		{"cad", "CAD", "Canadian Dollar"},
		{"chf", "CHF", "Swiss Franc"},
		{"cny", "CNY", "Chinese Yuan"},
	}

	for _, c := range currencies {
		if strings.HasPrefix(strings.ToLower(c.code), strings.ToLower(prefix)) && ac.currency.IsCurrency(c.code) {
			suggestions = append(suggestions, Suggestion{
				Text:        c.code,
				Display:     c.code,
				Category:    "currency",
				Description: c.name,
			})
		}
	}

	return suggestions
}

// isWordDelimiter checks if a rune is a word boundary character.
func isWordDelimiter(r rune) bool {
	return r == ' ' || r == '(' || r == ')' || r == ',' || r == '+' || r == '-' || r == '*' || r == '/'
}

// getLastWord extracts the last word or partial word from the input.
func getLastWord(input string) string {
	// If input ends with a delimiter, return empty (starting new word)
	if len(input) > 0 {
		lastChar := rune(input[len(input)-1])
		if isWordDelimiter(lastChar) {
			return ""
		}
	}
	
	// Find the last word boundary
	words := strings.FieldsFunc(input, isWordDelimiter)
	if len(words) == 0 {
		return ""
	}
	return words[len(words)-1]
}
