package parser

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

// TestParseUnitCommand tests parsing ":unit define" command form
func TestParseUnitCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd bool
		cmdName string
		args    []string
	}{
		{
			name:    "unit define command",
			input:   ":unit define spoon = 15 ml",
			wantCmd: true,
			cmdName: "unit",
			args:    []string{"define", "spoon", "=", "15", "ml"},
		},
		{
			name:    "unit list command",
			input:   ":unit list",
			wantCmd: true,
			cmdName: "unit",
			args:    []string{"list"},
		},
		{
			name:    "unit list custom command",
			input:   ":unit list custom",
			wantCmd: true,
			cmdName: "unit",
			args:    []string{"list", "custom"},
		},
		{
			name:    "unit show command",
			input:   ":unit show spoon",
			wantCmd: true,
			cmdName: "unit",
			args:    []string{"show", "spoon"},
		},
		{
			name:    "unit delete command",
			input:   ":unit delete spoon",
			wantCmd: true,
			cmdName: "unit",
			args:    []string{"delete", "spoon"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			
			if tt.wantCmd {
				cmd, ok := expr.(*CommandExpr)
				if !ok {
					t.Fatalf("expected CommandExpr, got %T", expr)
				}
				
				if cmd.Command != tt.cmdName {
					t.Errorf("expected command %q, got %q", tt.cmdName, cmd.Command)
				}
				
				if len(cmd.Args) != len(tt.args) {
					t.Errorf("expected %d args, got %d", len(tt.args), len(cmd.Args))
				}
				
				for i, arg := range tt.args {
					if i < len(cmd.Args) && cmd.Args[i] != arg {
						t.Errorf("arg[%d]: expected %q, got %q", i, arg, cmd.Args[i])
					}
				}
			}
		})
	}
}

// TestParseUnitDirective tests parsing ":unit name = value" directive form
func TestParseUnitDirective(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantDirective bool
		unitName     string
	}{
		{
			name:         "unit directive shorthand",
			input:        ":unit spoon = 15 ml",
			wantDirective: true,
			unitName:     "spoon",
		},
		{
			name:         "unit directive with expression",
			input:        ":unit bowl = 350 ml",
			wantDirective: true,
			unitName:     "bowl",
		},
		{
			name:         "unit directive with multiplication",
			input:        ":unit yard = 3 foot",
			wantDirective: true,
			unitName:     "yard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			
			if tt.wantDirective {
				dir, ok := expr.(*UnitDirectiveExpr)
				if !ok {
					t.Fatalf("expected UnitDirectiveExpr, got %T", expr)
				}
				
				if dir.Name != tt.unitName {
					t.Errorf("expected unit name %q, got %q", tt.unitName, dir.Name)
				}
				
				if dir.Value == nil {
					t.Error("expected non-nil value expression")
				}
			}
		})
	}
}

// TestParseUnitCommandVsDirective ensures commands and directives are parsed correctly
func TestParseUnitCommandVsDirective(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string // "command" or "directive"
	}{
		{
			name:     "command: unit define",
			input:    ":unit define spoon = 15 ml",
			wantType: "command",
		},
		{
			name:     "directive: unit shorthand",
			input:    ":unit spoon = 15 ml",
			wantType: "directive",
		},
		{
			name:     "command: unit list",
			input:    ":unit list",
			wantType: "command",
		},
		{
			name:     "command: unit show",
			input:    ":unit show spoon",
			wantType: "command",
		},
		{
			name:     "command: unit delete",
			input:    ":unit delete spoon",
			wantType: "command",
		},
		{
			name:     "directive: unit with complex expression",
			input:    ":unit dozen = 12",
			wantType: "directive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := l.AllTokens()
			p := New(tokens)
			
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			
			switch tt.wantType {
			case "command":
				if _, ok := expr.(*CommandExpr); !ok {
					t.Errorf("expected CommandExpr, got %T", expr)
				}
			case "directive":
				if _, ok := expr.(*UnitDirectiveExpr); !ok {
					t.Errorf("expected UnitDirectiveExpr, got %T", expr)
				}
			}
		})
	}
}
