package display

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andrewneudegg/calc/pkg/commands"
	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/formatter"
	"github.com/andrewneudegg/calc/pkg/graph"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/settings"
)

// Line represents a single calculation line.
type Line struct {
	ID     int
	Input  string
	Result evaluator.Value
	Expr   parser.Expr
}

// REPL manages the read-eval-print loop.
type REPL struct {
	lines      map[int]*Line
	nextID     int
	env        *evaluator.Environment
	eval       *evaluator.Evaluator
	formatter  *formatter.Formatter
	commands   *commands.Handler
	settings   *settings.Settings
	depGraph   *graph.Graph
}

// NewREPL creates a new REPL instance.
func NewREPL() *REPL {
	// Load settings
	homeDir, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.config/calc/settings.json", homeDir)
	
	sett, err := settings.Load(configPath)
	if err != nil {
		sett = settings.Default()
		sett.ConfigPath = configPath
	}
	
	env := evaluator.NewEnvironment()
	
	return &REPL{
		lines:     make(map[int]*Line),
		nextID:    1,
		env:       env,
		eval:      evaluator.New(env),
		formatter: formatter.New(sett),
		commands:  commands.New(sett),
		settings:  sett,
		depGraph:  graph.NewGraph(),
	}
}

// Run starts the REPL loop.
func (r *REPL) Run() {
	fmt.Println("Calc - A terminal notepad calculator")
	fmt.Println("Type :help for available commands, :quit to exit")
	fmt.Println()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Printf("%d> ", r.nextID)
		
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		result := r.EvaluateLine(input)
		if !result.IsError() || result.Error != "" {
			fmt.Printf("   = %s\n\n", r.formatter.Format(result))
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
	}
}

// EvaluateLine processes a single line of input.
func (r *REPL) EvaluateLine(input string) evaluator.Value {
	// Tokenise
	lex := lexer.New(input)
	tokens := lex.AllTokens()
	
	// Remove EOF token for parsing
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}
	
	// Parse
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return evaluator.NewError(err.Error())
	}
	
	// Check if it's a command
	if cmd, ok := expr.(*parser.CommandExpr); ok {
		msg := r.commands.Execute(cmd.Command, cmd.Args)
		fmt.Println(msg)
		return evaluator.Value{} // Empty value
	}
	
	// Evaluate
	result := r.eval.Eval(expr)
	
	// Store the line
	lineID := r.nextID
	r.nextID++
	
	r.lines[lineID] = &Line{
		ID:     lineID,
		Input:  input,
		Result: result,
		Expr:   expr,
	}
	
	return result
}

// GetLine retrieves a line by ID.
func (r *REPL) GetLine(id int) (*Line, bool) {
	line, ok := r.lines[id]
	return line, ok
}

// ListLines returns all lines in order.
func (r *REPL) ListLines() []*Line {
	var lines []*Line
	for i := 1; i < r.nextID; i++ {
		if line, ok := r.lines[i]; ok {
			lines = append(lines, line)
		}
	}
	return lines
}
