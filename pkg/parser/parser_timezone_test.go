package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParseTimezoneConversion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectType  string // "TimeConversionExpr" or other
	}{
		// Valid timezone conversion patterns
		{
			name:        "9 am EST in UTC",
			input:       "9 am EST in UTC",
			expectError: false,
			expectType:  "TimeConversionExpr",
		},
		{
			name:        "10:00 PST in EST",
			input:       "10:00 PST in EST",
			expectError: false,
			expectType:  "TimeConversionExpr",
		},
		{
			name:        "3 pm CET in PST",
			input:       "3 pm CET in PST",
			expectError: false,
			expectType:  "TimeConversionExpr",
		},
		{
			name:        "10:59 EST in UTC",
			input:       "10:59 EST in UTC",
			expectError: false,
			expectType:  "TimeConversionExpr",
		},
		
		// Non-matching patterns (should not be parsed as timezone conversion)
		{
			name:        "9 am + 1 (no timezone)",
			input:       "9 am + 1",
			expectError: false,
			expectType:  "NumberExpr", // "am" is treated as a dimensionless unit, result is just number
		},
		{
			name:        "9 am foo (invalid timezone)",
			input:       "9 am foo",
			expectError: false,
			expectType:  "NumberExpr", // Should parse as number with units
		},
		{
			name:        "9 am EST (missing 'in')",
			input:       "9 am EST",
			expectError: false,
			expectType:  "NumberExpr", // Should parse as number with units
		},
		{
			name:        "10:00 (just time, no timezone)",
			input:       "10:00",
			expectError: false,
			expectType:  "UnitExpr", // Should parse as time unit
		},
		{
			name:        "10:00 + 2 (time arithmetic)",
			input:       "10:00 + 2",
			expectError: false,
			expectType:  "BinaryExpr", // Should parse as addition
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := New(tokens)
			expr, err := p.Parse()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected parse error: %v", err)
			}

			// Check expression type
			exprType := getExprTypeName(expr)
			if exprType != tt.expectType {
				t.Errorf("Expected expression type %s, got %s", tt.expectType, exprType)
			}

			// For TimeConversionExpr, verify it has the expected fields
			if tt.expectType == "TimeConversionExpr" {
				timeConv, ok := expr.(*TimeConversionExpr)
				if !ok {
					t.Fatalf("Expected TimeConversionExpr but got %T", expr)
				}
				if timeConv.Time == nil {
					t.Error("TimeConversionExpr.Time should not be nil")
				}
				if timeConv.From == "" {
					t.Error("TimeConversionExpr.From should not be empty")
				}
				if timeConv.To == "" {
					t.Error("TimeConversionExpr.To should not be empty")
				}
			}
		})
	}
}

// Helper function to get expression type name
func getExprTypeName(expr Expr) string {
	switch expr.(type) {
	case *NumberExpr:
		return "NumberExpr"
	case *BinaryExpr:
		return "BinaryExpr"
	case *UnaryExpr:
		return "UnaryExpr"
	case *IdentExpr:
		return "IdentExpr"
	case *AssignExpr:
		return "AssignExpr"
	case *UnitExpr:
		return "UnitExpr"
	case *ConversionExpr:
		return "ConversionExpr"
	case *CurrencyExpr:
		return "CurrencyExpr"
	case *PercentExpr:
		return "PercentExpr"
	case *PercentOfExpr:
		return "PercentOfExpr"
	case *PercentChangeExpr:
		return "PercentChangeExpr"
	case *WhatPercentExpr:
		return "WhatPercentExpr"
	case *FunctionCallExpr:
		return "FunctionCallExpr"
	case *StringExpr:
		return "StringExpr"
	case *DateExpr:
		return "DateExpr"
	case *TimeExpr:
		return "TimeExpr"
	case *DateArithmeticExpr:
		return "DateArithmeticExpr"
	case *FuzzyExpr:
		return "FuzzyExpr"
	case *CommandExpr:
		return "CommandExpr"
	case *RateExpr:
		return "RateExpr"
	case *WeekdayExpr:
		return "WeekdayExpr"
	case *TimeInLocationExpr:
		return "TimeInLocationExpr"
	case *TimeDifferenceExpr:
		return "TimeDifferenceExpr"
	case *TimeConversionExpr:
		return "TimeConversionExpr"
	case *MonthExpr:
		return "MonthExpr"
	case *PrevExpr:
		return "PrevExpr"
	case *ArgDirectiveExpr:
		return "ArgDirectiveExpr"
	default:
		return "Unknown"
	}
}
