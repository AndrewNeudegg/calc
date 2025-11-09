package integration

import (
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestCompoundUnitsWithCurrencySymbols tests currency values with compound units using / and per
func TestCompoundUnitsWithCurrencySymbols(t *testing.T) {
	tests := []struct {
		category    string
		description string
		input       string
		expectType  evaluator.ValueType
		expectUnit  string // expected unit string (e.g., "$/hr")
	}{
		// Dollar with / notation
		{
			category:    "Currency Compound Units",
			description: "$2.93/hr",
			input:       "$2.93/hr",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/hr",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/hour",
			input:       "$2.93/hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/hour",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/h",
			input:       "$2.93/h",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/h",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/s",
			input:       "$2.93/s",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/s",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/second",
			input:       "$2.93/second",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/second",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/ms",
			input:       "$2.93/ms",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/ms",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/millisecond",
			input:       "$2.93/millisecond",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/millisecond",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/min",
			input:       "$2.93/min",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/min",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/minute",
			input:       "$2.93/minute",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/minute",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/day",
			input:       "$2.93/day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/day",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/week",
			input:       "$2.93/week",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/week",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/month",
			input:       "$2.93/month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/month",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/year",
			input:       "$2.93/year",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/year",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93/y (year shorthand)",
			input:       "$2.93/y",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/y",
		},

		// Dollar with per notation
		{
			category:    "Currency Compound Units",
			description: "$2.93 per hour",
			input:       "$2.93 per hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/hour",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per second",
			input:       "$2.93 per second",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/second",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per millisecond",
			input:       "$2.93 per millisecond",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/millisecond",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per minute",
			input:       "$2.93 per minute",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/minute",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per day",
			input:       "$2.93 per day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/day",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per week",
			input:       "$2.93 per week",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/week",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per month",
			input:       "$2.93 per month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/month",
		},
		{
			category:    "Currency Compound Units",
			description: "$2.93 per year",
			input:       "$2.93 per year",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "$/year",
		},

		// Other currencies
		{
			category:    "Currency Compound Units",
			description: "£50/hour",
			input:       "£50/hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "£/hour",
		},
		{
			category:    "Currency Compound Units",
			description: "€100/day",
			input:       "€100/day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "€/day",
		},
		{
			category:    "Currency Compound Units",
			description: "¥1000/month",
			input:       "¥1000/month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "¥/month",
		},
		{
			category:    "Currency Compound Units",
			description: "£50 per hour",
			input:       "£50 per hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "£/hour",
		},
		{
			category:    "Currency Compound Units",
			description: "€100 per day",
			input:       "€100 per day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "€/day",
		},
		{
			category:    "Currency Compound Units",
			description: "¥1000 per month",
			input:       "¥1000 per month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "¥/month",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			result := e.Eval(expr)

			if result.IsError() {
				t.Fatalf("Evaluation error for %q: %s", tt.input, result.Error)
			}

			if result.Type != tt.expectType {
				t.Errorf("Expected type %v for %q, got %v", tt.expectType, tt.input, result.Type)
			}

			if tt.expectType == evaluator.ValueUnit {
				if result.Unit != tt.expectUnit {
					t.Errorf("Expected unit %q for %q, got %q", tt.expectUnit, tt.input, result.Unit)
				}
			}
		})
	}
}

// TestCompoundUnitsArithmetic tests arithmetic operations with compound units
func TestCompoundUnitsArithmetic(t *testing.T) {
	tests := []struct {
		description string
		lines       []string // multiple lines to execute in sequence
		checkFinal  func(evaluator.Value) bool
	}{
		{
			description: "Calculate daily cost from hourly rate",
			lines: []string{
				"ab_cost = $2.93/hr",
				"ab_cost * 24",
			},
			checkFinal: func(v evaluator.Value) bool {
				return v.Type == evaluator.ValueUnit && v.Number > 70 && v.Number < 71
			},
		},
		{
			description: "Calculate weekly cost from hourly rate",
			lines: []string{
				"hourly = $50/hour",
				"hourly * 40", // 40 hour work week
			},
			checkFinal: func(v evaluator.Value) bool {
				return v.Type == evaluator.ValueUnit && v.Number == 2000
			},
		},
		{
			description: "Calculate monthly cost from daily rate",
			lines: []string{
				"daily = £100/day",
				"daily * 30",
			},
			checkFinal: func(v evaluator.Value) bool {
				return v.Type == evaluator.ValueUnit && v.Number == 3000
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			var result evaluator.Value

			for _, line := range tt.lines {
				l := lexer.New(line)
				tokens := l.AllTokens()
				p := parser.New(tokens)
				expr, err := p.Parse()
				if err != nil {
					t.Fatalf("Parse error for %q: %v", line, err)
				}
				result = e.Eval(expr)
				if result.IsError() {
					t.Fatalf("Evaluation error for %q: %s", line, result.Error)
				}
			}

			if !tt.checkFinal(result) {
				t.Errorf("Final value check failed for %q: got %+v", tt.description, result)
			}
		})
	}
}

// TestIssueExamples tests the specific examples from the GitHub issue
func TestIssueExamples(t *testing.T) {
	tests := []struct {
		description string
		input       string
		shouldWork  bool
	}{
		{
			description: "ab_cost = $2.93/hr",
			input:       "$2.93/hr",
			shouldWork:  true,
		},
		{
			description: "ab_cost = $2.93/hour",
			input:       "$2.93/hour",
			shouldWork:  true,
		},
		{
			description: "ab_cost = $2.93 per hour",
			input:       "$2.93 per hour",
			shouldWork:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected parse success for %q, got error: %v", tt.input, err)
					return
				}

				env := evaluator.NewEnvironment()
				e := evaluator.New(env)
				result := e.Eval(expr)

				if result.IsError() {
					t.Errorf("Expected evaluation success for %q, got error: %s", tt.input, result.Error)
					return
				}

				if result.Type != evaluator.ValueUnit {
					t.Errorf("Expected ValueUnit for %q, got %v", tt.input, result.Type)
				}

				// Check that the unit contains a slash (compound unit)
				if !strings.Contains(result.Unit, "/") {
					t.Errorf("Expected compound unit (with /) for %q, got %q", tt.input, result.Unit)
				}
			}
		})
	}
}
