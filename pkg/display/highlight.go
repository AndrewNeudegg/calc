package display

import (
	"strings"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

// Highlighter colorizes input using a Theme and the project's lexer.
type Highlighter struct {
	theme *Theme
}

func NewHighlighter(theme *Theme) *Highlighter {
	return &Highlighter{theme: theme}
}

// Colorize returns an ANSI-colored version of the input for display.
// It preserves original spacing by interleaving tokenized segments into the raw input.
func (h *Highlighter) Colorize(input string) string {
	if input == "" {
		return input
	}

	// Special-case commands: first non-space ':' then an ident
	trimmed := strings.TrimLeft(input, " \t")
	if strings.HasPrefix(trimmed, ":") {
		// Color only the ":<ident>" prefix; leave args plain
		// Find command word after ':'
		i := strings.Index(input, ":")
		start := i
		j := i + 1
		for j < len(input) {
			c := input[j]
			if !(isAlphaNum(c) || c == '_') {
				break
			}
			j++
		}
		cmd := input[start:j]
		colored := input[:start] + h.theme.wrap(cmd, h.theme.Command) + input[j:]
		return colored
	}

	// Tokenize with the lexer and stitch back with colors applied per token type
	l := lexer.New(input)
	toks := l.AllTokens()
	var b strings.Builder
	pos := 0
	for _, tok := range toks {
		if tok.Type == lexer.TokenEOF {
			break
		}
		lit := tok.Literal
		if lit == "" {
			continue
		}
		// Find next occurrence of token literal from current position
		idx := strings.Index(input[pos:], lit)
		if idx < 0 {
			// If not found (shouldn't happen), dump rest as-is and stop
			b.WriteString(input[pos:])
			pos = len(input)
			break
		}
		// Write any gap (whitespace or punctuation) unchanged
		b.WriteString(input[pos : pos+idx])
		pos += idx
		// Write colored token
		b.WriteString(h.colorToken(tok.Type, lit))
		pos += len(lit)
	}
	if pos < len(input) {
		b.WriteString(input[pos:])
	}
	return b.String()
}

func (h *Highlighter) colorToken(tt lexer.TokenType, s string) string {
	t := h.theme
	switch tt {
	case lexer.TokenNumber:
		return t.wrap(s, t.Number)
	case lexer.TokenUnit:
		return t.wrap(s, t.Unit)
	case lexer.TokenCurrency:
		return t.wrap(s, t.Currency)
	case lexer.TokenPlus, lexer.TokenMinus, lexer.TokenMultiply, lexer.TokenDivide, lexer.TokenPercent, lexer.TokenEquals,
		lexer.TokenLParen, lexer.TokenRParen, lexer.TokenComma:
		return t.wrap(s, t.Operator)
	case lexer.TokenIn, lexer.TokenOf, lexer.TokenPer, lexer.TokenBy, lexer.TokenWhat, lexer.TokenIs,
		lexer.TokenIncrease, lexer.TokenDecrease, lexer.TokenSum, lexer.TokenAverage, lexer.TokenMean, lexer.TokenTotal,
		lexer.TokenHalf, lexer.TokenDouble, lexer.TokenTwice, lexer.TokenQuarters, lexer.TokenThree, lexer.TokenAfter,
		lexer.TokenBefore, lexer.TokenFrom, lexer.TokenAgo, lexer.TokenNow, lexer.TokenToday, lexer.TokenTomorrow,
		lexer.TokenYesterday, lexer.TokenNext, lexer.TokenLast:
		return t.wrap(s, t.Keyword)
	case lexer.TokenDate:
		return t.wrap(s, t.Date)
	case lexer.TokenTime, lexer.TokenTimeValue:
		return t.wrap(s, t.Time)
	case lexer.TokenIdent:
		return t.wrap(s, t.Ident)
	default:
		return s
	}
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
