package main

import (
	"bufio"
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

// argsMap is a custom flag type for repeated --arg flags
type argsMap map[string]string

func (a argsMap) String() string {
	return fmt.Sprintf("%v", map[string]string(a))
}

func (a argsMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("argument must be in format name=value")
	}
	a[parts[0]] = parts[1]
	return nil
}

const helpText = `Calc - Terminal Notepad Calculator

A local-only terminal calculator.
Mix arithmetic, units, dates, currencies, and variables with natural language.

USAGE:
	calc                Start interactive REPL mode
	calc -c "expr"      Execute a single calculation and exit
	calc -f file.calc    Execute all lines from a file and print results

OPTIONS:
	-c string           Execute calculation and exit
	-f string           Execute a .calc file and print results
	-a, --arg name=value  Pass argument to script (can be repeated)
	--arg-file path     Read arguments from a file (key=value format)
	-h, --help          Show this help message

EXAMPLES:
	calc -c "half of 40"
	calc -c "10 m in cm"
	calc -c "11:00 - 09:00"
	calc -c "20% of 100"
	calc -c "£100 + £50"
	calc -f examples/k8s-cluster.calc
	calc -f script.calc --arg count=5 --arg rate=10
	calc -f script.calc --arg-file args.env

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
  • Script arguments (:arg directive with --arg flags)

REPL COMMANDS:
  :help              Show available commands
  :set precision N   Set decimal precision
  :set currency C    Set default currency (GBP, USD, EUR, JPY)
	:quiet [on|off]    Toggle or set quiet mode (suppress assignment output)
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
	argFile := flag.String("arg-file", "", "Read arguments from a file")
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message")
	
	// Custom argsMap for repeated --arg flags
	args := make(argsMap)
	flag.Var(&args, "arg", "Pass argument to script (name=value)")
	flag.Var(&args, "a", "Pass argument to script (name=value)")
	
	flag.Parse()

	// Show help if requested
	if *showHelp {
		fmt.Print(helpText)
		os.Exit(0)
	}

	// Load arguments from file if specified
	if *argFile != "" {
		fileArgs, err := loadArgsFromFile(*argFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading arg file: %v\n", err)
			os.Exit(1)
		}
		// Merge file args with CLI args (CLI args override file args)
		for k, v := range fileArgs {
			if _, exists := args[k]; !exists {
				args[k] = v
			}
		}
	}

	// If -f flag is provided, execute file and exit
	if *filePath != "" {
		if err := executeFile(*filePath, args); err != nil {
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

// loadArgsFromFile reads arguments from a .env-style file
func loadArgsFromFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	args := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			args[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return args, nil
}

// executeFile runs a .calc script file line-by-line, printing results to stdout.
// Commands (lines starting with :) are executed and their messages printed; comment-only lines are ignored.
func executeFile(path string, providedArgs map[string]string) error {
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

	// First pass: collect all :arg directives
	requiredArgs := make(map[string]string) // name -> prompt
	lines := strings.Split(string(b), "\n")
	
	for _, ln := range lines {
		input := strings.TrimSpace(ln)
		if input == "" || strings.HasPrefix(input, "#") {
			continue
		}
		
		// Parse to check if it's an :arg directive
		lex := lexer.New(input)
		tokens := lex.AllTokens()
		if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
			tokens = tokens[:len(tokens)-1]
		}
		if len(tokens) == 0 {
			continue
		}
		
		p := parser.New(tokens)
		expr, parseErr := p.Parse()
		if parseErr != nil {
			continue
		}
		
		if argDir, ok := expr.(*parser.ArgDirectiveExpr); ok {
			requiredArgs[argDir.Name] = argDir.Prompt
		}
	}

	// Process arguments: use provided args or prompt for missing ones
	for name, prompt := range requiredArgs {
		if val, exists := providedArgs[name]; exists {
			// Parse the provided value through lexer/parser for rich input
			if err := setArgVariable(repl, name, val); err != nil {
				return fmt.Errorf("error setting argument %s: %v", name, err)
			}
		} else {
			// Prompt user for the argument
			if prompt == "" {
				prompt = fmt.Sprintf("Enter value for %s:", name)
			}
			fmt.Printf("%s ", prompt)
			
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading argument %s: %v", name, err)
			}
			response = strings.TrimSpace(response)
			
			// Parse the response through lexer/parser for rich input
			if err := setArgVariable(repl, name, response); err != nil {
				return fmt.Errorf("error setting argument %s: %v", name, err)
			}
		}
	}

	// Second pass: execute the script
	for _, ln := range lines {
		input := strings.TrimSpace(ln)
		if input == "" || strings.HasPrefix(input, "#") {
			continue
		}
		
		// Parse to check if it's an :arg directive (skip execution)
		lex := lexer.New(input)
		tokens := lex.AllTokens()
		if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
			tokens = tokens[:len(tokens)-1]
		}
		if len(tokens) == 0 {
			continue
		}
		
		p := parser.New(tokens)
		expr, parseErr := p.Parse()
		if parseErr == nil {
			if _, ok := expr.(*parser.ArgDirectiveExpr); ok {
				// Skip :arg directives in execution phase
				continue
			}
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

// setArgVariable parses a string value and sets it as a variable in the REPL environment
func setArgVariable(repl *display.REPL, name, value string) error {
	// Parse the value through lexer/parser to support units, currency, expressions, etc.
	lex := lexer.New(value)
	tokens := lex.AllTokens()
	if len(tokens) > 0 && tokens[len(tokens)-1].Type == lexer.TokenEOF {
		tokens = tokens[:len(tokens)-1]
	}
	
	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return err
	}
	
	// Evaluate the expression
	result := repl.Env().Eval(expr)
	if result.IsError() {
		return fmt.Errorf("%s", result.Error)
	}
	
	// Set the variable
	repl.Env().SetVariable(name, result)
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
