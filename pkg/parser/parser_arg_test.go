package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

func TestParseArgDirective(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantName   string
		wantPrompt string
		wantErr    bool
	}{
		{
			name:       "arg with variable name only",
			input:      ":arg count",
			wantName:   "count",
			wantPrompt: "",
			wantErr:    false,
		},
		{
			name:       "arg with variable name and prompt",
			input:      `:arg total "What is the total?"`,
			wantName:   "total",
			wantPrompt: "What is the total?",
			wantErr:    false,
		},
		{
			name:       "arg with multi-word prompt",
			input:      `:arg nights "How many nights are you staying?"`,
			wantName:   "nights",
			wantPrompt: "How many nights are you staying?",
			wantErr:    false,
		},
		{
			name:    "arg without variable name",
			input:   ":arg",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			
			// Remove EOF token
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}
			
			p := New(tokens)
			expr, err := p.Parse()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			
			argExpr, ok := expr.(*ArgDirectiveExpr)
			if !ok {
				t.Fatalf("expected *ArgDirectiveExpr, got %T", expr)
			}
			
			if argExpr.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", argExpr.Name, tt.wantName)
			}
			
			if argExpr.Prompt != tt.wantPrompt {
				t.Errorf("Prompt = %q, want %q", argExpr.Prompt, tt.wantPrompt)
			}
		})
	}
}
