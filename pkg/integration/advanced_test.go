package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestAdvancedExpressions tests complex real-world calculation scenarios
func TestAdvancedExpressions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectType evaluator.ValueType
		checkValue func(evaluator.Value) bool
	}{
		{
			name:       "Compound rate with currency conversion (dollars/day to gbp/month)",
			input:      "32 dollars per day in gbp per month",
			expectType: evaluator.ValueCurrency,
			checkValue: func(v evaluator.Value) bool {
				// 32 * 30.417 (avg days per month) = 973.33 USD
				// 973.33 USD ≈ £764 (at ~1.27 exchange rate)
				// Allow range for exchange rate fluctuation
				return v.Currency == "£" && v.Number > 700 && v.Number < 850
			},
		},
		{
			name:       "Volume flow rate conversion (liters/min to m³/hour)",
			input:      "500 liters per minute in m3 per hour",
			expectType: evaluator.ValueUnit,
			checkValue: func(v evaluator.Value) bool {
				// 500 L/min = 0.5 m³/min = 30 m³/hour (allow tiny float tolerance)
				return v.Number > 29.999999 && v.Number < 30.000001 && v.Unit == "m3/hour"
			},
		},
		{
			name:       "Currency conversion with addition",
			input:      "250 dollars in eur + 100 euros",
			expectType: evaluator.ValueCurrency,
			checkValue: func(v evaluator.Value) bool {
				// 250 USD ≈ €227.27, then + 100 = €327.27
				// Allow range for exchange rate fluctuation
				return v.Currency == "€" && v.Number > 300 && v.Number < 350
			},
		},
		{
			name:       "Percentage increase with currency conversion",
			input:      "increase 320 usd by 12% in gbp",
			expectType: evaluator.ValueCurrency,
			checkValue: func(v evaluator.Value) bool {
				// 320 * 1.12 = 358.40 USD
				// 358.40 USD ≈ £281 (at ~1.27 exchange rate)
				// Allow range for exchange rate fluctuation
				return v.Currency == "£" && v.Number > 250 && v.Number < 320
			},
		},
		{
			name:       "Rate calculation with unit conversion",
			input:      "100 km / 2 hours in mph",
			expectType: evaluator.ValueUnit,
			checkValue: func(v evaluator.Value) bool {
				// 100 km / 2 hours = 50 km/h
				// 50 km/h ≈ 31.07 mph
				return v.Number > 31 && v.Number < 32 && v.Unit == "mph"
			},
		},
		{
			name:       "Fuzzy phrase with currency conversion",
			input:      "half of 960 dollars in gbp",
			expectType: evaluator.ValueCurrency,
			checkValue: func(v evaluator.Value) bool {
				// half of 960 = 480 USD
				// 480 USD ≈ £377.95 (at ~1.27 exchange rate)
				// Allow range for exchange rate fluctuation
				return v.Currency == "£" && v.Number > 350 && v.Number < 420
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error for input %q: %v", tt.input, err)
			}

			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			result := e.Eval(expr)

			if result.IsError() {
				t.Fatalf("Evaluation error for input %q: %s", tt.input, result.Error)
			}

			if result.Type != tt.expectType {
				t.Errorf("Input: %q\nExpected type %v, got %v", tt.input, tt.expectType, result.Type)
			}

			if !tt.checkValue(result) {
				t.Errorf("Input: %q\nValue check failed. Got: %+v", tt.input, result)
			}
		})
	}
}
