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
		// New mass units
		{1, "lbs", "kg", 0.453592},
		{10, "stone", "kg", 63.5029},
		{1, "stone", "lbs", 14},
		{1, "tonne", "kg", 1000},
		{1, "ton", "kg", 907.185},
		{1, "ton", "lbs", 2000},
		{1, "tonne", "ton", 1.10231},
		{150, "lbs", "stone", 10.7143},
		{100, "kg", "stone", 15.7473},
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

func TestAreaConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value     float64
		from      string
		to        string
		expected  float64
		tolerance float64
	}{
		{1, "sqm", "sqft", 10.7639, 0.01},
		{100, "sqft", "sqm", 9.2903, 0.01},
		{1, "hectare", "sqm", 10000, 0.1},
		{1, "acre", "sqm", 4046.86, 0.1},
		{1, "hectare", "acres", 2.47105, 0.01},
		{1, "acre", "hectares", 0.404686, 0.01},
		{1, "sqkm", "hectares", 100, 0.1},
		{1, "sqmi", "acres", 640, 1},
		{10000, "sqcm", "sqm", 1, 0.01},
		{1, "sqyd", "sqft", 9, 0.1},
		{144, "sqin", "sqft", 1, 0.01},
		{1, "sqkm", "sqmi", 0.386102, 0.001},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > tt.tolerance {
			t.Errorf("%f %s in %s: expected %.5f, got %.5f", tt.value, tt.from, tt.to, tt.expected, result)
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
		// Celsius to Fahrenheit
		{0, "c", "f", 32},
		{100, "c", "f", 212},
		{-40, "c", "f", -40},

		// Fahrenheit to Celsius
		{32, "f", "c", 0},
		{212, "f", "c", 100},
		{-40, "f", "c", -40},

		// Celsius to Kelvin
		{0, "c", "k", 273.15},
		{100, "c", "k", 373.15},
		{-273.15, "c", "k", 0},

		// Kelvin to Celsius
		{273.15, "k", "c", 0},
		{373.15, "k", "c", 100},
		{0, "k", "c", -273.15},

		// Fahrenheit to Kelvin
		{32, "f", "k", 273.15},
		{212, "f", "k", 373.15},
		{-459.67, "f", "k", 0},

		// Kelvin to Fahrenheit
		{273.15, "k", "f", 32},
		{373.15, "k", "f", 212},
		{0, "k", "f", -459.67},

		// Same unit conversions
		{25, "c", "celsius", 25},
		{77, "f", "fahrenheit", 77},
		{300, "k", "kelvin", 300},
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

func TestKelvinConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
	}{
		// Absolute zero
		{"Absolute zero K to C", 0, "k", "c", -273.15},
		{"Absolute zero K to F", 0, "k", "f", -459.67},

		// Water freezing point
		{"Water freezing K to C", 273.15, "k", "c", 0},
		{"Water freezing K to F", 273.15, "k", "f", 32},

		// Water boiling point
		{"Water boiling K to C", 373.15, "k", "c", 100},
		{"Water boiling K to F", 373.15, "k", "f", 212},

		// Room temperature (~ 20°C)
		{"Room temp K to C", 293.15, "k", "c", 20},
		{"Room temp K to F", 293.15, "k", "f", 68},

		// Reverse conversions
		{"C to K absolute zero", -273.15, "c", "k", 0},
		{"F to K absolute zero", -459.67, "f", "k", 0},

		// Using full names
		{"Kelvin full name to Celsius", 300, "kelvin", "celsius", 26.85},
		{"Celsius full name to Kelvin", 25, "celsius", "kelvin", 298.15},

		// Scientific temperatures
		{"Liquid nitrogen K to C", 77, "k", "c", -196.15},
		{"Human body temp C to K", 37, "c", "k", 310.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Fatalf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			}

			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("%f %s in %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
			}
		})
	}
}

func TestTemperatureUnitRecognition(t *testing.T) {
	s := NewSystem()

	units := []string{"c", "celsius", "f", "fahrenheit", "k", "kelvin"}

	for _, unit := range units {
		if !s.IsUnit(unit) {
			t.Errorf("Unit %q should be recognized", unit)
		}
	}

	// Test case insensitivity
	upperUnits := []string{"C", "CELSIUS", "F", "FAHRENHEIT", "K", "KELVIN"}
	for _, unit := range upperUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Unit %q (uppercase) should be recognized", unit)
		}
	}
}

func TestParseCompoundUnit(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		unitStr  string
		wantErr  bool
		checkNum string // expected numerator unit
		checkDen string // expected denominator unit
	}{
		{"km per hour", "km/h", false, "km", "h"},
		{"meters per second", "m/s", false, "m", "s"},
		{"miles per hour", "mi/hr", false, "mi", "hr"},
		{"feet per second", "ft/s", false, "ft", "s"},
		{"invalid numerator", "xyz/h", true, "", ""},
		{"invalid denominator", "km/xyz", true, "", ""},
		{"no slash", "kmh", true, "", ""},
		{"empty string", "", true, "", ""},
		{"only slash", "/", true, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cu, err := s.ParseCompoundUnit(tt.unitStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCompoundUnit(%q) expected error, got nil", tt.unitStr)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCompoundUnit(%q) unexpected error: %v", tt.unitStr, err)
				return
			}

			if cu.Numerator == nil || cu.Denominator == nil {
				t.Errorf("ParseCompoundUnit(%q) returned nil numerator or denominator", tt.unitStr)
				return
			}

			// Check that conversion factors are set
			if cu.ToBaseNum == 0 || cu.ToBaseDen == 0 {
				t.Errorf("ParseCompoundUnit(%q) conversion factors not set properly", tt.unitStr)
			}
		})
	}
}

func TestConvertCompoundUnit(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		value     float64
		from      string
		to        string
		expected  float64
		tolerance float64
	}{
		{"km/h to m/s", 50, "km/h", "m/s", 13.8889, 0.01},
		{"m/s to km/h", 13.8889, "m/s", "km/h", 50, 0.1},
		{"mi/hr to km/h", 60, "mi/hr", "km/h", 96.56, 0.1},
		{"km/h to mi/hr", 100, "km/h", "mi/hr", 62.14, 0.1},
		{"ft/s to m/s", 100, "ft/s", "m/s", 30.48, 0.01},
		{"m/s to ft/s", 30.48, "m/s", "ft/s", 100, 0.1},
		{"specification example", 50, "km/h", "m/s", 13.8889, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.ConvertCompoundUnit(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("ConvertCompoundUnit(%f %s to %s) failed: %s",
					tt.value, tt.from, tt.to, err)
				return
			}

			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("ConvertCompoundUnit(%f %s to %s): expected %.4f, got %.4f",
					tt.value, tt.from, tt.to, tt.expected, result)
			}
		})
	}
}

func TestConvertCompoundUnitErrors(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name  string
		value float64
		from  string
		to    string
	}{
		{"invalid from unit", 50, "xyz/h", "m/s"},
		{"invalid to unit", 50, "km/h", "xyz/s"},
		{"incompatible dimensions numerator", 50, "kg/h", "m/s"},
		{"incompatible dimensions denominator", 50, "km/kg", "m/s"},
		{"empty from", 50, "", "m/s"},
		{"empty to", 50, "km/h", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.ConvertCompoundUnit(tt.value, tt.from, tt.to)
			if err == nil {
				t.Errorf("ConvertCompoundUnit(%f %s to %s) expected error, got nil",
					tt.value, tt.from, tt.to)
			}
		})
	}
}

func TestIsCompoundUnit(t *testing.T) {
	tests := []struct {
		unitStr  string
		expected bool
	}{
		{"km/h", true},
		{"m/s", true},
		{"ft/s", true},
		{"mi/hr", true},
		{"km", false},
		{"m", false},
		{"hours", false},
		{"", false},
		{"/", true}, // contains slash even though invalid
	}

	for _, tt := range tests {
		result := IsCompoundUnit(tt.unitStr)
		if result != tt.expected {
			t.Errorf("IsCompoundUnit(%q) = %v, expected %v",
				tt.unitStr, result, tt.expected)
		}
	}
}
