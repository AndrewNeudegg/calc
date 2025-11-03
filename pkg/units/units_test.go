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
