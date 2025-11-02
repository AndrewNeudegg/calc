package lexer

import (
	"testing"
)

func TestParseNumberWords(t *testing.T) {
	tests := []struct {
		input    []string
		expected float64
		locale   string
	}{
		{[]string{"zero"}, 0, "en_GB"},
		{[]string{"one"}, 1, "en_GB"},
		{[]string{"five"}, 5, "en_GB"},
		{[]string{"ten"}, 10, "en_GB"},
		{[]string{"twenty"}, 20, "en_GB"},
		{[]string{"thirty"}, 30, "en_GB"},
		{[]string{"hundred"}, 100, "en_GB"},
		{[]string{"thousand"}, 1000, "en_GB"},
		
		// Compound numbers
		{[]string{"twenty", "one"}, 21, "en_GB"},
		{[]string{"thirty", "five"}, 35, "en_GB"},
		{[]string{"ninety", "nine"}, 99, "en_GB"},
		
		// With "and"
		{[]string{"twenty", "and", "one"}, 21, "en_GB"},
		{[]string{"fifty", "and", "seven"}, 57, "en_GB"},
		
		// Hundreds
		{[]string{"one", "hundred"}, 100, "en_GB"},
		{[]string{"two", "hundred"}, 200, "en_GB"},
		{[]string{"five", "hundred"}, 500, "en_GB"},
		{[]string{"one", "hundred", "and", "one"}, 101, "en_GB"},
		{[]string{"three", "hundred", "and", "forty", "two"}, 342, "en_GB"},
		
		// Thousands
		{[]string{"one", "thousand"}, 1000, "en_GB"},
		{[]string{"five", "thousand"}, 5000, "en_GB"},
		{[]string{"ten", "thousand"}, 10000, "en_GB"},
		{[]string{"twenty", "thousand"}, 20000, "en_GB"},
		
		// Complex numbers
		{[]string{"three", "thousand", "five", "hundred"}, 3500, "en_GB"},
		{[]string{"five", "thousand", "and", "twenty"}, 5020, "en_GB"},
		{[]string{"one", "hundred", "thousand"}, 100000, "en_GB"},
		
		// Million
		{[]string{"one", "million"}, 1000000, "en_GB"},
		{[]string{"five", "million"}, 5000000, "en_GB"},
	}
	
	for _, tt := range tests {
		result, ok := ParseNumberWords(tt.input, tt.locale)
		if !ok {
			t.Errorf("ParseNumberWords(%v) failed to parse", tt.input)
			continue
		}
		if result != tt.expected {
			t.Errorf("ParseNumberWords(%v) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}

func TestParseNumberWordsFail(t *testing.T) {
	tests := []struct {
		input []string
	}{
		{[]string{"not", "a", "number"}},
		{[]string{"hello"}},
		{[]string{""}},
	}
	
	for _, tt := range tests {
		_, ok := ParseNumberWords(tt.input, "en_GB")
		if ok {
			t.Errorf("ParseNumberWords(%v) should have failed", tt.input)
		}
	}
}

func TestIsNumberWord(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"one", true},
		{"two", true},
		{"three", true},
		{"hundred", true},
		{"thousand", true},
		{"and", true},
		{"a", true},
		{"hello", false},
		{"world", false},
	}
	
	for _, tt := range tests {
		result := IsNumberWord(tt.word, "en_GB")
		if result != tt.expected {
			t.Errorf("IsNumberWord(%q) = %v, want %v", tt.word, result, tt.expected)
		}
	}
}
