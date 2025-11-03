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
		// Fine-grained mass units
		{1, "mg", "µg", 1000},
		{1, "g", "mg", 1000},
		{1, "mg", "ug", 1000}, // alternative spelling
		{1000, "µg", "mg", 1},
		{1000000, "micrograms", "g", 1},
		// Jewellery units
		{1, "carat", "g", 0.2},
		{5, "carats", "mg", 1000},
		{1, "ct", "mg", 200},
		{100, "carats", "oz", 0.705479},
		// Troy measures
		{1, "troyounce", "g", 31.1035},
		{1, "ozt", "g", 31.1035},
		{1, "troyoz", "oz", 1.09714}, // troy oz is heavier than regular oz
		{10, "troyounces", "kg", 0.311035},
		{1, "ozt", "carats", 155.517}, // troy ounce to carats
		// Plural forms
		{2, "stones", "kg", 12.7006},
		{100, "micrograms", "µg", 100},
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
		// Fine-grained time units
		{1, "s", "ms", 1000},
		{1, "ms", "µs", 1000},
		{1, "µs", "ns", 1000},
		{1, "millisecond", "microseconds", 1000},
		{1000, "ns", "µs", 1},
		{1000000, "ns", "ms", 1},
		{1000000000, "nanoseconds", "s", 1},
		{0.5, "s", "ms", 500},
		{250, "ms", "s", 0.25},
		{100, "µs", "ns", 100000},
		{1, "us", "ns", 1000}, // alternative spelling
		// Informal time spans
		{1, "fortnight", "days", 14},
		{1, "fortnight", "weeks", 2},
		{2, "fortnights", "days", 28},
		{1, "quarter", "months", 3},
		{1, "quarter", "days", 91.3125},
		{4, "quarters", "year", 1},
		{1, "semester", "months", 6},
		{1, "semester", "days", 182.625},
		{2, "semesters", "year", 1},
		{1, "year", "quarters", 4},
		{1, "year", "semesters", 2},
		// Mixed conversions
		{1, "week", "fortnight", 0.5},
		{1, "month", "quarter", 0.333333},
		{6, "months", "semester", 1},
		{1, "fortnight", "hours", 336},
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

func TestSpeedConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
		delta    float64 // allowed difference for precision
	}{
		// mph conversions
		{50, "mph", "kph", 80.4672, 0.001},
		{60, "mph", "mps", 26.8224, 0.001},
		{25, "mph", "fps", 36.6667, 0.001},
		{100, "mph", "knots", 86.8976, 0.001},

		// kph conversions
		{100, "kph", "mph", 62.1371, 0.001},
		{80, "kph", "mps", 22.2222, 0.001},
		{50, "kmh", "kph", 50, 0.001}, // kmh alias test

		// mps conversions (base unit)
		{10, "mps", "kph", 36, 0.001},
		{10, "mps", "mph", 22.3694, 0.001},
		{10, "mps", "fps", 32.8084, 0.001},

		// knots conversions
		{20, "knots", "mph", 23.0156, 0.001},
		{30, "knot", "kph", 55.5556, 0.01}, // slightly larger delta for rounding
		{15, "kn", "mps", 7.71666, 0.001},

		// fps conversions
		{100, "fps", "mps", 30.48, 0.001},
		{50, "fps", "mph", 34.0909, 0.001},
	}

	for _, tt := range tests {
		result, err := s.Convert(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > tt.delta {
			t.Errorf("%f %s in %s: expected %.4f, got %.4f", tt.value, tt.from, tt.to, tt.expected, result)
		}
	}
}

func TestSpeedUnitRecognition(t *testing.T) {
	s := NewSystem()

	units := []string{"mps", "mph", "kph", "kmh", "fps", "knot", "knots", "kn"}

	for _, unit := range units {
		if !s.IsUnit(unit) {
			t.Errorf("Speed unit %q should be recognized", unit)
		}
	}

	// Test case insensitivity
	upperUnits := []string{"MPS", "MPH", "KPH", "KMH", "FPS", "KNOT", "KNOTS", "KN"}
	for _, unit := range upperUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Speed unit %q (uppercase) should be recognized", unit)
		}
	}
}

func TestMassUnitRecognition(t *testing.T) {
	s := NewSystem()

	units := []string{
		// Fine-grained
		"µg", "ug", "microgram", "micrograms",
		// Jewellery
		"carat", "carats", "ct",
		"troyounce", "troyounces", "troyoz", "ozt",
		// Plural
		"stones",
	}

	for _, unit := range units {
		if !s.IsUnit(unit) {
			t.Errorf("Mass unit %q should be recognized", unit)
		}
	}

	// Test case insensitivity
	upperUnits := []string{"UG", "MICROGRAM", "CARAT", "CT", "OZT", "STONES"}
	for _, unit := range upperUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Mass unit %q (uppercase) should be recognized", unit)
		}
	}
}

func TestTimeUnitRecognition(t *testing.T) {
	s := NewSystem()

	units := []string{
		// Fine-grained
		"ns", "nanosecond", "nanoseconds",
		"µs", "us", "microsecond", "microseconds",
		"ms", "millisecond", "milliseconds",
		// Informal spans
		"fortnight", "fortnights",
		"quarter", "quarters",
		"semester", "semesters",
	}

	for _, unit := range units {
		if !s.IsUnit(unit) {
			t.Errorf("Time unit %q should be recognized", unit)
		}
	}

	// Test case insensitivity
	upperUnits := []string{"NS", "US", "MS", "FORTNIGHT", "QUARTER", "SEMESTER"}
	for _, unit := range upperUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Time unit %q (uppercase) should be recognized", unit)
		}
	}
}

func TestSpeedAbbreviationToCompoundUnit(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Abbreviation to compound unit
		{50, "kph", "km/h", 50, 0.01},
		{100, "mph", "mi/hr", 100, 0.01},
		{30, "knots", "km/h", 55.56, 0.01},
		{60, "fps", "m/s", 18.29, 0.01},

		// Compound unit to abbreviation
		{50, "km/h", "kph", 50, 0.01},
		{100, "mi/hr", "mph", 100, 0.01},
		{55.56, "km/h", "knots", 30, 0.01},
		{18.29, "m/s", "fps", 60, 0.01},

		// Mixed conversions
		{50, "km/h", "mph", 31.07, 0.01},
		{100, "mph", "km/h", 160.93, 0.01},
	}

	for _, tt := range tests {
		result, err := s.ConvertCompoundUnit(tt.value, tt.from, tt.to)
		if err != nil {
			t.Errorf("conversion %f %s to %s failed: %s", tt.value, tt.from, tt.to, err)
			continue
		}

		if math.Abs(result-tt.expected) > tt.delta {
			t.Errorf("%f %s to %s: expected %.2f, got %.2f", tt.value, tt.from, tt.to, tt.expected, result)
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

func TestDigitalStorageConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		value     float64
		from      string
		to        string
		expected  float64
		tolerance float64
	}{
		// Bytes conversions
		{"bytes to kilobytes", 2048, "bytes", "kb", 2, 0.01},
		{"kilobytes to bytes", 1, "kb", "bytes", 1024, 1},
		{"megabytes to bytes", 1, "mb", "bytes", 1048576, 1},
		{"gigabytes to megabytes", 1, "gb", "mb", 1024, 1},
		{"terabytes to gigabytes", 1, "tb", "gb", 1024, 1},
		{"petabytes to terabytes", 1, "pb", "tb", 1024, 1},

		// Bits conversions
		{"bits to bytes", 8, "bits", "bytes", 1, 0.01},
		{"bytes to bits", 1, "bytes", "bits", 8, 0.01},
		{"kilobits to kilobytes", 8, "kbit", "kb", 1, 0.01},
		{"megabits to megabytes", 8, "mbit", "mb", 1, 0.01},
		{"gigabits to gigabytes", 8, "gbit", "gb", 1, 0.01},

		// Mixed conversions
		{"megabits to kilobytes", 1, "mbit", "kb", 128, 1},
		{"gigabytes to megabits", 1, "gb", "mbit", 8192, 1},
		{"100 megabits to megabytes", 100, "mbit", "mb", 12.5, 0.1},

		// Large conversions
		{"5 terabytes to gigabytes", 5, "tb", "gb", 5120, 1},
		{"10 petabytes to terabytes", 10, "pb", "tb", 10240, 1},

		// Practical examples
		{"4k video size", 100, "gb", "tb", 0.09765625, 0.001},
		{"internet speed conversion", 1000, "mbit", "mb", 125, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion %f %s to %s failed: %s",
					tt.value, tt.from, tt.to, err)
				return
			}

			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("%s: %f %s in %s: expected %.6f, got %.6f",
					tt.name, tt.value, tt.from, tt.to, tt.expected, result)
			}
		})
	}
}

func TestDataRateConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		value     float64
		from      string
		to        string
		expected  float64
		tolerance float64
	}{
		// Bytes per second conversions
		{"bps to kbps", 2048, "bps", "kbps", 2, 0.01},
		{"kbps to mbps", 1024, "kbps", "mbps", 1, 0.01},
		{"mbps to gbps", 1024, "mbps", "gbps", 1, 0.01},
		{"gbps to tbps", 1024, "gbps", "tbps", 1, 0.01},

		// Uppercase variants (common in networking)
		{"Bps to KBps", 2048, "Bps", "KBps", 2, 0.01},
		{"MBps to GBps", 1024, "MBps", "GBps", 1, 0.01},

		// Bits per second conversions
		{"bitps to kbitps", 2048, "bitps", "kbitps", 2, 0.01},
		{"kbitps to mbitps", 1024, "kbitps", "mbitps", 1, 0.01},
		{"mbitps to gbitps", 1024, "mbitps", "gbitps", 1, 0.01},

		// Mixed bits/bytes conversions
		{"8 bitps to bps", 8, "bitps", "bps", 1, 0.01},
		{"1 kbps to bitps", 1, "kbps", "bitps", 8192, 1},
		{"100 mbps to kbps", 100, "mbps", "kbps", 102400, 1},

		// Practical networking examples
		{"gigabit ethernet", 1, "gbps", "mbps", 1024, 1},
		{"100 megabit to bytes/sec", 100, "mbps", "bps", 104857600, 1000},
		{"1 gbps to MBps", 1, "gbps", "MBps", 1024, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion %f %s to %s failed: %s",
					tt.value, tt.from, tt.to, err)
				return
			}

			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("%s: %f %s in %s: expected %.6f, got %.6f",
					tt.name, tt.value, tt.from, tt.to, tt.expected, result)
			}
		})
	}
}

func TestDigitalUnitRecognition(t *testing.T) {
	s := NewSystem()

	// Test digital storage units
	storageUnits := []string{
		"b", "byte", "bytes",
		"kb", "kilobyte", "kilobytes",
		"mb", "megabyte", "megabytes",
		"gb", "gigabyte", "gigabytes",
		"tb", "terabyte", "terabytes",
		"pb", "petabyte", "petabytes",
		"bit", "bits",
		"kbit", "kilobit", "kilobits",
		"mbit", "megabit", "megabits",
		"gbit", "gigabit", "gigabits",
		"tbit", "terabit", "terabits",
		"pbit", "petabit", "petabits",
	}

	for _, unit := range storageUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Storage unit %q should be recognized", unit)
		}
	}

	// Test data rate units
	rateUnits := []string{
		"bps", "kbps", "mbps", "gbps", "tbps",
		"Bps", "KBps", "MBps", "GBps", "TBps",
		"bitps", "kbitps", "mbitps", "gbitps", "tbitps",
	}

	for _, unit := range rateUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Data rate unit %q should be recognized", unit)
		}
	}

	// Test case insensitivity for common variants
	upperUnits := []string{"KB", "MB", "GB", "TB", "PB"}
	for _, unit := range upperUnits {
		if !s.IsUnit(unit) {
			t.Errorf("Unit %q (uppercase) should be recognized", unit)
		}
	}
}

func TestDigitalStorageEdgeCases(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		value     float64
		from      string
		to        string
		expected  float64
		tolerance float64
	}{
		// Very small values
		{"1 byte to kilobytes", 1, "byte", "kb", 0.0009765625, 0.0000001},
		{"1 bit to bytes", 1, "bit", "byte", 0.125, 0.001},

		// Very large values
		{"1000 petabytes to terabytes", 1000, "pb", "tb", 1024000, 1000},

		// Precision tests
		{"128 bytes to bits", 128, "bytes", "bits", 1024, 1},
		{"1.5 gigabytes to megabytes", 1.5, "gb", "mb", 1536, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion %f %s to %s failed: %s",
					tt.value, tt.from, tt.to, err)
				return
			}

			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("%s: expected %.10f, got %.10f",
					tt.name, tt.expected, result)
			}
		})
	}
}

func TestJewelleryMassConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Carat conversions (precious stones)
		{"1 carat to milligrams", 1, "carat", "mg", 200, 0.01},
		{"5 carats to grams", 5, "carats", "g", 1, 0.01},
		{"100 ct to oz", 100, "ct", "oz", 0.7055, 0.001},
		{"1 gram to carats", 1, "g", "carats", 5, 0.01},

		// Troy ounce conversions (precious metals)
		{"1 troy ounce to grams", 1, "troyounce", "g", 31.1035, 0.001},
		{"1 ozt to regular oz", 1, "ozt", "oz", 1.09714, 0.001},
		{"10 troy oz to kg", 10, "troyoz", "kg", 0.311035, 0.0001},
		{"1 ozt to carats", 1, "ozt", "carats", 155.517, 0.01},
		{"100 grams to ozt", 100, "g", "ozt", 3.2151, 0.001},

		// Microgram conversions (laboratory)
		{"1000 µg to mg", 1000, "µg", "mg", 1, 0.001},
		{"1 mg to micrograms", 1, "mg", "micrograms", 1000, 0.01},
		{"1 g to ug", 1, "g", "ug", 1000000, 1},
		{"500 µg to g", 500, "µg", "g", 0.0005, 0.000001},

		// Mixed conversions
		{"1 carat to µg", 1, "carat", "µg", 200000, 1},
		{"1 troy oz to carats", 1, "troyounce", "carats", 155.517, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}

func TestFineGrainedTimeConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Nanosecond conversions
		{"1000 ns to µs", 1000, "ns", "µs", 1, 0.001},
		{"1000000 ns to ms", 1000000, "ns", "ms", 1, 0.001},
		{"1000000000 ns to s", 1000000000, "nanoseconds", "s", 1, 0.001},
		{"1 second to ns", 1, "s", "ns", 1000000000, 1},

		// Microsecond conversions
		{"1000 µs to ms", 1000, "µs", "ms", 1, 0.001},
		{"1000000 us to s", 1000000, "us", "s", 1, 0.001},
		{"1 ms to microseconds", 1, "ms", "microseconds", 1000, 0.01},

		// Millisecond conversions
		{"1000 ms to s", 1000, "ms", "s", 1, 0.001},
		{"500 milliseconds to s", 500, "milliseconds", "s", 0.5, 0.001},
		{"1 s to ms", 1, "s", "ms", 1000, 0.01},
		{"1 min to ms", 1, "min", "ms", 60000, 1},

		// Cross-scale conversions
		{"1 µs to ns", 1, "µs", "ns", 1000, 0.01},
		{"1 ms to µs", 1, "ms", "µs", 1000, 0.01},
		{"1 hour to ms", 1, "hour", "ms", 3600000, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}

func TestInformalTimeSpans(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Fortnight conversions
		{"1 fortnight to days", 1, "fortnight", "days", 14, 0.01},
		{"1 fortnight to weeks", 1, "fortnight", "weeks", 2, 0.01},
		{"2 fortnights to days", 2, "fortnights", "days", 28, 0.01},
		{"1 fortnight to hours", 1, "fortnight", "hours", 336, 0.01},

		// Quarter conversions
		{"1 quarter to months", 1, "quarter", "months", 3, 0.01},
		{"1 quarter to days", 1, "quarter", "days", 91.3125, 0.01},
		{"4 quarters to year", 4, "quarters", "year", 1, 0.01},
		{"1 year to quarters", 1, "year", "quarters", 4, 0.01},

		// Semester conversions
		{"1 semester to months", 1, "semester", "months", 6, 0.01},
		{"1 semester to days", 1, "semester", "days", 182.625, 0.01},
		{"2 semesters to year", 2, "semesters", "year", 1, 0.01},
		{"1 year to semesters", 1, "year", "semesters", 2, 0.01},

		// Mixed conversions
		{"1 week to fortnight", 1, "week", "fortnight", 0.5, 0.01},
		{"1 month to quarter", 1, "month", "quarter", 0.333333, 0.001},
		{"6 months to semester", 6, "months", "semester", 1, 0.01},
		{"1 quarter to fortnights", 1, "quarter", "fortnights", 6.522, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}

func TestExtendedVolumeConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Metric subunits
		{"100 cl to litres", 100, "cl", "l", 1, 0.01},
		{"10 dl to litres", 10, "dl", "l", 1, 0.01},
		{"1 litre to cl", 1, "l", "cl", 100, 0.01},
		{"1 litre to dl", 1, "l", "dl", 10, 0.01},

		// Cubic measures
		{"1 m³ to litres", 1, "m3", "l", 1000, 0.1},
		{"1000 cm³ to litre", 1000, "cm3", "l", 1, 0.01},
		{"1 cc to ml", 1, "cc", "ml", 1, 0.01},
		{"1 litre to cc", 1, "l", "cc", 1000, 0.1},
		{"1 ft³ to litres", 1, "ft3", "l", 28.3168, 0.01},
		{"1000 mm³ to ml", 1000, "mm3", "ml", 1, 0.01},

		// Kitchen measures (US)
		{"1 cup to ml", 1, "cup", "ml", 236.588, 0.01},
		{"16 tbsp to cup", 16, "tbsp", "cup", 1, 0.01},
		{"3 tsp to tbsp", 3, "tsp", "tbsp", 1, 0.01},
		{"8 floz to cup", 8, "floz", "cup", 1, 0.01},
		{"2 cups to pint", 2, "cups", "pint", 1, 0.01},
		{"2 pints to quart", 2, "pints", "quart", 1, 0.01},
		{"4 quarts to gallon", 4, "quarts", "gallon", 1, 0.01},

		// US vs UK gallons/pints
		{"1 usgal to ukgal", 1, "usgal", "ukgal", 0.832674, 0.001},
		{"1 ukgal to usgal", 1, "ukgal", "usgal", 1.20095, 0.001},
		{"1 uspint to ukpint", 1, "uspint", "ukpint", 0.832674, 0.001},
		{"1 ukpint to l", 1, "ukpint", "l", 0.568261, 0.001},
		{"1 uspint to l", 1, "uspint", "l", 0.473176, 0.001},
		{"1 imperialgallon to l", 1, "imperialgallon", "l", 4.54609, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}

func TestExtendedAreaConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Spelled phrases
		{"1 squaremetre to sqm", 1, "squaremetre", "sqm", 1, 0.01},
		{"1 squarefoot to sqft", 1, "squarefoot", "sqft", 1, 0.01},
		{"100 squaremetres to are", 100, "squaremetres", "are", 1, 0.01},

		// Scientific units
		{"1 are to sqm", 1, "are", "sqm", 100, 0.01},
		{"10 ares to decare", 10, "ares", "decare", 1, 0.01},
		{"10 decares to hectare", 10, "decares", "hectare", 1, 0.01},
		{"100 ares to hectare", 100, "ares", "hectare", 1, 0.01},

		// Small units
		{"1 sqmm to mm²", 1, "sqmm", "mm²", 1, 0.01},
		{"1000000 mm² to m²", 1000000, "mm²", "m²", 1, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}

func TestRankineTemperature(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Rankine to other scales
		{"Absolute zero R to F", 0, "r", "f", -459.67, 0.01},
		{"Absolute zero R to C", 0, "r", "c", -273.15, 0.01},
		{"Absolute zero R to K", 0, "r", "k", 0, 0.01},
		{"Water freezing R to F", 491.67, "r", "f", 32, 0.01},
		{"Water boiling R to F", 671.67, "r", "f", 212, 0.01},

		// Other scales to Rankine
		{"Absolute zero F to R", -459.67, "f", "r", 0, 0.01},
		{"Water freezing F to R", 32, "f", "r", 491.67, 0.01},
		{"Water boiling F to R", 212, "f", "r", 671.67, 0.01},
		{"0 C to R", 0, "c", "r", 491.67, 0.01},
		{"100 C to R", 100, "c", "r", 671.67, 0.01},
		{"273.15 K to R", 273.15, "k", "r", 491.67, 0.01},

		// Using degree symbol
		{"500 °R to F", 500, "°r", "f", 40.33, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestPressureConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Basic Pascal conversions
		{"1 kPa to Pa", 1, "kpa", "pa", 1000, 1},
		{"1 MPa to kPa", 1, "mpa", "kpa", 1000, 1},
		{"1 bar to Pa", 1, "bar", "pa", 100000, 1},
		{"1000 mbar to bar", 1000, "mbar", "bar", 1, 0.01},

		// Atmospheric pressure
		{"1 atm to Pa", 1, "atm", "pa", 101325, 1},
		{"1 atm to bar", 1, "atm", "bar", 1.01325, 0.001},
		{"1 atm to psi", 1, "atm", "psi", 14.6959, 0.01},

		// PSI conversions (tyres)
		{"30 psi to bar", 30, "psi", "bar", 2.06843, 0.001},
		{"2.5 bar to psi", 2.5, "bar", "psi", 36.2594, 0.01},
		{"35 psi to kPa", 35, "psi", "kpa", 241.317, 0.1},

		// Torr and mmHg
		{"760 torr to atm", 760, "torr", "atm", 1, 0.01},
		{"760 mmhg to atm", 760, "mmhg", "atm", 1, 0.01},
		{"1 inhg to mmhg", 1, "inhg", "mmhg", 25.4, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestForceConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Newton conversions
		{"1 kN to N", 1, "kilonewton", "n", 1000, 1},
		{"1 MN to kN", 1, "mn", "kilonewton", 1000, 1},

		// Pound-force
		{"1 lbf to N", 1, "lbf", "n", 4.44822, 0.001},
		{"100 N to lbf", 100, "n", "lbf", 22.4809, 0.01},

		// Kilogram-force
		{"1 kgf to N", 1, "kgf", "n", 9.80665, 0.001},
		{"10 kgf to lbf", 10, "kgf", "lbf", 22.0462, 0.01},

		// Dyne (CGS unit)
		{"100000 dyne to N", 100000, "dyne", "n", 1, 0.01},
		{"1 N to dynes", 1, "n", "dynes", 100000, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestAngleConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Degrees to radians
		{"180 deg to rad", 180, "deg", "rad", 3.14159, 0.001},
		{"90 degrees to radians", 90, "degrees", "radians", 1.5708, 0.001},
		{"360 deg to rad", 360, "deg", "rad", 6.28319, 0.001},

		// Radians to degrees
		{"π rad to deg", 3.14159, "rad", "deg", 180, 0.01},
		{"1 radian to degrees", 1, "radian", "degrees", 57.2958, 0.01},

		// Gradians
		{"100 grad to deg", 100, "grad", "deg", 90, 0.01},
		{"400 gradians to degrees", 400, "gradians", "degrees", 360, 0.01},
		{"200 gon to deg", 200, "gon", "deg", 180, 0.01},

		// Turns/revolutions
		{"1 turn to deg", 1, "turn", "deg", 360, 0.01},
		{"0.5 turns to degrees", 0.5, "turns", "degrees", 180, 0.01},
		{"1 revolution to deg", 1, "revolution", "deg", 360, 0.01},
		{"2 revolutions to rad", 2, "revolutions", "rad", 12.5664, 0.001},

		// Degree symbol
		{"90 ° to rad", 90, "°", "rad", 1.5708, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestFrequencyConversions(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name     string
		value    float64
		from     string
		to       string
		expected float64
		delta    float64
	}{
		// Basic frequency conversions
		{"1 kHz to Hz", 1, "khz", "hz", 1000, 1},
		{"1 MHz to kHz", 1, "mhz", "khz", 1000, 1},
		{"1 GHz to MHz", 1, "ghz", "mhz", 1000, 1},
		{"1 THz to GHz", 1, "thz", "ghz", 1000, 1},

		// Practical examples
		{"2.4 GHz WiFi", 2.4, "ghz", "mhz", 2400, 1},
		{"5 GHz WiFi", 5, "ghz", "mhz", 5000, 1},
		{"100 MHz FM radio", 100, "mhz", "hz", 100000000, 100},

		// RPM conversions
		{"60 rpm to Hz", 60, "rpm", "hz", 1, 0.01},
		{"3000 rpm to Hz", 3000, "rpm", "hz", 50, 0.1},
		{"1 Hz to rpm", 1, "hz", "rpm", 60, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Errorf("conversion failed: %s", err)
				return
			}

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestAllNewUnitsRecognition(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		category string
		units    []string
	}{
		{"Volume metric", []string{"cl", "centilitre", "dl", "decilitre"}},
		{"Volume cubic", []string{"m3", "m³", "cm3", "cm³", "mm3", "mm³", "ft3", "ft³", "in3", "in³", "cc"}},
		{"Volume kitchen", []string{"cup", "cups", "floz", "tbsp", "tsp", "tablespoon", "teaspoon"}},
		{"Volume US", []string{"usgal", "usgallon", "usquart", "uspint", "quart", "pint", "qt", "pt"}},
		{"Volume UK", []string{"ukgal", "ukgallon", "impgal", "imperialgallon", "ukquart", "ukpint", "imppint"}},
		{"Area spelled", []string{"squaremetre", "squarefoot", "squareinch", "squareyard", "squaremile"}},
		{"Area scientific", []string{"are", "ares", "decare", "decares", "sqmm", "mm2", "mm²"}},
		{"Temperature", []string{"r", "rankine", "°r"}},
		{"Pressure", []string{"pa", "pascal", "kpa", "bar", "mbar", "atm", "psi", "torr", "mmhg", "inhg"}},
		{"Force", []string{"n", "newton", "kilonewton", "mn", "lbf", "kgf", "dyne"}},
		{"Angle", []string{"deg", "degree", "rad", "radian", "grad", "gradian", "gon", "turn", "revolution", "°"}},
		{"Frequency", []string{"hz", "khz", "mhz", "ghz", "thz", "rpm"}},
	}

	for _, category := range tests {
		for _, unit := range category.units {
			if !s.IsUnit(unit) {
				t.Errorf("%s unit %q should be recognized", category.category, unit)
			}
		}
	}
}
