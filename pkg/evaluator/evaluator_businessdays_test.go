package evaluator

import (
	"strings"
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

// TestBusinessDaysConversion tests converting date differences to business days
func TestBusinessDaysConversion(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedBusinessDays float64
		locale            string
	}{
		{
			name:              "simple date difference to business days",
			input:             "08/01/2024 - 01/01/2024 in business days", // Monday to Monday (1 week)
			expectedBusinessDays: 5,
			locale:            "en_GB",
		},
		{
			name:              "over weekend to business days",
			input:             "08/01/2024 - 05/01/2024 in business days", // Friday to Monday
			expectedBusinessDays: 1,
			locale:            "en_GB",
		},
		{
			name:              "two weeks to business days",
			input:             "15/01/2024 - 01/01/2024 in business days", // Monday to Monday (2 weeks)
			expectedBusinessDays: 10,
			locale:            "en_GB",
		},
		{
			name:              "negative difference to business days",
			input:             "01/01/2024 - 08/01/2024 in business days", // Monday to Monday backwards
			expectedBusinessDays: -5,
			locale:            "en_GB",
		},
		{
			name:              "Middle Eastern locale date difference",
			input:             "11/01/2024 - 07/01/2024 in business days", // Sunday to Thursday (4 business days in ar_SA)
			expectedBusinessDays: 4,
			locale:            "ar_SA",
		},
		{
			name:              "same day difference",
			input:             "01/01/2024 - 01/01/2024 in business days",
			expectedBusinessDays: 0,
			locale:            "en_GB",
		},
		{
			name:              "single business day",
			input:             "02/01/2024 - 01/01/2024 in business day", // Monday to Tuesday
			expectedBusinessDays: 1,
			locale:            "en_GB",
		},
		{
			name:              "mixed case unit",
			input:             "08/01/2024 - 01/01/2024 in Business Days",
			expectedBusinessDays: 5,
			locale:            "en_GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			if tt.locale != "" {
				env.SetLocale(tt.locale)
			}
			
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

			if result.Type != ValueUnit {
				t.Fatalf("Expected ValueUnit, got %v", result.Type)
			}

			if result.Number != tt.expectedBusinessDays {
				t.Errorf("Expected %.0f business days, got %.0f", tt.expectedBusinessDays, result.Number)
			}
			
			// Verify the unit is preserved (either "business days" or "business day")
			lowerUnit := strings.ToLower(result.Unit)
			if lowerUnit != "business days" && lowerUnit != "business day" {
				t.Errorf("Expected unit 'business days' or 'business day', got '%s'", result.Unit)
			}
		})
	}
}

// TestBusinessDaysConversionWithKeywords tests business days conversion with date keywords
func TestBusinessDaysConversionWithKeywords(t *testing.T) {
	// Test with a fixed future date minus today
	tests := []struct {
		name  string
		input string
		minDays int
		maxDays int
	}{
		{
			name:  "future date minus today in business days",
			input: "31/12/2026 - today in business days",
			minDays: 100, // At least 100 business days from now to end of 2026
			maxDays: 500, // But less than 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			
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

			if result.Type != ValueUnit {
				t.Fatalf("Expected ValueUnit, got %v", result.Type)
			}
			
			// Verify it's a reasonable number
			if result.Number < float64(tt.minDays) || result.Number > float64(tt.maxDays) {
				t.Errorf("Expected between %d and %d business days, got %.0f", tt.minDays, tt.maxDays, result.Number)
			}
		})
	}
}

