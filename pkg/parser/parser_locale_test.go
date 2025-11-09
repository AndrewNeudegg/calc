package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

// TestLocaleParsingUSFormat tests US format parsing with en_US locale
func TestLocaleParsingUSFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Simple decimal", "2.115", 2.115},
		{"Two decimal places", "2.11", 2.11},
		{"Three decimal places", "3.142", 3.142},
		{"Four decimal places", "1.2345", 1.2345},
		{"Thousands with decimal", "1,234.56", 1234.56},
		{"Millions", "1,000,000.50", 1000000.50},
		{"Small decimal", "0.001", 0.001},
		{"Large number", "123,456.789", 123456.789},
		{"No decimal part", "1,234", 1234},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := NewWithLocale(tokens, "en_US")

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			numExpr, ok := expr.(*NumberExpr)
			if !ok {
				t.Fatalf("Expected *NumberExpr, got %T", expr)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, numExpr.Value)
			}
		})
	}
}

// TestLocaleParsingEuropeanFormat tests European format parsing with de_DE locale
func TestLocaleParsingEuropeanFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Simple decimal", "2,115", 2.115},
		{"Two decimal places", "2,11", 2.11},
		{"Three decimal places", "3,142", 3.142},
		{"Four decimal places", "1,2345", 1.2345},
		{"Thousands with decimal", "1.234,56", 1234.56},
		{"Millions", "1.000.000,50", 1000000.50},
		{"Small decimal", "0,001", 0.001},
		{"Large number", "123.456,789", 123456.789},
		{"No decimal part", "1.234", 1234},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := NewWithLocale(tokens, "de_DE")

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			numExpr, ok := expr.(*NumberExpr)
			if !ok {
				t.Fatalf("Expected *NumberExpr, got %T", expr)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, numExpr.Value)
			}
		})
	}
}

// TestLocaleAvoidsMisidentification ensures numbers with many decimal places
// are not misidentified as thousand separators
func TestLocaleAvoidsMisidentification(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		input    string
		expected float64
		desc     string
	}{
		{
			name:     "US: 2.115 is two point one one five",
			locale:   "en_US",
			input:    "2.115",
			expected: 2.115,
			desc:     "Not two thousand one hundred fifteen",
		},
		{
			name:     "US: 2.11 is two point one one",
			locale:   "en_US",
			input:    "2.11",
			expected: 2.11,
			desc:     "Not two thousand eleven",
		},
		{
			name:     "US: 3.14159 is pi",
			locale:   "en_US",
			input:    "3.14159",
			expected: 3.14159,
			desc:     "Not three thousand one hundred forty-one point five nine",
		},
		{
			name:     "EU: 2,115 is two point one one five",
			locale:   "de_DE",
			input:    "2,115",
			expected: 2.115,
			desc:     "Not two thousand one hundred fifteen",
		},
		{
			name:     "EU: 2,11 is two point one one",
			locale:   "de_DE",
			input:    "2,11",
			expected: 2.11,
			desc:     "Not two thousand eleven",
		},
		{
			name:     "EU: 3,14159 is pi",
			locale:   "de_DE",
			input:    "3,14159",
			expected: 3.14159,
			desc:     "Not three thousand one hundred forty-one point five nine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := NewWithLocale(tokens, tt.locale)

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			numExpr, ok := expr.(*NumberExpr)
			if !ok {
				t.Fatalf("Expected *NumberExpr, got %T", expr)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("%s: Expected %f, got %f", tt.desc, tt.expected, numExpr.Value)
			}
		})
	}
}

// TestLocaleCurrencyParsing tests currency parsing with different locales
func TestLocaleCurrencyParsing(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		input    string
		expected float64
		currency string
	}{
		// US locale with USD
		{"US: $2.115", "en_US", "$2.115", 2.115, "$"},
		{"US: $2.11", "en_US", "$2.11", 2.11, "$"},
		{"US: $1,234.56", "en_US", "$1,234.56", 1234.56, "$"},

		// European locale with EUR
		{"EU: €2,115", "de_DE", "€2,115", 2.115, "€"},
		{"EU: €2,11", "de_DE", "€2,11", 2.11, "€"},
		{"EU: €1.234,56", "de_DE", "€1.234,56", 1234.56, "€"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := NewWithLocale(tokens, tt.locale)

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			currExpr, ok := expr.(*CurrencyExpr)
			if !ok {
				t.Fatalf("Expected *CurrencyExpr, got %T", expr)
			}

			numExpr, ok := currExpr.Value.(*NumberExpr)
			if !ok {
				t.Fatalf("Expected *NumberExpr inside currency, got %T", currExpr.Value)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("Expected value %f, got %f", tt.expected, numExpr.Value)
			}

			if currExpr.Currency != tt.currency {
				t.Errorf("Expected currency %s, got %s", tt.currency, currExpr.Currency)
			}
		})
	}
}

// TestLocaleArithmeticOperations tests arithmetic with locale-aware parsing
func TestLocaleArithmeticOperations(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		input    string
		expected float64
	}{
		// US locale
		{"US: 2.115 / 2", "en_US", "2.115 / 2", 1.0575},
		{"US: 2.11 / 2", "en_US", "2.11 / 2", 1.055},
		{"US: 1,234.56 + 5,678.90", "en_US", "1,234.56 + 5,678.90", 6913.46},

		// European locale
		{"EU: 2,115 / 2", "de_DE", "2,115 / 2", 1.0575},
		{"EU: 2,11 / 2", "de_DE", "2,11 / 2", 1.055},
		{"EU: 1.234,56 + 5.678,90", "de_DE", "1.234,56 + 5.678,90", 6913.46},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := NewWithLocale(tokens, tt.locale)

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			// We can't evaluate here without the evaluator, but we can check it parses
			if expr == nil {
				t.Fatal("Expected expression, got nil")
			}
		})
	}
}

// TestDefaultLocale tests that New() uses en_US as default
func TestDefaultLocale(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"2.115", 2.115},
		{"1,234.56", 1234.56},
		{"3.14159", 3.14159},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := New(tokens) // Uses default en_GB locale

			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			numExpr, ok := expr.(*NumberExpr)
			if !ok {
				t.Fatalf("Expected *NumberExpr, got %T", expr)
			}

			if numExpr.Value != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, numExpr.Value)
			}
		})
	}
}
