package formatter

import (
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/settings"
)

// TestLocalisationNumberFormatting tests number formatting across multiple locales
func TestLocalisationNumberFormatting(t *testing.T) {
	tests := []struct {
		locale    string
		number    float64
		precision int
		expected  string
	}{
		// British English (commas for thousands, period for decimal)
		{"en_GB", 1234.56, 2, "1,234.56"},
		{"en_GB", 1000000, 2, "1,000,000.00"},
		{"en_GB", 42.123456, 4, "42.1235"},

		// US English (same as UK for numbers - both use comma thousands separator)
		{"en_US", 1234.56, 2, "1,234.56"},
		{"en_US", 1000000, 2, "1,000,000.00"},

		// French (space for thousands, comma for decimal - not yet implemented)
		// {"fr_FR", 1234.56, 2, "1 234,56"},

		// German (period for thousands, comma for decimal - not yet implemented)
		// {"de_DE", 1234.56, 2, "1.234,56"},

		// Precision variations
		{"en_GB", 3.14159265, 0, "3"},
		{"en_GB", 3.14159265, 2, "3.14"},
		{"en_GB", 3.14159265, 5, "3.14159"},
	}

	for _, tt := range tests {
		s := settings.Default()
		s.Locale = tt.locale
		s.Precision = tt.precision
		f := New(s)

		val := evaluator.Value{Type: evaluator.ValueNumber, Number: tt.number}
		result := f.Format(val)

		if result != tt.expected {
			t.Errorf("Locale %s: Format(%f, precision=%d) = %q, want %q",
				tt.locale, tt.number, tt.precision, result, tt.expected)
		}
	}
}

// TestLocalisationCurrencyFormatting tests currency formatting across locales
func TestLocalisationCurrencyFormatting(t *testing.T) {
	tests := []struct {
		locale   string
		amount   float64
		currency string
		expected string
	}{
		// British English - symbol before, comma separators
		{"en_GB", 1234.56, "£", "£1,234.56"},
		{"en_GB", 100, "£", "£100.00"},

		// US English - symbol before, comma separators (standard US format)
		{"en_US", 1234.56, "$", "$1,234.56"},
		{"en_US", 100, "$", "$100.00"},

		// Euro formatting (not yet fully implemented for different locales)
		{"en_GB", 1234.56, "€", "€1,234.56"},
	}

	for _, tt := range tests {
		s := settings.Default()
		s.Locale = tt.locale
		s.Precision = 2
		f := New(s)

		val := evaluator.Value{
			Type:     evaluator.ValueCurrency,
			Number:   tt.amount,
			Currency: tt.currency,
		}
		result := f.Format(val)

		if result != tt.expected {
			t.Errorf("Locale %s: Format(currency %f %s) = %q, want %q",
				tt.locale, tt.amount, tt.currency, result, tt.expected)
		}
	}
}

// TestLocalisationDateFormatting tests date formatting across locales
func TestLocalisationDateFormatting(t *testing.T) {
	testDate := time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		locale     string
		dateFormat string
		expected   string
	}{
		// British English - DD/MM/YYYY
		{"en_GB", "02/01/2006", "15/11/2025"},
		{"en_GB", "2 Jan 2006", "15 Nov 2025"},

		// US English - MM/DD/YYYY
		{"en_US", "01/02/2006", "11/15/2025"},
		{"en_US", "Jan 2, 2006", "Nov 15, 2025"},

		// ISO format (universal)
		{"en_GB", "2006-01-02", "2025-11-15"},
		{"en_US", "2006-01-02", "2025-11-15"},
	}

	for _, tt := range tests {
		s := settings.Default()
		s.Locale = tt.locale
		s.DateFormat = tt.dateFormat
		f := New(s)

		val := evaluator.Value{Type: evaluator.ValueDate, Date: testDate}
		result := f.Format(val)

		if result != tt.expected {
			t.Errorf("Locale %s with format %q: Format(date) = %q, want %q",
				tt.locale, tt.dateFormat, result, tt.expected)
		}
	}
}

// TestLocalisationPercentFormatting tests percent formatting
func TestLocalisationPercentFormatting(t *testing.T) {
	tests := []struct {
		locale    string
		percent   float64
		precision int
		expected  string
	}{
		{"en_GB", 50, 0, "50%"},
		{"en_GB", 33.333, 2, "33.33%"},
		{"en_US", 50, 0, "50%"},
		{"en_US", 33.333, 2, "33.33%"},
	}

	for _, tt := range tests {
		s := settings.Default()
		s.Locale = tt.locale
		s.Precision = tt.precision
		f := New(s)

		val := evaluator.Value{Type: evaluator.ValuePercent, Number: tt.percent}
		result := f.Format(val)

		if result != tt.expected {
			t.Errorf("Locale %s: Format(percent %f) = %q, want %q",
				tt.locale, tt.percent, result, tt.expected)
		}
	}
}

// TestLocalisationUnitFormatting tests unit formatting with different locales
func TestLocalisationUnitFormatting(t *testing.T) {
	tests := []struct {
		locale    string
		value     float64
		unit      string
		precision int
		expected  string
	}{
		// British English
		{"en_GB", 1234.56, "m", 2, "1,234.56 m"},
		{"en_GB", 100, "kg", 2, "100.00 kg"},

		// US English (also uses comma thousands separator)
		{"en_US", 1234.56, "m", 2, "1,234.56 m"},
		{"en_US", 100, "kg", 2, "100.00 kg"},
	}

	for _, tt := range tests {
		s := settings.Default()
		s.Locale = tt.locale
		s.Precision = tt.precision
		f := New(s)

		val := evaluator.Value{
			Type:   evaluator.ValueUnit,
			Number: tt.value,
			Unit:   tt.unit,
		}
		result := f.Format(val)

		if result != tt.expected {
			t.Errorf("Locale %s: Format(%f %s) = %q, want %q",
				tt.locale, tt.value, tt.unit, result, tt.expected)
		}
	}
}

// TestLocalisationEdgeCases tests edge cases in localisation
func TestLocalisationEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		locale    string
		number    float64
		precision int
		expected  string
	}{
		{"Zero value UK", "en_GB", 0, 2, "0.00"},
		{"Zero value US", "en_US", 0, 2, "0.00"},
		{"Negative UK", "en_GB", -1234.56, 2, "-1,234.56"},
		{"Negative US", "en_US", -1234.56, 2, "-1,234.56"},
		{"Very large UK", "en_GB", 1234567890.12, 2, "1,234,567,890.12"},
		{"Very small UK", "en_GB", 0.000001, 6, "0.000001"},
		{"No precision UK", "en_GB", 1234, 0, "1,234"},
		{"No precision US", "en_US", 1234, 0, "1,234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings.Default()
			s.Locale = tt.locale
			s.Precision = tt.precision
			f := New(s)

			val := evaluator.Value{Type: evaluator.ValueNumber, Number: tt.number}
			result := f.Format(val)

			if result != tt.expected {
				t.Errorf("%s: got %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

// TestLocalisationDefaultFallback tests that unknown locales fall back gracefully
func TestLocalisationDefaultFallback(t *testing.T) {
	s := settings.Default()
	s.Locale = "unknown_LOCALE"
	s.Precision = 2
	f := New(s)

	val := evaluator.Value{Type: evaluator.ValueNumber, Number: 1234.56}
	result := f.Format(val)

	// Should fall back to basic formatting (no comma separators)
	expected := "1234.56"
	if result != expected {
		t.Errorf("Unknown locale fallback: got %q, want %q", result, expected)
	}
}
