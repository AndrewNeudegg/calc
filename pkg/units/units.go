package units

import (
	"fmt"
	"strings"
)

// Dimension represents a physical dimension.
type Dimension int

const (
	DimensionNone Dimension = iota
	DimensionLength
	DimensionMass
	DimensionTime
	DimensionTemperature
	DimensionVolume
	DimensionArea
)

// Unit represents a unit of measurement.
type Unit struct {
	Name      string
	Dimension Dimension
	ToBase    float64 // conversion factor to base unit
	BaseUnit  string
	IsCustom  bool
}

// System manages all units and conversions.
type System struct {
	units  map[string]*Unit
	custom map[string]*Unit
}

// NewSystem creates a new unit system.
func NewSystem() *System {
	s := &System{
		units:  make(map[string]*Unit),
		custom: make(map[string]*Unit),
	}
	s.initStandardUnits()
	return s
}

func (s *System) initStandardUnits() {
	// Length units (base: metre)
	s.addUnit("m", DimensionLength, 1.0, "m")
	s.addUnit("metre", DimensionLength, 1.0, "m")
	s.addUnit("metres", DimensionLength, 1.0, "m")
	s.addUnit("meter", DimensionLength, 1.0, "m")
	s.addUnit("meters", DimensionLength, 1.0, "m")
	s.addUnit("cm", DimensionLength, 0.01, "m")
	s.addUnit("mm", DimensionLength, 0.001, "m")
	s.addUnit("km", DimensionLength, 1000.0, "m")
	s.addUnit("ft", DimensionLength, 0.3048, "m")
	s.addUnit("foot", DimensionLength, 0.3048, "m")
	s.addUnit("feet", DimensionLength, 0.3048, "m")
	s.addUnit("in", DimensionLength, 0.0254, "m")
	s.addUnit("inch", DimensionLength, 0.0254, "m")
	s.addUnit("inches", DimensionLength, 0.0254, "m")
	s.addUnit("yd", DimensionLength, 0.9144, "m")
	s.addUnit("yard", DimensionLength, 0.9144, "m")
	s.addUnit("yards", DimensionLength, 0.9144, "m")
	s.addUnit("mi", DimensionLength, 1609.344, "m")
	s.addUnit("mile", DimensionLength, 1609.344, "m")
	s.addUnit("miles", DimensionLength, 1609.344, "m")

	// Mass units (base: kilogram)
	s.addUnit("kg", DimensionMass, 1.0, "kg")
	s.addUnit("kilogram", DimensionMass, 1.0, "kg")
	s.addUnit("kilograms", DimensionMass, 1.0, "kg")
	s.addUnit("g", DimensionMass, 0.001, "kg")
	s.addUnit("gram", DimensionMass, 0.001, "kg")
	s.addUnit("grams", DimensionMass, 0.001, "kg")
	s.addUnit("mg", DimensionMass, 0.000001, "kg")
	s.addUnit("lb", DimensionMass, 0.453592, "kg")
	s.addUnit("pound", DimensionMass, 0.453592, "kg")
	s.addUnit("pounds", DimensionMass, 0.453592, "kg")
	s.addUnit("oz", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounce", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounces", DimensionMass, 0.0283495, "kg")

	// Time units (base: second)
	s.addUnit("s", DimensionTime, 1.0, "s")
	s.addUnit("sec", DimensionTime, 1.0, "s")
	s.addUnit("second", DimensionTime, 1.0, "s")
	s.addUnit("seconds", DimensionTime, 1.0, "s")
	s.addUnit("min", DimensionTime, 60.0, "s")
	s.addUnit("minute", DimensionTime, 60.0, "s")
	s.addUnit("minutes", DimensionTime, 60.0, "s")
	s.addUnit("h", DimensionTime, 3600.0, "s")
	s.addUnit("hr", DimensionTime, 3600.0, "s")
	s.addUnit("hour", DimensionTime, 3600.0, "s")
	s.addUnit("hours", DimensionTime, 3600.0, "s")
	s.addUnit("day", DimensionTime, 86400.0, "s")
	s.addUnit("days", DimensionTime, 86400.0, "s")
	s.addUnit("week", DimensionTime, 604800.0, "s")
	s.addUnit("weeks", DimensionTime, 604800.0, "s")
	s.addUnit("year", DimensionTime, 31557600.0, "s") // 365.25 days
	s.addUnit("years", DimensionTime, 31557600.0, "s")

	// Volume units (base: litre)
	s.addUnit("l", DimensionVolume, 1.0, "l")
	s.addUnit("litre", DimensionVolume, 1.0, "l")
	s.addUnit("litres", DimensionVolume, 1.0, "l")
	s.addUnit("liter", DimensionVolume, 1.0, "l")
	s.addUnit("liters", DimensionVolume, 1.0, "l")
	s.addUnit("ml", DimensionVolume, 0.001, "l")
	s.addUnit("gal", DimensionVolume, 3.78541, "l")
	s.addUnit("gallon", DimensionVolume, 3.78541, "l")
	s.addUnit("gallons", DimensionVolume, 3.78541, "l")

	// Temperature units (special handling needed)
	s.addUnit("c", DimensionTemperature, 1.0, "c")
	s.addUnit("celsius", DimensionTemperature, 1.0, "c")
	s.addUnit("f", DimensionTemperature, 1.0, "f")
	s.addUnit("fahrenheit", DimensionTemperature, 1.0, "f")
}

func (s *System) addUnit(name string, dim Dimension, toBase float64, baseUnit string) {
	s.units[strings.ToLower(name)] = &Unit{
		Name:      name,
		Dimension: dim,
		ToBase:    toBase,
		BaseUnit:  baseUnit,
		IsCustom:  false,
	}
}

// AddCustomUnit adds a custom unit definition.
func (s *System) AddCustomUnit(name string, value float64, baseUnit string) error {
	name = strings.ToLower(name)

	// Check if base unit exists
	base, exists := s.units[strings.ToLower(baseUnit)]
	if !exists {
		return fmt.Errorf("unknown base unit: %s", baseUnit)
	}

	// Check for circular definition
	if name == strings.ToLower(baseUnit) {
		return fmt.Errorf("circular unit definition")
	}

	s.custom[name] = &Unit{
		Name:      name,
		Dimension: base.Dimension,
		ToBase:    value * base.ToBase,
		BaseUnit:  base.BaseUnit,
		IsCustom:  true,
	}

	s.units[name] = s.custom[name]

	return nil
}

// Convert converts a value from one unit to another.
func (s *System) Convert(value float64, fromUnit, toUnit string) (float64, error) {
	fromUnit = strings.ToLower(fromUnit)
	toUnit = strings.ToLower(toUnit)

	from, ok := s.units[fromUnit]
	if !ok {
		return 0, fmt.Errorf("unknown unit '%s'", fromUnit)
	}

	to, ok := s.units[toUnit]
	if !ok {
		return 0, fmt.Errorf("unknown unit '%s'", toUnit)
	}

	// Check dimension compatibility
	if from.Dimension != to.Dimension {
		return 0, fmt.Errorf("cannot convert %s to %s", fromUnit, toUnit)
	}

	// Handle temperature specially
	if from.Dimension == DimensionTemperature {
		return s.convertTemperature(value, fromUnit, toUnit)
	}

	// Convert to base unit, then to target unit
	baseValue := value * from.ToBase
	result := baseValue / to.ToBase

	return result, nil
}

func (s *System) convertTemperature(value float64, from, to string) (float64, error) {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	// Normalise to celsius first
	var celsius float64
	switch from {
	case "c", "celsius":
		celsius = value
	case "f", "fahrenheit":
		celsius = (value - 32) * 5 / 9
	default:
		return 0, fmt.Errorf("unknown temperature unit: %s", from)
	}

	// Convert from celsius to target
	switch to {
	case "c", "celsius":
		return celsius, nil
	case "f", "fahrenheit":
		return celsius*9/5 + 32, nil
	default:
		return 0, fmt.Errorf("unknown temperature unit: %s", to)
	}
}

// IsUnit checks if a string is a known unit.
func (s *System) IsUnit(name string) bool {
	_, ok := s.units[strings.ToLower(name)]
	return ok
}

// GetDimension returns the dimension of a unit.
func (s *System) GetDimension(name string) (Dimension, error) {
	unit, ok := s.units[strings.ToLower(name)]
	if !ok {
		return DimensionNone, fmt.Errorf("unknown unit: %s", name)
	}
	return unit.Dimension, nil
}
