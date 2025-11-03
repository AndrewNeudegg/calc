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
	DimensionData      // Digital storage (bytes, bits)
	DimensionDataRate  // Data transfer rate (bytes/s, bits/s)
	DimensionSpeed     // Speed/velocity (m/s, mph, kph, etc.)
	DimensionPressure  // Pressure (Pa, bar, atm, psi)
	DimensionForce     // Force (N, lbf)
	DimensionAngle     // Angle (degrees, radians, gradians)
	DimensionFrequency // Frequency (Hz, kHz, MHz, GHz)
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
	s.addUnit("µg", DimensionMass, 0.000000001, "kg") // microgram = 1e-9 kg
	s.addUnit("ug", DimensionMass, 0.000000001, "kg")
	s.addUnit("microgram", DimensionMass, 0.000000001, "kg")
	s.addUnit("micrograms", DimensionMass, 0.000000001, "kg")
	s.addUnit("lb", DimensionMass, 0.453592, "kg")
	s.addUnit("lbs", DimensionMass, 0.453592, "kg")
	s.addUnit("pound", DimensionMass, 0.453592, "kg")
	s.addUnit("pounds", DimensionMass, 0.453592, "kg")
	s.addUnit("oz", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounce", DimensionMass, 0.0283495, "kg")
	s.addUnit("ounces", DimensionMass, 0.0283495, "kg")
	s.addUnit("stone", DimensionMass, 6.35029, "kg") // 14 pounds
	s.addUnit("stones", DimensionMass, 6.35029, "kg")
	s.addUnit("st", DimensionMass, 6.35029, "kg")
	s.addUnit("carat", DimensionMass, 0.0002, "kg") // metric carat = 200 mg
	s.addUnit("carats", DimensionMass, 0.0002, "kg")
	s.addUnit("ct", DimensionMass, 0.0002, "kg")
	s.addUnit("troyounce", DimensionMass, 0.0311035, "kg") // troy ounce ≈ 31.1035 g
	s.addUnit("troyounces", DimensionMass, 0.0311035, "kg")
	s.addUnit("troyoz", DimensionMass, 0.0311035, "kg")
	s.addUnit("ozt", DimensionMass, 0.0311035, "kg")
	s.addUnit("tonne", DimensionMass, 1000.0, "kg") // metric ton
	s.addUnit("tonnes", DimensionMass, 1000.0, "kg")
	s.addUnit("ton", DimensionMass, 907.185, "kg") // US short ton (2000 lbs)
	s.addUnit("tons", DimensionMass, 907.185, "kg")

	// Time units (base: second)
	// Fine-grained units
	s.addUnit("ns", DimensionTime, 0.000000001, "s") // nanosecond = 1e-9 s
	s.addUnit("nanosecond", DimensionTime, 0.000000001, "s")
	s.addUnit("nanoseconds", DimensionTime, 0.000000001, "s")
	s.addUnit("µs", DimensionTime, 0.000001, "s") // microsecond = 1e-6 s
	s.addUnit("us", DimensionTime, 0.000001, "s")
	s.addUnit("microsecond", DimensionTime, 0.000001, "s")
	s.addUnit("microseconds", DimensionTime, 0.000001, "s")
	s.addUnit("ms", DimensionTime, 0.001, "s") // millisecond = 1e-3 s
	s.addUnit("millisecond", DimensionTime, 0.001, "s")
	s.addUnit("milliseconds", DimensionTime, 0.001, "s")
	// Standard units
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
	s.addUnit("fortnight", DimensionTime, 1209600.0, "s") // 2 weeks = 14 days
	s.addUnit("fortnights", DimensionTime, 1209600.0, "s")
	s.addUnit("month", DimensionTime, 2629800.0, "s") // average month: 30.4375 days
	s.addUnit("months", DimensionTime, 2629800.0, "s")
	s.addUnit("quarter", DimensionTime, 7889400.0, "s") // 3 months = 91.3125 days
	s.addUnit("quarters", DimensionTime, 7889400.0, "s")
	s.addUnit("semester", DimensionTime, 15778800.0, "s") // 6 months = 182.625 days
	s.addUnit("semesters", DimensionTime, 15778800.0, "s")
	s.addUnit("year", DimensionTime, 31557600.0, "s") // 365.25 days
	s.addUnit("years", DimensionTime, 31557600.0, "s")
	// Special time-of-day unit (stores decimal hours as HH:MM format)
	s.addUnit("time", DimensionTime, 3600.0, "s") // time unit for HH:MM format

	// Volume units (base: litre)
	// Metric
	s.addUnit("l", DimensionVolume, 1.0, "l")
	s.addUnit("litre", DimensionVolume, 1.0, "l")
	s.addUnit("litres", DimensionVolume, 1.0, "l")
	s.addUnit("liter", DimensionVolume, 1.0, "l")
	s.addUnit("liters", DimensionVolume, 1.0, "l")
	s.addUnit("ml", DimensionVolume, 0.001, "l")
	s.addUnit("millilitre", DimensionVolume, 0.001, "l")
	s.addUnit("millilitres", DimensionVolume, 0.001, "l")
	s.addUnit("milliliter", DimensionVolume, 0.001, "l")
	s.addUnit("milliliters", DimensionVolume, 0.001, "l")
	s.addUnit("cl", DimensionVolume, 0.01, "l")
	s.addUnit("centilitre", DimensionVolume, 0.01, "l")
	s.addUnit("centilitres", DimensionVolume, 0.01, "l")
	s.addUnit("centiliter", DimensionVolume, 0.01, "l")
	s.addUnit("centiliters", DimensionVolume, 0.01, "l")
	s.addUnit("dl", DimensionVolume, 0.1, "l")
	s.addUnit("decilitre", DimensionVolume, 0.1, "l")
	s.addUnit("decilitres", DimensionVolume, 0.1, "l")
	s.addUnit("deciliter", DimensionVolume, 0.1, "l")
	s.addUnit("deciliters", DimensionVolume, 0.1, "l")
	// Cubic measures
	s.addUnit("m3", DimensionVolume, 1000.0, "l") // 1 m³ = 1000 litres
	s.addUnit("m³", DimensionVolume, 1000.0, "l")
	s.addUnit("cm3", DimensionVolume, 0.001, "l") // 1 cm³ = 1 ml
	s.addUnit("cm³", DimensionVolume, 0.001, "l")
	s.addUnit("cc", DimensionVolume, 0.001, "l")     // cubic centimeter
	s.addUnit("mm3", DimensionVolume, 0.000001, "l") // 1 mm³ = 0.001 ml
	s.addUnit("mm³", DimensionVolume, 0.000001, "l")
	s.addUnit("ft3", DimensionVolume, 28.3168, "l") // cubic foot
	s.addUnit("ft³", DimensionVolume, 28.3168, "l")
	s.addUnit("in3", DimensionVolume, 0.0163871, "l") // cubic inch
	s.addUnit("in³", DimensionVolume, 0.0163871, "l")
	// US customary
	s.addUnit("usgal", DimensionVolume, 3.78541, "l") // US gallon
	s.addUnit("usgallon", DimensionVolume, 3.78541, "l")
	s.addUnit("usgallons", DimensionVolume, 3.78541, "l")
	s.addUnit("gal", DimensionVolume, 3.78541, "l") // default to US gallon
	s.addUnit("gallon", DimensionVolume, 3.78541, "l")
	s.addUnit("gallons", DimensionVolume, 3.78541, "l")
	s.addUnit("usquart", DimensionVolume, 0.946353, "l") // US quart
	s.addUnit("usquarts", DimensionVolume, 0.946353, "l")
	s.addUnit("quart", DimensionVolume, 0.946353, "l")
	s.addUnit("quarts", DimensionVolume, 0.946353, "l")
	s.addUnit("qt", DimensionVolume, 0.946353, "l")
	s.addUnit("uspint", DimensionVolume, 0.473176, "l") // US pint
	s.addUnit("uspints", DimensionVolume, 0.473176, "l")
	s.addUnit("pint", DimensionVolume, 0.473176, "l")
	s.addUnit("pints", DimensionVolume, 0.473176, "l")
	s.addUnit("pt", DimensionVolume, 0.473176, "l")
	s.addUnit("cup", DimensionVolume, 0.236588, "l") // US cup
	s.addUnit("cups", DimensionVolume, 0.236588, "l")
	s.addUnit("floz", DimensionVolume, 0.0295735, "l") // US fluid ounce
	s.addUnit("fluidounce", DimensionVolume, 0.0295735, "l")
	s.addUnit("fluidounces", DimensionVolume, 0.0295735, "l")
	s.addUnit("tbsp", DimensionVolume, 0.0147868, "l") // US tablespoon
	s.addUnit("tablespoon", DimensionVolume, 0.0147868, "l")
	s.addUnit("tablespoons", DimensionVolume, 0.0147868, "l")
	s.addUnit("tsp", DimensionVolume, 0.00492892, "l") // US teaspoon
	s.addUnit("teaspoon", DimensionVolume, 0.00492892, "l")
	s.addUnit("teaspoons", DimensionVolume, 0.00492892, "l")
	// UK/Imperial
	s.addUnit("ukgal", DimensionVolume, 4.54609, "l") // UK gallon (larger)
	s.addUnit("ukgallon", DimensionVolume, 4.54609, "l")
	s.addUnit("ukgallons", DimensionVolume, 4.54609, "l")
	s.addUnit("impgal", DimensionVolume, 4.54609, "l")
	s.addUnit("imperialgallon", DimensionVolume, 4.54609, "l")
	s.addUnit("ukquart", DimensionVolume, 1.13652, "l") // UK quart
	s.addUnit("ukquarts", DimensionVolume, 1.13652, "l")
	s.addUnit("ukpint", DimensionVolume, 0.568261, "l") // UK pint (larger than US)
	s.addUnit("ukpints", DimensionVolume, 0.568261, "l")
	s.addUnit("imppint", DimensionVolume, 0.568261, "l")
	s.addUnit("imperialpint", DimensionVolume, 0.568261, "l")

	// Area units (base: square meter)
	// Numeric notation
	s.addUnit("sqm", DimensionArea, 1.0, "sqm")
	s.addUnit("m2", DimensionArea, 1.0, "sqm")
	s.addUnit("m²", DimensionArea, 1.0, "sqm")
	s.addUnit("sqmm", DimensionArea, 0.000001, "sqm")
	s.addUnit("mm2", DimensionArea, 0.000001, "sqm")
	s.addUnit("mm²", DimensionArea, 0.000001, "sqm")
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
	// Spelled phrases
	s.addUnit("squaremetre", DimensionArea, 1.0, "sqm")
	s.addUnit("squaremetres", DimensionArea, 1.0, "sqm")
	s.addUnit("squaremeter", DimensionArea, 1.0, "sqm")
	s.addUnit("squaremeters", DimensionArea, 1.0, "sqm")
	s.addUnit("squarefoot", DimensionArea, 0.092903, "sqm")
	s.addUnit("squarefeet", DimensionArea, 0.092903, "sqm")
	s.addUnit("squareinch", DimensionArea, 0.00064516, "sqm")
	s.addUnit("squareinches", DimensionArea, 0.00064516, "sqm")
	s.addUnit("squareyard", DimensionArea, 0.836127, "sqm")
	s.addUnit("squareyards", DimensionArea, 0.836127, "sqm")
	s.addUnit("squaremile", DimensionArea, 2589988.11, "sqm")
	s.addUnit("squaremiles", DimensionArea, 2589988.11, "sqm")
	s.addUnit("squarekilometre", DimensionArea, 1000000.0, "sqm")
	s.addUnit("squarekilometres", DimensionArea, 1000000.0, "sqm")
	s.addUnit("squarekilometer", DimensionArea, 1000000.0, "sqm")
	s.addUnit("squarekilometers", DimensionArea, 1000000.0, "sqm")
	// Land area units
	s.addUnit("acre", DimensionArea, 4046.86, "sqm")
	s.addUnit("acres", DimensionArea, 4046.86, "sqm")
	s.addUnit("hectare", DimensionArea, 10000.0, "sqm")
	s.addUnit("hectares", DimensionArea, 10000.0, "sqm")
	s.addUnit("ha", DimensionArea, 10000.0, "sqm")
	s.addUnit("are", DimensionArea, 100.0, "sqm") // 10m x 10m
	s.addUnit("ares", DimensionArea, 100.0, "sqm")
	s.addUnit("decare", DimensionArea, 1000.0, "sqm") // 10 ares
	s.addUnit("decares", DimensionArea, 1000.0, "sqm")

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
	s.addUnit("r", DimensionTemperature, 1.0, "r")
	s.addUnit("rankine", DimensionTemperature, 1.0, "r")
	s.addUnit("°r", DimensionTemperature, 1.0, "r")

	// Speed units (base: meters per second)
	// These are shortcuts for compound units to avoid needing slashes
	s.addUnit("mps", DimensionSpeed, 1.0, "mps")       // meters per second (base)
	s.addUnit("kph", DimensionSpeed, 0.277778, "mps")  // kilometers per hour: 1 kph = 1000m/3600s
	s.addUnit("kmh", DimensionSpeed, 0.277778, "mps")  // alternative: km/h without slash
	s.addUnit("mph", DimensionSpeed, 0.44704, "mps")   // miles per hour: 1 mph = 1609.34m/3600s
	s.addUnit("fps", DimensionSpeed, 0.3048, "mps")    // feet per second
	s.addUnit("knot", DimensionSpeed, 0.514444, "mps") // nautical miles per hour
	s.addUnit("knots", DimensionSpeed, 0.514444, "mps")
	s.addUnit("kn", DimensionSpeed, 0.514444, "mps")

	// Pressure units (base: Pascal)
	s.addUnit("pa", DimensionPressure, 1.0, "pa")
	s.addUnit("pascal", DimensionPressure, 1.0, "pa")
	s.addUnit("pascals", DimensionPressure, 1.0, "pa")
	s.addUnit("kpa", DimensionPressure, 1000.0, "pa")
	s.addUnit("kilopascal", DimensionPressure, 1000.0, "pa")
	s.addUnit("kilopascals", DimensionPressure, 1000.0, "pa")
	s.addUnit("mpa", DimensionPressure, 1000000.0, "pa")
	s.addUnit("megapascal", DimensionPressure, 1000000.0, "pa")
	s.addUnit("megapascals", DimensionPressure, 1000000.0, "pa")
	s.addUnit("bar", DimensionPressure, 100000.0, "pa") // 1 bar = 100,000 Pa
	s.addUnit("bars", DimensionPressure, 100000.0, "pa")
	s.addUnit("mbar", DimensionPressure, 100.0, "pa") // millibar
	s.addUnit("millibar", DimensionPressure, 100.0, "pa")
	s.addUnit("millibars", DimensionPressure, 100.0, "pa")
	s.addUnit("atm", DimensionPressure, 101325.0, "pa") // standard atmosphere
	s.addUnit("atmosphere", DimensionPressure, 101325.0, "pa")
	s.addUnit("atmospheres", DimensionPressure, 101325.0, "pa")
	s.addUnit("psi", DimensionPressure, 6894.76, "pa")  // pounds per square inch
	s.addUnit("torr", DimensionPressure, 133.322, "pa") // 1/760 atm
	s.addUnit("mmhg", DimensionPressure, 133.322, "pa") // millimeters of mercury
	s.addUnit("inhg", DimensionPressure, 3386.39, "pa") // inches of mercury

	// Force units (base: Newton)
	s.addUnit("n", DimensionForce, 1.0, "n")
	s.addUnit("newton", DimensionForce, 1.0, "n")
	s.addUnit("newtons", DimensionForce, 1.0, "n")
	// Note: "kn" is reserved for knots (speed), use "kilonewton" for force
	s.addUnit("kilonewton", DimensionForce, 1000.0, "n")
	s.addUnit("kilonewtons", DimensionForce, 1000.0, "n")
	s.addUnit("mn", DimensionForce, 1000000.0, "n")
	s.addUnit("meganewton", DimensionForce, 1000000.0, "n")
	s.addUnit("meganewtons", DimensionForce, 1000000.0, "n")
	s.addUnit("lbf", DimensionForce, 4.44822, "n") // pound-force
	s.addUnit("poundforce", DimensionForce, 4.44822, "n")
	s.addUnit("poundsforce", DimensionForce, 4.44822, "n")
	s.addUnit("kgf", DimensionForce, 9.80665, "n") // kilogram-force
	s.addUnit("kilogramforce", DimensionForce, 9.80665, "n")
	s.addUnit("dyne", DimensionForce, 0.00001, "n") // cgs unit
	s.addUnit("dynes", DimensionForce, 0.00001, "n")

	// Angle units (base: degrees)
	s.addUnit("deg", DimensionAngle, 1.0, "deg")
	s.addUnit("degree", DimensionAngle, 1.0, "deg")
	s.addUnit("degrees", DimensionAngle, 1.0, "deg")
	s.addUnit("°", DimensionAngle, 1.0, "deg")
	s.addUnit("rad", DimensionAngle, 57.2958, "deg") // 180/π
	s.addUnit("radian", DimensionAngle, 57.2958, "deg")
	s.addUnit("radians", DimensionAngle, 57.2958, "deg")
	s.addUnit("grad", DimensionAngle, 0.9, "deg") // 360/400
	s.addUnit("gradian", DimensionAngle, 0.9, "deg")
	s.addUnit("gradians", DimensionAngle, 0.9, "deg")
	s.addUnit("gon", DimensionAngle, 0.9, "deg")    // same as gradian
	s.addUnit("turn", DimensionAngle, 360.0, "deg") // full rotation
	s.addUnit("turns", DimensionAngle, 360.0, "deg")
	s.addUnit("revolution", DimensionAngle, 360.0, "deg")
	s.addUnit("revolutions", DimensionAngle, 360.0, "deg")

	// Frequency units (base: Hertz)
	s.addUnit("hz", DimensionFrequency, 1.0, "hz")
	s.addUnit("hertz", DimensionFrequency, 1.0, "hz")
	s.addUnit("khz", DimensionFrequency, 1000.0, "hz")
	s.addUnit("kilohertz", DimensionFrequency, 1000.0, "hz")
	s.addUnit("mhz", DimensionFrequency, 1000000.0, "hz")
	s.addUnit("megahertz", DimensionFrequency, 1000000.0, "hz")
	s.addUnit("ghz", DimensionFrequency, 1000000000.0, "hz")
	s.addUnit("gigahertz", DimensionFrequency, 1000000000.0, "hz")
	s.addUnit("thz", DimensionFrequency, 1000000000000.0, "hz")
	s.addUnit("terahertz", DimensionFrequency, 1000000000000.0, "hz")
	s.addUnit("rpm", DimensionFrequency, 0.0166667, "hz") // revolutions per minute
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
	case "r", "rankine", "°r":
		celsius = (value - 491.67) * 5 / 9 // R to C: (R - 491.67) × 5/9
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
	case "r", "rankine", "°r":
		return (celsius + 273.15) * 9 / 5, nil // C to R: (C + 273.15) × 9/5
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
// This also handles conversions between speed abbreviations (kph, mph) and compound units (km/h, mi/h).
func (s *System) ConvertCompoundUnit(value float64, fromUnit, toUnit string) (float64, error) {
	// If one or both units are speed abbreviations, convert through the base unit (mps)
	fromIsSimple := !IsCompoundUnit(fromUnit)
	toIsSimple := !IsCompoundUnit(toUnit)

	// Case 1: Both are compound units - use normal compound unit conversion
	if !fromIsSimple && !toIsSimple {
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
		result := value * (from.ToBaseNum / from.ToBaseDen) * (to.ToBaseDen / to.ToBaseNum)
		return result, nil
	}

	// Case 2: Simple to compound, or compound to simple - convert through base unit
	// For speed units, the base is mps (meters per second)
	if fromIsSimple && !toIsSimple {
		// Convert from simple (e.g., kph) to base (mps), then to compound (e.g., km/h)
		// First convert simple unit to mps
		inMps, err := s.Convert(value, fromUnit, "mps")
		if err != nil {
			return 0, err
		}
		// Then convert mps to compound unit
		to, err := s.ParseCompoundUnit(toUnit)
		if err != nil {
			return 0, err
		}
		// mps is 1.0 m/s, so we need to convert to the target compound unit
		// result = mps * (targetDen / targetNum)
		result := inMps * (to.ToBaseDen / to.ToBaseNum)
		return result, nil
	}

	if !fromIsSimple && toIsSimple {
		// Convert from compound (e.g., km/h) to base (mps), then to simple (e.g., mph)
		from, err := s.ParseCompoundUnit(fromUnit)
		if err != nil {
			return 0, err
		}
		// Convert compound to mps
		inMps := value * (from.ToBaseNum / from.ToBaseDen)
		// Then convert mps to simple unit
		result, err := s.Convert(inMps, "mps", toUnit)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	// Case 3: Both simple - shouldn't happen in this function, but handle it
	return s.Convert(value, fromUnit, toUnit)
}

// IsCompoundUnit checks if a string looks like a compound unit (contains /).
func IsCompoundUnit(unitStr string) bool {
	return strings.Contains(unitStr, "/")
}
