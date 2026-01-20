package integration

import (
	"fmt"
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

// TestBusinessDaysConversionIntegration tests converting date differences to business days end-to-end
func TestBusinessDaysConversionIntegration(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedBusinessDays float64
		locale            string
	}{
		{
			name:              "date difference to business days - one week",
			input:             "08/01/2024 - 01/01/2024 in business days",
			expectedBusinessDays: 5,
			locale:            "en_GB",
		},
		{
			name:              "date difference to business days - over weekend",
			input:             "08/01/2024 - 05/01/2024 in business days",
			expectedBusinessDays: 1,
			locale:            "en_GB",
		},
		{
			name:              "date difference to business days - two weeks",
			input:             "15/01/2024 - 01/01/2024 in business days",
			expectedBusinessDays: 10,
			locale:            "en_GB",
		},
		{
			name:              "backwards date difference to business days",
			input:             "01/01/2024 - 15/01/2024 in business days",
			expectedBusinessDays: -10,
			locale:            "en_GB",
		},
		{
			name:              "same day difference to business days",
			input:             "01/01/2024 - 01/01/2024 in business days",
			expectedBusinessDays: 0,
			locale:            "en_GB",
		},
		{
			name:              "Middle Eastern schedule - Sunday to Thursday",
			input:             "11/01/2024 - 07/01/2024 in business days",
			expectedBusinessDays: 4,
			locale:            "ar_SA",
		},
		{
			name:              "singular form - business day",
			input:             "02/01/2024 - 01/01/2024 in business day",
			expectedBusinessDays: 1,
			locale:            "en_GB",
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
			if tt.locale != "" {
				env.SetLocale(tt.locale)
			}
			eval := evaluator.New(env)
			result := eval.Eval(expr)

			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != evaluator.ValueUnit {
				t.Fatalf("Expected ValueUnit, got %v", result.Type)
			}

			if result.Number != tt.expectedBusinessDays {
				t.Errorf("Expected %.0f business days, got %.0f", tt.expectedBusinessDays, result.Number)
			}
		})
	}
}

// TestBusinessDaysRoundTrip tests both adding and subtracting business days
func TestBusinessDaysRoundTrip(t *testing.T) {
	// Test that adding and then subtracting business days gives us back the original date
	tests := []struct {
		name       string
		startDate  string
		businessDays int
	}{
		{
			name:       "5 business days round trip",
			startDate:  "01/01/2024",
			businessDays: 5,
		},
		{
			name:       "10 business days round trip",
			startDate:  "01/01/2024",
			businessDays: 10,
		},
		{
			name:       "3 business days from Friday",
			startDate:  "05/01/2024",
			businessDays: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			
			// Add business days
			addInput := fmt.Sprintf("%s + %d business days", tt.startDate, tt.businessDays)
			l := lexer.New(addInput)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error for add: %v", err)
			}
			
			eval := evaluator.New(env)
			addResult := eval.Eval(expr)
			if addResult.IsError() {
				t.Fatalf("Eval error for add: %s", addResult.Error)
			}
			if addResult.Type != evaluator.ValueDate {
				t.Fatalf("Expected ValueDate from add, got %v", addResult.Type)
			}
			
			// Now count business days between start and result
			countInput := fmt.Sprintf("%s - %s in business days", addResult.Date.Format("02/01/2006"), tt.startDate)
			l2 := lexer.New(countInput)
			toks2 := l2.AllTokens()
			p2 := parser.New(toks2)
			expr2, err := p2.Parse()
			if err != nil {
				t.Fatalf("Parse error for count: %v", err)
			}
			
			eval2 := evaluator.New(env)
			countResult := eval2.Eval(expr2)
			if countResult.IsError() {
				t.Fatalf("Eval error for count: %s", countResult.Error)
			}
			if countResult.Type != evaluator.ValueUnit {
				t.Fatalf("Expected ValueUnit from count, got %v", countResult.Type)
			}
			
			if int(countResult.Number) != tt.businessDays {
				t.Errorf("Expected %d business days, got %.0f", tt.businessDays, countResult.Number)
			}
		})
	}
}

