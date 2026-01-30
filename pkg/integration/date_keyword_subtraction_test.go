package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
)

// TestDateKeywordSubtractionIntegration tests date subtraction with date keywords like "today"
// This is an integration test for the issue: "today - 19/09/2025" should work the same as "19/09/2025 - today"
func TestDateKeywordSubtractionIntegration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectType  evaluator.ValueType
		expectUnit  string
	}{
		{
			name:        "today minus date literal",
			input:       "today - 19/09/2025",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "date literal minus today",
			input:       "19/09/2025 - today",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "tomorrow minus date literal",
			input:       "tomorrow - 01/01/2025",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "yesterday minus date literal",
			input:       "yesterday - 15/06/2024",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "today minus today",
			input:       "today - today",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "tomorrow minus today",
			input:       "tomorrow - today",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
		{
			name:        "today minus yesterday",
			input:       "today - yesterday",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)

			if tt.expectError {
				if !result.IsError() {
					t.Errorf("Expected error, but got success with value: %v", result)
				}
			} else {
				if result.IsError() {
					t.Fatalf("Unexpected error: %s", result.Error)
				}

				if result.Type != tt.expectType {
					t.Errorf("Expected type %v, got %v", tt.expectType, result.Type)
				}

				if tt.expectUnit != "" && result.Unit != tt.expectUnit {
					t.Errorf("Expected unit '%s', got '%s'", tt.expectUnit, result.Unit)
				}
			}
		})
	}
}

// TestDateKeywordSubtractionReversibilityIntegration tests that "today - date" and "date - today" give opposite results
func TestDateKeywordSubtractionReversibilityIntegration(t *testing.T) {
	testPairs := []struct {
		name    string
		forward string
		reverse string
	}{
		{
			name:    "today vs 19 Sep 2025",
			forward: "today - 19/09/2025",
			reverse: "19/09/2025 - today",
		},
		{
			name:    "tomorrow vs 10 Oct 2024",
			forward: "tomorrow - 10/10/2024",
			reverse: "10/10/2024 - tomorrow",
		},
		{
			name:    "yesterday vs 25 Dec 2025",
			forward: "yesterday - 25/12/2025",
			reverse: "25/12/2025 - yesterday",
		},
	}

	for _, pair := range testPairs {
		t.Run(pair.name, func(t *testing.T) {
			resultForward := evalExpr(pair.forward)
			resultReverse := evalExpr(pair.reverse)

			// Both should succeed
			if resultForward.IsError() {
				t.Fatalf("Forward expression failed: %s", resultForward.Error)
			}
			if resultReverse.IsError() {
				t.Fatalf("Reverse expression failed: %s", resultReverse.Error)
			}

			// Both should be unit values with "days"
			if resultForward.Type != evaluator.ValueUnit {
				t.Errorf("Forward result expected ValueUnit, got %v", resultForward.Type)
			}
			if resultReverse.Type != evaluator.ValueUnit {
				t.Errorf("Reverse result expected ValueUnit, got %v", resultReverse.Type)
			}

			if resultForward.Unit != "days" {
				t.Errorf("Forward result expected unit 'days', got '%s'", resultForward.Unit)
			}
			if resultReverse.Unit != "days" {
				t.Errorf("Reverse result expected unit 'days', got '%s'", resultReverse.Unit)
			}

			// Values should be opposite
			if resultForward.Number != -resultReverse.Number {
				t.Errorf("Expected opposite values, got forward=%v and reverse=%v",
					resultForward.Number, resultReverse.Number)
			}

			// Verify that the result is non-zero to catch bugs where both are 0
			if resultForward.Number == 0 {
				t.Errorf("Expected non-zero result for date difference, got 0")
			}
		})
	}
}

// TestDateKeywordSubtractionWithConversion tests date subtraction with business day conversion
func TestDateKeywordSubtractionWithConversion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectType  evaluator.ValueType
		expectUnit  string
	}{
		{
			name:        "today minus date in business days",
			input:       "today - 19/09/2025 in business days",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "business days",
		},
		{
			name:        "date minus today in business days",
			input:       "19/09/2025 - today in business days",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "business days",
		},
		{
			name:        "tomorrow minus date in business days",
			input:       "tomorrow - 01/01/2025 in business days",
			expectError: false,
			expectType:  evaluator.ValueUnit,
			expectUnit:  "business days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(tt.input)

			if tt.expectError {
				if !result.IsError() {
					t.Errorf("Expected error, but got success")
				}
			} else {
				if result.IsError() {
					t.Fatalf("Unexpected error: %s", result.Error)
				}

				if result.Type != tt.expectType {
					t.Errorf("Expected type %v, got %v", tt.expectType, result.Type)
				}

				if result.Unit != tt.expectUnit {
					t.Errorf("Expected unit '%s', got '%s'", tt.expectUnit, result.Unit)
				}
			}
		})
	}
}
