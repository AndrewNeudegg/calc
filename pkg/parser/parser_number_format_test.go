package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParserNormalizeNumber(t *testing.T) {
	p := New([]lexer.Token{})

	tests := []struct {
		input    string
		expected string
	}{
		// US format (comma as thousand separator, period as decimal)
		{"1,234", "1234"},
		{"1,234.56", "1234.56"},
		{"55,101.10", "55101.10"},
		{"31,432", "31432"},
		{"1,500,000", "1500000"},
		{"12,345,678.90", "12345678.90"},

		// European format (period as thousand separator, comma as decimal)
		{"1.234,56", "1234.56"},
		{"65.342,10", "65342.10"},
		{"12.345.678,90", "12345678.90"},

		// Edge cases with single separator
		{"100,000", "100000"},     // US thousands
		{"100.000", "100000"},     // European thousands (multiple periods)
		{"50,000.50", "50000.50"}, // US format
		{"50.000,50", "50000.50"}, // European format

		// Simple numbers (no separators)
		{"42", "42"},
		{"3.14", "3.14"},
		{"100.5", "100.5"},

		// Single separator near end
		{"1,5", "1.5"}, // European decimal
	}

	for _, tt := range tests {
		result := p.normalizeNumber(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeNumber(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParserFormattedNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		// US format
		{"1,234", 1234},
		{"1,234.56", 1234.56},
		{"55,101.10", 55101.10},
		{"31,432", 31432},
		{"1,500,000", 1500000},

		// European format
		{"1.234,56", 1234.56},
		{"65.342,10", 65342.10},

		// Edge cases
		{"100,000", 100000},
		{"100.000", 100000},
		{"50,000.50", 50000.50},
		{"50.000,50", 50000.50},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := New(tokens)

		expr, err := p.Parse()
		if err != nil {
			t.Errorf("input %q: parse error: %v", tt.input, err)
			continue
		}

		numExpr, ok := expr.(*NumberExpr)
		if !ok {
			t.Errorf("input %q: expected *NumberExpr, got %T", tt.input, expr)
			continue
		}

		if numExpr.Value != tt.expected {
			t.Errorf("input %q: expected value %f, got %f", tt.input, tt.expected, numExpr.Value)
		}
	}
}

func TestParserCurrencyWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input          string
		expectedValue  float64
		expectedCurrency string
	}{
		{"$31,432", 31432, "$"},
		{"$55,101.10", 55101.10, "$"},
		{"€65.342,10", 65342.10, "€"},
		{"£1,234,567.89", 1234567.89, "£"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := New(tokens)

		expr, err := p.Parse()
		if err != nil {
			t.Errorf("input %q: parse error: %v", tt.input, err)
			continue
		}

		currExpr, ok := expr.(*CurrencyExpr)
		if !ok {
			t.Errorf("input %q: expected *CurrencyExpr, got %T", tt.input, expr)
			continue
		}

		numExpr, ok := currExpr.Value.(*NumberExpr)
		if !ok {
			t.Errorf("input %q: expected *NumberExpr inside currency, got %T", tt.input, currExpr.Value)
			continue
		}

		if numExpr.Value != tt.expectedValue {
			t.Errorf("input %q: expected value %f, got %f", tt.input, tt.expectedValue, numExpr.Value)
		}

		if currExpr.Currency != tt.expectedCurrency {
			t.Errorf("input %q: expected currency %q, got %q", tt.input, tt.expectedCurrency, currExpr.Currency)
		}
	}
}

func TestParserExpressionsWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{"1,234 + 5,678", "addition with US format"},
		{"10,000 * 2.5", "multiplication with US format"},
		{"1.234,56 + 2.345,67", "addition with European format"},
		{"100,000 / 2", "division with US format"},
		{"50.000,50 - 25.000,25", "subtraction with European format"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := New(tokens)

		_, err := p.Parse()
		if err != nil {
			t.Errorf("%s - input %q: parse error: %v", tt.description, tt.input, err)
		}
	}
}
