package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/formatter"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/settings"
)

// TestLocaleIntegrationUSFormat tests complete workflow with US format
func TestLocaleIntegrationUSFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Issue examples - these were incorrectly parsed before the fix
		{"Issue example 1: $2.115 / 2", "$2.115 / 2", "$1.06"},
		{"Issue example 2: $2.11 / 2", "$2.11 / 2", "$1.06"},

		// Additional US format tests
		{"Decimal with 3 places", "3.142", "3.14"},
		{"Decimal with 4 places", "1.2345", "1.23"},
		{"Currency with 3 decimals", "$5.678", "$5.68"},
		{"Arithmetic with decimals", "10.5 + 20.3", "30.80"},
		{"Division with decimals", "7.5 / 2.5", "3.00"},
		{"Thousands separator", "1,234 + 5,678", "6,912.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings.Default() // Default is en_US
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			f := formatter.New(s)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.NewWithLocale(tokens, s.Locale)

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

// TestLocaleIntegrationEuropeanFormat tests complete workflow with European format
func TestLocaleIntegrationEuropeanFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // Output is always in US format (with commas)
	}{
		{"Decimal with comma", "3,142", "3.14"},
		{"Currency with comma", "€5,678", "€5.68"},
		{"Thousands with comma decimal", "1.234,56", "1,234.56"},
		{"Arithmetic with European format", "10,5 + 20,3", "30.80"},
		{"Division with comma decimal", "7,5 / 2,5", "3.00"},
		{"Large number", "1.234.567,89", "1,234,567.89"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings.Default()
			s.Locale = "de_DE" // European locale for parsing
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			
			// Use US locale for formatting (output)
			sFormat := settings.Default()
			sFormat.Locale = "en_US"
			f := formatter.New(sFormat)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.NewWithLocale(tokens, "de_DE")

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

// TestLocaleWithVariables tests locale handling with variable assignments
func TestLocaleWithVariables(t *testing.T) {
	tests := []struct {
		name       string
		locale     string
		statements []string    // Sequence of statements to execute
		checks     []struct {  // Values to check after executing statements
			varName  string
			expected float64
		}
	}{
		{
			name:   "US format with variables",
			locale: "en_US",
			statements: []string{
				"x = 2.115",
				"y = x / 2",
			},
			checks: []struct {
				varName  string
				expected float64
			}{
				{"x", 2.115},
				{"y", 1.0575},
			},
		},
		{
			name:   "European format with variables",
			locale: "de_DE",
			statements: []string{
				"x = 2,115",
				"y = x / 2",
			},
			checks: []struct {
				varName  string
				expected float64
			}{
				{"x", 2.115},
				{"y", 1.0575},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)

			// Execute statements
			for _, stmt := range tt.statements {
				l := lexer.New(stmt)
				tokens := l.AllTokens()
				p := parser.NewWithLocale(tokens, tt.locale)

				expr, err := p.Parse()
				if err != nil {
					t.Fatalf("Parse error for %q: %v", stmt, err)
				}

				val := e.Eval(expr)
				if val.IsError() {
					t.Fatalf("Eval error for %q: %s", stmt, val.Error)
				}
			}

			// Check variable values
			for _, check := range tt.checks {
				val, exists := e.GetVariable(check.varName)
				if !exists {
					t.Fatalf("Variable %q not found", check.varName)
				}

				if val.Number != check.expected {
					t.Errorf("Variable %q: expected %f, got %f", check.varName, check.expected, val.Number)
				}
			}
		})
	}
}

// TestLocaleEdgeCases tests edge cases for locale handling
func TestLocaleEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		input    string
		expected string
	}{
		// Very small decimals
		{"US: Small decimal", "en_US", "0.001", "0.00"},
		{"EU: Small decimal", "de_DE", "0,001", "0.00"},

		// Zero
		{"US: Zero", "en_US", "0.00", "0.00"},
		{"EU: Zero", "de_DE", "0,00", "0.00"},

		// Negative numbers
		{"US: Negative", "en_US", "-2.115", "-2.12"},
		{"EU: Negative", "de_DE", "-2,115", "-2.12"},

		// Large numbers
		{"US: Million", "en_US", "1,000,000.50", "1,000,000.50"},
		{"EU: Million", "de_DE", "1.000.000,50", "1,000,000.50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings.Default()
			s.Locale = "en_US" // Output formatting
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)
			f := formatter.New(s)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.NewWithLocale(tokens, tt.locale)

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

// TestLocaleWithUnits tests locale handling with unit conversions
func TestLocaleWithUnits(t *testing.T) {
	tests := []struct {
		name   string
		locale string
		input  string
	}{
		{"US: Decimal with units", "en_US", "2.115 kg"},
		{"EU: Comma decimal with units", "de_DE", "2,115 kg"},
		{"US: Thousands with units", "en_US", "1,234.56 m"},
		{"EU: Thousands with units", "de_DE", "1.234,56 m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			e := evaluator.New(env)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.NewWithLocale(tokens, tt.locale)

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			val := e.Eval(expr)
			if val.IsError() {
				t.Fatalf("Eval error: %s", val.Error)
			}

			// Just check that it evaluates without error
			if val.Type != evaluator.ValueUnit {
				t.Errorf("Expected ValueUnit, got %v", val.Type)
			}
		})
	}
}
