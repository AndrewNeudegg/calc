package constants

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/units"
)

func TestNewSystem(t *testing.T) {
	sys := NewSystem()
	if sys == nil {
		t.Fatal("NewSystem returned nil")
	}
	
	if len(sys.constants) == 0 {
		t.Error("NewSystem created empty constants map")
	}
}

func TestIsConstant(t *testing.T) {
	sys := NewSystem()
	
	tests := []struct {
		name     string
		expected bool
	}{
		{"c", true},
		{"speed_of_light", true},
		{"h", true},
		{"planck", true},
		{"G", true},
		{"gravitational_constant", true},
		{"e", true},
		{"elementary_charge", true},
		{"σ", true},
		{"stefan_boltzmann", true},
		{"unknown", false},
		{"xyz", false},
		{"C", true}, // case insensitive
		{"H", true}, // case insensitive
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sys.IsConstant(tt.name)
			if result != tt.expected {
				t.Errorf("IsConstant(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetConstant(t *testing.T) {
	sys := NewSystem()
	
	tests := []struct {
		name        string
		shouldError bool
		checkValue  bool
		expected    float64
	}{
		{"c", false, true, 299792458.0},
		{"speed_of_light", false, true, 299792458.0},
		{"h", false, true, 6.62607015e-34},
		{"planck", false, true, 6.62607015e-34},
		{"G", false, true, 6.67430e-11},
		{"e", false, true, 1.602176634e-19},
		{"unknown", true, false, 0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := sys.GetConstant(tt.name)
			
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("GetConstant(%q) returned error: %v", tt.name, err)
			}
			
			if c == nil {
				t.Fatal("GetConstant returned nil constant")
			}
			
			if tt.checkValue && c.Value != tt.expected {
				t.Errorf("GetConstant(%q).Value = %e, want %e", tt.name, c.Value, tt.expected)
			}
		})
	}
}

func TestConstantFields(t *testing.T) {
	sys := NewSystem()
	
	c, err := sys.GetConstant("c")
	if err != nil {
		t.Fatalf("GetConstant('c') failed: %v", err)
	}
	
	if c.Name != "speed_of_light" {
		t.Errorf("c.Name = %q, want 'speed_of_light'", c.Name)
	}
	
	if c.Symbol != "c" {
		t.Errorf("c.Symbol = %q, want 'c'", c.Symbol)
	}
	
	if c.Value != 299792458.0 {
		t.Errorf("c.Value = %e, want 299792458.0", c.Value)
	}
	
	if c.Unit != "m/s" {
		t.Errorf("c.Unit = %q, want 'm/s'", c.Unit)
	}
	
	if c.Dimension != units.DimensionSpeed {
		t.Errorf("c.Dimension = %v, want DimensionSpeed", c.Dimension)
	}
	
	if c.Category != "fundamental" {
		t.Errorf("c.Category = %q, want 'fundamental'", c.Category)
	}
	
	if c.Description == "" {
		t.Error("c.Description is empty")
	}
}

func TestListConstants(t *testing.T) {
	sys := NewSystem()
	
	constants := sys.ListConstants()
	
	if len(constants) == 0 {
		t.Error("ListConstants returned empty list")
	}
	
	// Check that we have at least the major constants
	names := make(map[string]bool)
	for _, c := range constants {
		names[c.Name] = true
	}
	
	required := []string{
		"speed_of_light",
		"planck",
		"gravitational_constant",
		"elementary_charge",
		"stefan_boltzmann",
	}
	
	for _, name := range required {
		if !names[name] {
			t.Errorf("ListConstants missing required constant: %s", name)
		}
	}
}

func TestListByCategory(t *testing.T) {
	sys := NewSystem()
	
	tests := []struct {
		category string
		minCount int
	}{
		{"fundamental", 5},
		{"electromagnetic", 3},
		{"universal", 5},
	}
	
	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			constants := sys.ListByCategory(tt.category)
			
			if len(constants) < tt.minCount {
				t.Errorf("ListByCategory(%q) returned %d constants, want at least %d",
					tt.category, len(constants), tt.minCount)
			}
			
			// Verify all returned constants are in the correct category
			for _, c := range constants {
				if c.Category != tt.category {
					t.Errorf("Constant %s has category %q, want %q",
						c.Name, c.Category, tt.category)
				}
			}
		})
	}
}

func TestGetCategories(t *testing.T) {
	sys := NewSystem()
	
	categories := sys.GetCategories()
	
	if len(categories) == 0 {
		t.Error("GetCategories returned empty list")
	}
	
	// Check for expected categories
	expected := map[string]bool{
		"fundamental":     false,
		"electromagnetic": false,
		"universal":       false,
	}
	
	for _, cat := range categories {
		if _, ok := expected[cat]; ok {
			expected[cat] = true
		}
	}
	
	for cat, found := range expected {
		if !found {
			t.Errorf("GetCategories missing expected category: %s", cat)
		}
	}
}

func TestCaseInsensitivity(t *testing.T) {
	sys := NewSystem()
	
	// Test that both uppercase and lowercase work
	tests := [][]string{
		{"c", "C"},
		{"h", "H"},
		{"g", "G"},
		{"planck", "PLANCK", "Planck"},
	}
	
	for _, variants := range tests {
		var firstConstant *Constant
		
		for i, variant := range variants {
			c, err := sys.GetConstant(variant)
			if err != nil {
				t.Errorf("GetConstant(%q) failed: %v", variant, err)
				continue
			}
			
			if i == 0 {
				firstConstant = c
			} else {
				// Verify all variants return the same constant
				if c.Name != firstConstant.Name {
					t.Errorf("GetConstant(%q).Name = %q, want %q (same as %q)",
						variant, c.Name, firstConstant.Name, variants[0])
				}
			}
		}
	}
}

func TestSymbolAccess(t *testing.T) {
	sys := NewSystem()
	
	// Test access by symbol vs name returns same constant
	tests := []struct {
		symbol string
		name   string
	}{
		{"c", "speed_of_light"},
		{"h", "planck"},
		{"G", "gravitational_constant"},
		{"e", "elementary_charge"},
		{"σ", "stefan_boltzmann"},
	}
	
	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			bySymbol, err := sys.GetConstant(tt.symbol)
			if err != nil {
				t.Fatalf("GetConstant(%q) failed: %v", tt.symbol, err)
			}
			
			byName, err := sys.GetConstant(tt.name)
			if err != nil {
				t.Fatalf("GetConstant(%q) failed: %v", tt.name, err)
			}
			
			if bySymbol.Name != byName.Name {
				t.Errorf("Symbol %q and name %q returned different constants: %q vs %q",
					tt.symbol, tt.name, bySymbol.Name, byName.Name)
			}
		})
	}
}
