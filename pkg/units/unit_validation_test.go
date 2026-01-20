package units

import (
	"testing"
)

// TestEvaluatorUnitsRegistration ensures all time units used in the evaluator
// switch statements are properly registered in the units system.
// This test prevents issues like the business days bug where units were handled
// in date arithmetic but not registered for general conversions.
func TestEvaluatorUnitsRegistration(t *testing.T) {
	s := NewSystem()

	// These are all the time unit variants used in evaluator switch statements
	// for date arithmetic (from evaluator.go lines 231-243, 702-714, 1216-1224)
	evaluatorTimeUnits := []string{
		// Full forms
		"day", "days",
		"week", "weeks",
		"month", "months",
		"year", "years",
		"hour", "hours",
		"minute", "minutes",
		"second", "seconds",
		// Short forms
		"d",   // day
		"w",   // week
		"mo",  // month
		"y",   // year
		"h",   // hour
		"hr",  // hour
		"min", // minute
		"m",   // minute (used in line 1221)
		"s",   // second
		"sec", // second
	}

	missing := []string{}
	for _, unit := range evaluatorTimeUnits {
		if !s.IsUnit(unit) {
			missing = append(missing, unit)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Found %d time units used in evaluator that are not registered in the units system: %v", len(missing), missing)
		t.Error("This will cause conversion errors. Please register these units in initStandardUnits()")
	}
}

// TestBusinessUnitsRegistration ensures business time units are properly registered.
// Business units should support conversion to standard time units.
func TestBusinessUnitsRegistration(t *testing.T) {
	s := NewSystem()

	businessUnits := []string{
		"business day",
		"business days",
		"business week",
		"business weeks",
		"business month",
		"business months",
	}

	missing := []string{}
	for _, unit := range businessUnits {
		if !s.IsUnit(unit) {
			missing = append(missing, unit)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Found %d business time units that are not registered: %v", len(missing), missing)
		t.Error("Business units must be registered to support conversion to standard time units")
	}
}

// TestUnitConversionConsistency verifies that units can be converted both ways.
// This ensures that if a unit is registered, it can participate in conversions.
func TestUnitConversionConsistency(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		unit1     string
		unit2     string
		value     float64
		tolerance float64
	}{
		// Short form time units should work
		{name: "day short form", unit1: "d", unit2: "hours", value: 1, tolerance: 0.01},
		{name: "week short form", unit1: "w", unit2: "days", value: 1, tolerance: 0.01},
		{name: "month short form", unit1: "mo", unit2: "days", value: 1, tolerance: 0.5},
		// Business time units should work
		{name: "business day to hours", unit1: "business day", unit2: "hours", value: 1, tolerance: 0.01},
		{name: "business week to hours", unit1: "business week", unit2: "hours", value: 1, tolerance: 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.unit1, tt.unit2)
			if err != nil {
				t.Errorf("Failed to convert %f %s to %s: %v", tt.value, tt.unit1, tt.unit2, err)
			}

			// Also test reverse conversion
			reverseResult, err := s.Convert(result, tt.unit2, tt.unit1)
			if err != nil {
				t.Errorf("Failed to convert back %f %s to %s: %v", result, tt.unit2, tt.unit1, err)
			}

			// Check round-trip conversion gives us back the original value
			if result > 0 && tt.value-reverseResult > tt.tolerance {
				t.Errorf("Round-trip conversion %f %s -> %s -> %s gave %f, expected %f",
					tt.value, tt.unit1, tt.unit2, tt.unit1, reverseResult, tt.value)
			}
		})
	}
}
