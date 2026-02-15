package units

import (
	"testing"
)

func TestAddCustomUnit(t *testing.T) {
	s := NewSystem()
	
	// Test adding a custom unit
	err := s.AddCustomUnit("spoon", 15.0, "ml")
	if err != nil {
		t.Fatalf("failed to add custom unit: %v", err)
	}
	
	// Check if unit exists
	if !s.UnitExists("spoon") {
		t.Errorf("custom unit 'spoon' should exist")
	}
	
	// Test conversion
	result, err := s.Convert(2.0, "spoon", "ml")
	if err != nil {
		t.Fatalf("failed to convert: %v", err)
	}
	
	expected := 30.0
	if result != expected {
		t.Errorf("expected 2 spoons = %.2f ml, got %.2f ml", expected, result)
	}
}

func TestAddCustomUnitWithUnknownBase(t *testing.T) {
	s := NewSystem()
	
	err := s.AddCustomUnit("foo", 10.0, "unknownunit")
	if err == nil {
		t.Errorf("expected error when adding custom unit with unknown base unit")
	}
}

func TestListCustomUnits(t *testing.T) {
	s := NewSystem()
	
	s.AddCustomUnit("spoon", 15.0, "ml")
	s.AddCustomUnit("cup", 240.0, "ml")
	
	units := s.ListCustomUnits()
	if len(units) != 2 {
		t.Errorf("expected 2 custom units, got %d", len(units))
	}
}

func TestDeleteCustomUnit(t *testing.T) {
	s := NewSystem()
	
	s.AddCustomUnit("spoon", 15.0, "ml")
	
	// Delete the unit
	err := s.DeleteCustomUnit("spoon")
	if err != nil {
		t.Fatalf("failed to delete custom unit: %v", err)
	}
	
	// Check it's gone
	if s.UnitExists("spoon") {
		t.Errorf("custom unit 'spoon' should not exist after deletion")
	}
}

func TestDeleteNonExistentCustomUnit(t *testing.T) {
	s := NewSystem()
	
	err := s.DeleteCustomUnit("nonexistent")
	if err == nil {
		t.Errorf("expected error when deleting non-existent custom unit")
	}
}
