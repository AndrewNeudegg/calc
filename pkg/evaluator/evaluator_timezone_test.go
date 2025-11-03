package evaluator

import (
	"strings"
	"testing"
	"time"

	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

func TestTimeInLocation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ValueType
		checkFn  func(*testing.T, Value)
	}{
		{
			name:     "time in london",
			input:    "time in london",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// Check that we get a time close to now
				now := time.Now().UTC()
				diff := v.Date.Sub(now)
				if diff < 0 {
					diff = -diff
				}
				if diff > time.Minute {
					t.Errorf("Time difference too large: %v", diff)
				}
			},
		},
		{
			name:     "time in Tokyo",
			input:    "time in Tokyo",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
			},
		},
		{
			name:     "time in New York",
			input:    "time in New York",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
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

func TestTimeDifference(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantType   ValueType
		wantNumber float64
		wantUnit   string
	}{
		{
			name:       "time difference London Sydney",
			input:      "time difference London Sydney",
			wantType:   ValueUnit,
			wantNumber: 10,
			wantUnit:   "hours",
		},
		{
			name:       "time difference Sydney London",
			input:      "time difference Sydney London",
			wantType:   ValueUnit,
			wantNumber: -10,
			wantUnit:   "hours",
		},
		{
			name:       "time difference New York London",
			input:      "time difference New York London",
			wantType:   ValueUnit,
			wantNumber: 5,
			wantUnit:   "hours",
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

			if result.Number != tt.wantNumber {
				t.Errorf("got number %v, want %v", result.Number, tt.wantNumber)
			}

			if result.Unit != tt.wantUnit {
				t.Errorf("got unit %v, want %v", result.Unit, tt.wantUnit)
			}
		})
	}
}

func TestTimeDifferenceWithUnitConversion(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantType   ValueType
		wantNumber float64
		wantUnit   string
	}{
		{
			name:       "time difference in hours (explicit)",
			input:      "time difference London Kabul in hours",
			wantType:   ValueUnit,
			wantNumber: 4,
			wantUnit:   "hours",
		},
		{
			name:       "time difference in days",
			input:      "time difference London Kabul in days",
			wantType:   ValueUnit,
			wantNumber: 0.16666666666666666, // 4/24
			wantUnit:   "days",
		},
		{
			name:       "time difference in minutes",
			input:      "time difference London Kabul in minutes",
			wantType:   ValueUnit,
			wantNumber: 240, // 4*60
			wantUnit:   "minutes",
		},
		{
			name:       "time difference in seconds",
			input:      "time difference London Kabul in seconds",
			wantType:   ValueUnit,
			wantNumber: 14400, // 4*3600
			wantUnit:   "seconds",
		},
		{
			name:       "time difference New York London in days",
			input:      "time difference New York London in days",
			wantType:   ValueUnit,
			wantNumber: 0.20833333333333334, // 5/24
			wantUnit:   "days",
		},
		{
			name:       "time difference Sydney Tokyo in minutes",
			input:      "time difference Sydney Tokyo in minutes",
			wantType:   ValueUnit,
			wantNumber: -60, // -1*60
			wantUnit:   "minutes",
		},
		{
			name:       "time difference default (no unit)",
			input:      "time difference London Paris",
			wantType:   ValueUnit,
			wantNumber: 1, // Paris is UTC+1, London is UTC+0
			wantUnit:   "hours",
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

			if result.Number != tt.wantNumber {
				t.Errorf("got number %v, want %v", result.Number, tt.wantNumber)
			}

			if result.Unit != tt.wantUnit {
				t.Errorf("got unit %v, want %v", result.Unit, tt.wantUnit)
			}
		})
	}
}

func TestTimeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ValueType
		checkFn  func(*testing.T, Value)
	}{
		{
			name:     "time in Sydney plus 3 hours in London",
			input:    "time in Sydney plus 3 hours in London",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// Sydney is UTC+10, add 3 hours, convert to London (UTC+0)
				// Current time in Sydney: now UTC + 10
				// Add 3 hours: now UTC + 10 + 3 = now + 13 in Sydney
				// Convert to London: We want the same moment in time, just expressed in London's timezone
				// London is UTC+0, Sydney is UTC+10, so London is 10 hours behind
				// The same moment (now + 13 in Sydney time) = (now + 3) in UTC/London time
				now := time.Now().UTC()
				expected := now.Add(3 * time.Hour)
				diff := v.Date.Sub(expected)
				if diff < 0 {
					diff = -diff
				}
				// Allow 1 minute tolerance
				if diff > time.Minute {
					t.Errorf("Time difference too large: %v, got %v, expected around %v", diff, v.Date, expected)
				}
			},
		},
		{
			name:     "time in New York minus 2 hours in Sydney",
			input:    "time in New York minus 2 hours in Sydney",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
				}
				// New York is UTC-5, subtract 2 hours, convert to Sydney (UTC+10)
				// Offset: 10 - (-5) = 15
				// Result: now - 5 - 2 + 15 = now + 8 hours
				now := time.Now().UTC()
				expected := now.Add(8 * time.Hour)
				diff := v.Date.Sub(expected)
				if diff < 0 {
					diff = -diff
				}
				// Allow 1 minute tolerance
				if diff > time.Minute {
					t.Errorf("Time difference too large: %v, got %v, expected around %v", diff, v.Date, expected)
				}
			},
		},
		{
			name:     "time in London plus 5 hours in New York",
			input:    "time in London plus 5 hours in New York",
			wantType: ValueDate,
			checkFn: func(t *testing.T, v Value) {
				if v.Date.IsZero() {
					t.Error("Expected non-zero date")
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

func TestTimezoneErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{
			name:      "unknown timezone",
			input:     "time in Atlantis",
			wantError: "unknown timezone",
		},
		{
			name:      "time difference with unknown location",
			input:     "time difference London Atlantis",
			wantError: "unknown timezone",
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
			if !result.IsError() {
				t.Fatal("Expected error but got success")
			}

			if !strings.Contains(result.Error, tt.wantError) {
				t.Errorf("got error %q, want it to contain %q", result.Error, tt.wantError)
			}
		})
	}
}
