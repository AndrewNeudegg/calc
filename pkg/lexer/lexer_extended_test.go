package lexer

import (
	"testing"
)

// TestLexerMalformedNumbers tests edge cases with numbers
func TestLexerMalformedNumbers(t *testing.T) {
	tests := []struct {
		input       string
		description string
	}{
		{"123.456.789", "multiple decimals"},
		{".", "just a decimal"},
		{".123", "leading decimal"},
		{"123.", "trailing decimal"},
		{"0.0", "zero with decimal"},
		{"00.00", "leading zeros"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		if len(tokens) == 0 {
			t.Errorf("%s (%q): got no tokens", tt.description, tt.input)
		}
		// Just check it tokenizes without panic
	}
}

// TestLexerAllKeywords tests all keyword tokenization
func TestLexerAllKeywords(t *testing.T) {
	keywords := []struct {
		input    string
		expected TokenType
	}{
		{"in", TokenIn},
		{"of", TokenOf},
		{"per", TokenPer},
		{"by", TokenBy},
		{"what", TokenWhat},
		{"is", TokenIs},
		{"increase", TokenIncrease},
		{"decrease", TokenDecrease},
		{"sum", TokenSum},
		{"average", TokenAverage},
		{"mean", TokenMean},
		{"total", TokenTotal},
		{"half", TokenHalf},
		{"double", TokenDouble},
		{"twice", TokenTwice},
		{"quarters", TokenQuarters},
		{"three", TokenThree},
		{"after", TokenAfter},
		{"before", TokenBefore},
		{"from", TokenFrom},
		{"ago", TokenAgo},
		{"now", TokenNow},
		{"today", TokenToday},
		{"tomorrow", TokenTomorrow},
		{"yesterday", TokenYesterday},
		{"next", TokenNext},
		{"last", TokenLast},
		{"monday", TokenMonday},
		{"tuesday", TokenTuesday},
		{"wednesday", TokenWednesday},
		{"thursday", TokenThursday},
		{"friday", TokenFriday},
		{"saturday", TokenSaturday},
		{"sunday", TokenSunday},
	}

	for _, tt := range keywords {
		l := New(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expected {
			t.Errorf("keyword %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

// TestLexerCaseInsensitiveKeywords tests keywords are case-insensitive
func TestLexerCaseInsensitiveKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"TODAY", TokenToday},
		{"Today", TokenToday},
		{"ToDay", TokenToday},
		{"HALF", TokenHalf},
		{"Half", TokenHalf},
		{"MONDAY", TokenMonday},
		{"Monday", TokenMonday},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expected {
			t.Errorf("case variant %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

// TestLexerAllCurrencySymbols tests all currency symbols
func TestLexerAllCurrencySymbols(t *testing.T) {
	tests := []string{
		"$", "£", "€", "¥",
	}

	for _, input := range tests {
		l := New(input + "100")
		tok := l.NextToken()
		if tok.Type != TokenCurrency {
			t.Errorf("currency %q: expected TokenCurrency, got %s", input, tok.Type)
		}
	}
}

// TestLexerAllOperators tests all operator tokens
func TestLexerAllOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"+", TokenPlus},
		{"-", TokenMinus},
		{"*", TokenMultiply},
		{"/", TokenDivide},
		{"%", TokenPercent},
		{"=", TokenEquals},
		{"(", TokenLParen},
		{")", TokenRParen},
		{",", TokenComma},
		{":", TokenColon},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()
		if tok.Type != tt.expected {
			t.Errorf("operator %q: expected %s, got %s", tt.input, tt.expected, tok.Type)
		}
	}
}

// TestLexerComplexExpressions tests tokenization of complex expressions
func TestLexerComplexExpressions(t *testing.T) {
	tests := []struct {
		input         string
		expectedCount int
	}{
		{"10 + 20 * 30", 6}, // number, plus, number, multiply, number, EOF
		{"(100 + 50) / 2", 8},
		{"x = 10 + 5", 6},
		{"sum(1, 2, 3)", 9}, // sum, lparen, 1, comma, 2, comma, 3, rparen, EOF
		{"$100 + £50", 6},   // currency, number, plus, currency, number, EOF
		{"10 m in cm", 5},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		if len(tokens) != tt.expectedCount {
			t.Errorf("%q: expected %d tokens, got %d", tt.input, tt.expectedCount, len(tokens))
			for i, tok := range tokens {
				t.Logf("  [%d] %s: %q", i, tok.Type, tok.Literal)
			}
		}
	}
}

// TestLexerWhitespace tests whitespace handling
func TestLexerWhitespace(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"  10  +  20  "},
		{"\t10\t+\t20\t"},
		{"10\n+\n20"},
		{"   "},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()
		// Check it doesn't panic and produces tokens
		if len(tokens) == 0 {
			t.Errorf("%q: got no tokens", tt.input)
		}
	}
}

// TestLexerLineAndColumn tests line and column tracking
func TestLexerLineAndColumn(t *testing.T) {
	input := "10\n+ 20"
	l := New(input)

	tok1 := l.NextToken() // 10
	if tok1.Line != 1 || tok1.Column != 1 {
		t.Errorf("first token: expected line=1, col=1, got line=%d, col=%d", tok1.Line, tok1.Column)
	}

	tok2 := l.NextToken() // +
	if tok2.Line != 2 {
		t.Errorf("plus token: expected line=2, got line=%d", tok2.Line)
	}
}

// TestLexerUnknownCharacters tests error handling for unknown characters
func TestLexerUnknownCharacters(t *testing.T) {
	tests := []string{
		"@",
		"#",
		"&",
		"[",
		"]",
		"{",
		"}",
	}

	for _, input := range tests {
		l := New(input)
		tok := l.NextToken()
		if tok.Type != TokenError {
			t.Errorf("unknown char %q: expected TokenError, got %s", input, tok.Type)
		}
	}
}

// TestLexerIdentifiers tests identifier tokenization
func TestLexerIdentifiers(t *testing.T) {
	tests := []string{
		"x",
		"alpha",
		"_var",
		"var123",
		"CamelCase",
	}

	for _, input := range tests {
		l := New(input)
		tok := l.NextToken()
		if tok.Type != TokenIdent {
			t.Errorf("identifier %q: expected TokenIdent, got %s", input, tok.Type)
		}
		if tok.Literal != input {
			t.Errorf("identifier %q: literal mismatch, got %q", input, tok.Literal)
		}
	}
}

// TestLexerAllUnits tests all known unit types
func TestLexerAllUnits(t *testing.T) {
	units := []string{
		// Length
		"m", "cm", "mm", "km", "ft", "yd", "mi",
		"mile", "metre", "meter", "foot", "feet", "inch", "yard",
		// Mass
		"g", "kg", "mg", "lb", "lbs", "oz",
		"gram", "kilogram", "pound", "ounce",
		"stone", "st", "tonne", "tonnes", "ton", "tons",
		// Time
		"s", "sec", "second", "min", "minute", "h", "hr", "hour",
		"day", "week", "month", "year",
		// Volume
		"l", "ml", "litre", "liter", "gal", "gallon",
		// Temperature
		"c", "f", "celsius", "fahrenheit",
	}

	for _, unit := range units {
		l := New("10 " + unit)
		tokens := l.AllTokens()

		foundUnit := false
		for _, tok := range tokens {
			if tok.Type == TokenUnit && tok.Literal == unit {
				foundUnit = true
				break
			}
		}

		if !foundUnit {
			t.Errorf("unit %q: not recognized as TokenUnit", unit)
		}
	}
}

// TestLexerEmptyInput tests empty input handling
func TestLexerEmptyInput(t *testing.T) {
	l := New("")
	tok := l.NextToken()
	if tok.Type != TokenEOF {
		t.Errorf("empty input: expected TokenEOF, got %s", tok.Type)
	}
}

// TestLexerNumbersWithUnits tests combined numbers and units
func TestLexerNumbersWithUnits(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"10m"},
		{"10 m"},
		{"  10   m  "},
		{"3.14 cm"},
		{"100.5 kg"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.AllTokens()

		hasNumber := false
		hasUnit := false
		for _, tok := range tokens {
			if tok.Type == TokenNumber {
				hasNumber = true
			}
			if tok.Type == TokenUnit {
				hasUnit = true
			}
		}

		if !hasNumber || !hasUnit {
			t.Errorf("%q: hasNumber=%v, hasUnit=%v (expected both true)",
				tt.input, hasNumber, hasUnit)
		}
	}
}

// TestLexerMultipleTokens tests AllTokens method
func TestLexerMultipleTokens(t *testing.T) {
	input := "10 + 20 * 30"
	l := New(input)
	tokens := l.AllTokens()

	if len(tokens) == 0 {
		t.Fatal("AllTokens returned empty slice")
	}

	lastToken := tokens[len(tokens)-1]
	if lastToken.Type != TokenEOF {
		t.Errorf("last token should be EOF, got %s", lastToken.Type)
	}
}
