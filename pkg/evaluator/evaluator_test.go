package evaluator

import (
	"math"
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

func parseAndEval(input string) Value {
	lex := lexer.New(input)
	tokens := lex.AllTokens()
	
	// Remove EOF
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}
	
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return NewError(err.Error())
	}
	
	env := NewEnvironment()
	eval := New(env)
	return eval.Eval(expr)
}

func TestArithmeticBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3 + 4 * 5", 23},
		{"(10 + 5) / 2", 7.5},
		{"10 - 3", 7},
		{"6 * 7", 42},
		{"100 / 4", 25},
	}
	
	for _, tt := range tests {
		result := parseAndEval(tt.input)
		if result.IsError() {
			t.Errorf("input %q: unexpected error: %s", tt.input, result.Error)
			continue
		}
		
		if math.Abs(result.Number-tt.expected) > 0.01 {
			t.Errorf("input %q: expected %.2f, got %.2f", tt.input, tt.expected, result.Number)
		}
	}
}

func TestVariableAssignment(t *testing.T) {
	env := NewEnvironment()
	eval := New(env)
	
	// x = 10
	lex1 := lexer.New("x = 10")
	tokens1 := lex1.AllTokens()
	tokens1 = tokens1[:len(tokens1)-1]
	p1 := parser.New(tokens1)
	expr1, _ := p1.Parse()
	result1 := eval.Eval(expr1)
	
	if result1.IsError() {
		t.Fatalf("assignment failed: %s", result1.Error)
	}
	
	if result1.Number != 10 {
		t.Errorf("expected 10, got %.2f", result1.Number)
	}
	
	// y = x + 3
	lex2 := lexer.New("y = x + 3")
	tokens2 := lex2.AllTokens()
	tokens2 = tokens2[:len(tokens2)-1]
	p2 := parser.New(tokens2)
	expr2, _ := p2.Parse()
	result2 := eval.Eval(expr2)
	
	if result2.IsError() {
		t.Fatalf("y assignment failed: %s", result2.Error)
	}
	
	if result2.Number != 13 {
		t.Errorf("expected 13, got %.2f", result2.Number)
	}
}

func TestPercentages(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"30 + 20%", 36},     // 30 + (30 * 0.20)
		{"20% of 50", 10},    // 50 * 0.20
		{"increase 100 by 10%", 110},
		{"decrease 100 by 10%", 90},
	}
	
	for _, tt := range tests {
		result := parseAndEval(tt.input)
		if result.IsError() {
			t.Errorf("input %q: unexpected error: %s", tt.input, result.Error)
			continue
		}
		
		if math.Abs(result.Number-tt.expected) > 0.01 {
			t.Errorf("input %q: expected %.2f, got %.2f", tt.input, tt.expected, result.Number)
		}
	}
}

func TestFuzzyPhrases(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"half of 80", 40},
		{"double 15", 30},
		{"twice 4", 8},
		{"three quarters of 200", 150},
	}
	
	for _, tt := range tests {
		result := parseAndEval(tt.input)
		if result.IsError() {
			t.Errorf("input %q: unexpected error: %s", tt.input, result.Error)
			continue
		}
		
		if math.Abs(result.Number-tt.expected) > 0.01 {
			t.Errorf("input %q: expected %.2f, got %.2f", tt.input, tt.expected, result.Number)
		}
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"sum(10, 20, 30)", 60},
		{"average(3, 4, 5)", 4},
	}
	
	for _, tt := range tests {
		result := parseAndEval(tt.input)
		if result.IsError() {
			t.Errorf("input %q: unexpected error: %s", tt.input, result.Error)
			continue
		}
		
		if math.Abs(result.Number-tt.expected) > 0.01 {
			t.Errorf("input %q: expected %.2f, got %.2f", tt.input, tt.expected, result.Number)
		}
	}
}

func TestCurrency(t *testing.T) {
	result := parseAndEval("$50")
	if result.IsError() {
		t.Fatalf("unexpected error: %s", result.Error)
	}
	
	if result.Type != ValueCurrency {
		t.Errorf("expected currency type, got %v", result.Type)
	}
	
	if result.Number != 50 {
		t.Errorf("expected 50, got %.2f", result.Number)
	}
	
	if result.Currency != "$" {
		t.Errorf("expected $, got %s", result.Currency)
	}
}

func TestUnits(t *testing.T) {
	result := parseAndEval("10 m")
	if result.IsError() {
		t.Fatalf("unexpected error: %s", result.Error)
	}
	
	if result.Type != ValueUnit {
		t.Errorf("expected unit type, got %v", result.Type)
	}
	
	if result.Number != 10 {
		t.Errorf("expected 10, got %.2f", result.Number)
	}
	
	if result.Unit != "m" {
		t.Errorf("expected m, got %s", result.Unit)
	}
}
