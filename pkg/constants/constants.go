package constants

import (
	"fmt"
	"strings"

	"github.com/andrewneudegg/calc/pkg/units"
)

// Constant represents a physical constant with a value, unit, and description.
type Constant struct {
	Name        string
	Symbol      string
	Value       float64
	Unit        string
	Dimension   units.Dimension
	Description string
	Category    string
}

// System manages all physical constants.
type System struct {
	constants map[string]*Constant
}

// NewSystem creates a new constants system with all standard constants loaded.
func NewSystem() *System {
	s := &System{
		constants: make(map[string]*Constant),
	}
	// Initialize all constant categories
	s.initFundamental()
	s.initElectromagnetic()
	s.initUniversal()
	return s
}

// AddConstant adds a constant to the system.
func (s *System) addConstant(name, symbol string, value float64, unit string, dim units.Dimension, description, category string) {
	c := &Constant{
		Name:        name,
		Symbol:      symbol,
		Value:       value,
		Unit:        unit,
		Dimension:   dim,
		Description: description,
		Category:    category,
	}
	
	// Register by both name and symbol (case-insensitive)
	key := strings.ToLower(name)
	s.constants[key] = c
	
	if symbol != "" && symbol != name {
		symKey := strings.ToLower(symbol)
		s.constants[symKey] = c
	}
}

// IsConstant checks if a string is a known constant.
func (s *System) IsConstant(name string) bool {
	_, ok := s.constants[strings.ToLower(name)]
	return ok
}

// GetConstant retrieves a constant by name or symbol.
func (s *System) GetConstant(name string) (*Constant, error) {
	c, ok := s.constants[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown constant: %s", name)
	}
	return c, nil
}

// ListConstants returns all constants in the system.
func (s *System) ListConstants() []*Constant {
	// Use a map to deduplicate (same constant may have multiple keys)
	seen := make(map[*Constant]bool)
	result := make([]*Constant, 0)
	
	for _, c := range s.constants {
		if !seen[c] {
			seen[c] = true
			result = append(result, c)
		}
	}
	
	return result
}

// ListByCategory returns all constants in a specific category.
func (s *System) ListByCategory(category string) []*Constant {
	seen := make(map[*Constant]bool)
	result := make([]*Constant, 0)
	
	for _, c := range s.constants {
		if !seen[c] && strings.EqualFold(c.Category, category) {
			seen[c] = true
			result = append(result, c)
		}
	}
	
	return result
}

// GetCategories returns all available categories.
func (s *System) GetCategories() []string {
	categorySet := make(map[string]bool)
	
	for _, c := range s.constants {
		categorySet[c.Category] = true
	}
	
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	
	return categories
}
