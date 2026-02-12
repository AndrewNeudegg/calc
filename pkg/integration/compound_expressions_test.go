package integration

import (
	"math"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestCompoundExpressionsWithConversions tests compound expressions that include
// conversions followed by arithmetic operations
func TestCompoundExpressionsWithConversions(t *testing.T) {
	tests := []struct {
		description string
		input       string
		checkResult func(evaluator.Value) bool
	}{
		{
			description: "Unit conversion followed by division",
			input:       "100 cm in m / 2",
			checkResult: func(v evaluator.Value) bool {
				// 100 cm = 1 m, then 1 / 2 = 0.5
				return !v.IsError() && math.Abs(v.Number-0.5) < 0.01
			},
		},
		{
			description: "Unit conversion followed by multiplication",
			input:       "100 cm in m * 2",
			checkResult: func(v evaluator.Value) bool {
				// 100 cm = 1 m, then 1 * 2 = 2
				return !v.IsError() && math.Abs(v.Number-2.0) < 0.01
			},
		},
		{
			description: "Unit conversion followed by addition",
			input:       "100 cm in m + 1",
			checkResult: func(v evaluator.Value) bool {
				// 100 cm = 1 m, then 1 + 1 = 2
				return !v.IsError() && math.Abs(v.Number-2.0) < 0.01
			},
		},
		{
			description: "Unit conversion followed by subtraction",
			input:       "100 cm in m - 0.5",
			checkResult: func(v evaluator.Value) bool {
				// 100 cm = 1 m, then 1 - 0.5 = 0.5
				return !v.IsError() && math.Abs(v.Number-0.5) < 0.01
			},
		},
		{
			description: "Chained conversions with arithmetic",
			input:       "1000 m in km in m / 10",
			checkResult: func(v evaluator.Value) bool {
				// 1000 m = 1 km = 1000 m, then 1000 / 10 = 100
				return !v.IsError() && math.Abs(v.Number-100) < 0.01
			},
		},
		{
			description: "Complex expression with conversion and multiple operations",
			input:       "500 g in kg * 2 + 0.5",
			checkResult: func(v evaluator.Value) bool {
				// 500 g = 0.5 kg, then 0.5 * 2 = 1, then 1 + 0.5 = 1.5
				return !v.IsError() && math.Abs(v.Number-1.5) < 0.01
			},
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

			if !tt.checkResult(result) {
				if result.IsError() {
					t.Errorf("Test failed for %q: %s", tt.input, result.Error)
				} else {
					t.Errorf("Test failed for %q: got value %v", tt.input, result.Number)
				}
			}
		})
	}
}

// TestCompoundExpressionsWithBusinessDays tests compound expressions involving
// business days calculations with arithmetic operations
func TestCompoundExpressionsWithBusinessDays(t *testing.T) {
	tests := []struct {
		description string
		lines       []string // multiple lines to execute in sequence
		checkFinal  func(evaluator.Value) bool
	}{
		{
			description: "Date difference in business days divided by 7 (weeks)",
			lines: []string{
				"future = 03/04/2026",
				"past = 01/02/2026",
				"future - past in business days / 7",
			},
			checkFinal: func(v evaluator.Value) bool {
				// Should calculate business days between dates, then divide by 7
				// The exact value depends on the business day calculation, but it should be a reasonable number
				// We just check it's not an error and is positive
				return !v.IsError() && v.Number > 0 && v.Number < 100
			},
		},
		{
			description: "Business days calculation with multiplication",
			lines: []string{
				"bdays = 5 business days in hours",
				"bdays / 8", // work hours per day
			},
			checkFinal: func(v evaluator.Value) bool {
				// 5 business days = 5 * 8 hours = 40 hours, then 40 / 8 = 5
				return !v.IsError() && math.Abs(v.Number-5) < 0.1
			},
		},
		{
			description: "Business days conversion then arithmetic",
			lines: []string{
				"10 business days in hours / 24",
			},
			checkFinal: func(v evaluator.Value) bool {
				// 10 business days = 80 hours, then 80 / 24 = 3.33...
				return !v.IsError() && math.Abs(v.Number-3.333) < 0.1
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
					t.Logf("Evaluation error for %q: %s", line, result.Error)
				}
			}

			if !tt.checkFinal(result) {
				if result.IsError() {
					t.Errorf("Final check failed for %q: %s", tt.description, result.Error)
				} else {
					t.Errorf("Final check failed for %q: got value %v", tt.description, result.Number)
				}
			}
		})
	}
}

// TestCompoundExpressionsOperatorPrecedence tests that operator precedence is
// correctly maintained in compound expressions
func TestCompoundExpressionsOperatorPrecedence(t *testing.T) {
	tests := []struct {
		description string
		input       string
		checkResult func(evaluator.Value) bool
	}{
		{
			description: "Multiplication before addition after conversion",
			input:       "100 cm in m * 2 + 1",
			checkResult: func(v evaluator.Value) bool {
				// (100 cm in m) * 2 + 1 = 1 * 2 + 1 = 3
				return !v.IsError() && math.Abs(v.Number-3.0) < 0.01
			},
		},
		{
			description: "Division before subtraction after conversion",
			input:       "1000 g in kg / 2 - 0.25",
			checkResult: func(v evaluator.Value) bool {
				// (1000 g in kg) / 2 - 0.25 = 1 / 2 - 0.25 = 0.25
				return !v.IsError() && math.Abs(v.Number-0.25) < 0.01
			},
		},
		{
			description: "Complex precedence with parentheses",
			input:       "(100 cm in m + 1) * 2",
			checkResult: func(v evaluator.Value) bool {
				// (100 cm in m + 1) * 2 = (1 + 1) * 2 = 4
				return !v.IsError() && math.Abs(v.Number-4.0) < 0.01
			},
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

			if !tt.checkResult(result) {
				if result.IsError() {
					t.Errorf("Test failed for %q: %s", tt.input, result.Error)
				} else {
					t.Errorf("Test failed for %q: got value %v (expected check to pass)", tt.input, result.Number)
				}
			}
		})
	}
}

// TestIssueExample tests the exact example from the GitHub issue
func TestIssueExample(t *testing.T) {
	// Note: This test uses relative dates, so we can't check exact values
	// We just ensure it parses and evaluates without error
	tests := []struct {
		description   string
		input         string
		shouldSucceed bool
	}{
		{
			description:   "Date arithmetic with business days and division",
			input:         "03/04/2026 - today in business days / 7",
			shouldSucceed: true,
		},
		{
			description:   "Date arithmetic with business days (no division)",
			input:         "03/04/2026 - today in business days",
			shouldSucceed: true,
		},
		{
			description:   "Parenthesized version should also work",
			input:         "(03/04/2026 - today in business days) / 7",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()

			if tt.shouldSucceed {
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

				// Check that we got a numeric result
				if result.Type != evaluator.ValueNumber && result.Type != evaluator.ValueUnit {
					t.Errorf("Expected numeric result for %q, got type %v", tt.input, result.Type)
				}
			} else {
				if err == nil {
					t.Errorf("Expected parse error for %q, but it succeeded", tt.input)
				}
			}
		})
	}
}
