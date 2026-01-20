package integration

import (
	"math"
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestBusinessTimeUnitConversions tests the conversion of business time units to other time units.
// This addresses the issue where business days couldn't be chained into hours.
func TestBusinessTimeUnitConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		unit     string
		tolerance float64
	}{
		// Business day to time unit conversions
		{
			name:      "business days to hours - direct conversion",
			input:     "5 business days in hours",
			expected:  40,
			unit:      "hours",
			tolerance: 0.01,
		},
		{
			name:      "business day (singular) to hours",
			input:     "1 business day in hours",
			expected:  8,
			unit:      "hours",
			tolerance: 0.01,
		},
		{
			name:      "business days to minutes",
			input:     "2 business days in minutes",
			expected:  960,
			unit:      "minutes",
			tolerance: 0.01,
		},
		{
			name:      "business days to seconds",
			input:     "1 business day in seconds",
			expected:  28800,
			unit:      "seconds",
			tolerance: 0.01,
		},
		
		// Business week to time unit conversions
		{
			name:      "business week to hours",
			input:     "1 business week in hours",
			expected:  40,
			unit:      "hours",
			tolerance: 0.01,
		},
		{
			name:      "business weeks to business days",
			input:     "2 business weeks in business days",
			expected:  10,
			unit:      "business days",
			tolerance: 0.01,
		},
		{
			name:      "business week to seconds",
			input:     "1 business week in seconds",
			expected:  144000,
			unit:      "seconds",
			tolerance: 0.01,
		},
		
		// Business month conversions
		{
			name:      "business month to hours",
			input:     "1 business month in hours",
			expected:  173.33,
			unit:      "hours",
			tolerance: 0.5,
		},
		{
			name:      "business month to business days",
			input:     "1 business month in business days",
			expected:  21.67,
			unit:      "business days",
			tolerance: 0.1,
		},
		{
			name:      "business months to business weeks",
			input:     "3 business months in business weeks",
			expected:  13,
			unit:      "business weeks",
			tolerance: 0.5,
		},
		
		// Chained conversions (the main issue from the bug report)
		{
			name:      "business days to hours to seconds - chained",
			input:     "5 business days in hours in seconds",
			expected:  144000,
			unit:      "seconds",
			tolerance: 0.01,
		},
		{
			name:      "business week to hours to minutes - chained",
			input:     "1 business week in hours in minutes",
			expected:  2400,
			unit:      "minutes",
			tolerance: 0.01,
		},
		{
			name:      "business month to business days to hours - chained",
			input:     "1 business month in business days in hours",
			expected:  173.33,
			unit:      "hours",
			tolerance: 0.5,
		},
		
		// Time units to business units
		{
			name:      "hours to business days",
			input:     "24 hours in business days",
			expected:  3,
			unit:      "business days",
			tolerance: 0.01,
		},
		{
			name:      "hours to business weeks",
			input:     "80 hours in business weeks",
			expected:  2,
			unit:      "business weeks",
			tolerance: 0.01,
		},
		
		// Calendar time to business time
		{
			name:      "calendar days to business days",
			input:     "7 days in business days",
			expected:  21,
			unit:      "business days",
			tolerance: 0.01,
		},
		{
			name:      "calendar week to business weeks",
			input:     "2 weeks in business weeks",
			expected:  8.4,
			unit:      "business weeks",
			tolerance: 0.1,
		},
		
		// Real-world business calculations
		{
			name:      "business year in business days (52.2 weeks)",
			input:     "52.2 business weeks in business days",
			expected:  261,
			unit:      "business days",
			tolerance: 0.1,
		},
		{
			name:      "business year in hours",
			input:     "260 business days in hours",
			expected:  2080,
			unit:      "hours",
			tolerance: 0.1,
		},
		{
			name:      "business quarter approximation",
			input:     "65 business days in business weeks",
			expected:  13,
			unit:      "business weeks",
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			
			if err != nil {
				t.Fatalf("parser error: %v", err)
			}
			
			env := evaluator.NewEnvironment()
			ev := evaluator.New(env)
			result := ev.Eval(expr)
			
			if result.Type == evaluator.ValueError {
				t.Fatalf("evaluation error: %s", result.Error)
			}
			
			if result.Type != evaluator.ValueUnit {
				t.Fatalf("expected unit value, got type %v", result.Type)
			}
			
			// Check the result value
			if math.Abs(result.Number-tt.expected) > tt.tolerance {
				t.Errorf("expected %.2f, got %.2f (tolerance: %.2f)", 
					tt.expected, result.Number, tt.tolerance)
			}
			
			// Check the unit name (case-insensitive)
			if !strings.EqualFold(result.Unit, tt.unit) {
				t.Errorf("expected unit '%s', got '%s'", tt.unit, result.Unit)
			}
		})
	}
}

// TestBusinessTimeUnitWithDateDifferences tests business time units with date calculations.
// This tests the complete workflow: date subtraction → business days → hours.
func TestBusinessTimeUnitWithDateDifferences(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectMin float64
		expectMax float64
		unit      string
	}{
		{
			name:      "date difference to business days to hours",
			input:     "15/01/2024 - 01/01/2024 in business days in hours",
			expectMin: 70,   // 10 business days * 8 hours - some margin
			expectMax: 90,   // 10 business days * 8 hours + some margin
			unit:      "hours",
		},
		{
			name:      "business week calculation from dates",
			input:     "08/01/2024 - 01/01/2024 in business days in business weeks",
			expectMin: 0.9,
			expectMax: 1.1,
			unit:      "business weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			
			if err != nil {
				t.Fatalf("parser error: %v", err)
			}
			
			env := evaluator.NewEnvironment()
			ev := evaluator.New(env)
			result := ev.Eval(expr)
			
			if result.Type == evaluator.ValueError {
				t.Fatalf("evaluation error: %s", result.Error)
			}
			
			if result.Type != evaluator.ValueUnit {
				t.Fatalf("expected unit value, got type %v", result.Type)
			}
			
			// Check the result is within expected range
			if result.Number < tt.expectMin || result.Number > tt.expectMax {
				t.Errorf("expected value between %.2f and %.2f, got %.2f", 
					tt.expectMin, tt.expectMax, result.Number)
			}
			
			// Check the unit name (case-insensitive)
			if !strings.EqualFold(result.Unit, tt.unit) {
				t.Errorf("expected unit '%s', got '%s'", tt.unit, result.Unit)
			}
		})
	}
}

// TestBusinessTimeUnitVariants tests different spelling variants of business time units.
func TestBusinessTimeUnitVariants(t *testing.T) {
	tests := []struct {
		name  string
		input string
		expected float64
	}{
		// Space variants (these are the main requirement)
		{name: "business day (with space)", input: "1 business day in hours", expected: 8},
		{name: "business days (with space, plural)", input: "1 business days in hours", expected: 8},
		
		// Week variants
		{name: "business week (with space)", input: "1 business week in hours", expected: 40},
		{name: "business weeks (with space, plural)", input: "1 business weeks in hours", expected: 40},
		
		// Month variants
		{name: "business month (with space)", input: "1 business month in business days", expected: 21.67},
		{name: "business months (with space, plural)", input: "1 business months in business days", expected: 21.67},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			
			if err != nil {
				t.Fatalf("parser error: %v", err)
			}
			
			env := evaluator.NewEnvironment()
			ev := evaluator.New(env)
			result := ev.Eval(expr)
			
			if result.Type == evaluator.ValueError {
				t.Fatalf("evaluation error: %s", result.Error)
			}
			
			// Allow 5% tolerance for business month calculations
			tolerance := 0.1
			if strings.Contains(tt.input, "month") {
				tolerance = 0.5
			}
			
			if math.Abs(result.Number-tt.expected) > tolerance {
				t.Errorf("expected %.2f, got %.2f", tt.expected, result.Number)
			}
		})
	}
}

// TestIssueReproduction tests the exact scenarios from the bug report.
func TestIssueReproduction(t *testing.T) {
	// These are the exact scenarios from the issue that were failing
	tests := []struct {
		name          string
		expressions   []string
		shouldSucceed bool
	}{
		{
			name: "original failing case: date diff to business days to hours (chained)",
			expressions: []string{
				"03/04/2026 - today in business days in hours",
			},
			shouldSucceed: true,
		},
		{
			name: "original failing case: variable then convert to hours",
			expressions: []string{
				"t = 03/04/2026 - today in business days",
				"t in hours",
			},
			shouldSucceed: true,
		},
		{
			name: "working case: regular days to hours to seconds",
			expressions: []string{
				"204 days in hours in seconds",
			},
			shouldSucceed: true,
		},
		{
			name: "business year calculation",
			expressions: []string{
				"01/01/2026 - 01/01/2027 in business days",
			},
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			
			for i, input := range tt.expressions {
				l := lexer.New(input)
				toks := l.AllTokens()
				p := parser.New(toks)
				expr, err := p.Parse()
				
				if err != nil {
					if tt.shouldSucceed {
						t.Fatalf("expression %d: parser error: %v", i+1, err)
					}
					continue
				}
				
				ev := evaluator.New(env)
				result := ev.Eval(expr)
				
				if result.Type == evaluator.ValueError {
					if tt.shouldSucceed {
						t.Fatalf("expression %d: evaluation error: %s (input: %s)", i+1, result.Error, input)
					}
					continue
				}
				
				if !tt.shouldSucceed {
					t.Fatalf("expression %d: expected error but got result: %v", i+1, result)
				}
				
				// If this was a variable assignment, store it
				if strings.HasPrefix(input, "t =") {
					env.SetVariable("t", result)
				}
			}
		})
	}
}
