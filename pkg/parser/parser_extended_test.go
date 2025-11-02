package parser

import (
	"testing"
)

// TestParserOperatorPrecedence tests operator precedence
func TestParserOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{"1 + 2 * 3", "multiplication before addition"},
		{"2 * 3 + 4", "multiplication before addition (reversed)"},
		{"10 - 2 * 3", "multiplication before subtraction"},
		{"10 / 2 + 3", "division before addition"},
		{"(1 + 2) * 3", "parentheses override precedence"},
		{"2 * (3 + 4)", "parentheses override precedence (reversed)"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)
		if err != nil {
			t.Errorf("%s (%q): parse error %v", tt.description, tt.input, err)
		}
		if expr == nil {
			t.Errorf("%s (%q): got nil expression", tt.description, tt.input)
		}
	}
}

// TestParserComplexExpressions tests complex nested expressions
func TestParserComplexExpressions(t *testing.T) {
	tests := []string{
		"((1 + 2) * 3) / 4",
		"1 + 2 + 3 + 4",
		"10 - 5 - 2",
		"2 * 3 * 4",
		"100 / 10 / 2",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserUnaryOperators tests unary minus
func TestParserUnaryOperators(t *testing.T) {
	tests := []string{
		"-5",
		"-10 + 5",
		"10 + -5",
		"-(3 + 4)",
		"--5",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserAssignment tests variable assignment
func TestParserAssignment(t *testing.T) {
	tests := []string{
		"x = 10",
		"alpha = 20 + 5",
		"result = 10 * 2",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*AssignExpr)
		if !ok {
			t.Errorf("%q: expected AssignExpr, got %T", input, expr)
		}
	}
}

// TestParserFunctions tests function calls
func TestParserFunctions(t *testing.T) {
	tests := []string{
		"sum(1, 2, 3)",
		"average(10, 20, 30)",
		"total(5)",
		"mean(1, 2)",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*FunctionCallExpr)
		if !ok {
			t.Errorf("%q: expected FunctionCallExpr, got %T", input, expr)
		}
	}
}

// TestParserUnitConversions tests unit conversion expressions
func TestParserUnitConversions(t *testing.T) {
	tests := []string{
		"10 m in cm",
		"5 km in m",
		"100 g in kg",
		"2 hours in minutes",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*ConversionExpr)
		if !ok {
			t.Errorf("%q: expected ConversionExpr, got %T", input, expr)
		}
	}
}

// TestParserCurrency tests currency expressions
func TestParserCurrency(t *testing.T) {
	tests := []string{
		"$100",
		"£50",
		"€75",
		"¥1000",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*CurrencyExpr)
		if !ok {
			t.Errorf("%q: expected CurrencyExpr, got %T", input, expr)
		}
	}
}

// TestParserDateKeywords tests date keyword expressions
func TestParserDateKeywords(t *testing.T) {
	tests := []string{
		"today",
		"tomorrow",
		"yesterday",
		"now",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserWeekdays tests weekday expressions
func TestParserWeekdays(t *testing.T) {
	tests := []string{
		"next monday",
		"last friday",
		"next wednesday",
		"last tuesday",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*WeekdayExpr)
		if !ok {
			t.Errorf("%q: expected WeekdayExpr, got %T", input, expr)
		}
	}
}

// TestParserFuzzyExpressions tests fuzzy language expressions
func TestParserFuzzyExpressions(t *testing.T) {
	tests := []string{
		"half of 100",
		"double 50",
		"twice 25",
		"three quarters of 200",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*FuzzyExpr)
		if !ok {
			t.Errorf("%q: expected FuzzyExpr, got %T", input, expr)
		}
	}
}

// TestParserPercentageVariants tests all percentage expression types
func TestParserPercentageVariants(t *testing.T) {
	tests := []struct {
		input        string
		expectedType string
	}{
		{"50%", "PercentExpr"},
		{"10% of 100", "PercentOfExpr"},
		{"increase 100 by 10%", "PercentChangeExpr"},
		{"decrease 100 by 10%", "PercentChangeExpr"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)
		if err != nil {
			t.Errorf("%q: parse error %v", tt.input, err)
			continue
		}

		var typeName string
		switch expr.(type) {
		case *PercentExpr:
			typeName = "PercentExpr"
		case *PercentOfExpr:
			typeName = "PercentOfExpr"
		case *PercentChangeExpr:
			typeName = "PercentChangeExpr"
		default:
			typeName = "unknown"
		}

		if typeName != tt.expectedType {
			t.Errorf("%q: expected %s, got %s", tt.input, tt.expectedType, typeName)
		}
	}
}

// TestParserErrorRecovery tests error handling
func TestParserErrorRecovery(t *testing.T) {
	tests := []string{
		"10 +",    // Missing operand
		"* 10",    // Missing left operand
		"(10 + 5", // Unclosed parenthesis
		// Removed: "10 20", - valid as separate statements
		"", // Empty input
	}

	for _, input := range tests {
		_, err := parseInput(input)
		if err == nil {
			t.Errorf("%q: expected error, got none", input)
		}
	}
}

// TestParserDateArithmetic tests date arithmetic expressions
func TestParserDateArithmetic(t *testing.T) {
	tests := []string{
		"today + 1 day",
		"tomorrow - 2 days",
		"yesterday + 1 week",
		"today + 1 month",
		"today - 1 year",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*DateArithmeticExpr)
		if !ok {
			t.Errorf("%q: expected DateArithmeticExpr, got %T", input, expr)
		}
	}
}

// TestParserIdentifiers tests identifier parsing
func TestParserIdentifiers(t *testing.T) {
	tests := []string{
		"x",
		"alpha",
		"my_var",
		"var123",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*IdentExpr)
		if !ok {
			t.Errorf("%q: expected IdentExpr, got %T", input, expr)
		}
	}
}

// TestParserMixedOperations tests mixed operation types
func TestParserMixedOperations(t *testing.T) {
	tests := []string{
		"10 m + 5 m",
		"$100 + $50",
		"£100 - £25",
		"10 kg * 2",
		"100 cm / 2",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserNestedFunctions tests nested function calls
func TestParserNestedFunctions(t *testing.T) {
	tests := []string{
		"sum(10, 20)",
		"average(sum(1, 2), sum(3, 4))",
	}

	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserNumberWords tests parsing English number words
func TestParserNumberWords(t *testing.T) {
	tests := []struct {
		input       string
		expectedNum float64
	}{
		{"one", 1},
		{"five", 5},
		{"ten", 10},
		{"twenty", 20},
		{"thirty five", 35},
		{"forty two", 42},
		{"ninety nine", 99},
		{"one hundred", 100},
		{"two hundred", 200},
		{"three hundred and forty two", 342},
		{"one thousand", 1000},
		{"five thousand", 5000},
		{"three thousand five hundred", 3500},
		{"one million", 1000000},
	}
	
	for _, tt := range tests {
		expr, err := parseInput(tt.input)
		if err != nil {
			t.Errorf("%q: parse error %v", tt.input, err)
			continue
		}
		numExpr, ok := expr.(*NumberExpr)
		if !ok {
			t.Errorf("%q: expected NumberExpr, got %T", tt.input, expr)
			continue
		}
		if numExpr.Value != tt.expectedNum {
			t.Errorf("%q: got %f, want %f", tt.input, numExpr.Value, tt.expectedNum)
		}
	}
}

// TestParserNumberWordsInExpressions tests number words in complex expressions
func TestParserNumberWordsInExpressions(t *testing.T) {
	tests := []string{
		"five plus ten",
		"twenty minus seven",
		"three times four",
		"one hundred divided by ten",
		"three hundred and forty two divided by seventeen",
		"one thousand plus five hundred",
	}
	
	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
		}
		if expr == nil {
			t.Errorf("%q: got nil expression", input)
		}
	}
}

// TestParserNumberWordsWithUnits tests number words with units
func TestParserNumberWordsWithUnits(t *testing.T) {
	tests := []string{
		"three years in hours",
		"five meters in cm",
		"ten kg in grams",
		"twenty minutes in seconds",
	}
	
	for _, input := range tests {
		expr, err := parseInput(input)
		if err != nil {
			t.Errorf("%q: parse error %v", input, err)
			continue
		}
		_, ok := expr.(*ConversionExpr)
		if !ok {
			t.Errorf("%q: expected ConversionExpr, got %T", input, expr)
		}
	}
}
