package evaluator

import (
	"fmt"
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

func TestEvaluator_Prev(t *testing.T) {
	tests := []struct {
		name     string
		history  []Value
		input    string
		expected float64
	}{
		{
			name: "prev with one item",
			history: []Value{
				NewNumber(25),
			},
			input:    "prev",
			expected: 25,
		},
		{
			name: "prev with expression",
			history: []Value{
				NewNumber(25),
			},
			input:    "10 + prev",
			expected: 35,
		},
		{
			name: "prev~1 refers to second-to-last",
			history: []Value{
				NewNumber(25),
				NewNumber(35),
			},
			input:    "prev~1",
			expected: 25,
		},
		{
			name: "prev~5 with multiple items",
			history: []Value{
				NewNumber(10),
				NewNumber(20),
				NewNumber(30),
				NewNumber(40),
				NewNumber(50),
				NewNumber(60),
			},
			input:    "prev~5",
			expected: 10,
		},
		{
			name: "multiple prev in expression",
			history: []Value{
				NewNumber(25),
				NewNumber(35),
			},
			input:    "prev + prev~1",
			expected: 60, // 35 + 25
		},
		{
			name: "prev with currency",
			history: []Value{
				NewCurrency(100, "$"),
			},
			input:    "prev * 2",
			expected: 200,
		},
		{
			name: "prev with unit",
			history: []Value{
				NewUnit(10, "m"),
			},
			input:    "prev * 2",
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			
			// Set up history function
			env.SetHistoryFunc(func(offset int) (Value, error) {
				idx := len(tt.history) - 1 - offset
				if idx < 0 || idx >= len(tt.history) {
					return NewError(""), nil
				}
				return tt.history[idx], nil
			})

			eval := New(env)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error: %v", result.Error)
			}

			if result.Number != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result.Number)
			}
		})
	}
}

func TestEvaluator_PrevWithoutHistory(t *testing.T) {
	env := NewEnvironment()
	eval := New(env)

	l := lexer.New("prev")
	tokens := l.AllTokens()
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Eval(expr)
	if !result.IsError() {
		t.Fatal("Expected error when using prev without history function")
	}

	if result.Error != "prev is only available in REPL mode" {
		t.Errorf("Expected specific error message, got: %s", result.Error)
	}
}

func TestEvaluator_PrevOutOfRange(t *testing.T) {
	env := NewEnvironment()
	
	// Set up history function with only one item
	env.SetHistoryFunc(func(offset int) (Value, error) {
		if offset > 0 {
			return NewError(""), nil
		}
		return NewNumber(25), nil
	})

	eval := New(env)

	l := lexer.New("prev~5")
	tokens := l.AllTokens()
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Eval(expr)
	if !result.IsError() {
		t.Fatal("Expected error when accessing prev out of range")
	}
}

func TestEvaluator_PrevAbsolute(t *testing.T) {
	tests := []struct {
		name     string
		history  map[int]Value // lineID -> Value
		input    string
		expected float64
	}{
		{
			name: "prev#1",
			history: map[int]Value{
				1: NewNumber(10),
				2: NewNumber(20),
				3: NewNumber(30),
			},
			input:    "prev#1",
			expected: 10,
		},
		{
			name: "prev#15",
			history: map[int]Value{
				15: NewNumber(110),
				16: NewNumber(120),
			},
			input:    "prev#15",
			expected: 110,
		},
		{
			name: "prev#15 * 42",
			history: map[int]Value{
				15: NewNumber(110),
			},
			input:    "prev#15 * 42",
			expected: 4620,
		},
		{
			name: "prev#10 with currency",
			history: map[int]Value{
				10: NewCurrency(100, "$"),
			},
			input:    "prev#10 * 2",
			expected: 200,
		},
		{
			name: "prev#5 with unit",
			history: map[int]Value{
				5: NewUnit(10, "m"),
			},
			input:    "prev#5 * 3",
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			
			// Set up absolute history function
			env.SetAbsoluteHistoryFunc(func(lineID int) (Value, error) {
				val, ok := tt.history[lineID]
				if !ok {
					return Value{}, fmt.Errorf("no result found for line %d", lineID)
				}
				return val, nil
			})

			eval := New(env)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := eval.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error: %v", result.Error)
			}

			if result.Number != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result.Number)
			}
		})
	}
}

func TestEvaluator_PrevAbsoluteWithoutHistory(t *testing.T) {
	env := NewEnvironment()
	eval := New(env)

	l := lexer.New("prev#15")
	tokens := l.AllTokens()
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Eval(expr)
	if !result.IsError() {
		t.Fatal("Expected error when using prev#N without history function")
	}

	if result.Error != "prev is only available in REPL mode" {
		t.Errorf("Expected specific error message, got: %s", result.Error)
	}
}

func TestEvaluator_PrevAbsoluteOutOfRange(t *testing.T) {
	env := NewEnvironment()
	
	// Set up absolute history function with limited data
	env.SetAbsoluteHistoryFunc(func(lineID int) (Value, error) {
		if lineID == 5 {
			return NewNumber(50), nil
		}
		return Value{}, fmt.Errorf("no result found for line %d", lineID)
	})

	eval := New(env)

	l := lexer.New("prev#100")
	tokens := l.AllTokens()
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := eval.Eval(expr)
	if !result.IsError() {
		t.Fatal("Expected error when accessing non-existent line")
	}
}
