package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParserNormalizeNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		locale   string
	}{
		// US format (comma as thousand separator, period as decimal)
		{"1,234", "1234", "en_US"},
		{"1,234.56", "1234.56", "en_US"},
		{"55,101.10", "55101.10", "en_US"},
		{"31,432", "31432", "en_US"},
		{"1,500,000", "1500000", "en_US"},
		{"12,345,678.90", "12345678.90", "en_US"},

		// European format (period as thousand separator, comma as decimal)
		{"1.234,56", "1234.56", "de_DE"},
		{"65.342,10", "65342.10", "de_DE"},
		{"12.345.678,90", "12345678.90", "de_DE"},

		// Edge cases with single separator - US locale
		{"100,000", "100000", "en_US"},     // US thousands
		{"100.000", "100.000", "en_US"},    // US decimal with 3 places
		{"50,000.50", "50000.50", "en_US"}, // US format
		{"1.5", "1.5", "en_US"},            // US decimal

		// Edge cases with single separator - European locale
		{"100.000", "100000", "de_DE"},   // European thousands
		{"50.000,50", "50000.50", "de_DE"}, // European format
		{"1,5", "1.5", "de_DE"},            // European decimal

		// Simple numbers (no separators) work in both locales
		{"42", "42", "en_US"},
		{"3.14", "3.14", "en_US"},
		{"100.5", "100.5", "en_US"},
		{"42", "42", "de_DE"},
	}

	for _, tt := range tests {
		p := NewWithLocale([]lexer.Token{}, tt.locale)
		result := p.normalizeNumber(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeNumber(%q) with locale %s = %q, want %q", tt.input, tt.locale, result, tt.expected)
		}
	}
}

func TestParserFormattedNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		locale   string
	}{
		// US format
		{"1,234", 1234, "en_US"},
		{"1,234.56", 1234.56, "en_US"},
		{"55,101.10", 55101.10, "en_US"},
		{"31,432", 31432, "en_US"},
		{"1,500,000", 1500000, "en_US"},
		{"100,000", 100000, "en_US"},
		{"50,000.50", 50000.50, "en_US"},

		// European format
		{"1.234,56", 1234.56, "de_DE"},
		{"65.342,10", 65342.10, "de_DE"},
		{"100.000", 100000, "de_DE"},
		{"50.000,50", 50000.50, "de_DE"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := NewWithLocale(tokens, tt.locale)

		expr, err := p.Parse()
		if err != nil {
			t.Errorf("input %q (locale %s): parse error: %v", tt.input, tt.locale, err)
			continue
		}

		numExpr, ok := expr.(*NumberExpr)
		if !ok {
			t.Errorf("input %q (locale %s): expected *NumberExpr, got %T", tt.input, tt.locale, expr)
			continue
		}

		if numExpr.Value != tt.expected {
			t.Errorf("input %q (locale %s): expected value %f, got %f", tt.input, tt.locale, tt.expected, numExpr.Value)
		}
	}
}

func TestParserCurrencyWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input            string
		expectedValue    float64
		expectedCurrency string
		locale           string
	}{
		{"$31,432", 31432, "$", "en_US"},
		{"$55,101.10", 55101.10, "$", "en_US"},
		{"€65.342,10", 65342.10, "€", "de_DE"},
		{"£1,234,567.89", 1234567.89, "£", "en_US"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := NewWithLocale(tokens, tt.locale)

		expr, err := p.Parse()
		if err != nil {
			t.Errorf("input %q (locale %s): parse error: %v", tt.input, tt.locale, err)
			continue
		}

		currExpr, ok := expr.(*CurrencyExpr)
		if !ok {
			t.Errorf("input %q (locale %s): expected *CurrencyExpr, got %T", tt.input, tt.locale, expr)
			continue
		}

		numExpr, ok := currExpr.Value.(*NumberExpr)
		if !ok {
			t.Errorf("input %q (locale %s): expected *NumberExpr inside currency, got %T", tt.input, tt.locale, currExpr.Value)
			continue
		}

		if numExpr.Value != tt.expectedValue {
			t.Errorf("input %q (locale %s): expected value %f, got %f", tt.input, tt.locale, tt.expectedValue, numExpr.Value)
		}

		if currExpr.Currency != tt.expectedCurrency {
			t.Errorf("input %q (locale %s): expected currency %q, got %q", tt.input, tt.locale, tt.expectedCurrency, currExpr.Currency)
		}
	}
}

func TestParserExpressionsWithFormattedNumbers(t *testing.T) {
	tests := []struct {
		input       string
		description string
		locale      string
	}{
		{"1,234 + 5,678", "addition with US format", "en_US"},
		{"10,000 * 2.5", "multiplication with US format", "en_US"},
		{"1.234,56 + 2.345,67", "addition with European format", "de_DE"},
		{"100,000 / 2", "division with US format", "en_US"},
		{"50.000,50 - 25.000,25", "subtraction with European format", "de_DE"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.AllTokens()
		p := NewWithLocale(tokens, tt.locale)

		_, err := p.Parse()
		if err != nil {
			t.Errorf("%s - input %q (locale %s): parse error: %v", tt.description, tt.input, tt.locale, err)
		}
	}
}
