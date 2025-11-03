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
	DimensionData     // Digital storage (bytes, bits)
	DimensionDataRate // Data transfer rate (bytes/s, bits/s)
)

// Unit represents a unit of measurement.
type Unit struct {
	Name      string
	Dimension Dimension
	ToBase    float64 // conversion factor to base unit
	BaseUnit  string
	IsCustom  bool
}

// CompoundUnit represents a compound unit like km/h or m/s.
type CompoundUnit struct {
	Numerator   *Unit
	Denominator *Unit
	ToBaseNum   float64 // conversion factor for numerator to base
	ToBaseDen   float64 // conversion factor for denominator to base
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
	s.addUnit("milligram", DimensionMass, 0.000001, "kg")
	s.addUnit("milligrams", DimensionMass, 0.000001, "kg")
	s.addUnit("lb", DimensionMass, 0.453592, "kg")
	s.addUnit("lbs", DimensionMass, 0.453592, "kg")
	s.addUnit("pound", DimensionMass, 0.453592, "kg")
	s.addUnit("pounds", DimensionMass, 0.453592, "kg")
	s.addUnit("oz", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounce", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounces", DimensionMass, 0.0283495, "kg")
	s.addUnit("stone", DimensionMass, 6.35029, "kg") // 14 pounds
	s.addUnit("st", DimensionMass, 6.35029, "kg")
	s.addUnit("tonne", DimensionMass, 1000.0, "kg") // metric ton
	s.addUnit("tonnes", DimensionMass, 1000.0, "kg")
	s.addUnit("ton", DimensionMass, 907.185, "kg") // US short ton (2000 lbs)
	s.addUnit("tons", DimensionMass, 907.185, "kg")

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

	// Area units (base: square meter)
	s.addUnit("sqm", DimensionArea, 1.0, "sqm")
	s.addUnit("m2", DimensionArea, 1.0, "sqm")
	s.addUnit("m²", DimensionArea, 1.0, "sqm")
	s.addUnit("sqcm", DimensionArea, 0.0001, "sqm")
	s.addUnit("cm2", DimensionArea, 0.0001, "sqm")
	s.addUnit("cm²", DimensionArea, 0.0001, "sqm")
	s.addUnit("sqkm", DimensionArea, 1000000.0, "sqm")
	s.addUnit("km2", DimensionArea, 1000000.0, "sqm")
	s.addUnit("km²", DimensionArea, 1000000.0, "sqm")
	s.addUnit("sqft", DimensionArea, 0.092903, "sqm")
	s.addUnit("ft2", DimensionArea, 0.092903, "sqm")
	s.addUnit("ft²", DimensionArea, 0.092903, "sqm")
	s.addUnit("sqin", DimensionArea, 0.00064516, "sqm")
	s.addUnit("in2", DimensionArea, 0.00064516, "sqm")
	s.addUnit("in²", DimensionArea, 0.00064516, "sqm")
	s.addUnit("sqyd", DimensionArea, 0.836127, "sqm")
	s.addUnit("yd2", DimensionArea, 0.836127, "sqm")
	s.addUnit("yd²", DimensionArea, 0.836127, "sqm")
	s.addUnit("sqmi", DimensionArea, 2589988.11, "sqm")
	s.addUnit("mi2", DimensionArea, 2589988.11, "sqm")
	s.addUnit("mi²", DimensionArea, 2589988.11, "sqm")
	s.addUnit("acre", DimensionArea, 4046.86, "sqm")
	s.addUnit("acres", DimensionArea, 4046.86, "sqm")
	s.addUnit("hectare", DimensionArea, 10000.0, "sqm")
	s.addUnit("hectares", DimensionArea, 10000.0, "sqm")
	s.addUnit("ha", DimensionArea, 10000.0, "sqm")

	// Digital storage units (base: byte)
	// Bytes
	s.addUnit("b", DimensionData, 1.0, "b")
	s.addUnit("byte", DimensionData, 1.0, "b")
	s.addUnit("bytes", DimensionData, 1.0, "b")
	s.addUnit("kb", DimensionData, 1024.0, "b")
	s.addUnit("kilobyte", DimensionData, 1024.0, "b")
	s.addUnit("kilobytes", DimensionData, 1024.0, "b")
	s.addUnit("mb", DimensionData, 1048576.0, "b") // 1024^2
	s.addUnit("megabyte", DimensionData, 1048576.0, "b")
	s.addUnit("megabytes", DimensionData, 1048576.0, "b")
	s.addUnit("gb", DimensionData, 1073741824.0, "b") // 1024^3
	s.addUnit("gigabyte", DimensionData, 1073741824.0, "b")
	s.addUnit("gigabytes", DimensionData, 1073741824.0, "b")
	s.addUnit("tb", DimensionData, 1099511627776.0, "b") // 1024^4
	s.addUnit("terabyte", DimensionData, 1099511627776.0, "b")
	s.addUnit("terabytes", DimensionData, 1099511627776.0, "b")
	s.addUnit("pb", DimensionData, 1125899906842624.0, "b") // 1024^5
	s.addUnit("petabyte", DimensionData, 1125899906842624.0, "b")
	s.addUnit("petabytes", DimensionData, 1125899906842624.0, "b")

	// Bits
	s.addUnit("bit", DimensionData, 0.125, "b") // 1 bit = 1/8 byte
	s.addUnit("bits", DimensionData, 0.125, "b")
	s.addUnit("kbit", DimensionData, 128.0, "b") // 1024 bits / 8
	s.addUnit("kilobit", DimensionData, 128.0, "b")
	s.addUnit("kilobits", DimensionData, 128.0, "b")
	s.addUnit("mbit", DimensionData, 131072.0, "b") // 1024^2 bits / 8
	s.addUnit("megabit", DimensionData, 131072.0, "b")
	s.addUnit("megabits", DimensionData, 131072.0, "b")
	s.addUnit("gbit", DimensionData, 134217728.0, "b") // 1024^3 bits / 8
	s.addUnit("gigabit", DimensionData, 134217728.0, "b")
	s.addUnit("gigabits", DimensionData, 134217728.0, "b")
	s.addUnit("tbit", DimensionData, 137438953472.0, "b") // 1024^4 bits / 8
	s.addUnit("terabit", DimensionData, 137438953472.0, "b")
	s.addUnit("terabits", DimensionData, 137438953472.0, "b")
	s.addUnit("pbit", DimensionData, 140737488355328.0, "b") // 1024^5 bits / 8
	s.addUnit("petabit", DimensionData, 140737488355328.0, "b")
	s.addUnit("petabits", DimensionData, 140737488355328.0, "b")

	// Data rate units (base: bytes per second)
	// Bytes per second
	s.addUnit("bps", DimensionDataRate, 1.0, "bps")
	s.addUnit("kbps", DimensionDataRate, 1024.0, "bps")
	s.addUnit("mbps", DimensionDataRate, 1048576.0, "bps")
	s.addUnit("gbps", DimensionDataRate, 1073741824.0, "bps")
	s.addUnit("tbps", DimensionDataRate, 1099511627776.0, "bps")

	// Uppercase variants (common in networking)
	s.addUnit("Bps", DimensionDataRate, 1.0, "bps")
	s.addUnit("KBps", DimensionDataRate, 1024.0, "bps")
	s.addUnit("MBps", DimensionDataRate, 1048576.0, "bps")
	s.addUnit("GBps", DimensionDataRate, 1073741824.0, "bps")
	s.addUnit("TBps", DimensionDataRate, 1099511627776.0, "bps")

	// Bits per second (lowercase)
	s.addUnit("bitps", DimensionDataRate, 0.125, "bps")
	s.addUnit("kbitps", DimensionDataRate, 128.0, "bps")
	s.addUnit("mbitps", DimensionDataRate, 131072.0, "bps")
	s.addUnit("gbitps", DimensionDataRate, 134217728.0, "bps")
	s.addUnit("tbitps", DimensionDataRate, 137438953472.0, "bps")

	// Temperature units (special handling needed)
	s.addUnit("c", DimensionTemperature, 1.0, "c")
	s.addUnit("celsius", DimensionTemperature, 1.0, "c")
	s.addUnit("f", DimensionTemperature, 1.0, "f")
	s.addUnit("fahrenheit", DimensionTemperature, 1.0, "f")
	s.addUnit("k", DimensionTemperature, 1.0, "k")
	s.addUnit("kelvin", DimensionTemperature, 1.0, "k")
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
	case "k", "kelvin":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unknown temperature unit: %s", from)
	}

	// Convert from celsius to target
	switch to {
	case "c", "celsius":
		return celsius, nil
	case "f", "fahrenheit":
		return celsius*9/5 + 32, nil
	case "k", "kelvin":
		return celsius + 273.15, nil
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

// ParseCompoundUnit parses a compound unit string like "km/h" or "m/s".
func (s *System) ParseCompoundUnit(unitStr string) (*CompoundUnit, error) {
	unitStr = strings.ToLower(unitStr)
	parts := strings.Split(unitStr, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid compound unit format: %s", unitStr)
	}

	numStr := strings.TrimSpace(parts[0])
	denStr := strings.TrimSpace(parts[1])

	num, ok := s.units[numStr]
	if !ok {
		return nil, fmt.Errorf("unknown numerator unit: %s", numStr)
	}

	den, ok := s.units[denStr]
	if !ok {
		return nil, fmt.Errorf("unknown denominator unit: %s", denStr)
	}

	return &CompoundUnit{
		Numerator:   num,
		Denominator: den,
		ToBaseNum:   num.ToBase,
		ToBaseDen:   den.ToBase,
	}, nil
}

// ConvertCompoundUnit converts a value from one compound unit to another.
// Example: Convert 50 km/h to m/s
func (s *System) ConvertCompoundUnit(value float64, fromUnit, toUnit string) (float64, error) {
	from, err := s.ParseCompoundUnit(fromUnit)
	if err != nil {
		return 0, err
	}

	to, err := s.ParseCompoundUnit(toUnit)
	if err != nil {
		return 0, err
	}

	// Check dimension compatibility
	if from.Numerator.Dimension != to.Numerator.Dimension {
		return 0, fmt.Errorf("incompatible numerator dimensions: %s vs %s",
			fromUnit, toUnit)
	}

	if from.Denominator.Dimension != to.Denominator.Dimension {
		return 0, fmt.Errorf("incompatible denominator dimensions: %s vs %s",
			fromUnit, toUnit)
	}

	// Convert: value * (fromNum/fromDen) * (toDen/toNum)
	// Example: 50 km/h = 50 * (1000m/3600s) * (1s/1m) = 50 * 1000/3600 = 13.89 m/s
	result := value * (from.ToBaseNum / from.ToBaseDen) * (to.ToBaseDen / to.ToBaseNum)

	return result, nil
}

// IsCompoundUnit checks if a string looks like a compound unit (contains /).
func IsCompoundUnit(unitStr string) bool {
	return strings.Contains(unitStr, "/")
}
