package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andrewneudegg/calc/pkg/display"
	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/formatter"
	"github.com/andrewneudegg/calc/pkg/lexer"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/settings"
)

const helpText = `Calc - Terminal Notepad Calculator

A local-only, dependency-free terminal calculator inspired by Soulver.
Mix arithmetic, units, dates, currencies, and variables with natural language.

USAGE:
	calc                Start interactive REPL mode
	calc -c "expr"      Execute a single calculation and exit
	calc -f file.calc    Execute all lines from a file and print results

OPTIONS:
	-c string           Execute calculation and exit
	-f string           Execute a .calc file and print results
	-h, --help          Show this help message

EXAMPLES:
	calc -c "half of 40"
	calc -c "10 m in cm"
	calc -c "11:00 - 09:00"
	calc -c "20% of 100"
	calc -c "£100 + £50"
	calc -f examples/k8s-cluster.calc

FEATURES:
  • Arithmetic with operator precedence and parentheses
  • Variables and assignments (x = 10, y = x * 2)
  • Unit conversions (10 m in cm, 70 kg in lb)
  • Time arithmetic (11:00 - 09:00, 14:00 + 2)
  • Currency operations (£100 + $50, $100 in GBP)
  • Percentages (20% of 50, increase 100 by 10%)
  • Fuzzy phrases (half of X, double Y, three quarters of Z)
  • Date arithmetic (today + 3 weeks, tomorrow - 2 days)
  • Functions (sum(1,2,3), average(10,20,30))

REPL COMMANDS:
  :help              Show available commands
  :set precision N   Set decimal precision
  :set currency C    Set default currency (GBP, USD, EUR, JPY)
  :save file.txt     Save workspace to file
  :open file.txt     Load workspace from file
  :quit              Exit calculator

For more information, visit: https://github.com/AndrewNeudegg/calc
`

func main() {
	// Custom usage function
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, helpText)
	}

	// Define flags
	calcExpr := flag.String("c", "", "Execute a single calculation and exit")
	filePath := flag.String("f", "", "Execute a .calc file and print results")
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message")
	flag.Parse()

	// Show help if requested
	if *showHelp {
		fmt.Print(helpText)
		os.Exit(0)
	}

	// If -f flag is provided, execute file and exit
	if *filePath != "" {
		if err := executeFile(*filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// If -c flag is provided, execute and exit
	if *calcExpr != "" {
		executeAndExit(*calcExpr)
		return
	}

	// Otherwise, start the REPL
	repl := display.NewREPL()
	repl.Run()
}

// executeFile runs a .calc script file line-by-line, printing results to stdout.
// Commands (lines starting with :) are executed and their messages printed; comment-only lines are ignored.
func executeFile(path string) error {
	var b []byte
	var err error

	if path == "-" {
		// Read from stdin
		b, err = io.ReadAll(os.Stdin)
	} else {
		b, err = os.ReadFile(path)
	}
	if err != nil {
		return err
	}

	repl := display.NewREPL()
	repl.SetSilent(true)

	// Iterate over lines to preserve REPL semantics (variables, settings, commands)
	lines := strings.Split(string(b), "\n")
	for _, ln := range lines {
		input := strings.TrimSpace(ln)
		if input == "" || strings.HasPrefix(input, "#") {
			continue
		}
		v := repl.EvaluateLine(input)
		// Skip sentinel no-op (commands or comment-only handled by EvaluateLine)
		if v.IsError() {
			if v.Error == "" {
				continue
			}
			// Print errors to stderr to mimic typical CLI behavior
			fmt.Fprintln(os.Stderr, repl.Formatter().Format(v))
			continue
		}
		// Print formatted value to stdout
		fmt.Println(repl.Formatter().Format(v))
	}

	return nil
}

func executeAndExit(input string) {
	// Create lexer and tokenise input
	l := lexer.New(input)
	tokens := l.AllTokens()

	// Parse tokens into AST
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create evaluator and evaluate expression
	env := evaluator.NewEnvironment()
	eval := evaluator.New(env)
	result := eval.Eval(expr)

	// Format and print result
	s := settings.Default()
	f := formatter.New(s)
	output := f.Format(result)

	if result.IsError() {
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Println(output)
}
