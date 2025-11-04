package integration

import (
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// Helper function to evaluate an expression
func evalExpr(input string) evaluator.Value {
	l := lexer.New(input)
	tokens := l.AllTokens()
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return evaluator.NewError(err.Error())
	}
	env := evaluator.NewEnvironment()
	e := evaluator.New(env)
	return e.Eval(expr)
}

// TestComprehensiveExamples tests EVERY possible statement type in calc
func TestComprehensiveExamples(t *testing.T) {
	tests := []struct {
		category    string
		description string
		input       string
		expectType  evaluator.ValueType
		checkValue  func(evaluator.Value) bool
	}{
		// ==================== BASIC ARITHMETIC ====================
		{
			category:    "Basic Arithmetic",
			description: "Simple addition",
			input:       "5 + 3",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Simple subtraction",
			input:       "10 - 4",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 6 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Simple multiplication",
			input:       "6 * 7",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 42 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Simple division",
			input:       "20 / 4",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Operator precedence (multiplication before addition)",
			input:       "2 + 3 * 4",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 14 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Parentheses override precedence",
			input:       "(2 + 3) * 4",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Nested parentheses",
			input:       "((2 + 3) * (4 + 5)) / 3",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 15 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Unary minus",
			input:       "-5",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == -5 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Unary minus in expression",
			input:       "10 + -5",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 },
		},
		{
			category:    "Basic Arithmetic",
			description: "Decimal numbers",
			input:       "3.14 * 2",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 6.2 && v.Number < 6.3 },
		},

		// ==================== VARIABLE ASSIGNMENT ====================
		{
			category:    "Variables",
			description: "Simple assignment",
			input:       "x = 10",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},
		{
			category:    "Variables",
			description: "Assignment with expression",
			input:       "y = 5 + 3",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 },
		},

		// ==================== PERCENTAGES ====================
		{
			category:    "Percentages",
			description: "Percentage as decimal",
			input:       "20%",
			expectType:  evaluator.ValuePercent,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 },
		},
		{
			category:    "Percentages",
			description: "Add percentage to number",
			input:       "100 + 10%",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 110 },
		},
		{
			category:    "Percentages",
			description: "Subtract percentage from number",
			input:       "100 - 10%",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 90 },
		},
		{
			category:    "Percentages",
			description: "X% of Y",
			input:       "20% of 50",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},
		{
			category:    "Percentages",
			description: "Increase by percentage",
			input:       "increase 100 by 10%",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 109.9 && v.Number < 110.1 },
		},
		{
			category:    "Percentages",
			description: "Decrease by percentage",
			input:       "decrease 100 by 10%",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 90 },
		},
		{
			category:    "Percentages",
			description: "What percent is X of Y",
			input:       "20 is what % of 50",
			expectType:  evaluator.ValuePercent,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 40 },
		},

		// ==================== CURRENCY ====================
		{
			category:    "Currency",
			description: "Currency symbol prefix (GBP)",
			input:       "£12",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 12 && v.Currency == "£" },
		},
		{
			category:    "Currency",
			description: "Currency symbol prefix (USD)",
			input:       "$50",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 && v.Currency == "$" },
		},
		{
			category:    "Currency",
			description: "Currency symbol prefix (EUR)",
			input:       "€100",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Currency == "€" },
		},
		{
			category:    "Currency",
			description: "Currency symbol prefix (JPY)",
			input:       "¥1000",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Currency == "¥" },
		},
		{
			category:    "Currency",
			description: "Currency code postfix (GBP)",
			input:       "12 gbp",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 12 && v.Currency == "£" },
		},
		{
			category:    "Currency",
			description: "Currency code postfix (USD)",
			input:       "50 usd",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 && v.Currency == "$" },
		},
		{
			category:    "Currency",
			description: "Currency name postfix (dollars)",
			input:       "50 dollars",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 && v.Currency == "$" },
		},
		{
			category:    "Currency",
			description: "Currency name postfix (euros)",
			input:       "25 euros",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 25 && v.Currency == "€" },
		},
		{
			category:    "Currency",
			description: "Currency name postfix (yen)",
			input:       "1000 yen",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Currency == "¥" },
		},
		{
			category:    "Currency",
			description: "Currency conversion",
			input:       "12 gbp in usd",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Currency == "$" && v.Number > 0 },
		},
		{
			category:    "Currency",
			description: "Currency addition",
			input:       "£10 + £5",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 15 && v.Currency == "£" },
		},
		{
			category:    "Currency",
			description: "Currency subtraction",
			input:       "$100 - $20",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 80 && v.Currency == "$" },
		},
		{
			category:    "Currency",
			description: "Currency multiplication by number",
			input:       "£10 * 3",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 30 && v.Currency == "£" },
		},
		{
			category:    "Currency",
			description: "Currency division by number",
			input:       "$100 / 4",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 25 && v.Currency == "$" },
		},

		// ==================== UNITS - LENGTH ====================
		{
			category:    "Units - Length",
			description: "Metres with unit",
			input:       "10 m",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 && v.Unit == "m" },
		},
		{
			category:    "Units - Length",
			description: "Convert metres to centimetres",
			input:       "10 m in cm",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Unit == "cm" },
		},
		{
			category:    "Units - Length",
			description: "Convert feet to metres",
			input:       "10 feet in metres",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 3 && v.Number < 4 && v.Unit == "metres" },
		},
		{
			category:    "Units - Length",
			description: "Convert inches to centimetres",
			input:       "12 inches in cm",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 30 && v.Number < 31 && v.Unit == "cm" },
		},
		{
			category:    "Units - Length",
			description: "Convert miles to kilometres",
			input:       "5 miles in km",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 8 && v.Number < 9 && v.Unit == "km" },
		},
		{
			category:    "Units - Length",
			description: "Add same units",
			input:       "5 m + 3 m",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 && v.Unit == "m" },
		},
		{
			category:    "Units - Length",
			description: "Subtract same units",
			input:       "10 km - 3 km",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 7 && v.Unit == "km" },
		},

		// ==================== UNITS - MASS ====================
		{
			category:    "Units - Mass",
			description: "Kilograms",
			input:       "5 kg",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "kg" },
		},
		{
			category:    "Units - Mass",
			description: "Convert kg to grams",
			input:       "2 kg in g",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2000 && v.Unit == "g" },
		},
		{
			category:    "Units - Mass",
			description: "Convert pounds to kg",
			input:       "10 lb in kg",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 4 && v.Number < 5 && v.Unit == "kg" },
		},
		{
			category:    "Units - Mass",
			description: "Ounces",
			input:       "16 oz",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 16 && v.Unit == "oz" },
		},
		{
			category:    "Units - Mass",
			description: "Stone (British)",
			input:       "10 stone",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 && v.Unit == "stone" },
		},
		{
			category:    "Units - Mass",
			description: "Carats (jewellery)",
			input:       "5 carats in grams",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 && v.Unit == "grams" },
		},

		// ==================== UNITS - TIME ====================
		{
			category:    "Units - Time",
			description: "Seconds",
			input:       "60 seconds",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 60 && v.Unit == "seconds" },
		},
		{
			category:    "Units - Time",
			description: "Convert seconds to minutes",
			input:       "120 seconds in minutes",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2 && v.Unit == "minutes" },
		},
		{
			category:    "Units - Time",
			description: "Convert hours to minutes",
			input:       "2 hours in minutes",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 120 && v.Unit == "minutes" },
		},
		{
			category:    "Units - Time",
			description: "Days",
			input:       "7 days",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 7 && v.Unit == "days" },
		},
		{
			category:    "Units - Time",
			description: "Weeks",
			input:       "2 weeks in days",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 14 && v.Unit == "days" },
		},
		{
			category:    "Units - Time",
			description: "Fortnights",
			input:       "1 fortnight in days",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 14 && v.Unit == "days" },
		},
		{
			category:    "Units - Time",
			description: "Milliseconds",
			input:       "1000 ms in seconds",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 && v.Unit == "seconds" },
		},
		{
			category:    "Units - Time",
			description: "Microseconds",
			input:       "1000 µs in ms",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 && v.Unit == "ms" },
		},
		{
			category:    "Units - Time",
			description: "Nanoseconds",
			input:       "1000 ns in µs",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 0.99 && v.Number < 1.01 && v.Unit == "µs" },
		},

		// ==================== UNITS - VOLUME ====================
		{
			category:    "Units - Volume",
			description: "Litres",
			input:       "5 litres",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "litres" },
		},
		{
			category:    "Units - Volume",
			description: "Convert litres to millilitres",
			input:       "2 l in ml",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2000 && v.Unit == "ml" },
		},
		{
			category:    "Units - Volume",
			description: "Gallons (US)",
			input:       "5 usgallon",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 },
		},
		{
			category:    "Units - Volume",
			description: "Gallons (UK/Imperial)",
			input:       "5 ukgallon",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 },
		},
		{
			category:    "Units - Volume",
			description: "Pints",
			input:       "8 pints",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 && v.Unit == "pints" },
		},
		{
			category:    "Units - Volume",
			description: "Cups",
			input:       "2 cups",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2 && v.Unit == "cups" },
		},
		{
			category:    "Units - Volume",
			description: "Tablespoons",
			input:       "3 tbsp",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 3 && v.Unit == "tbsp" },
		},
		{
			category:    "Units - Volume",
			description: "Teaspoons",
			input:       "5 tsp",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "tsp" },
		},

		// ==================== UNITS - TEMPERATURE ====================
		{
			category:    "Units - Temperature",
			description: "Celsius",
			input:       "20 celsius",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 && v.Unit == "celsius" },
		},
		{
			category:    "Units - Temperature",
			description: "Convert Celsius to Fahrenheit",
			input:       "0 c in f",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 32 && v.Unit == "f" },
		},
		{
			category:    "Units - Temperature",
			description: "Convert Fahrenheit to Celsius",
			input:       "32 f in c",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 0 && v.Unit == "c" },
		},
		{
			category:    "Units - Temperature",
			description: "Kelvin",
			input:       "273 kelvin",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 273 && v.Unit == "kelvin" },
		},
		{
			category:    "Units - Temperature",
			description: "Convert Kelvin to Celsius",
			input:       "273 k in c",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number < 1 && v.Number > -1 && v.Unit == "c" },
		},

		// ==================== UNITS - SPEED ====================
		{
			category:    "Units - Speed",
			description: "Metres per second",
			input:       "10 mps",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},
		{
			category:    "Units - Speed",
			description: "Kilometres per hour",
			input:       "100 kph",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 },
		},
		{
			category:    "Units - Speed",
			description: "Miles per hour",
			input:       "60 mph",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 60 },
		},
		{
			category:    "Units - Speed",
			description: "Convert mph to kph",
			input:       "60 mph in kph",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 96 && v.Number < 97 },
		},

		// ==================== UNITS - AREA ====================
		{
			category:    "Units - Area",
			description: "Square metres",
			input:       "100 sqm",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 },
		},
		{
			category:    "Units - Area",
			description: "Square feet",
			input:       "500 sqft",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 500 },
		},
		{
			category:    "Units - Area",
			description: "Hectares",
			input:       "2 hectares",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2 && v.Unit == "hectares" },
		},
		{
			category:    "Units - Area",
			description: "Acres",
			input:       "10 acres",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 && v.Unit == "acres" },
		},

		// ==================== UNITS - PRESSURE ====================
		{
			category:    "Units - Pressure",
			description: "Pascals",
			input:       "1000 pa",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Unit == "pa" },
		},
		{
			category:    "Units - Pressure",
			description: "Bar",
			input:       "2 bar",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2 && v.Unit == "bar" },
		},
		{
			category:    "Units - Pressure",
			description: "PSI (pounds per square inch)",
			input:       "30 psi",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 30 && v.Unit == "psi" },
		},
		{
			category:    "Units - Pressure",
			description: "Convert bar to psi",
			input:       "2 bar in psi",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 29 && v.Number < 30 && v.Unit == "psi" },
		},

		// ==================== UNITS - FORCE ====================
		{
			category:    "Units - Force",
			description: "Newtons",
			input:       "100 n",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "n" },
		},
		{
			category:    "Units - Force",
			description: "Kilonewtons",
			input:       "5 kn",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "kn" },
		},

		// ==================== UNITS - ANGLE ====================
		{
			category:    "Units - Angle",
			description: "Degrees",
			input:       "180 degrees",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 180 },
		},
		{
			category:    "Units - Angle",
			description: "Radians",
			input:       "3.14 radians",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 3.1 && v.Number < 3.2 },
		},
		{
			category:    "Units - Angle",
			description: "Convert degrees to radians",
			input:       "180 degrees in radians",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 3.1 && v.Number < 3.2 },
		},

		// ==================== UNITS - FREQUENCY ====================
		{
			category:    "Units - Frequency",
			description: "Hertz",
			input:       "1000 hz",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Unit == "hz" },
		},
		{
			category:    "Units - Frequency",
			description: "Kilohertz",
			input:       "2.4 khz",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2.4 && v.Unit == "khz" },
		},
		{
			category:    "Units - Frequency",
			description: "Megahertz",
			input:       "100 mhz",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "mhz" },
		},
		{
			category:    "Units - Frequency",
			description: "Gigahertz",
			input:       "2.4 ghz",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2.4 && v.Unit == "ghz" },
		},

		// ==================== UNITS - DATA STORAGE ====================
		{
			category:    "Units - Data Storage",
			description: "Bytes",
			input:       "1024 bytes",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1024 },
		},
		{
			category:    "Units - Data Storage",
			description: "Kilobytes",
			input:       "100 kb",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "kb" },
		},
		{
			category:    "Units - Data Storage",
			description: "Megabytes",
			input:       "500 mb",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 500 && v.Unit == "mb" },
		},
		{
			category:    "Units - Data Storage",
			description: "Gigabytes",
			input:       "100 gb",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "gb" },
		},
		{
			category:    "Units - Data Storage",
			description: "Terabytes",
			input:       "5 tb",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "tb" },
		},
		{
			category:    "Units - Data Storage",
			description: "Convert GB to MB",
			input:       "1 gb in mb",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1024 && v.Unit == "mb" },
		},
		{
			category:    "Units - Data Storage",
			description: "Bits",
			input:       "8 bits",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 && v.Unit == "bits" },
		},
		{
			category:    "Units - Data Storage",
			description: "Convert bits to bytes",
			input:       "8 bits in bytes",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 && v.Unit == "bytes" },
		},

		// ==================== UNITS - DATA RATE ====================
		{
			category:    "Units - Data Rate",
			description: "Bits per second",
			input:       "1000 bps",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1000 && v.Unit == "bps" },
		},
		{
			category:    "Units - Data Rate",
			description: "Kilobits per second",
			input:       "100 kbps",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "kbps" },
		},
		{
			category:    "Units - Data Rate",
			description: "Megabits per second",
			input:       "100 mbps",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Unit == "mbps" },
		},
		{
			category:    "Units - Data Rate",
			description: "Gigabits per second",
			input:       "1 gbps",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 && v.Unit == "gbps" },
		},

		// ==================== RATE EXPRESSIONS ====================
		{
			category:    "Rate Expressions",
			description: "Simple rate (division)",
			input:       "100 km / 2 hours",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 && strings.Contains(v.Unit, "/") },
		},
		{
			category:    "Rate Expressions",
			description: "Metres per second",
			input:       "100 metres / 10 seconds",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},

		// ==================== TIME VALUES ====================
		{
			category:    "Time Values",
			description: "Time in HH:MM format",
			input:       "14:00",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 14 && v.Unit == "time" },
		},
		{
			category:    "Time Values",
			description: "Add hours to time",
			input:       "14:00 + 2",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 16 && v.Unit == "time" },
		},
		{
			category:    "Time Values",
			description: "Subtract times",
			input:       "17:30 - 09:15",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 8.2 && v.Number < 8.3 && v.Unit == "time" },
		},

		// ==================== DATE LITERALS ====================
		{
			category:    "Date Literals",
			description: "Date in DD/MM/YYYY format",
			input:       "21/10/2024",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return v.Date.Day() == 21 && v.Date.Month() == 10 },
		},
		{
			category:    "Date Literals",
			description: "Add months to date",
			input:       "21/10/2024 + 3 months",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return v.Date.Month() == 1 && v.Date.Year() == 2025 },
		},
		{
			category:    "Date Literals",
			description: "Add years to date",
			input:       "01/01/2024 + 1 year",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return v.Date.Year() == 2025 },
		},
		{
			category:    "Date Literals",
			description: "Subtract days from date",
			input:       "15/02/2024 - 7 days",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return v.Date.Day() == 8 },
		},
		{
			category:    "Date Literals",
			description: "Subtract two dates (returns days)",
			input:       "25/12/2025 - 4/11/2025",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 51 && v.Unit == "days" },
		},

		// ==================== DATE KEYWORDS ====================
		{
			category:    "Date Keywords",
			description: "Today",
			input:       "today",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Tomorrow",
			input:       "tomorrow",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Yesterday",
			input:       "yesterday",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Today plus days",
			input:       "today + 7 days",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Today plus weeks",
			input:       "today + 2 weeks",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Today plus months",
			input:       "today + 3 months",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Date Keywords",
			description: "Today minus months",
			input:       "today - 1 month",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},

		// ==================== WEEKDAY EXPRESSIONS ====================
		{
			category:    "Weekday Expressions",
			description: "Next Monday",
			input:       "next monday",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Weekday Expressions",
			description: "Last Friday",
			input:       "last friday",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},
		{
			category:    "Weekday Expressions",
			description: "This Tuesday",
			input:       "tuesday",
			expectType:  evaluator.ValueDate,
			checkValue:  func(v evaluator.Value) bool { return !v.Date.IsZero() },
		},

		// ==================== FUZZY PHRASES ====================
		{
			category:    "Fuzzy Phrases",
			description: "Half of number",
			input:       "half of 100",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 },
		},
		{
			category:    "Fuzzy Phrases",
			description: "Double number",
			input:       "double 25",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 },
		},
		{
			category:    "Fuzzy Phrases",
			description: "Twice number",
			input:       "twice 30",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 60 },
		},
		{
			category:    "Fuzzy Phrases",
			description: "Three quarters of number",
			input:       "three quarters of 80",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 60 },
		},
		{
			category:    "Fuzzy Phrases",
			description: "Half of currency",
			input:       "half of £100",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 50 && v.Currency == "£" },
		},
		{
			category:    "Fuzzy Phrases",
			description: "Double currency",
			input:       "double $50",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 && v.Currency == "$" },
		},

		// ==================== FUNCTIONS ====================
		{
			category:    "Functions",
			description: "Sum of numbers",
			input:       "sum(10, 20, 30)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 60 },
		},
		{
			category:    "Functions",
			description: "Sum with two arguments",
			input:       "sum(5, 10)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 15 },
		},
		{
			category:    "Functions",
			description: "Average of numbers",
			input:       "average(10, 20, 30)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 },
		},
		{
			category:    "Functions",
			description: "Mean (alias for average)",
			input:       "mean(5, 10, 15)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},
		{
			category:    "Functions",
			description: "Total (alias for sum)",
			input:       "total(1, 2, 3, 4)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},

		// ==================== NUMBER WORDS ====================
		{
			category:    "Number Words",
			description: "Word: one",
			input:       "one",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1 },
		},
		{
			category:    "Number Words",
			description: "Word: two",
			input:       "two",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 2 },
		},
		{
			category:    "Number Words",
			description: "Word: five",
			input:       "five",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 },
		},
		{
			category:    "Number Words",
			description: "Word: ten",
			input:       "ten",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 10 },
		},
		{
			category:    "Number Words",
			description: "Word: twenty",
			input:       "twenty",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 },
		},
		{
			category:    "Number Words",
			description: "Word: hundred",
			input:       "one hundred",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 100 },
		},
		{
			category:    "Number Words",
			description: "Number word with unit",
			input:       "five metres",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 5 && v.Unit == "metres" },
		},
		{
			category:    "Number Words",
			description: "Number word arithmetic",
			input:       "five + three",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 8 },
		},

		// ==================== COMPLEX EXPRESSIONS ====================
		{
			category:    "Complex Expressions",
			description: "Mixed units and currency",
			input:       "£10 * 5 + £20",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 70 && v.Currency == "£" },
		},
		{
			category:    "Complex Expressions",
			description: "Percentage of currency",
			input:       "20% of £100",
			expectType:  evaluator.ValueCurrency,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 20 && v.Currency == "£" },
		},
		{
			category:    "Complex Expressions",
			description: "Unit conversion in expression",
			input:       "(10 m + 5 m) in cm",
			expectType:  evaluator.ValueUnit,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 1500 && v.Unit == "cm" },
		},
		{
			category:    "Complex Expressions",
			description: "Nested percentages",
			input:       "increase 100 by 10% + 5",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number > 114.9 && v.Number < 115.1 },
		},
		{
			category:    "Complex Expressions",
			description: "Function with units",
			input:       "sum(10 m, 5 m)",
			expectType:  evaluator.ValueNumber,
			checkValue:  func(v evaluator.Value) bool { return v.Number == 15 },
		},
	}

	// Group by category for better output
	categories := make(map[string][]int)
	for i, test := range tests {
		categories[test.category] = append(categories[test.category], i)
	}

	// Run tests in category order
	categoryOrder := []string{
		"Basic Arithmetic",
		"Variables",
		"Percentages",
		"Currency",
		"Units - Length",
		"Units - Mass",
		"Units - Time",
		"Units - Volume",
		"Units - Temperature",
		"Units - Speed",
		"Units - Area",
		"Units - Pressure",
		"Units - Force",
		"Units - Angle",
		"Units - Frequency",
		"Units - Data Storage",
		"Units - Data Rate",
		"Rate Expressions",
		"Time Values",
		"Date Literals",
		"Date Keywords",
		"Weekday Expressions",
		"Fuzzy Phrases",
		"Functions",
		"Number Words",
		"Complex Expressions",
	}

	passCount := 0
	failCount := 0

	for _, category := range categoryOrder {
		indices := categories[category]
		if len(indices) == 0 {
			continue
		}

		t.Run(category, func(t *testing.T) {
			for _, i := range indices {
				test := tests[i]
				t.Run(test.description, func(t *testing.T) {
					result := evalExpr(test.input)

					if result.IsError() {
						t.Errorf("Input: %q\nGot error: %s", test.input, result.Error)
						failCount++
						return
					}

					if result.Type != test.expectType {
						t.Errorf("Input: %q\nExpected type %v, got %v", test.input, test.expectType, result.Type)
						failCount++
						return
					}

					if !test.checkValue(result) {
						t.Errorf("Input: %q\nValue check failed. Got: %+v", test.input, result)
						failCount++
						return
					}

					passCount++
				})
			}
		})
	}

	t.Logf("\n\n=== SUMMARY ===")
	t.Logf("Total tests: %d", len(tests))
	t.Logf("Passed: %d", passCount)
	t.Logf("Failed: %d", failCount)
	t.Logf("Coverage: %d categories", len(categoryOrder))
}
