package display

// Theme defines ANSI color styles for different token classes in the REPL.
// Values should be full SGR sequences like "\x1b[33m"; Reset should be "\x1b[0m".
type Theme struct {
	// Prompt color (e.g., grey)
	Prompt string
	// Command (":help", ":save")
	Command string
	// Numbers (123, 3.14)
	Number string
	// Units (m, km, day, GBP, etc.)
	Unit string
	// Currency symbols and codes
	Currency string
	// Operators (+ - * / % = ( ) ,)
	Operator string
	// Keywords (in, of, per, by, what, is, increase, decrease, next, last, etc.)
	Keyword string
	// Date/time literals
	Date string
	Time string
	// Identifiers (variables)
	Ident string
	// Error text
	Error string
	// Reset sequence
	Reset string
}

// DefaultTheme returns a subtle, readable default theme.
func DefaultTheme() *Theme {
	return &Theme{
		Prompt:   "\x1b[90m", // bright black (grey)
		Command:  "\x1b[33m", // yellow
		Number:   "\x1b[36m", // cyan
		Unit:     "\x1b[34m", // blue
		Currency: "\x1b[32m", // green
		Operator: "\x1b[90m", // grey
		Keyword:  "\x1b[35m", // magenta
		Date:     "\x1b[36m", // cyan
		Time:     "\x1b[36m", // cyan
		Ident:    "\x1b[37m", // white (default-ish)
		Error:    "\x1b[31m", // red
		Reset:    "\x1b[0m",
	}
}

func (t *Theme) wrap(s, style string) string {
	if s == "" || style == "" {
		return s
	}
	return style + s + t.Reset
}
