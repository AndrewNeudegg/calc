package parser

import (
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func parseInput(input string) (Expr, error) {
	l := lexer.New(input)
	tokens := l.AllTokens()
	p := New(tokens)
	return p.Parse()
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42", 42},
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"1000", 1000},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors: %v", err)
			continue
		}

		numExpr, ok := expr.(*NumberExpr)
		if !ok {
			t.Errorf("Expected NumberExpr, got %T", expr)
			continue
		}

		if numExpr.Value != tt.expected {
			t.Errorf("Expected %f, got %f", tt.expected, numExpr.Value)
		}
	}
}

func TestParseBinaryOperations(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"3 + 4", "+"},
		{"10 - 5", "-"},
		{"6 * 7", "*"},
		{"20 / 4", "/"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors: %v", err)
			continue
		}

		binExpr, ok := expr.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr, got %T", expr)
			continue
		}

		if binExpr.Operator != tt.operator {
			t.Errorf("Expected operator %s, got %s", tt.operator, binExpr.Operator)
		}
	}
}

func TestParsePercentages(t *testing.T) {
	tests := []struct {
		input      string
		expectType string
	}{
		{"20%", "PercentExpr"},
		{"10% of 100", "PercentOfExpr"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		var typeName string
		switch expr.(type) {
		case *PercentExpr:
			typeName = "PercentExpr"
		case *PercentOfExpr:
			typeName = "PercentOfExpr"
		case *BinaryExpr:
			typeName = "BinaryExpr"
		default:
			typeName = "Unknown"
		}

		if typeName != tt.expectType {
			t.Errorf("For %q: expected %s, got %s", tt.input, tt.expectType, typeName)
		}
	}
}

func TestParseUnitConversion(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"10 m in cm"},
		{"5 kg in g"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		_, ok := expr.(*ConversionExpr)
		if !ok {
			t.Errorf("Expected ConversionExpr for %q, got %T", tt.input, expr)
		}
	}
}

func TestParseFuzzyPhrases(t *testing.T) {
	tests := []struct {
		input   string
		pattern string
	}{
		{"half of 100", "half"},
		{"double 50", "double"},
		{"twice 25", "twice"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		fuzzyExpr, ok := expr.(*FuzzyExpr)
		if !ok {
			t.Errorf("Expected FuzzyExpr for %q, got %T", tt.input, expr)
			continue
		}

		if fuzzyExpr.Pattern != tt.pattern {
			t.Errorf("Expected pattern %q, got %q", tt.pattern, fuzzyExpr.Pattern)
		}
	}
}

func TestParseDateKeywords(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"today"},
		{"tomorrow"},
		{"yesterday"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		_, ok := expr.(*DateExpr)
		if !ok {
			t.Errorf("Expected DateExpr for %q, got %T", tt.input, expr)
		}
	}
}

func TestParseWeekdays(t *testing.T) {
	tests := []struct {
		input    string
		modifier string
		weekday  time.Weekday
	}{
		{"next monday", "next", time.Monday},
		{"last friday", "last", time.Friday},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		weekdayExpr, ok := expr.(*WeekdayExpr)
		if !ok {
			t.Errorf("Expected WeekdayExpr for %q, got %T", tt.input, expr)
			continue
		}

		if weekdayExpr.Modifier != tt.modifier {
			t.Errorf("Expected modifier %q, got %q", tt.modifier, weekdayExpr.Modifier)
		}

		if weekdayExpr.Weekday != tt.weekday {
			t.Errorf("Expected weekday %v, got %v", tt.weekday, weekdayExpr.Weekday)
		}
	}
}

func TestParseFunctions(t *testing.T) {
	tests := []struct {
		input    string
		funcName string
		argCount int
	}{
		{"sum(1, 2, 3)", "sum", 3},
		{"average(10, 20)", "average", 2},
		{"min(3, 1, 2)", "min", 3},
		{"max(3, 1, 2)", "max", 3},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		funcExpr, ok := expr.(*FunctionCallExpr)
		if !ok {
			t.Errorf("Expected FunctionCallExpr for %q, got %T", tt.input, expr)
			continue
		}

		if funcExpr.Name != tt.funcName {
			t.Errorf("Expected function %q, got %q", tt.funcName, funcExpr.Name)
		}

		if len(funcExpr.Args) != tt.argCount {
			t.Errorf("Expected %d args, got %d", tt.argCount, len(funcExpr.Args))
		}
	}
}

func TestParseAssignment(t *testing.T) {
	tests := []struct {
		input    string
		variable string
	}{
		{"x = 10", "x"},
		{"price = 100", "price"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		assignExpr, ok := expr.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr for %q, got %T", tt.input, expr)
			continue
		}

		if assignExpr.Name != tt.variable {
			t.Errorf("Expected variable %q, got %q", tt.variable, assignExpr.Name)
		}
	}
}

func TestParseCommands(t *testing.T) {
	tests := []struct {
		input   string
		command string
	}{
		{":help", "help"},
		{":save", "save"},
		{":set", "set"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
			continue
		}

		cmdExpr, ok := expr.(*CommandExpr)
		if !ok {
			t.Errorf("Expected CommandExpr for %q, got %T", tt.input, expr)
			continue
		}

		if cmdExpr.Command != tt.command {
			t.Errorf("Expected command %q, got %q", tt.command, cmdExpr.Command)
		}
	}
}

func TestParseComplexExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"3 + 4 * 5"},
		{"(10 + 20) * 3"},
		{"10 m + 50 cm"},
		{"£100 + £50"},
	}

	for _, tt := range tests {
		expr, err := parseInput(tt.input)

		if err != nil {
			t.Errorf("Parser errors for %q: %v", tt.input, err)
		}

		if expr == nil {
			t.Errorf("Expected expression for %q, got nil", tt.input)
		}
	}
}

// TestParseFuzzyPhraseAssignments tests that fuzzy phrases can be assigned to variables
func TestParseFuzzyPhraseAssignments(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		variableName string
		fuzzyPattern string
	}{
		{
			name:         "assign half of number",
			input:        "fo = half of 99",
			variableName: "fo",
			fuzzyPattern: "half",
		},
		{
			name:         "assign double variable",
			input:        "y = double x",
			variableName: "y",
			fuzzyPattern: "double",
		},
		{
			name:         "assign twice variable",
			input:        "z = twice x",
			variableName: "z",
			fuzzyPattern: "twice",
		},
		{
			name:         "assign three quarters",
			input:        "part = three quarters of amount",
			variableName: "part",
			fuzzyPattern: "three quarters",
		},
		{
			name:         "assign half of variable",
			input:        "result = half of fo",
			variableName: "result",
			fuzzyPattern: "half",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parse error for %q: %v", tt.input, err)
			}

			// Should be an assignment expression
			assignExpr, ok := expr.(*AssignExpr)
			if !ok {
				t.Fatalf("expected AssignExpr for %q, got %T", tt.input, expr)
			}

			// Check variable name
			if assignExpr.Name != tt.variableName {
				t.Errorf("expected variable name %q, got %q", tt.variableName, assignExpr.Name)
			}

			// The value should be a FuzzyExpr
			fuzzyExpr, ok := assignExpr.Value.(*FuzzyExpr)
			if !ok {
				t.Fatalf("expected FuzzyExpr as value, got %T", assignExpr.Value)
			}

			// Check the pattern
			if fuzzyExpr.Pattern != tt.fuzzyPattern {
				t.Errorf("expected pattern %q, got %q", tt.fuzzyPattern, fuzzyExpr.Pattern)
			}

			// Verify that the fuzzy expression has a valid Value field
			if fuzzyExpr.Value == nil {
				t.Errorf("fuzzy expression Value should not be nil")
			}
		})
	}
}

// TestParseFuzzyPhraseWithVariableReference tests fuzzy phrases that reference variables
func TestParseFuzzyPhraseWithVariableReference(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedPattern string
		expectedVarName string
	}{
		{
			name:            "half of variable",
			input:           "half of fo",
			expectedPattern: "half",
			expectedVarName: "fo",
		},
		{
			name:            "double variable",
			input:           "double x",
			expectedPattern: "double",
			expectedVarName: "x",
		},
		{
			name:            "twice variable",
			input:           "twice myvar",
			expectedPattern: "twice",
			expectedVarName: "myvar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parse error for %q: %v", tt.input, err)
			}

			// Should be a FuzzyExpr
			fuzzyExpr, ok := expr.(*FuzzyExpr)
			if !ok {
				t.Fatalf("expected FuzzyExpr for %q, got %T", tt.input, expr)
			}

			// Check pattern
			if fuzzyExpr.Pattern != tt.expectedPattern {
				t.Errorf("expected pattern %q, got %q", tt.expectedPattern, fuzzyExpr.Pattern)
			}

			// The value should be an IdentExpr (variable reference)
			identExpr, ok := fuzzyExpr.Value.(*IdentExpr)
			if !ok {
				t.Fatalf("expected IdentExpr as value, got %T", fuzzyExpr.Value)
			}

			// Check variable name
			if identExpr.Name != tt.expectedVarName {
				t.Errorf("expected variable name %q, got %q", tt.expectedVarName, identExpr.Name)
			}
		})
	}
}
