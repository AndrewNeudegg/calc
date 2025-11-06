package evaluator

import (
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
