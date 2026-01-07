package evaluator

import (
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestBusinessDaysArithmetic tests business day calculations
func TestBusinessDaysArithmetic(t *testing.T) {
	// Use a fixed date for testing: Monday, January 1, 2024
	// This is a Monday, so we can test week transitions
	
	tests := []struct {
		name        string
		input       string
		expectedDay time.Weekday // Expected day of week
		skipDays    int          // Number of calendar days from start
	}{
		{
			name:        "today + 1 business day",
			input:       "01/01/2024 + 1 business day",
			expectedDay: time.Tuesday,
			skipDays:    1, // Monday -> Tuesday
		},
		{
			name:        "today + 5 business days",
			input:       "01/01/2024 + 5 business days",
			expectedDay: time.Monday,
			skipDays:    7, // Monday -> next Monday (full work week)
		},
		{
			name:        "Friday + 1 business day",
			input:       "05/01/2024 + 1 business day",
			expectedDay: time.Monday,
			skipDays:    3, // Friday -> Monday (skips weekend)
		},
		{
			name:        "Friday + 3 business days",
			input:       "05/01/2024 + 3 business days",
			expectedDay: time.Wednesday,
			skipDays:    5, // Friday -> Wednesday (skips weekend)
		},
		{
			name:        "today - 1 business day",
			input:       "02/01/2024 - 1 business day",
			expectedDay: time.Monday,
			skipDays:    -1, // Tuesday -> Monday
		},
		{
			name:        "Monday - 1 business day",
			input:       "08/01/2024 - 1 business day",
			expectedDay: time.Friday,
			skipDays:    -3, // Monday -> previous Friday (skips weekend)
		},
		{
			name:        "today + 10 business days",
			input:       "01/01/2024 + 10 business days",
			expectedDay: time.Monday,
			skipDays:    14, // Monday -> Monday two weeks later
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)
			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != ValueDate {
				t.Fatalf("Expected ValueDate, got %v", result.Type)
			}

			// Check the day of the week
			if result.Date.Weekday() != tt.expectedDay {
				t.Errorf("Expected weekday %s, got %s", tt.expectedDay, result.Date.Weekday())
			}

			// Also verify the calendar day count is as expected
			// Parse the start date from the input
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			// Extract the start date
			var startDate time.Time
			if binExpr, ok := expr.(*parser.BinaryExpr); ok {
				if dateExpr, ok := binExpr.Left.(*parser.DateExpr); ok {
					startDate = dateExpr.Date
				}
			}

			expectedDate := startDate.AddDate(0, 0, tt.skipDays)
			// Normalize to start of day
			expectedDate = time.Date(expectedDate.Year(), expectedDate.Month(), expectedDate.Day(), 0, 0, 0, 0, expectedDate.Location())
			resultDate := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), 0, 0, 0, 0, result.Date.Location())

			if !resultDate.Equal(expectedDate) {
				t.Errorf("Expected date %s, got %s", expectedDate.Format("2006-01-02"), resultDate.Format("2006-01-02"))
			}
		})
	}
}

// TestBusinessDaysWithDifferentLocales tests business days with different locale schedules
func TestBusinessDaysWithDifferentLocales(t *testing.T) {
	tests := []struct {
		name        string
		locale      string
		input       string
		expectedDay time.Weekday
	}{
		{
			name:        "Western locale (Monday-Friday): Thursday + 1 business day",
			locale:      "en_GB",
			input:       "04/01/2024 + 1 business day",
			expectedDay: time.Friday, // Thursday -> Friday
		},
		{
			name:        "Middle Eastern locale (Sunday-Thursday): Thursday + 1 business day",
			locale:      "ar_SA",
			input:       "04/01/2024 + 1 business day",
			expectedDay: time.Sunday, // Thursday -> Sunday (skips Friday-Saturday weekend)
		},
		{
			name:        "Middle Eastern locale (Sunday-Thursday): Sunday + 5 business days",
			locale:      "ar_AE",
			input:       "07/01/2024 + 5 business days",
			expectedDay: time.Sunday, // Sunday -> next Sunday (skips Friday-Saturday)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.SetLocale(tt.locale)
			
			l := lexer.New(tt.input)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			evaluator := New(env)
			result := evaluator.Eval(expr)

			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != ValueDate {
				t.Fatalf("Expected ValueDate, got %v", result.Type)
			}

			if result.Date.Weekday() != tt.expectedDay {
				t.Errorf("Expected weekday %s, got %s", tt.expectedDay, result.Date.Weekday())
			}
		})
	}
}

// TestBusinessDaysKeywordWithToday tests "today + N business days"
func TestBusinessDaysKeywordWithToday(t *testing.T) {
	// This test uses the actual current date, so we just verify it doesn't error
	// and returns a date in the future
	inputs := []string{
		"today + 1 business day",
		"today + 5 business days",
		"today + 10 business days",
		"today - 1 business day",
		"today - 5 business days",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			result := evalExpr(input)
			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != ValueDate {
				t.Fatalf("Expected ValueDate, got %v", result.Type)
			}

			// Just verify it's a valid date - we can't check the exact value
			// since it depends on today's date
			if result.Date.IsZero() {
				t.Error("Expected non-zero date")
			}
		})
	}
}
