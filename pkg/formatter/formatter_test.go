package formatter

import (
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/settings"
)

func TestFormatNumber(t *testing.T) {
	s := settings.Default()
	s.Precision = 2
	s.Locale = "en_GB"
	f := New(s)

	tests := []struct {
		input    float64
		expected string
	}{
		{1000, "1,000.00"},
		{1000000, "1,000,000.00"},
		{3.14159, "3.14"},
		{0.5, "0.50"},
		{42, "42.00"},
	}

	for _, tt := range tests {
		val := evaluator.Value{Type: evaluator.ValueNumber, Number: tt.input}
		result := f.Format(val)
		if result != tt.expected {
			t.Errorf("Format(%f) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatCurrency(t *testing.T) {
	s := settings.Default()
	s.Precision = 2
	s.Locale = "en_GB"
	f := New(s)

	tests := []struct {
		amount   float64
		code     string
		expected string
	}{
		{100, "£", "£100.00"},
		{50.5, "£", "£50.50"},
		{1000, "$", "$1,000.00"},
	}

	for _, tt := range tests {
		val := evaluator.Value{Type: evaluator.ValueCurrency, Number: tt.amount, Currency: tt.code}
		result := f.Format(val)
		if result != tt.expected {
			t.Errorf("Format(currency %f, %q) = %q, want %q", tt.amount, tt.code, result, tt.expected)
		}
	}
}

func TestFormatDate(t *testing.T) {
	s := settings.Default()
	s.DateFormat = "2 Jan 2006"
	f := New(s)

	date := time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC)
	val := evaluator.Value{Type: evaluator.ValueDate, Date: date}
	result := f.Format(val)
	expected := "15 Nov 2025"

	if result != expected {
		t.Errorf("Format(date) = %q, want %q", result, expected)
	}
}

func TestFormatUnit(t *testing.T) {
	s := settings.Default()
	s.Precision = 2
	f := New(s)

	val := evaluator.Value{Type: evaluator.ValueUnit, Number: 1000, Unit: "cm"}
	result := f.Format(val)
	expected := "1,000.00 cm"

	if result != expected {
		t.Errorf("Format(unit) = %q, want %q", result, expected)
	}
}

func TestFormatPercent(t *testing.T) {
	s := settings.Default()
	s.Precision = 0
	f := New(s)

	tests := []struct {
		value    float64
		expected string
	}{
		{50, "50%"},
		{25, "25%"},
		{100, "100%"},
		{10, "10%"},
	}

	for _, tt := range tests {
		val := evaluator.Value{Type: evaluator.ValuePercent, Number: tt.value}
		result := f.Format(val)
		if result != tt.expected {
			t.Errorf("Format(percent %f) = %q, want %q", tt.value, result, tt.expected)
		}
	}
}

func TestFormatError(t *testing.T) {
	s := settings.Default()
	f := New(s)

	val := evaluator.Value{Type: evaluator.ValueError, Error: "test error"}
	result := f.Format(val)

	if result != "Error: test error" {
		t.Errorf("Format(error) = %q, want %q", result, "Error: test error")
	}
}

func TestFormatDateWithTime(t *testing.T) {
	s := settings.Default()
	s.DateFormat = "2 Jan 2006"
	f := New(s)

	// Date with time component should show time
	dateWithTime := time.Date(2025, 11, 15, 14, 30, 45, 0, time.UTC)
	val := evaluator.Value{Type: evaluator.ValueDate, Date: dateWithTime}
	result := f.Format(val)
	expected := "15 Nov 2025 14:30:45 UTC"

	if result != expected {
		t.Errorf("Format(date with time) = %q, want %q", result, expected)
	}

	// Date without time component should just show date
	dateWithoutTime := time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC)
	val2 := evaluator.Value{Type: evaluator.ValueDate, Date: dateWithoutTime}
	result2 := f.Format(val2)
	expected2 := "15 Nov 2025"

	if result2 != expected2 {
		t.Errorf("Format(date without time) = %q, want %q", result2, expected2)
	}
}
