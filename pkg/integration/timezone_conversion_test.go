package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

// TestTimezoneConversion tests timezone conversion with various formats
func TestTimezoneConversion(t *testing.T) {
	tests := []struct {
		input        string
		expectNoErr  bool
		expectedHour int // Expected hour component of result
	}{
		// Basic timezone abbreviation conversions
		{"9 am EST in UTC", true, 14},       // 9 AM EST (UTC-5) = 14:00 UTC
		{"10:00 EST in UTC", true, 15},      // 10:00 EST = 15:00 UTC
		{"09:00 PST in UTC", true, 17},      // 09:00 PST (UTC-8) = 17:00 UTC
		{"10 am PST in EST", true, 13},      // 10 AM PST = 1 PM EST (displayed as 13:00)
		{"3 pm CET in PST", true, 6},        // 3 PM CET (UTC+1) = 14:00 UTC = 6 AM PST
		
		// More timezone abbreviations
		{"12 pm GMT in JST", true, 21},      // 12 PM GMT = 21:00 JST (UTC+9)
		{"8 am AEST in GMT", true, 22},      // 8 AM AEST (UTC+10) = 22:00 GMT (previous day)
		{"2 pm MST in CST", true, 15},       // 2 PM MST (UTC-7) = 21:00 UTC = 3 PM CST (UTC-6)
		
		// HH:MM format with timezones
		{"14:30 UTC in EST", true, 9},       // 14:30 UTC = 9:30 EST
		{"18:00 PST in CET", true, 3},       // 18:00 PST = 2:00 UTC+1day = 3:00 CET
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Create environment
			env := evaluator.NewEnvironment()
			eval := evaluator.New(env)

			// Lex
			l := lexer.New(tt.input)
			tokens := l.AllTokens()

			// Parse
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				if tt.expectNoErr {
					t.Fatalf("Parse error: %v", err)
				}
				return
			}

			// Evaluate
			result := eval.Eval(expr)
			if result.IsError() {
				if tt.expectNoErr {
					t.Fatalf("Eval error: %s", result.Error)
				}
				return
			}

			if !tt.expectNoErr {
				t.Fatalf("Expected error but got success")
			}

			// Check result type
			if result.Type != evaluator.ValueDate {
				t.Errorf("Expected ValueDate, got %v", result.Type)
			}

			// Check hour component
			hour := result.Date.Hour()
			if hour != tt.expectedHour {
				t.Errorf("Expected hour %d, got %d", tt.expectedHour, hour)
			}
		})
	}
}
