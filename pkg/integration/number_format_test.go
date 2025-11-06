package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/formatter"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/settings"
)

func TestNumberFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// US format (comma as thousand separator, period as decimal)
		{
			name:     "US format - simple thousands",
			input:    "31,432",
			expected: "31,432.00",
		},
		{
			name:     "US format - with decimals",
			input:    "55,101.10",
			expected: "55,101.10",
		},
		{
			name:     "US format - millions",
			input:    "1,500,000",
			expected: "1,500,000.00",
		},
		{
			name:     "US format - large number with decimals",
			input:    "12,345,678.90",
			expected: "12,345,678.90",
		},
		
		// European format (period as thousand separator, comma as decimal)
		{
			name:     "European format - with comma decimal",
			input:    "65.342,10",
			expected: "65,342.10",
		},
		{
			name:     "European format - multiple periods",
			input:    "12.345.678,90",
			expected: "12,345,678.90",
		},
		
		// Currency with formatted numbers
		{
			name:     "USD with thousands",
			input:    "$31,432",
			expected: "$31,432.00",
		},
		{
			name:     "USD with thousands and decimals",
			input:    "$55,101.10",
			expected: "$55,101.10",
		},
		{
			name:     "EUR with European format",
			input:    "€65.342,10",
			expected: "€65,342.10",
		},
		
		// Arithmetic with formatted numbers
		{
			name:     "Addition with formatted numbers",
			input:    "1,234 + 5,678",
			expected: "6,912.00",
		},
		{
			name:     "Multiplication with formatted numbers",
			input:    "10,000 * 2.5",
			expected: "25,000.00",
		},
		{
			name:     "Division with formatted numbers",
			input:    "100,000 / 4",
			expected: "25,000.00",
		},
	}

	s := settings.Default()
	env := evaluator.NewEnvironment()
	e := evaluator.New(env)
	f := formatter.New(s)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			val := e.Eval(expr)
			if val.IsError() {
				t.Fatalf("Eval error: %s", val.Error)
			}

			result := f.Format(val)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFuzzyPhrasesWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "half of formatted number",
			input:    "half of $31,432",
			expected: "$15,716.00",
		},
		{
			name:     "double formatted number",
			input:    "double 10,000",
			expected: "20,000.00",
		},
		{
			name:     "percentage of formatted number",
			input:    "20% of 50,000",
			expected: "10,000.00",
		},
		{
			name:     "three quarters of formatted number",
			input:    "three quarters of 100,000",
			expected: "75,000.00",
		},
	}

	s := settings.Default()
	env := evaluator.NewEnvironment()
	e := evaluator.New(env)
	f := formatter.New(s)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			val := e.Eval(expr)
			if val.IsError() {
				t.Fatalf("Eval error: %s", val.Error)
			}

			result := f.Format(val)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestUnitConversionsWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		// We just check that it parses and evaluates without error
	}{
		{
			name:  "formatted number with unit conversion",
			input: "10,000 kg in lb",
		},
		{
			name:  "formatted number with distance conversion",
			input: "1,234.56 km in miles",
		},
		{
			name:  "formatted currency conversion",
			input: "55,101.10 dollars in euros",
		},
	}

	env := evaluator.NewEnvironment()
	e := evaluator.New(env)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			val := e.Eval(expr)
			if val.IsError() {
				t.Fatalf("Eval error: %s", val.Error)
			}
		})
	}
}

func TestVariablesWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		checks []struct {
			input    string
			expected string
		}
	}{
		{
			name: "assign and use formatted number",
			inputs: []string{
				"price = 1,234.56",
			},
			checks: []struct {
				input    string
				expected string
			}{
				{"price", "1,234.56"},
				{"price * 2", "2,469.12"},
			},
		},
		{
			name: "assign and use currency with formatted number",
			inputs: []string{
				"budget = $55,101.10",
			},
			checks: []struct {
				input    string
				expected string
			}{
				{"budget", "$55,101.10"},
				{"half of budget", "$27,550.55"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings.Default()
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			f := formatter.New(s)

			// Process input assignments
			for _, input := range tt.inputs {
				l := lexer.New(input)
				tokens := l.AllTokens()
				p := parser.New(tokens)
				
				expr, err := p.Parse()
				if err != nil {
					t.Fatalf("Parse error for %q: %v", input, err)
				}

				val := e.Eval(expr)
				if val.IsError() {
					t.Fatalf("Eval error for %q: %s", input, val.Error)
				}
			}

			// Check results
			for _, check := range tt.checks {
				l := lexer.New(check.input)
				tokens := l.AllTokens()
				p := parser.New(tokens)
				
				expr, err := p.Parse()
				if err != nil {
					t.Fatalf("Parse error for %q: %v", check.input, err)
				}

				val := e.Eval(expr)
				if val.IsError() {
					t.Fatalf("Eval error for %q: %s", check.input, val.Error)
				}

				result := f.Format(val)
				if result != check.expected {
					t.Errorf("For %q: expected %q, got %q", check.input, check.expected, result)
				}
			}
		})
	}
}
