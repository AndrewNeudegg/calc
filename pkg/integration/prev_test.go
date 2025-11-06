package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/display"
)

func TestREPL_PrevFeature(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []string
		expected []string
	}{
		{
			name: "basic prev usage",
			inputs: []string{
				"5 * 5",
				"10 + prev",
			},
			expected: []string{
				"25.00",
				"35.00",
			},
		},
		{
			name: "prev with assignment",
			inputs: []string{
				"T = 30 / 5",
				"T",
				"prev",
			},
			expected: []string{
				"6.00",
				"6.00",
				"6.00",
			},
		},
		{
			name: "prev with tilde",
			inputs: []string{
				"10",
				"20",
				"30",
				"prev~",
			},
			expected: []string{
				"10.00",
				"20.00",
				"30.00",
				"20.00",
			},
		},
		{
			name: "prev~1",
			inputs: []string{
				"10",
				"20",
				"30",
				"prev~1",
			},
			expected: []string{
				"10.00",
				"20.00",
				"30.00",
				"20.00",
			},
		},
		{
			name: "prev~5",
			inputs: []string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"prev~5",
			},
			expected: []string{
				"1.00",
				"2.00",
				"3.00",
				"4.00",
				"5.00",
				"6.00",
				"1.00",
			},
		},
		{
			name: "multiple prev in expression",
			inputs: []string{
				"10",
				"20",
				"prev + prev~1",
			},
			expected: []string{
				"10.00",
				"20.00",
				"30.00",
			},
		},
		{
			name: "prev with currency",
			inputs: []string{
				"$100",
				"prev * 1.5",
			},
			expected: []string{
				"$100.00",
				"$150.00",
			},
		},
		{
			name: "prev with units",
			inputs: []string{
				"10 m",
				"prev * 2",
			},
			expected: []string{
				"10.00 m",
				"20.00 m",
			},
		},
		{
			name: "prev in complex expression",
			inputs: []string{
				"5",
				"10",
				"(prev + prev~1) * 2",
			},
			expected: []string{
				"5.00",
				"10.00",
				"30.00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repl := display.NewREPL()
			repl.SetSilent(true)

			for i, input := range tt.inputs {
				result := repl.EvaluateLine(input)
				if result.IsError() && result.Error != "" {
					t.Fatalf("Line %d error: %v", i+1, result.Error)
				}

				formatted := repl.Formatter().Format(result)
				if formatted != tt.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i+1, tt.expected[i], formatted)
				}
			}
		})
	}
}

func TestREPL_PrevErrors(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []string
		expectError int // 0-indexed line that should error
	}{
		{
			name: "prev on first line",
			inputs: []string{
				"prev",
			},
			expectError: 0,
		},
		{
			name: "prev~5 out of range",
			inputs: []string{
				"10",
				"20",
				"prev~5",
			},
			expectError: 2,
		},
		{
			name: "prev~1 on second line",
			inputs: []string{
				"10",
				"prev~1",
			},
			expectError: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repl := display.NewREPL()
			repl.SetSilent(true)

			for i, input := range tt.inputs {
				result := repl.EvaluateLine(input)
				if i == tt.expectError {
					if !result.IsError() || result.Error == "" {
						t.Errorf("Line %d: expected error, got success", i+1)
					}
				} else {
					if result.IsError() && result.Error != "" {
						t.Errorf("Line %d: unexpected error: %v", i+1, result.Error)
					}
				}
			}
		})
	}
}
