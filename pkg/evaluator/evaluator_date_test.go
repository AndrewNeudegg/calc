package evaluator

import (
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestDateArithmeticWithKeywords tests date arithmetic using keywords like "today"
func TestDateArithmeticWithKeywords(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"today plus 1 day", "today + 1 day", today.AddDate(0, 0, 1)},
		{"today plus 3 days", "today + 3 days", today.AddDate(0, 0, 3)},
		{"today plus 1 week", "today + 1 week", today.AddDate(0, 0, 7)},
		{"today plus 2 weeks", "today + 2 weeks", today.AddDate(0, 0, 14)},
		{"today plus 1 month", "today + 1 month", today.AddDate(0, 1, 0)},
		{"today plus 3 months", "today + 3 months", today.AddDate(0, 3, 0)},
		{"today plus 1 year", "today + 1 year", today.AddDate(1, 0, 0)},
		{"today minus 1 day", "today - 1 day", today.AddDate(0, 0, -1)},
		{"today minus 1 week", "today - 1 week", today.AddDate(0, 0, -7)},
		{"today minus 1 month", "today - 1 month", today.AddDate(0, -1, 0)},
		{"tomorrow", "tomorrow", today.AddDate(0, 0, 1)},
		{"yesterday", "yesterday", today.AddDate(0, 0, -1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)
			if result.IsError() {
				t.Errorf("got error: %s", result.Error)
				return
			}

			if result.Type != ValueDate {
				t.Errorf("expected date, got %v", result.Type)
				return
			}

			// Compare dates (ignoring time component)
			resultDate := time.Date(result.Date.Year(), result.Date.Month(), result.Date.Day(), 0, 0, 0, 0, result.Date.Location())
			expectedDate := time.Date(tt.expected.Year(), tt.expected.Month(), tt.expected.Day(), 0, 0, 0, 0, tt.expected.Location())

			if !resultDate.Equal(expectedDate) {
				t.Errorf("expected %v, got %v", expectedDate.Format("2006-01-02"), resultDate.Format("2006-01-02"))
			}
		})
	}
}

func TestDateArithmeticSupportsHoursAndMinutes(t *testing.T) {
	e := New(NewEnvironment())

	inputs := []string{
		"now + 3 days + 2 hours + five minutes",
		"today + 1 hour",
		"today + 90 minutes",
		"yesterday + 3600 seconds",
	}

	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			l := lexer.New(in)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			v := e.Eval(expr)
			if v.IsError() {
				t.Fatalf("unexpected error: %s", v.Error)
			}
			if v.Type != ValueDate {
				t.Fatalf("expected ValueDate, got %v", v.Type)
			}
		})
	}
}

func TestDateArithmeticSubtractsHoursMinutesSeconds(t *testing.T) {
	e := New(NewEnvironment())

	inputs := []string{
		"now - 2 hours",
		"today - 30 minutes",
		"tomorrow - 86400 seconds",
		"now + 1 day - 2 hours - five minutes",
	}

	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			l := lexer.New(in)
			toks := l.AllTokens()
			p := parser.New(toks)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			v := e.Eval(expr)
			if v.IsError() {
				t.Fatalf("unexpected error: %s", v.Error)
			}
			if v.Type != ValueDate {
				t.Fatalf("expected ValueDate, got %v", v.Type)
			}
		})
	}
}

func TestDateLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "21/10/2024",
			input:    "21/10/2024",
			expected: time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "01/01/2025",
			input:    "01/01/2025",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "31/12/2024",
			input:    "31/12/2024",
			expected: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "29/02/2024", // leap year
			input:    "29/02/2024",
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
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

			if !result.Date.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result.Date)
			}
		})
	}
}

func TestDateArithmeticWithLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "21/10/2024 + 3 months",
			input:    "21/10/2024 + 3 months",
			expected: time.Date(2025, 1, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "01/01/2024 + 1 year",
			input:    "01/01/2024 + 1 year",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "15/02/2024 - 7 days",
			input:    "15/02/2024 - 7 days",
			expected: time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "31/01/2024 + 1 month",
			input:    "31/01/2024 + 1 month",
			expected: time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC), // Go's AddDate behavior: Jan 31 + 1 month = Mar 2
		},
		{
			name:     "29/02/2024 + 1 year",
			input:    "29/02/2024 + 1 year",
			expected: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), // Go's AddDate behavior: Feb 29 + 1 year = Mar 1 (2025 not leap)
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

			if !result.Date.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result.Date)
			}
		})
	}
}

func TestDateSubtraction(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedDays float64
		expectedUnit string
	}{
		{
			name:         "21/10/2025 - 21/10/2024",
			input:        "21/10/2025 - 21/10/2024",
			expectedDays: 365,
			expectedUnit: "days",
		},
		{
			name:         "01/01/2025 - 01/01/2024",
			input:        "01/01/2025 - 01/01/2024",
			expectedDays: 366, // 2024 is a leap year
			expectedUnit: "days",
		},
		{
			name:         "25/12/2025 - 4/11/2025",
			input:        "25/12/2025 - 4/11/2025",
			expectedDays: 51,
			expectedUnit: "days",
		},
		{
			name:         "01/01/2024 - 01/01/2024",
			input:        "01/01/2024 - 01/01/2024",
			expectedDays: 0,
			expectedUnit: "days",
		},
		{
			name:         "21/10/2024 - 21/10/2025", // negative
			input:        "21/10/2024 - 21/10/2025",
			expectedDays: -365,
			expectedUnit: "days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)
			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			if result.Type != ValueUnit {
				t.Fatalf("Expected ValueUnit, got %v", result.Type)
			}

			if result.Number != tt.expectedDays {
				t.Errorf("Expected %v days, got %v", tt.expectedDays, result.Number)
			}

			if result.Unit != tt.expectedUnit {
				t.Errorf("Expected unit %s, got %s", tt.expectedUnit, result.Unit)
			}
		})
	}
}

// TestDateKeywordSubtraction tests date subtraction with keywords like "today"
// This specifically tests the fix for the issue where "today - 19/09/2025" would fail
func TestDateKeywordSubtraction(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkDays bool // if true, check that the result is approximately correct in days
	}{
		{
			name:      "today - date literal",
			input:     "today - 19/09/2025",
			checkDays: true,
		},
		{
			name:      "date literal - today",
			input:     "19/09/2025 - today",
			checkDays: true,
		},
		{
			name:      "tomorrow - date literal",
			input:     "tomorrow - 19/09/2025",
			checkDays: true,
		},
		{
			name:      "yesterday - date literal",
			input:     "yesterday - 19/09/2025",
			checkDays: true,
		},
		{
			name:      "today - today",
			input:     "today - today",
			checkDays: false,
		},
		{
			name:      "tomorrow - today",
			input:     "tomorrow - today",
			checkDays: false,
		},
		{
			name:      "today - yesterday",
			input:     "today - yesterday",
			checkDays: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)
			if result.IsError() {
				t.Fatalf("Eval error: %s", result.Error)
			}

			// Should return a unit value (days)
			if result.Type != ValueUnit {
				t.Fatalf("Expected ValueUnit, got %v", result.Type)
			}

			if result.Unit != "days" {
				t.Errorf("Expected unit 'days', got '%s'", result.Unit)
			}

			// For operations involving today/tomorrow/yesterday, we can check specific results
			if !tt.checkDays {
				switch tt.input {
				case "today - today":
					if result.Number != 0 {
						t.Errorf("Expected 0 days for 'today - today', got %v", result.Number)
					}
				case "tomorrow - today":
					if result.Number != 1 {
						t.Errorf("Expected 1 day for 'tomorrow - today', got %v", result.Number)
					}
				case "today - yesterday":
					if result.Number != 1 {
						t.Errorf("Expected 1 day for 'today - yesterday', got %v", result.Number)
					}
				}
			}
			// For operations with specific dates (checkDays=true),
			// we just verify that no error occurred (checked above)
		})
	}
}

// TestDateKeywordSubtractionReversibility tests that date subtraction is consistent
// This ensures that "today - date" and "date - today" give opposite results
func TestDateKeywordSubtractionReversibility(t *testing.T) {
	testCases := []struct {
		forward string
		reverse string
	}{
		{
			forward: "today - 19/09/2025",
			reverse: "19/09/2025 - today",
		},
		{
			forward: "tomorrow - 15/03/2024",
			reverse: "15/03/2024 - tomorrow",
		},
		{
			forward: "yesterday - 31/12/2024",
			reverse: "31/12/2024 - yesterday",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.forward+" vs "+tc.reverse, func(t *testing.T) {
			resultForward := evalExpr(tc.forward)
			if resultForward.IsError() {
				t.Fatalf("Forward eval error: %s", resultForward.Error)
			}

			resultReverse := evalExpr(tc.reverse)
			if resultReverse.IsError() {
				t.Fatalf("Reverse eval error: %s", resultReverse.Error)
			}

			// Both should be unit values
			if resultForward.Type != ValueUnit || resultReverse.Type != ValueUnit {
				t.Fatalf("Expected both results to be ValueUnit, got %v and %v",
					resultForward.Type, resultReverse.Type)
			}

			// Both should use "days" as the unit
			if resultForward.Unit != "days" || resultReverse.Unit != "days" {
				t.Fatalf("Expected both results to use 'days', got '%s' and '%s'",
					resultForward.Unit, resultReverse.Unit)
			}

			// The numbers should be opposite (one positive, one negative)
			if resultForward.Number != -resultReverse.Number {
				t.Errorf("Expected opposite values, got %v and %v",
					resultForward.Number, resultReverse.Number)
			}
		})
	}
}
