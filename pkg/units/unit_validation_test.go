package units

import (
	"math"
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
	// Note: "m" is NOT included as it conflicts with metres and is only used
	// in a specific context (evalTimeDifference) where it's disambiguated
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
			// Use absolute difference to handle both positive and negative differences
			if math.Abs(tt.value-reverseResult) > tt.tolerance {
				t.Errorf("Round-trip conversion %f %s -> %s -> %s gave %f, expected %f (diff: %f, tolerance: %f)",
					tt.value, tt.unit1, tt.unit2, tt.unit1, reverseResult, tt.value,
					math.Abs(tt.value-reverseResult), tt.tolerance)
			}
		})
	}
}

// TestRoundTripConversionThorough provides comprehensive testing of round-trip conversions
// to ensure the math.Abs fix works correctly in all scenarios.
// This test addresses the review comment about incorrect round-trip logic.
func TestRoundTripConversionThorough(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name      string
		unit1     string
		unit2     string
		values    []float64 // Test multiple values including edge cases
		tolerance float64
	}{
		{
			name:      "positive values - time units",
			unit1:     "hours",
			unit2:     "minutes",
			values:    []float64{1, 2.5, 10, 100, 0.1},
			tolerance: 0.0001,
		},
		{
			name:      "positive values - length units",
			unit1:     "m",
			unit2:     "cm",
			values:    []float64{1, 5.5, 0.01, 1000},
			tolerance: 0.0001,
		},
		{
			name:      "positive values - mass units",
			unit1:     "kg",
			unit2:     "g",
			values:    []float64{1, 0.5, 10, 0.001},
			tolerance: 0.0001,
		},
		{
			name:      "business time to standard time",
			unit1:     "business days",
			unit2:     "hours",
			values:    []float64{1, 5, 10, 0.5},
			tolerance: 0.001,
		},
		{
			name:      "business week conversions",
			unit1:     "business weeks",
			unit2:     "business days",
			values:    []float64{1, 2, 10, 0.2},
			tolerance: 0.001,
		},
		{
			name:      "fractional conversions",
			unit1:     "days",
			unit2:     "hours",
			values:    []float64{0.25, 0.5, 0.75, 1.5, 2.25},
			tolerance: 0.0001,
		},
		{
			name:      "very small values",
			unit1:     "seconds",
			unit2:     "milliseconds",
			values:    []float64{0.001, 0.0001, 0.00001},
			tolerance: 0.000001,
		},
		{
			name:      "very large values",
			unit1:     "km",
			unit2:     "m",
			values:    []float64{1000, 10000, 100000},
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, value := range tt.values {
				// Forward conversion
				result, err := s.Convert(value, tt.unit1, tt.unit2)
				if err != nil {
					t.Errorf("Failed to convert %f %s to %s: %v", value, tt.unit1, tt.unit2, err)
					continue
				}

				// Reverse conversion
				reverseResult, err := s.Convert(result, tt.unit2, tt.unit1)
				if err != nil {
					t.Errorf("Failed to convert back %f %s to %s: %v", result, tt.unit2, tt.unit1, err)
					continue
				}

				// Check round-trip with absolute difference
				diff := math.Abs(value - reverseResult)
				if diff > tt.tolerance {
					t.Errorf("Round-trip conversion failed for value %f:\n"+
						"  %f %s -> %f %s -> %f %s\n"+
						"  Expected: %f, Got: %f, Diff: %f, Tolerance: %f",
						value,
						value, tt.unit1, result, tt.unit2, reverseResult, tt.unit1,
						value, reverseResult, diff, tt.tolerance)
				}

				// Verify the difference calculation is correct
				if diff < 0 {
					t.Errorf("BUG: Absolute difference should never be negative! Got: %f", diff)
				}
			}
		})
	}
}

// TestRoundTripEdgeCases tests edge cases that could expose issues with the
// round-trip conversion logic (e.g., the old bug with result > 0 check).
func TestRoundTripEdgeCases(t *testing.T) {
	s := NewSystem()

	tests := []struct {
		name        string
		unit1       string
		unit2       string
		value       float64
		tolerance   float64
		description string
	}{
		{
			name:        "zero value",
			unit1:       "hours",
			unit2:       "minutes",
			value:       0,
			tolerance:   0.0001,
			description: "Zero should round-trip correctly",
		},
		{
			name:        "very small positive",
			unit1:       "m",
			unit2:       "mm",
			value:       0.000001,
			tolerance:   0.000001,
			description: "Very small values should work",
		},
		{
			name:        "exact power of 10",
			unit1:       "km",
			unit2:       "m",
			value:       1.0,
			tolerance:   0.0001,
			description: "Powers of 10 should convert exactly",
		},
		{
			name:        "irrational result",
			unit1:       "mile",
			unit2:       "km",
			value:       1.0,
			tolerance:   0.001,
			description: "Non-exact conversions should still round-trip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.Convert(tt.value, tt.unit1, tt.unit2)
			if err != nil {
				t.Fatalf("Forward conversion failed: %v", err)
			}

			reverseResult, err := s.Convert(result, tt.unit2, tt.unit1)
			if err != nil {
				t.Fatalf("Reverse conversion failed: %v", err)
			}

			diff := math.Abs(tt.value - reverseResult)
			if diff > tt.tolerance {
				t.Errorf("%s: Round-trip failed\n"+
					"  Original: %f %s\n"+
					"  After conversion: %f %s\n"+
					"  After round-trip: %f %s\n"+
					"  Difference: %f (tolerance: %f)",
					tt.description,
					tt.value, tt.unit1,
					result, tt.unit2,
					reverseResult, tt.unit1,
					diff, tt.tolerance)
			}
		})
	}
}
