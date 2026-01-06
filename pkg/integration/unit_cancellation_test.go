package integration

import (
	"fmt"
	"math"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestNumberSlashUnit tests parsing of expressions like "730/month"
func TestNumberSlashUnit(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expectType  evaluator.ValueType
		expectUnit  string
		expectNum   float64
	}{
		{
			description: "730/month",
			input:       "730/month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/month",
			expectNum:   730,
		},
		{
			description: "50/day",
			input:       "50/day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/day",
			expectNum:   50,
		},
		{
			description: "100/hour",
			input:       "100/hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/hour",
			expectNum:   100,
		},
		{
			description: "1000/year",
			input:       "1000/year",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/year",
			expectNum:   1000,
		},
		{
			description: "24/hr",
			input:       "24/hr",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/hr",
			expectNum:   24,
		},
		{
			description: "60/min",
			input:       "60/min",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/min",
			expectNum:   60,
		},
		{
			description: "3600/s",
			input:       "3600/s",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/s",
			expectNum:   3600,
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

			if result.Unit != tt.expectUnit {
				t.Errorf("Expected unit %q for %q, got %q", tt.expectUnit, tt.input, result.Unit)
			}

			if result.Number != tt.expectNum {
				t.Errorf("Expected number %v for %q, got %v", tt.expectNum, tt.input, result.Number)
			}
		})
	}
}

// TestNumberPerUnit tests parsing of expressions like "730 per month"
func TestNumberPerUnit(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expectType  evaluator.ValueType
		expectUnit  string
		expectNum   float64
	}{
		{
			description: "730 per month",
			input:       "730 per month",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/month",
			expectNum:   730,
		},
		{
			description: "50 per day",
			input:       "50 per day",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/day",
			expectNum:   50,
		},
		{
			description: "100 per hour",
			input:       "100 per hour",
			expectType:  evaluator.ValueUnit,
			expectUnit:  "1/hour",
			expectNum:   100,
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

			if result.Unit != tt.expectUnit {
				t.Errorf("Expected unit %q for %q, got %q", tt.expectUnit, tt.input, result.Unit)
			}

			if result.Number != tt.expectNum {
				t.Errorf("Expected number %v for %q, got %v", tt.expectNum, tt.input, result.Number)
			}
		})
	}
}

// TestUnitCancellation tests unit cancellation in multiplication
func TestUnitCancellation(t *testing.T) {
	tests := []struct {
		description string
		lines       []string
		checkFinal  func(evaluator.Value) error
	}{
		{
			description: "$/hr * hours -> $",
			lines: []string{
				"rate = $0.192/hr",
				"hours = 730 hours",
				"rate * hours",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "$" {
					return errorf("Expected unit $, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 140.16, 0.01) {
					return errorf("Expected ~140.16, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "$/hour * hours -> $",
			lines: []string{
				"rate = $0.192/hour",
				"hours = 730 hours",
				"rate * hours",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "$" {
					return errorf("Expected unit $, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 140.16, 0.01) {
					return errorf("Expected ~140.16, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "km/h * h -> km",
			lines: []string{
				"speed = 100 km/h",
				"time = 2 h",
				"speed * time",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "km" {
					return errorf("Expected unit km, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 200, 0.01) {
					return errorf("Expected 200, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "m/s * s -> m",
			lines: []string{
				"velocity = 10 m/s",
				"time = 5 s",
				"velocity * time",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "m" {
					return errorf("Expected unit m, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 50, 0.01) {
					return errorf("Expected 50, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "$/hr * hours/month * number -> $/month",
			lines: []string{
				"instance_price = $0.192/hr",
				"hours_in_month = 730 hours per month",
				"num_instances = 20",
				"instance_price * hours_in_month * num_instances",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "$/month" {
					return errorf("Expected unit $/month, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 2803.2, 0.1) {
					return errorf("Expected ~2803.2, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "Multiple unit cancellations",
			lines: []string{
				"rate = $50/hour",
				"hours = 8 hours",
				"days = 5 days",
				"rate * hours",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "$" {
					return errorf("Expected unit $, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 400, 0.01) {
					return errorf("Expected 400, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "£ per hour * hours -> £",
			lines: []string{
				"rate = £25/hour",
				"hours = 40 hours",
				"rate * hours",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "£" {
					return errorf("Expected unit £, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 1000, 0.01) {
					return errorf("Expected 1000, got %v", v.Number)
				}
				return nil
			},
		},
		{
			description: "€ per day * days -> €",
			lines: []string{
				"rate = €100/day",
				"days = 30 days",
				"rate * days",
			},
			checkFinal: func(v evaluator.Value) error {
				if v.Type != evaluator.ValueUnit {
					return errorf("Expected type ValueUnit, got %v", v.Type)
				}
				if v.Unit != "€" {
					return errorf("Expected unit €, got %q", v.Unit)
				}
				if !approxEqual(v.Number, 3000, 0.01) {
					return errorf("Expected 3000, got %v", v.Number)
				}
				return nil
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

			if err := tt.checkFinal(result); err != nil {
				t.Errorf("Final value check failed for %q: %v", tt.description, err)
			}
		})
	}
}

// TestCompleteCancellation tests cases where all units cancel, reducing to dimensionless
func TestCompleteCancellation(t *testing.T) {
	tests := []struct {
		description string
		lines       []string
		expectNum   float64
	}{
		{
			description: "$ * hours / $ / hours -> dimensionless",
			lines: []string{
				"a = $10/hour",
				"b = 5 hours",
				"c = $50",
				"result = a * b / c",
			},
			expectNum: 1.0,
		},
		{
			description: "currency rate * time / currency -> dimensionless",
			lines: []string{
				"rate = $20/hour",
				"time = 3 hours",
				"cost = $60",
				"result = rate * time / cost",
			},
			expectNum: 1.0,
		},
		{
			description: "km/hour * hours / km -> dimensionless",
			lines: []string{
				"speed = 100 km/hour",
				"time = 2 hours",
				"dist = 200 km",
				"result = speed * time / dist",
			},
			expectNum: 1.0,
		},
		{
			description: "Multiple same units cancel completely",
			lines: []string{
				"val1 = 10 m",
				"val2 = 5 m",
				"val3 = 2 m",
				"val4 = 5 m",
				"result = val1 * val2 / val3 / val4",
			},
			expectNum: 5.0,
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

			// Check that result is dimensionless (plain number)
			if result.Type != evaluator.ValueNumber {
				t.Errorf("Expected dimensionless result (ValueNumber), got %v with unit %q", result.Type, result.Unit)
			}

			if !approxEqual(result.Number, tt.expectNum, 0.01) {
				t.Errorf("Expected %v, got %v", tt.expectNum, result.Number)
			}
		})
	}
}

// TestMultipleSlashBehavior tests behavior with multiple division operations
// This documents the current behavior where a/b/c is parsed as (a/b)/c
func TestMultipleSlashBehavior(t *testing.T) {
	tests := []struct {
		description string
		lines       []string
		expectType  evaluator.ValueType
		expectUnit  string
		expectNum   float64
	}{
		{
			description: "10 m / 2 m / 5 -> dimensionless 1",
			lines: []string{
				"result = 10 m / 2 m / 5",
			},
			expectType: evaluator.ValueNumber,
			expectUnit: "",
			expectNum:  1.0,
		},
		{
			description: "100 km / 2 hours / 5 -> 10 km/hours",
			lines: []string{
				"result = 100 km / 2 hours / 5",
			},
			expectType: evaluator.ValueUnit,
			expectUnit: "km/hours",
			expectNum:  10.0,
		},
		{
			description: "100 m / 10 m / 2 m -> 5 1/m (current behavior, not fully canceled)",
			lines: []string{
				"result = 100 m / 10 m / 2 m",
			},
			expectType: evaluator.ValueUnit,
			expectUnit: "1/m",
			expectNum:  5.0,
		},
		{
			description: "Grouped divisions for complete cancellation: (100 m / 10 m) / 2 m",
			lines: []string{
				"temp = 100 m / 10 m",
				"result = temp / 2 m",
			},
			expectType: evaluator.ValueUnit,
			expectUnit: "1/m",
			expectNum:  5.0,
		},
		{
			description: "Proper complete cancellation with explicit grouping",
			lines: []string{
				"val1 = 100 m",
				"val2 = 10 m",
				"val3 = 10 m",
				"result = val1 / val2 / val3",
			},
			expectType: evaluator.ValueUnit,
			expectUnit: "1/m",
			expectNum:  1.0,
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

			if result.Type != tt.expectType {
				t.Errorf("Expected type %v, got %v", tt.expectType, result.Type)
			}

			if tt.expectType == evaluator.ValueUnit && result.Unit != tt.expectUnit {
				t.Errorf("Expected unit %q, got %q", tt.expectUnit, result.Unit)
			}

			if !approxEqual(result.Number, tt.expectNum, 0.01) {
				t.Errorf("Expected %v, got %v", tt.expectNum, result.Number)
			}
		})
	}
}

// TestUnitCancellationVariations tests all variations of time units
func TestUnitCancellationVariations(t *testing.T) {
	tests := []struct {
		description string
		rateUnit    string // denominator in rate (e.g., "hr" in "$/hr")
		timeUnit    string // unit for time value (e.g., "hours")
		shouldCancel bool
	}{
		// hr variations
		{"hr and hours", "hr", "hours", true},
		{"hr and hour", "hr", "hour", true},
		{"hr and h", "hr", "h", true},
		{"hr and hr", "hr", "hr", true},
		
		// hour variations
		{"hour and hours", "hour", "hours", true},
		{"hour and hour", "hour", "hour", true},
		{"hour and h", "hour", "h", true},
		{"hour and hr", "hour", "hr", true},
		
		// h variations
		{"h and hours", "h", "hours", true},
		{"h and hour", "h", "hour", true},
		{"h and h", "h", "h", true},
		{"h and hr", "h", "hr", true},
		
		// hours variations
		{"hours and hours", "hours", "hours", true},
		{"hours and hour", "hours", "hour", true},
		{"hours and h", "hours", "h", true},
		{"hours and hr", "hours", "hr", true},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)

			// Parse and evaluate rate (e.g., "$10/hr")
			rateLine := "$10/" + tt.rateUnit
			l1 := lexer.New(rateLine)
			tokens1 := l1.AllTokens()
			p1 := parser.New(tokens1)
			expr1, err := p1.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", rateLine, err)
			}
			rateVal := e.Eval(expr1)
			if rateVal.IsError() {
				t.Fatalf("Evaluation error for %q: %s", rateLine, rateVal.Error)
			}

			// Parse and evaluate time (e.g., "5 hours")
			timeLine := "5 " + tt.timeUnit
			l2 := lexer.New(timeLine)
			tokens2 := l2.AllTokens()
			p2 := parser.New(tokens2)
			expr2, err := p2.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", timeLine, err)
			}
			timeVal := e.Eval(expr2)
			if timeVal.IsError() {
				t.Fatalf("Evaluation error for %q: %s", timeLine, timeVal.Error)
			}

			// Multiply rate * time
			multiplyLine := rateLine + " * " + timeLine
			l3 := lexer.New(multiplyLine)
			tokens3 := l3.AllTokens()
			p3 := parser.New(tokens3)
			expr3, err := p3.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", multiplyLine, err)
			}
			result := e.Eval(expr3)
			if result.IsError() {
				t.Fatalf("Evaluation error for %q: %s", multiplyLine, result.Error)
			}

			if tt.shouldCancel {
				// Units should cancel, leaving just currency
				if result.Unit != "$" {
					t.Errorf("Expected unit $ (cancelled), got %q for %s * %s", result.Unit, tt.rateUnit, tt.timeUnit)
				}
				if !approxEqual(result.Number, 50, 0.01) {
					t.Errorf("Expected 50, got %v", result.Number)
				}
			}
		})
	}
}

// Helper functions
func approxEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
