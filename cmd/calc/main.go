package main

import (
	"flag"
	"fmt"
	"os"

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
  calc              Start interactive REPL mode
  calc -c "expr"    Execute a single calculation and exit

OPTIONS:
  -c string         Execute calculation and exit
  -h, --help        Show this help message

EXAMPLES:
  calc -c "half of 40"
  calc -c "10 m in cm"
  calc -c "11:00 - 09:00"
  calc -c "20% of 100"
  calc -c "£100 + £50"

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
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message")
	flag.Parse()

	// Show help if requested
	if *showHelp {
		fmt.Print(helpText)
		os.Exit(0)
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
