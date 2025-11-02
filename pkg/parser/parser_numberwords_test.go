package parser

import (
	"testing"
)

// TestNumberWordsWithUnitsRegressions tests specific failing cases
func TestNumberWordsWithUnitsRegressions(t *testing.T) {
	tests := []struct {
		input       string
		description string
		shouldParse bool
	}{
		{
			"three years and 47 days in days",
			"mixed number words and digits with units",
			true,
		},
		{
			"7 meters divided by twelve cm",
			"division with number words in second operand",
			true,
		},
		{
			"twelve cm",
			"number word with unit",
			true,
		},
		{
			"three years",
			"number word with time unit",
			true,
		},
		{
			"forty two meters",
			"compound number word with unit",
			true,
		},
	}
	
	for _, tt := range tests {
		expr, err := parseInput(tt.input)
		if tt.shouldParse && err != nil {
			t.Errorf("%s: %q failed to parse: %v", tt.description, tt.input, err)
			continue
		}
		if tt.shouldParse && expr == nil {
			t.Errorf("%s: %q returned nil expression", tt.description, tt.input)
		}
		t.Logf("%s: %q -> %T", tt.description, tt.input, expr)
	}
}

// TestNumberWordFollowedByUnit tests that number words can be followed by units
func TestNumberWordFollowedByUnit(t *testing.T) {
	tests := []string{
		"twelve cm",
		"twenty meters",
		"five kg",
		"ten seconds",
		"three hundred km",
		"one thousand grams",
	}
	
	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		
		// Should be parsed as a UnitExpr
		if _, ok := expr.(*UnitExpr); !ok {
			t.Errorf("%q: expected UnitExpr, got %T", input, expr)
		}
	}
}
