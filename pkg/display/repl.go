package display

import (
	"bufio"
	"fmt"
	"io"
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
	lines     map[int]*Line
	nextID    int
	env       *evaluator.Environment
	eval      *evaluator.Evaluator
	formatter *formatter.Formatter
	commands  *commands.Handler
	settings  *settings.Settings
	depGraph  *graph.Graph
	theme     *Theme
	silent    bool
	quiet     bool
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

	r := &REPL{
		lines:     make(map[int]*Line),
		nextID:    1,
		env:       env,
		eval:      evaluator.New(env),
		formatter: formatter.New(sett),
		commands:  commands.New(sett),
		settings:  sett,
		depGraph:  graph.NewGraph(),
		theme:     DefaultTheme(),
	}
	
	// Set up history function for prev support
	env.SetHistoryFunc(r.getHistoryValue)
	
	// Wire workspace handlers for :save and :open
	r.commands.SaveWorkspace = r.saveWorkspace
	r.commands.LoadWorkspace = r.loadWorkspace
	// Wire clear handler for :clear
	r.commands.ClearWorkspace = r.clearWorkspace
	// Wire quiet controls
	r.commands.SetQuiet = r.SetQuiet
	r.commands.ToggleQuiet = r.ToggleQuiet
	r.commands.GetQuiet = r.IsQuiet
	return r
}

// Run starts the REPL loop.
func (r *REPL) Run() {
	fmt.Println("Calc - A terminal notepad calculator")
	fmt.Println("Type :help for available commands, :quit to exit")
	fmt.Println()

	// Try to use interactive line editor with control key support.
	// If it fails (e.g., not a TTY), fall back to simple Scanner.
	if isATTY(os.Stdin.Fd()) && isATTY(os.Stdout.Fd()) {
		r.runInteractive()
		return
	}

	// Fallback: basic line-by-line input
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

// runInteractive runs the REPL with a minimal line editor that supports control characters.
func (r *REPL) runInteractive() {
	reader := bufio.NewReader(os.Stdin)
	// Enable raw mode; ensure we restore on exit
	state, err := enableRawMode(int(os.Stdin.Fd()))
	if err != nil {
		// Fall back if raw mode cannot be enabled
		r.Run()
		return
	}
	defer restoreRawMode(int(os.Stdin.Fd()), state)

	for {
		rawPrompt := fmt.Sprintf("%d> ", r.nextID)
		prompt := r.theme.wrap(rawPrompt, r.theme.Prompt) + r.theme.Reset
		ed := NewEditor(prompt, r.collectHistory())
		// Install syntax highlighter for the buffer
		hl := NewHighlighter(r.theme)
		ed.SetHighlighter(hl.Colorize)
		line, aborted, eof := ed.ReadLine(reader, os.Stdout)
		if eof {
			fmt.Fprintln(os.Stdout)
			break
		}
		if aborted {
			// Show a helpful tip when Ctrl-C is pressed in raw mode
			printWithCRLF(os.Stdout, ctrlCTip())
			continue
		}
		input := strings.TrimSpace(line)
		if input == "" {
			fmt.Fprintln(os.Stdout)
			continue
		}
		result := r.EvaluateLine(input)
		if !result.IsError() || result.Error != "" {
			fmt.Fprintf(os.Stdout, "   = %s\n\n", r.formatter.Format(result))
		}
	}
}

func (r *REPL) collectHistory() []string {
	var h []string
	for i := 1; i < r.nextID; i++ {
		if line, ok := r.lines[i]; ok {
			if strings.TrimSpace(line.Input) != "" {
				h = append(h, line.Input)
			}
		}
	}
	return h
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

	// If the line reduces to nothing (e.g., comment-only or whitespace), treat as no-op
	if len(tokens) == 0 {
		return evaluator.NewError("")
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
		if !r.silent {
			printWithCRLF(os.Stdout, msg)
		}
		// Return a sentinel error value with empty message so caller skips printing a result line.
		return evaluator.NewError("")
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

	// Quiet mode: suppress printing for assignment lines
	if r.quiet {
		if _, isAssign := expr.(*parser.AssignExpr); isAssign {
			return evaluator.NewError("")
		}
	}

	return result
}

// clearWorkspace resets the current REPL session: history, variables, and evaluation state.
func (r *REPL) clearWorkspace() error {
	// Reset stored lines and prompt counter
	r.lines = make(map[int]*Line)
	r.nextID = 1

	// Reset evaluation environment and evaluator (clears variables and systems)
	r.env = evaluator.NewEnvironment()
	r.eval = evaluator.New(r.env)
	
	// Re-wire history function
	r.env.SetHistoryFunc(r.getHistoryValue)

	// Reset dependency graph
	r.depGraph = graph.NewGraph()

	return nil
}

// saveWorkspace writes the current REPL inputs to a file.
func (r *REPL) saveWorkspace(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	// Optional header
	fmt.Fprintln(f, "# calc workspace")
	for _, line := range r.ListLines() {
		if strings.TrimSpace(line.Input) == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line.Input), ":") {
			// Do not persist command lines
			continue
		}
		fmt.Fprintln(f, line.Input)
	}
	return nil
}

// loadWorkspace loads inputs from a file, replacing current session.
func (r *REPL) loadWorkspace(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// Reset state
	r.lines = make(map[int]*Line)
	r.nextID = 1
	r.env = evaluator.NewEnvironment()
	r.eval = evaluator.New(r.env)
	
	// Re-wire history function
	r.env.SetHistoryFunc(r.getHistoryValue)

	lines := strings.Split(string(b), "\n")
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		// Skip embedded commands in workspace files
		if strings.HasPrefix(t, ":") {
			continue
		}
		// Evaluate silently (no printing here)
		_ = r.EvaluateLine(t)
	}
	return nil
}

// printWithCRLF writes a possibly multi-line message ensuring lines start at column 0
// by converting bare "\n" linefeeds into "\r\n". This avoids ragged left margins when
// the REPL is in raw mode on terminals where LF does not imply carriage return.
func printWithCRLF(w io.Writer, s string) {
	if s == "" {
		return
	}
	// Normalize all standalone \n into \r\n
	// First, ensure we don't double-convert existing CRLF sequences: replace them temporarily
	// with a placeholder, convert remaining LFs, then restore CRLF.
	const ph = "\x00\x00CRLF_PLACEHOLDER\x00\x00"
	s = strings.ReplaceAll(s, "\r\n", ph)
	s = strings.ReplaceAll(s, "\n", "\r\n")
	s = strings.ReplaceAll(s, ph, "\r\n")
	// Ensure trailing newline for a clean break after the block
	if !strings.HasSuffix(s, "\r\n") {
		s += "\r\n"
	}
	fmt.Fprint(w, s)
}

// ctrlCTip returns the message shown when the user presses Ctrl-C in raw mode.
func ctrlCTip() string {
	return "Tip: Ctrl-C cancels the current line. Press Ctrl-D to exit, or type :help for commands."
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

// Formatter returns the formatter used by the REPL. This reflects live settings
// updates made via :set commands during the session.
func (r *REPL) Formatter() *formatter.Formatter {
	return r.formatter
}

// SetSilent toggles printing of command outputs during EvaluateLine. Useful for batch/script mode.
func (r *REPL) SetSilent(s bool) {
	r.silent = s
}

// SetQuiet enables or disables quiet mode (suppresses assignment output).
func (r *REPL) SetQuiet(q bool) {
	r.quiet = q
}

// ToggleQuiet flips quiet mode and returns the new state.
func (r *REPL) ToggleQuiet() bool {
	r.quiet = !r.quiet
	return r.quiet
}

// IsQuiet reports whether quiet mode is enabled.
func (r *REPL) IsQuiet() bool {
	return r.quiet
}

// getHistoryValue retrieves a previous result by offset.
// offset 0 means the most recent result (previous line),
// offset 1 means the result before that, etc.
func (r *REPL) getHistoryValue(offset int) (evaluator.Value, error) {
	// Calculate the line ID to retrieve
	targetID := r.nextID - 1 - offset
	
	if targetID < 1 {
		return evaluator.Value{}, fmt.Errorf("no previous result at offset %d", offset)
	}
	
	line, ok := r.lines[targetID]
	if !ok {
		return evaluator.Value{}, fmt.Errorf("no result found for prev~%d", offset)
	}
	
	// Return the result of that line
	return line.Result, nil
}
