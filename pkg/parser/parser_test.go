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
		input       string
		expectType  string
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
		input    string
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
