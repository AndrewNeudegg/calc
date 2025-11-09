package integration

import (
	"os"
	"testing"

	"github.com/andrewneudegg/calc/pkg/display"
)

// TestREPLCustomUnitCommands tests all :unit commands in the REPL
func TestREPLCustomUnitCommands(t *testing.T) {
	repl := display.NewREPL()
	repl.SetSilent(true)
	
	// Clean up any existing test units first
	repl.Env().Units().DeleteCustomUnit("customspoon")
	repl.Env().Units().DeleteCustomUnit("customtablespoon")
	repl.Env().Units().DeleteCustomUnit("megaspoon")
	
	// Test 1: Define a custom unit using command form
	result := repl.EvaluateLine(":unit define customspoon = 5 ml")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to define customspoon: %s", result.Error)
	}
	
	// Test 2: Verify unit was created
	if !repl.Env().Units().UnitExists("customspoon") {
		t.Error("customspoon unit should exist after definition")
	}
	
	// Test 3: Use the custom unit
	result = repl.EvaluateLine("2 customspoon in ml")
	if result.IsError() {
		t.Fatalf("failed to use custom unit: %s", result.Error)
	}
	if result.Number != 10.0 {
		t.Errorf("expected 10.0 ml, got %.2f", result.Number)
	}
	
	// Test 4: Define using shorthand directive form
	result = repl.EvaluateLine(":unit customtablespoon = 15 ml")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to define customtablespoon: %s", result.Error)
	}
	
	// Test 5: Verify shorthand form works
	if !repl.Env().Units().UnitExists("customtablespoon") {
		t.Error("customtablespoon unit should exist after definition")
	}
	
	// Test 6: Unit chaining - define unit based on another custom unit
	result = repl.EvaluateLine(":unit megaspoon = 100 customspoon")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to define megaspoon: %s", result.Error)
	}
	
	result = repl.EvaluateLine("1 megaspoon in ml")
	if result.IsError() {
		t.Fatalf("failed to convert megaspoon: %s", result.Error)
	}
	if result.Number != 500.0 {
		t.Errorf("expected 500.0 ml, got %.2f", result.Number)
	}
	
	// Test 7: Test :unit show command
	result = repl.EvaluateLine(":unit show customspoon")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to show unit: %s", result.Error)
	}
	
	// Test 8: Test :unit list command
	result = repl.EvaluateLine(":unit list custom")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to list custom units: %s", result.Error)
	}
	
	// Test 9: Delete a custom unit
	result = repl.EvaluateLine(":unit delete customspoon")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to delete unit: %s", result.Error)
	}
	
	// Test 10: Verify unit was deleted
	if repl.Env().Units().UnitExists("customspoon") {
		t.Error("customspoon should not exist after deletion")
	}
	
	// Clean up
	repl.Env().Units().DeleteCustomUnit("customtablespoon")
	repl.Env().Units().DeleteCustomUnit("megaspoon")
}

// TestREPLCustomUnitErrorHandling tests error cases
func TestREPLCustomUnitErrorHandling(t *testing.T) {
	repl := display.NewREPL()
	repl.SetSilent(true)
	
	// Test 1: Try to redefine a standard unit
	result := repl.EvaluateLine(":unit meter = 100 cm")
	if !result.IsError() {
		t.Error("should not allow redefining standard units")
	}
	
	// Test 2: Invalid base unit
	result = repl.EvaluateLine(":unit foo = 10 nonexistent")
	if !result.IsError() {
		t.Error("should error with invalid base unit")
	}
	
	// Test 3: Circular definition
	repl.Env().Units().AddCustomUnit("bar", 1.0, "m")
	result = repl.EvaluateLine(":unit bar = 5 bar")
	if !result.IsError() {
		t.Error("should error on circular definition")
	}
	
	// Test 4: Try to delete non-existent unit
	result = repl.EvaluateLine(":unit delete nonexistent")
	if result.IsError() && result.Error != "" {
		// Expected - command will return error message
	}
	
	// Clean up
	repl.Env().Units().DeleteCustomUnit("bar")
}

// TestREPLCustomUnitExpressions tests defining units with expressions
func TestREPLCustomUnitExpressions(t *testing.T) {
	repl := display.NewREPL()
	repl.SetSilent(true)
	
	// Clean up first
	repl.Env().Units().DeleteCustomUnit("customfoot")
	repl.Env().Units().DeleteCustomUnit("customyard")
	
	// Test 1: Define foot and yard with custom names
	result := repl.EvaluateLine(":unit define customfoot = 0.3048 m")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to define customfoot: %s", result.Error)
	}
	
	result = repl.EvaluateLine(":unit customyard = 3 customfoot")
	if result.IsError() && result.Error != "" {
		t.Errorf("failed to define customyard: %s", result.Error)
	}
	
	// Test 2: Verify customyard = 3 customfoot = 3 * 0.3048 m = 0.9144 m
	result = repl.EvaluateLine("1 customyard in m")
	if result.IsError() {
		t.Fatalf("failed to convert customyard: %s", result.Error)
	}
	expected := 0.9144
	if result.Number < expected-0.001 || result.Number > expected+0.001 {
		t.Errorf("expected approximately %.4f m, got %.4f m", expected, result.Number)
	}
	
	// Test 3: Verify conversions work
	result = repl.EvaluateLine("1 customyard in customfoot")
	if result.IsError() {
		t.Fatalf("failed to convert customyard to customfoot: %s", result.Error)
	}
	if result.Number < 2.99 || result.Number > 3.01 {
		t.Errorf("expected approximately 3.0 customfoot, got %.2f", result.Number)
	}
	
	// Clean up
	repl.Env().Units().DeleteCustomUnit("customfoot")
	repl.Env().Units().DeleteCustomUnit("customyard")
}

// TestREPLCustomUnitPersistence tests persistence functionality
func TestREPLCustomUnitPersistence(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := t.TempDir() + "/test_custom_units.json"
	
	// Create first REPL and add units
	repl1 := display.NewREPL()
	repl1.SetSilent(true)
	
	repl1.Env().Units().AddCustomUnit("testunit1", 25.0, "ml")
	repl1.Env().Units().AddCustomUnit("testunit2", 50.0, "g")
	
	// Save units
	err := repl1.Env().Units().SaveCustomUnits(tmpFile)
	if err != nil {
		t.Fatalf("failed to save custom units: %v", err)
	}
	
	// Create second REPL and load units
	repl2 := display.NewREPL()
	repl2.SetSilent(true)
	
	err = repl2.Env().Units().LoadCustomUnits(tmpFile)
	if err != nil {
		t.Fatalf("failed to load custom units: %v", err)
	}
	
	// Verify units were loaded
	if !repl2.Env().Units().UnitExists("testunit1") {
		t.Error("testunit1 should exist after loading")
	}
	if !repl2.Env().Units().UnitExists("testunit2") {
		t.Error("testunit2 should exist after loading")
	}
	
	// Test conversions work correctly after loading
	result := repl2.EvaluateLine("2 testunit1 in ml")
	if result.IsError() {
		t.Fatalf("conversion error: %s", result.Error)
	}
	if result.Number != 50.0 {
		t.Errorf("expected 50.0 ml, got %.2f", result.Number)
	}
	
	// Clean up
	os.Remove(tmpFile)
}

// TestREPLCustomUnitListFilters tests the list command filters
func TestREPLCustomUnitListFilters(t *testing.T) {
	repl := display.NewREPL()
	repl.SetSilent(true)
	
	// Add some custom units
	repl.Env().Units().AddCustomUnit("testunit", 10.0, "ml")
	
	// Test :unit list (should work)
	result := repl.EvaluateLine(":unit list")
	if result.IsError() && result.Error != "" {
		t.Errorf("':unit list' failed: %s", result.Error)
	}
	
	// Test :unit list custom
	result = repl.EvaluateLine(":unit list custom")
	if result.IsError() && result.Error != "" {
		t.Errorf("':unit list custom' failed: %s", result.Error)
	}
	
	// Test :unit list builtin
	result = repl.EvaluateLine(":unit list builtin")
	if result.IsError() && result.Error != "" {
		t.Errorf("':unit list builtin' failed: %s", result.Error)
	}
	
	// Clean up
	repl.Env().Units().DeleteCustomUnit("testunit")
}

// TestREPLCustomUnitWithDifferentDimensions tests units across different dimensions
func TestREPLCustomUnitWithDifferentDimensions(t *testing.T) {
	// Test different unit dimensions
	tests := []struct {
		name     string
		unitDef  string
		testExpr string
		expected float64
		unit     string
		unitName string
	}{
		{
			name:     "volume unit",
			unitDef:  ":unit mymug = 300 ml",
			testExpr: "2 mymug in l",
			expected: 0.6,
			unit:     "l",
			unitName: "mymug",
		},
		{
			name:     "length unit",
			unitDef:  ":unit myhand = 10.16 cm",
			testExpr: "1 myhand in m",
			expected: 0.1016,
			unit:     "m",
			unitName: "myhand",
		},
		{
			name:     "mass unit",
			unitDef:  ":unit myrock = 6350 g",
			testExpr: "1 myrock in kg",
			expected: 6.35,
			unit:     "kg",
			unitName: "myrock",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh REPL for each test to avoid interference
			repl := display.NewREPL()
			repl.SetSilent(true)
			
			// Clean up first
			repl.Env().Units().DeleteCustomUnit(tt.unitName)
			
			// Define the unit
			result := repl.EvaluateLine(tt.unitDef)
			if result.IsError() && result.Error != "" {
				t.Fatalf("failed to define unit: %s", result.Error)
			}
			
			// Test the unit
			result = repl.EvaluateLine(tt.testExpr)
			if result.IsError() {
				t.Fatalf("failed to evaluate %q: %s", tt.testExpr, result.Error)
			}
			
			if result.Number < tt.expected-0.01 || result.Number > tt.expected+0.01 {
				t.Errorf("expected approximately %.4f %s, got %.4f %s", 
					tt.expected, tt.unit, result.Number, result.Unit)
			}
			
			// Clean up
			repl.Env().Units().DeleteCustomUnit(tt.unitName)
		})
	}
}
