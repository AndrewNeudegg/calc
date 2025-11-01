package units

import (
	"math"
	"testing"
)

func TestLengthConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
	}{
		{10, "m", "cm", 1000},
		{100, "km", "m", 100000},
		{1, "km", "miles", 0.621371},
		{1, "mile", "km", 1.60934},
		{12, "inches", "cm", 30.48},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > 0.1 {
			t.Errorf("%f %s in %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
		}
	}
}

func TestMassConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
	}{
		{1, "kg", "g", 1000},
		{70, "kg", "lb", 154.324},
		{1, "lb", "g", 453.592},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > 0.1 {
			t.Errorf("%f %s in %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
		}
	}
}

func TestTimeConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
	}{
		{2, "hours", "minutes", 120},
		{1, "hour", "seconds", 3600},
		{1, "day", "hours", 24},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("%f %s in %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
		}
	}
}

func TestIncompatibleUnits(t *testing.T) {
	s := NewSystem()

	_, err := s.Convert(10, "kg", "metres")
	if err == nil {
		t.Error("expected error when converting incompatible units, got nil")
	}
}

func TestCustomUnits(t *testing.T) {
	s := NewSystem()

	// 1 box = 20 apples (using a base unit for reference)
	err := s.AddCustomUnit("box", 20, "kg")
	if err != nil {
		t.Fatalf("failed to add custom unit: %s", err)
	}

	result, err := s.Convert(2, "box", "kg")
	if err != nil {
		t.Fatalf("conversion failed: %s", err)
	}

	if math.Abs(result-40) > 0.01 {
		t.Errorf("expected 40, got %.2f", result)
	}
}

func TestTemperatureConversion(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
	}{
		{0, "c", "f", 32},
		{100, "c", "f", 212},
		{32, "f", "c", 0},
		{212, "f", "c", 100},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("%f %s in %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
		}
	}
}
