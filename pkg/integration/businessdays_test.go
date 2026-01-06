package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestBusinessDaysIntegration tests business days functionality end-to-end
func TestBusinessDaysIntegration(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		shouldError    bool
		checkWeekday   bool
		expectedWeekday time.Weekday
	}{
		{
			name:           "simple business day addition",
			input:          "01/01/2024 + 1 business day",
			shouldError:    false,
			checkWeekday:   true,
			expectedWeekday: time.Tuesday, // Monday + 1 business day = Tuesday
		},
		{
			name:           "business days over weekend",
			input:          "05/01/2024 + 3 business days",
			shouldError:    false,
			checkWeekday:   true,
			expectedWeekday: time.Wednesday, // Friday + 3 business days = Wednesday (skips weekend)
		},
		{
			name:           "subtract business days",
			input:          "08/01/2024 - 1 business day",
			shouldError:    false,
			checkWeekday:   true,
			expectedWeekday: time.Friday, // Monday - 1 business day = previous Friday
		},
		{
			name:        "today with business days",
			input:       "today + 5 business days",
			shouldError: false,
			checkWeekday: false, // Can't check exact weekday since it depends on today
		},
		{
			name:        "tomorrow with business days",
			input:       "tomorrow + 10 business days",
			shouldError: false,
			checkWeekday: false,
		},
		{
			name:        "yesterday with business days",
			input:       "yesterday - 3 business days",
			shouldError: false,
			checkWeekday: false,
		},
		{
			name:           "full work week",
			input:          "01/01/2024 + 5 business days",
			shouldError:    false,
			checkWeekday:   true,
			expectedWeekday: time.Monday, // Monday + 5 business days = next Monday
		},
		{
			name:           "two work weeks",
			input:          "01/01/2024 + 10 business days",
			shouldError:    false,
			checkWeekday:   true,
			expectedWeekday: time.Monday, // Monday + 10 business days = Monday 2 weeks later
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				if !tt.shouldError {
					t.Fatalf("Parse error: %v", err)
				}
				return
			}

			env := evaluator.NewEnvironment()
			eval := evaluator.New(env)
			result := eval.Eval(expr)

			if result.IsError() {
				if !tt.shouldError {
					t.Fatalf("Eval error: %s", result.Error)
				}
				return
			}

			if tt.shouldError {
				t.Fatal("Expected error but got success")
			}

			if result.Type != evaluator.ValueDate {
				t.Fatalf("Expected ValueDate, got %v", result.Type)
			}

			if tt.checkWeekday && result.Date.Weekday() != tt.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s", tt.expectedWeekday, result.Date.Weekday())
			}
		})
	}
}

// TestBusinessDaysWithLocales tests business days with different locale settings
func TestBusinessDaysWithLocales(t *testing.T) {
	tests := []struct {
		name           string
		locale         string
		input          string
		expectedWeekday time.Weekday
	}{
		{
			name:           "Western (Mon-Fri): Thursday + 1 bd",
			locale:         "en_GB",
			input:          "04/01/2024 + 1 business day",
			expectedWeekday: time.Friday,
		},
		{
			name:           "Western (Mon-Fri): Friday + 1 bd",
			locale:         "en_US",
			input:          "05/01/2024 + 1 business day",
			expectedWeekday: time.Monday,
		},
		{
			name:           "Middle Eastern (Sun-Thu): Thursday + 1 bd",
			locale:         "ar_SA",
			input:          "04/01/2024 + 1 business day",
			expectedWeekday: time.Sunday, // Skips Fri-Sat weekend
		},
		{
			name:           "Middle Eastern (Sun-Thu): Thursday + 2 bd",
			locale:         "ar_AE",
			input:          "04/01/2024 + 2 business days",
			expectedWeekday: time.Monday, // Thu -> Sun (skip Fri-Sat) -> Mon
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			env := evaluator.NewEnvironment()
			env.SetLocale(tt.locale)
			eval := evaluator.New(env)
			result := eval.Eval(expr)

			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != evaluator.ValueDate {
				t.Fatalf("Expected ValueDate, got %v", result.Type)
			}

			if result.Date.Weekday() != tt.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s (date: %s)",
					tt.expectedWeekday, result.Date.Weekday(), result.Date.Format("2006-01-02"))
			}
		})
	}
}

// TestBusinessDaysInExpressions tests business days in more complex expressions
func TestBusinessDaysInExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "business days with variable assignment",
			input: "deadline = today + 10 business days",
		},
		{
			name:  "business days subtraction with variable",
			input: "start = 15/01/2024\nend = start + 5 business days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(tt.input, "\n")
			env := evaluator.NewEnvironment()
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				l := lexer.New(line)
				toks := l.AllTokens()
				p := parser.New(toks)
				expr, err := p.Parse()
				if err != nil {
					t.Fatalf("Parse error on line '%s': %v", line, err)
				}

				eval := evaluator.New(env)
				result := eval.Eval(expr)

				if result.IsError() {
					t.Fatalf("Eval error on line '%s': %s", line, result.Error)
				}
			}
		})
	}
}
