package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestCustomUnitDefinitionBasic tests basic custom unit definition and usage
func TestCustomUnitDefinitionBasic(t *testing.T) {
	env := evaluator.NewEnvironment()
	
	// Define a custom unit: spoon = 15 ml
	err := env.Units().AddCustomUnit("spoon", 15.0, "ml")
	if err != nil {
		t.Fatalf("failed to add custom unit: %v", err)
	}
	
	// Test that unit exists
	if !env.Units().UnitExists("spoon") {
		t.Error("custom unit 'spoon' should exist")
	}
	
	// Test parsing "2 spoon in ml"
	tests := []struct {
		input    string
		expected float64
		unit     string
	}{
		{"1 spoon", 1.0, "spoon"},
		{"2 spoon in ml", 30.0, "ml"},
		{"60 ml in spoon", 4.0, "spoon"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			l.SetUnitChecker(env.Units().UnitExists)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error for %q: %v", tt.input, err)
			}
			
			result := env.Eval(expr)
			if result.IsError() {
				t.Fatalf("evaluation error for %q: %s", tt.input, result.Error)
			}
			
			if result.Number != tt.expected {
				t.Errorf("expected %.2f for %q, got %.2f", tt.expected, tt.input, result.Number)
			}
			
			if result.Unit != tt.unit {
				t.Errorf("expected unit %q for %q, got %q", tt.unit, tt.input, result.Unit)
			}
		})
	}
}

// TestCustomUnitChaining tests defining units using other custom units
func TestCustomUnitChaining(t *testing.T) {
	env := evaluator.NewEnvironment()
	
	// Define spoon = 15 ml
	err := env.Units().AddCustomUnit("spoon", 15.0, "ml")
	if err != nil {
		t.Fatalf("failed to add spoon unit: %v", err)
	}
	
	// Define bowl = 350 ml
	err = env.Units().AddCustomUnit("bowl", 350.0, "ml")
	if err != nil {
		t.Fatalf("failed to add bowl unit: %v", err)
	}
	
	// Test bowl / spoon (should be approximately 23.33)
	l := lexer.New("1 bowl in spoon")
	l.SetUnitChecker(env.Units().UnitExists)
	tokens := l.AllTokens()
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	result := env.Eval(expr)
	if result.IsError() {
		t.Fatalf("evaluation error: %s", result.Error)
	}
	
	expected := 350.0 / 15.0 // 23.33...
	if result.Number < expected-0.01 || result.Number > expected+0.01 {
		t.Errorf("expected approximately %.2f, got %.2f", expected, result.Number)
	}
}

// TestCustomUnitErrorHandling tests error cases
func TestCustomUnitErrorHandling(t *testing.T) {
	env := evaluator.NewEnvironment()
	
	// Test adding unit with unknown base
	err := env.Units().AddCustomUnit("foo", 10.0, "unknownunit")
	if err == nil {
		t.Error("expected error when adding unit with unknown base")
	}
	
	// Test circular definition
	err = env.Units().AddCustomUnit("bar", 1.0, "bar")
	if err == nil {
		t.Error("expected error for circular definition")
	}
}

// TestCustomUnitDelete tests deleting custom units
func TestCustomUnitDelete(t *testing.T) {
	env := evaluator.NewEnvironment()
	
	// Add a custom unit
	err := env.Units().AddCustomUnit("temp", 10.0, "ml")
	if err != nil {
		t.Fatalf("failed to add custom unit: %v", err)
	}
	
	// Verify it exists
	if !env.Units().UnitExists("temp") {
		t.Error("custom unit should exist")
	}
	
	// Delete it
	err = env.Units().DeleteCustomUnit("temp")
	if err != nil {
		t.Fatalf("failed to delete custom unit: %v", err)
	}
	
	// Verify it's gone
	if env.Units().UnitExists("temp") {
		t.Error("custom unit should not exist after deletion")
	}
}

// TestCustomUnitList tests listing custom units
func TestCustomUnitList(t *testing.T) {
	env := evaluator.NewEnvironment()
	
	// Add multiple custom units
	env.Units().AddCustomUnit("spoon", 15.0, "ml")
	env.Units().AddCustomUnit("cup", 240.0, "ml")
	env.Units().AddCustomUnit("bowl", 350.0, "ml")
	
	units := env.Units().ListCustomUnits()
	
	if len(units) != 3 {
		t.Errorf("expected 3 custom units, got %d", len(units))
	}
	
	// Check that all units are marked as custom
	for _, u := range units {
		if !u.IsCustom {
			t.Errorf("unit %s should be marked as custom", u.Name)
		}
	}
}

// TestCustomUnitPersistence tests saving and loading custom units
func TestCustomUnitPersistence(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := t.TempDir() + "/custom_units.json"
	
	// Create first system and add units
	s1 := evaluator.NewEnvironment().Units()
	s1.AddCustomUnit("spoon", 15.0, "ml")
	s1.AddCustomUnit("cup", 240.0, "ml")
	
	// Save units
	err := s1.SaveCustomUnits(tmpFile)
	if err != nil {
		t.Fatalf("failed to save custom units: %v", err)
	}
	
	// Create second system and load units
	s2 := evaluator.NewEnvironment().Units()
	err = s2.LoadCustomUnits(tmpFile)
	if err != nil {
		t.Fatalf("failed to load custom units: %v", err)
	}
	
	// Verify units were loaded
	if !s2.UnitExists("spoon") {
		t.Error("spoon unit should exist after loading")
	}
	if !s2.UnitExists("cup") {
		t.Error("cup unit should exist after loading")
	}
	
	// Test conversion works correctly after loading
	result, err := s2.Convert(2.0, "spoon", "ml")
	if err != nil {
		t.Fatalf("conversion error: %v", err)
	}
	if result != 30.0 {
		t.Errorf("expected 30.0, got %.2f", result)
	}
}
