package integration

import (
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
)

func TestPhysicalConstants(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		contains string // check if output contains this string
	}{
		{
			name:    "speed of light by symbol",
			input:   "c",
			wantErr: false,
		},
		{
			name:    "speed of light by name",
			input:   "speed_of_light",
			wantErr: false,
		},
		{
			name:    "planck constant by symbol",
			input:   "h",
			wantErr: false,
		},
		{
			name:    "planck constant by name",
			input:   "planck",
			wantErr: false,
		},
		{
			name:    "gravitational constant by symbol",
			input:   "G",
			wantErr: false,
		},
		{
			name:    "gravitational constant by name",
			input:   "gravitational_constant",
			wantErr: false,
		},
		{
			name:    "elementary charge",
			input:   "e",
			wantErr: false,
		},
		{
			name:    "stefan-boltzmann constant",
			input:   "σ",
			wantErr: false,
		},
		{
			name:    "constant in expression",
			input:   "c * 2",
			wantErr: false,
		},
		{
			name:    "constant with unit conversion",
			input:   "c in km/s",
			wantErr: false,
		},
		{
			name:    "assign constant to variable",
			input:   "speed = c",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			
			// Set up lexer with constant checker
			l := lexer.New(tt.input)
			l.SetConstantChecker(env.Constants().IsConstant)
			
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()

			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := env.Eval(expr)

			if tt.wantErr && !result.IsError() {
				t.Errorf("Expected error but got result: %v", result)
			}

			if !tt.wantErr && result.IsError() {
				t.Errorf("Unexpected error: %v", result.Error)
			}

			if tt.contains != "" {
				resultStr := result.String()
				if !strings.Contains(resultStr, tt.contains) {
					t.Errorf("Expected result to contain %q, got %q", tt.contains, resultStr)
				}
			}
		})
	}
}

func TestConstantArithmetic(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "multiply constants",
			input:   "c * h",
			wantErr: false,
		},
		{
			name:    "divide constants",
			input:   "c / h",
			wantErr: false,
		},
		{
			name:    "add constants with same dimension",
			input:   "electron_mass + proton_mass",
			wantErr: false,
		},
		{
			name:    "constant times number",
			input:   "3 * c",
			wantErr: false,
		},
		{
			name:    "number times constant",
			input:   "c * 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			l := lexer.New(tt.input)
			l.SetConstantChecker(env.Constants().IsConstant)
			
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()

			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := env.Eval(expr)

			if tt.wantErr && !result.IsError() {
				t.Errorf("Expected error but got result: %v", result)
			}

			if !tt.wantErr && result.IsError() {
				t.Errorf("Unexpected error: %v", result.Error)
			}
		})
	}
}

func TestConstantCaseInsensitive(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"lowercase c", "c"},
		{"uppercase C", "C"},
		{"mixed Planck", "Planck"},
		{"uppercase PLANCK", "PLANCK"},
		{"lowercase planck", "planck"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			l := lexer.New(tt.input)
			l.SetConstantChecker(env.Constants().IsConstant)
			
			tokens := l.AllTokens()
			p := parser.New(tokens)
			expr, err := p.Parse()

			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := env.Eval(expr)

			if result.IsError() {
				t.Errorf("Unexpected error: %v", result.Error)
			}
		})
	}
}
