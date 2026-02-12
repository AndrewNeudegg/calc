package evaluator

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

func TestTimezoneAbbreviationConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ValueType
		checkFn  func(*testing.T, Value)
	}{
		{
			name:     "9 am EST in UTC",
			input:    "9 am EST in UTC",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// 9 AM EST = 14:00 UTC (EST is UTC-5)
				hour := v.Date.Hour()
				if hour != 14 {
					t.Errorf("Expected hour 14, got %d", hour)
				}
			},
		},
		{
			name:     "10:00 EST in UTC",
			input:    "10:00 EST in UTC",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// 10:00 EST = 15:00 UTC
				hour := v.Date.Hour()
				if hour != 15 {
					t.Errorf("Expected hour 15, got %d", hour)
				}
			},
		},
		{
			name:     "10 am PST in EST",
			input:    "10 am PST in EST",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// Convert 10 AM PST to EST:
				// 10 AM PST (UTC-8) = 18:00 UTC (10 + 8)
				// 18:00 UTC in EST (UTC-5) = 18:00 - 5 = 13:00
				// The result time has EST offset applied, displayed as UTC
				hour := v.Date.Hour()
				if hour != 13 {
					t.Errorf("Expected hour 13, got %d", hour)
				}
			},
		},
		{
			name:     "3 pm CET in PST",
			input:    "3 pm CET in PST",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// 3 PM (15:00) CET = 14:00 UTC = 6 AM PST
				// Result in PST offset: 14:00 - 8 = 6:00
				hour := v.Date.Hour()
				if hour != 6 {
					t.Errorf("Expected hour 6, got %d", hour)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			e := New(env)

			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := e.Eval(expr)
			if result.IsError() {
				t.Fatalf("Eval error: %v", result.Error)
			}

			if result.Type != tt.wantType {
				t.Errorf("got type %v, want %v", result.Type, tt.wantType)
			}

			if tt.checkFn != nil {
				tt.checkFn(t, result)
			}
		})
	}
}
