package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestMixedNumericTextualRepresentations tests mixing numeric literals with textual number words
func TestMixedNumericTextualRepresentations(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		// Valid: numeric + scale word
		{"5 million", 5000000},
		{"10 million", 10000000},
		{"3.5 billion", 3500000000},
		{"2.5 thousand", 2500},
		{"10 thousand", 10000},
		{"100 thousand", 100000},
		
		// Valid: numeric + hundred + scale word
		{"5 hundred thousand", 500000},
		{"10 hundred thousand", 1000000},
		{"2 hundred million", 200000000},
		
		// Valid: all textual
		{"five hundred thousand", 500000},
		{"five million", 5000000},
		{"ten hundred thousand", 1000000},
		{"one billion", 1000000000},
		
		// Valid: decimal + scale word
		{"1.5 million", 1500000},
		{"0.5 billion", 500000000},
		{"12.7 thousand", 12700},
		
		// Valid: large numeric + scale word
		{"1000 thousand", 1000000},
		{"999 million", 999000000},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			lex := lexer.New(tt.input)
			lex.SetConstantChecker(env.Constants().IsConstant)
			tokens := lex.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			eval := evaluator.New(env)
			result := eval.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error for %q: %s", tt.input, result.Error)
			}

			if result.Number != tt.expected {
				t.Errorf("%q = %f, want %f", tt.input, result.Number, tt.expected)
			}
		})
	}
}

// TestInvalidMixedNumericTextualRepresentations tests that invalid patterns are rejected
func TestInvalidMixedNumericTextualRepresentations(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		// Invalid: numeric + connector + number word
		{"100000 and three", "numeric literal with 'and' followed by basic number word"},
		{"5000 and five", "numeric literal with 'and' followed by basic number word"},
		{"99 and one", "numeric literal with 'and' followed by basic number word"},
		
		// Invalid: numeric + connector + scale word (caught by our validation)
		{"100 and thousand", "numeric literal with 'and' followed by scale word"},
		{"5 and million", "numeric literal with 'and' followed by scale word"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			lex := lexer.New(tt.input)
			lex.SetConstantChecker(env.Constants().IsConstant)
			tokens := lex.AllTokens()
			p := parser.New(tokens)
			_, err := p.Parse()
			
			// We expect a parse error for these invalid patterns
			if err == nil {
				t.Errorf("%q (%s) should have been rejected but was accepted", tt.input, tt.description)
			} else {
				t.Logf("%q correctly rejected with error: %v", tt.input, err)
			}
		})
	}
}

// TestValidAdditionWithAnd tests that valid uses of "and" still work
func TestValidAdditionWithAnd(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		// Valid: numeric + numeric with "and"
		{"5 and 5", 10},
		{"10 and 20", 30},
		{"100 and 50", 150},
		
		// Valid: textual + textual with "and"
		{"five and three", 8},
		{"ten and twenty", 30},
		{"fifty and fifty", 100},
		
		// Valid: unit expressions with "and"
		{"5 meters and 3 meters", 8}, // Result in meters
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			lex := lexer.New(tt.input)
			lex.SetConstantChecker(env.Constants().IsConstant)
			tokens := lex.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			eval := evaluator.New(env)
			result := eval.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error for %q: %s", tt.input, result.Error)
			}

			if result.Number != tt.expected {
				t.Errorf("%q = %f, want %f", tt.input, result.Number, tt.expected)
			}
		})
	}
}

// TestEdgeCasesForMixedRepresentations tests edge cases
func TestEdgeCasesForMixedRepresentations(t *testing.T) {
	tests := []struct {
		input       string
		expected    float64
		shouldError bool
	}{
		// Edge case: zero with scale word
		{"0 million", 0, false},
		{"0.0 thousand", 0, false},
		
		// Edge case: very large numbers
		{"999999 million", 999999000000, false},
		
		// Edge case: very small decimals
		{"0.001 million", 1000, false},
		{"0.1 billion", 100000000, false},
		
		// Edge case: just scale word (treated as the number itself)
		{"million", 1000000, false},
		{"thousand", 1000, false},
		{"hundred", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			lex := lexer.New(tt.input)
			lex.SetConstantChecker(env.Constants().IsConstant)
			tokens := lex.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("%q should have errored but didn't", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			eval := evaluator.New(env)
			result := eval.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error for %q: %s", tt.input, result.Error)
			}

			if result.Number != tt.expected {
				t.Errorf("%q = %f, want %f", tt.input, result.Number, tt.expected)
			}
		})
	}
}
