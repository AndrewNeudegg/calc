package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andrewneudegg/calc/pkg/lexer"
)

// Parser parses tokens into an AST.
type Parser struct {
	tokens []lexer.Token
	pos    int
	locale string // Locale for number parsing (e.g., "en_GB", "en_US")
}

// New creates a new parser from tokens with default UK locale.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		locale: "en_GB", // Default to UK format
	}
}

// NewWithLocale creates a new parser from tokens with a specific locale.
func NewWithLocale(tokens []lexer.Token, locale string) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		locale: locale,
	}
}

// Parse parses the tokens and returns an expression.
func (p *Parser) Parse() (Expr, error) {
	return p.parseExpression()
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek(offset int) lexer.Token {
	pos := p.pos + offset
	if pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[pos]
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) expect(typ lexer.TokenType) (lexer.Token, error) {
	tok := p.current()
	if tok.Type != typ {
		return tok, fmt.Errorf("expected %s, got %s", typ, tok.Type)
	}
	p.advance()
	return tok, nil
}

// isCurrencyCode checks if a unit string is a currency code or name
func (p *Parser) isCurrencyCode(unit string) bool {
	lower := strings.ToLower(unit)
	switch lower {
	case "usd", "dollar", "dollars",
		"gbp",
		"eur", "euro", "euros",
		"jpy", "yen",
		// Expanded 3-letter codes
		"aud", "cad", "nzd",
		"chf", "sek", "nok", "dkk",
		"pln", "czk", "huf", "ron",
		"rub", "try",
		"aed", "sar", "ils",
		"cny", "hkd", "sgd", "inr", "krw", "twd", "thb", "myr", "idr", "php",
		"mxn", "brl", "zar":
		return true
	default:
		return false
	}
}

// normalizeNumber converts a number string with thousand separators to a valid float string.
// It uses the parser's locale setting to determine how to interpret commas and periods.
// UK/US format (en_GB, en_US): comma as thousand separator, period as decimal (1,234.56)
// European format (de_DE, fr_FR, etc.): period as thousand separator, comma as decimal (1.234,56)
func (p *Parser) normalizeNumber(s string) string {
	// If there are no commas or periods, return as-is
	if !strings.Contains(s, ",") && !strings.Contains(s, ".") {
		return s
	}

	// Determine if locale uses European format (period=thousands, comma=decimal)
	isEuropeanLocale := p.isEuropeanLocale()

	// Count commas and periods
	commaCount := strings.Count(s, ",")
	periodCount := strings.Count(s, ".")

	if isEuropeanLocale {
		// European format: period = thousands, comma = decimal
		if commaCount == 0 && periodCount > 0 {
			// Only periods - they are thousand separators
			return strings.ReplaceAll(s, ".", "")
		}
		if periodCount == 0 && commaCount > 0 {
			// Only commas - they are decimal separators
			return strings.ReplaceAll(s, ",", ".")
		}
		// Both present - period is thousands, comma is decimal
		s = strings.ReplaceAll(s, ".", "")  // Remove thousand separators
		s = strings.ReplaceAll(s, ",", ".") // Convert decimal separator
		return s
	}

	// UK/US format: comma = thousands, period = decimal
	if commaCount == 0 && periodCount > 0 {
		// Only periods - they are decimal separators
		return s
	}
	if periodCount == 0 && commaCount > 0 {
		// Only commas - they are thousand separators
		return strings.ReplaceAll(s, ",", "")
	}
	// Both present - comma is thousands, period is decimal
	s = strings.ReplaceAll(s, ",", "") // Remove thousand separators
	return s
}

// isEuropeanLocale returns true if the locale uses European number format
// (period as thousand separator, comma as decimal separator)
func (p *Parser) isEuropeanLocale() bool {
	// Common European locales that use comma as decimal separator
	switch p.locale {
	case "de_DE", "de_AT", "de_CH", // German
		"fr_FR", "fr_BE", "fr_CH", // French
		"es_ES", "es_MX", "es_AR", // Spanish
		"it_IT", // Italian
		"nl_NL", "nl_BE", // Dutch
		"pt_PT", "pt_BR", // Portuguese
		"pl_PL", // Polish
		"ru_RU", // Russian
		"cs_CZ", // Czech
		"da_DK", // Danish
		"fi_FI", // Finnish
		"sv_SE", // Swedish
		"no_NO", // Norwegian
		"tr_TR", // Turkish
		"hu_HU", // Hungarian
		"ro_RO": // Romanian
		return true
	default:
		return false
	}
}

func (p *Parser) parseExpression() (Expr, error) {
	// Check for command
	if p.current().Type == lexer.TokenColon {
		return p.parseCommand()
	}

	// Check for assignment (allow keywords and units as variable names)
	if (p.current().Type == lexer.TokenIdent ||
		p.isKeywordToken(p.current().Type) ||
		p.current().Type == lexer.TokenUnit) &&
		p.peek(1).Type == lexer.TokenEquals {
		return p.parseAssignment()
	}

	// Try parsing timezone queries
	if expr, ok := p.tryParseTimezoneQuery(); ok {
		return expr, nil
	}

	// Try parsing fuzzy phrases first
	if expr, ok := p.tryParseFuzzyPhrase(); ok {
		return expr, nil
	}

	// Parse standard expression
	return p.parseConversion()
}

func (p *Parser) parseCommand() (Expr, error) {
	p.advance() // skip ':'

	if p.current().Type != lexer.TokenIdent && p.current().Type != lexer.TokenArg {
		return nil, fmt.Errorf("expected command name")
	}

	command := p.current().Literal
	p.advance()

	// Special handling for :arg directive
	if command == "arg" {
		return p.parseArgDirective()
	}

	// Reconstruct the remainder of the line into a raw tail string while
	// preserving filename/path punctuation like '.', '/', and '-' by gluing
	// those to adjacent tokens without spaces. Then split on spaces to get args.
	var tailBuilder strings.Builder
	glue := map[string]bool{".": true, "/": true, "-": true}
	wrote := false
	for p.current().Type != lexer.TokenEOF {
		lit := p.current().Literal
		if wrote {
			// Add a space before this token unless either side is a glue token
			if !glue[lit] {
				// Also peek previous written rune to handle trailing glue
				// We can't easily get previous token here, so approximate by
				// checking last rune written.
				if tailBuilder.Len() > 0 {
					last, _ := utf8DecLastRune(&tailBuilder)
					if last != '.' && last != '/' && last != '-' {
						tailBuilder.WriteByte(' ')
					}
				}
			}
		}
		tailBuilder.WriteString(lit)
		wrote = true
		p.advance()
	}
	tail := strings.TrimSpace(tailBuilder.String())
	var args []string
	if tail != "" {
		args = strings.Fields(tail)
	}

	return &CommandExpr{
		Command: command,
		Args:    args,
	}, nil
}

// parseArgDirective parses ":arg var_name "prompt text"" directives.
func (p *Parser) parseArgDirective() (Expr, error) {
	// Expect variable name (can be an identifier, a keyword, or a unit token used as variable name)
	if p.current().Type != lexer.TokenIdent &&
		!p.isKeywordToken(p.current().Type) &&
		p.current().Type != lexer.TokenUnit {
		return nil, fmt.Errorf("expected variable name after :arg")
	}

	varName := p.current().Literal
	p.advance()

	// Optional prompt string
	var prompt string
	if p.current().Type == lexer.TokenString {
		prompt = p.current().Literal
		p.advance()
	}

	return &ArgDirectiveExpr{
		Name:   varName,
		Prompt: prompt,
	}, nil
}

// isKeywordToken checks if a token type is a keyword that can be used as a variable name or identifier.
// Not all keywords are included—only those allowed in this context.
func (p *Parser) isKeywordToken(t lexer.TokenType) bool {
	switch t {
	case lexer.TokenIn, lexer.TokenOf, lexer.TokenPer, lexer.TokenBy,
		lexer.TokenWhat, lexer.TokenIs, lexer.TokenIncrease, lexer.TokenDecrease,
		lexer.TokenSum, lexer.TokenAverage, lexer.TokenMean, lexer.TokenTotal,
		lexer.TokenHalf, lexer.TokenDouble, lexer.TokenTwice, lexer.TokenQuarters,
		lexer.TokenThree, lexer.TokenArg, lexer.TokenAfter, lexer.TokenBefore,
		lexer.TokenFrom, lexer.TokenAgo, lexer.TokenNow, lexer.TokenToday,
		lexer.TokenTomorrow, lexer.TokenYesterday, lexer.TokenNext, lexer.TokenLast,
		lexer.TokenPrev, lexer.TokenTime, lexer.TokenMonday, lexer.TokenTuesday,
		lexer.TokenWednesday, lexer.TokenThursday, lexer.TokenFriday, lexer.TokenSaturday,
		lexer.TokenSunday, lexer.TokenJanuary, lexer.TokenFebruary, lexer.TokenMarch,
		lexer.TokenApril, lexer.TokenMay, lexer.TokenJune, lexer.TokenJuly,
		lexer.TokenAugust, lexer.TokenSeptember, lexer.TokenOctober, lexer.TokenNovember,
		lexer.TokenDecember:
		return true
	default:
		return false
	}
}

// utf8DecLastRune returns the last rune written in a strings.Builder and a bool indicating success.
func utf8DecLastRune(b *strings.Builder) (rune, bool) {
	// strings.Builder doesn't expose bytes; use String(), acceptable for short command tails
	s := b.String()
	if s == "" {
		return 0, false
	}
	// Walk back to decode last rune
	// Bytes are ASCII for our glue checks; this is sufficient and safe.
	return rune(s[len(s)-1]), true
}

func (p *Parser) parseAssignment() (Expr, error) {
	name := p.current().Literal
	p.advance() // skip identifier
	p.advance() // skip '='

	// Try parsing fuzzy phrases first in assignments
	if expr, ok := p.tryParseFuzzyPhrase(); ok {
		return &AssignExpr{
			Name:  name,
			Value: expr,
		}, nil
	}

	// Allow timezone queries on the right-hand side of assignments
	if expr, ok := p.tryParseTimezoneQuery(); ok {
		return &AssignExpr{
			Name:  name,
			Value: expr,
		}, nil
	}

	value, err := p.parseConversion()
	if err != nil {
		return nil, err
	}

	return &AssignExpr{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) tryParseFuzzyPhrase() (Expr, bool) {
	tok := p.current()

	// "half of X"
	if tok.Type == lexer.TokenHalf {
		p.advance()
		if p.current().Type == lexer.TokenOf {
			p.advance()
		}
		value, err := p.parseConversion()
		if err != nil {
			return nil, false
		}
		return &FuzzyExpr{Pattern: "half", Value: value}, true
	}

	// "double X" or "twice X"
	if tok.Type == lexer.TokenDouble || tok.Type == lexer.TokenTwice {
		pattern := tok.Literal
		p.advance()
		value, err := p.parseConversion()
		if err != nil {
			return nil, false
		}
		return &FuzzyExpr{Pattern: pattern, Value: value}, true
	}

	// "three quarters of X"
	if tok.Type == lexer.TokenThree && p.peek(1).Type == lexer.TokenQuarters {
		p.advance() // skip 'three'
		p.advance() // skip 'quarters'
		if p.current().Type == lexer.TokenOf {
			p.advance()
		}
		value, err := p.parseConversion()
		if err != nil {
			return nil, false
		}
		return &FuzzyExpr{Pattern: "three quarters", Value: value}, true
	}

	// "increase X by Y%"
	if tok.Type == lexer.TokenIncrease {
		p.advance()
		base, err := p.parseAdditive()
		if err != nil {
			return nil, false
		}
		if p.current().Type == lexer.TokenBy {
			p.advance()
			percent, err := p.parseAdditive()
			if err != nil {
				return nil, false
			}
			expr := &PercentChangeExpr{Base: base, Percent: percent, Increase: true}
			// Optional trailing conversion: "in <unit>"
			if wrapped, ok := p.tryWrapWithConversion(expr); ok {
				return wrapped, true
			}
			return expr, true
		}
	}

	// "decrease X by Y%"
	if tok.Type == lexer.TokenDecrease {
		p.advance()
		base, err := p.parseAdditive()
		if err != nil {
			return nil, false
		}
		if p.current().Type == lexer.TokenBy {
			p.advance()
			percent, err := p.parseAdditive()
			if err != nil {
				return nil, false
			}
			expr := &PercentChangeExpr{Base: base, Percent: percent, Increase: false}
			if wrapped, ok := p.tryWrapWithConversion(expr); ok {
				return wrapped, true
			}
			return expr, true
		}
	}

	// "X is what % of Y"
	if p.pos+3 < len(p.tokens) {
		if p.peek(1).Type == lexer.TokenIs && p.peek(2).Type == lexer.TokenWhat && p.peek(3).Type == lexer.TokenPercent {
			part, err := p.parseAdditive()
			if err != nil {
				return nil, false
			}
			p.advance() // 'is'
			p.advance() // 'what'
			p.advance() // '%'
			if p.current().Type == lexer.TokenOf {
				p.advance()
			}
			whole, err := p.parseAdditive()
			if err != nil {
				return nil, false
			}
			expr := &WhatPercentExpr{Part: part, Whole: whole}
			if wrapped, ok := p.tryWrapWithConversion(expr); ok {
				return wrapped, true
			}
			return expr, true
		}
	}

	return nil, false
}

// tryWrapWithConversion checks for a trailing "in ..." conversion and wraps the given expr
func (p *Parser) tryWrapWithConversion(expr Expr) (Expr, bool) {
	if p.current().Type != lexer.TokenIn {
		return nil, false
	}
	// Parse one or more chained conversions
	for p.current().Type == lexer.TokenIn {
		p.advance()
		toUnit := p.current().Literal
		p.advance()
		if p.current().Type == lexer.TokenPer {
			p.advance()
			if p.current().Type == lexer.TokenUnit {
				toUnit = toUnit + "/" + p.current().Literal
				p.advance()
			}
		} else if p.current().Type == lexer.TokenDivide {
			if p.peek(1).Type == lexer.TokenUnit {
				p.advance()
				toUnit = toUnit + "/" + p.current().Literal
				p.advance()
			}
		}
		expr = &ConversionExpr{Value: expr, ToUnit: toUnit}
	}
	return expr, true
}

// tryParseTimezoneQuery attempts to parse timezone-related queries
func (p *Parser) tryParseTimezoneQuery() (Expr, bool) {
	tok := p.current()

	// "time in <location>"
	if tok.Type == lexer.TokenTime && p.peek(1).Type == lexer.TokenIn {
		p.advance() // skip 'time'
		p.advance() // skip 'in'

		// Get location name (could be multi-word like "New York")
		location := p.parseLocationName()

		// Check for offset: "time in Sydney plus 3 hours in London"
		if p.current().Type == lexer.TokenPlus || p.current().Type == lexer.TokenMinus ||
			(p.current().Type == lexer.TokenIdent && (p.current().Literal == "plus" || p.current().Literal == "minus")) {
			operator := "+"
			if p.current().Type == lexer.TokenMinus || p.current().Literal == "minus" {
				operator = "-"
			}
			p.advance()

			// Parse offset value and unit (e.g., "3 hours")
			offsetValue, err := p.parseAdditive()
			if err != nil {
				return nil, false
			}

			// Check for "in <target_location>"
			if p.current().Type == lexer.TokenIn {
				p.advance()
				targetLocation := p.parseLocationName()
				return &TimeConversionExpr{
					Time:     nil, // current time
					From:     location,
					To:       targetLocation,
					Offset:   offsetValue,
					Operator: operator,
				}, true
			}
		}

		return &TimeInLocationExpr{Location: location}, true
	}

	// "time difference between <loc1> and <loc2>"
	// or "time difference <loc1> <loc2>"
	if tok.Type == lexer.TokenTime && p.peek(1).Type == lexer.TokenIdent && p.peek(1).Literal == "difference" {
		p.advance() // skip 'time'
		p.advance() // skip 'difference'

		// Optional "between"
		if p.current().Type == lexer.TokenIdent && p.current().Literal == "between" {
			p.advance()
		}

		from := p.parseLocationName()

		// Optional "and"
		if p.current().Type == lexer.TokenIdent && p.current().Literal == "and" {
			p.advance()
		}

		to := p.parseLocationName()

		// Optional "in [unit]" for target unit
		var targetUnit string
		if p.current().Type == lexer.TokenIn {
			p.advance() // skip 'in'
			if p.current().Type == lexer.TokenUnit || p.current().Type == lexer.TokenIdent {
				targetUnit = p.current().Literal
				p.advance()
			}
		}

		return &TimeDifferenceExpr{From: from, To: to, TargetUnit: targetUnit}, true
	}

	return nil, false
}

// parseLocationName parses a location name (can be multi-word like "New York")
func (p *Parser) parseLocationName() string {
	var parts []string

	// Known multi-word cities
	multiWordCities := []string{"New York", "Los Angeles", "Hong Kong"}

	// Try to match multi-word cities first
	for _, city := range multiWordCities {
		cityParts := strings.Split(city, " ")
		matches := true

		for i, part := range cityParts {
			if p.peek(i).Type != lexer.TokenIdent ||
				!strings.EqualFold(p.peek(i).Literal, part) {
				matches = false
				break
			}
		}

		if matches {
			// Consume all parts of the city name
			for range cityParts {
				parts = append(parts, p.current().Literal)
				p.advance()
			}
			return strings.Join(parts, " ")
		}
	}

	// Single-word city - just take one identifier
	if p.current().Type == lexer.TokenIdent {
		city := p.current().Literal
		p.advance()
		return city
	}

	return ""
}

func (p *Parser) parseConversion() (Expr, error) {
	// Parse the left-hand side expression first
	expr, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	// Handle one or more postfix "in ..." conversions that apply to the current expr
	for p.current().Type == lexer.TokenIn {
		p.advance()
		toUnit := p.current().Literal
		p.advance()

		// Check if this is a compound unit (e.g., "m/s" or "km per hour")
		if p.current().Type == lexer.TokenPer {
			p.advance()
			if p.current().Type == lexer.TokenUnit {
				toUnit = toUnit + "/" + p.current().Literal
				p.advance()
			}
		} else if p.current().Type == lexer.TokenDivide {
			// Look ahead to see if next token is a unit
			if p.peek(1).Type == lexer.TokenUnit {
				p.advance() // consume /
				toUnit = toUnit + "/" + p.current().Literal
				p.advance()
			}
		}

		expr = &ConversionExpr{Value: expr, ToUnit: toUnit}
	}

	// After applying any conversions, allow additive tail (e.g., "(a in x) + b")
	for {
		tok := p.current()
		var op string

		if tok.Type == lexer.TokenPlus {
			op = "+"
		} else if tok.Type == lexer.TokenMinus {
			op = "-"
		} else if tok.Type == lexer.TokenIdent {
			if tok.Literal == "plus" {
				op = "+"
			} else if tok.Literal == "minus" || tok.Literal == "and" {
				op = map[bool]string{true: "+", false: "-"}[tok.Literal == "plus" || tok.Literal == "and"]
			} else {
				break
			}
		} else {
			break
		}

		p.advance()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseAdditive() (Expr, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.current()
		var op string

		// Check for symbolic operators
		if tok.Type == lexer.TokenPlus {
			op = "+"
		} else if tok.Type == lexer.TokenMinus {
			op = "-"
		} else if tok.Type == lexer.TokenIdent {
			// Check for word operators
			if tok.Literal == "plus" {
				op = "+"
			} else if tok.Literal == "minus" {
				op = "-"
			} else if tok.Literal == "and" {
				// "and" acts as addition when connecting units or numbers
				op = "+"
			} else {
				break
			}
		} else {
			break
		}

		p.advance()

		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseMultiplicative() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.current()
		var op string

		// Check for symbolic operators
		if tok.Type == lexer.TokenMultiply {
			op = "*"
		} else if tok.Type == lexer.TokenDivide {
			op = "/"
		} else if tok.Type == lexer.TokenIdent {
			// Check for word operators
			if tok.Literal == "times" || tok.Literal == "multiplied" {
				op = "*"
			} else if tok.Literal == "divided" {
				// Check if followed by "by"
				if p.peek(1).Type == lexer.TokenBy {
					op = "/"
					p.advance() // consume "divided"
					p.advance() // consume "by"
					right, err := p.parseUnary()
					if err != nil {
						return nil, err
					}
					left = &BinaryExpr{
						Left:     left,
						Operator: op,
						Right:    right,
					}
					continue
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}

		p.advance()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (Expr, error) {
	tok := p.current()

	if tok.Type == lexer.TokenMinus {
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Operator: "-",
			Operand:  operand,
		}, nil
	}

	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Check for unit
	if p.current().Type == lexer.TokenUnit {
		unit := p.current().Literal
		p.advance()

		// Check if this unit is actually a currency code
		if p.isCurrencyCode(unit) {
			// Convert to CurrencyExpr
			expr = &CurrencyExpr{
				Value:    expr,
				Currency: unit,
			}

			// Check for "per" (rate) after currency - e.g., "32 dollars per day"
			if p.current().Type == lexer.TokenPer {
				p.advance()
				if p.current().Type == lexer.TokenUnit {
					unit2 := p.current().Literal
					p.advance()
					// Convert to a unit expression with currency/time rate
					// Store as a UnitExpr with compound unit like "$/day"
					currencySymbol := p.getCurrencySymbol(unit)
					expr = &UnitExpr{Value: expr, Unit: currencySymbol + "/" + unit2}
				}
			} else if p.current().Type == lexer.TokenDivide {
				// Look ahead to see if this is a rate (/ followed by unit)
				if p.peek(1).Type == lexer.TokenUnit {
					p.advance() // consume the /
					unit2 := p.current().Literal
					p.advance()
					currencySymbol := p.getCurrencySymbol(unit)
					expr = &UnitExpr{Value: expr, Unit: currencySymbol + "/" + unit2}
				}
			}
		} else {
			// Regular unit
			expr = &UnitExpr{Value: expr, Unit: unit}

			// Check for "per" (rate) - only consume / if immediately followed by a unit
			// If followed by a number, leave the / for the binary operator parser
			if p.current().Type == lexer.TokenPer {
				p.advance()
				if p.current().Type == lexer.TokenUnit {
					unit2 := p.current().Literal
					p.advance()
					expr = &UnitExpr{Value: expr, Unit: unit + "/" + unit2}
				}
			} else if p.current().Type == lexer.TokenDivide {
				// Look ahead to see if this is a rate (/ followed by unit) or division (/ followed by number)
				if p.peek(1).Type == lexer.TokenUnit {
					p.advance() // consume the /
					unit2 := p.current().Literal
					p.advance()
					expr = &UnitExpr{Value: expr, Unit: unit + "/" + unit2}
				}
				// Otherwise, leave the / for the binary operator parser to handle
			}
		}
	}

	// Check for compound unit rate after currency expression (e.g., "$2.93/hr" or "$2.93 per hour")
	// This handles prefix currency symbols like $, £, €, ¥
	if currExpr, ok := expr.(*CurrencyExpr); ok {
		if p.current().Type == lexer.TokenPer {
			p.advance()
			if p.current().Type == lexer.TokenUnit {
				unit := p.current().Literal
				p.advance()
				expr = &UnitExpr{Value: currExpr, Unit: currExpr.Currency + "/" + unit}
			}
		} else if p.current().Type == lexer.TokenDivide {
			// Look ahead to see if this is a rate (/ followed by unit)
			if p.peek(1).Type == lexer.TokenUnit {
				p.advance() // consume the /
				unit := p.current().Literal
				p.advance()
				expr = &UnitExpr{Value: currExpr, Unit: currExpr.Currency + "/" + unit}
			}
		}
	}

	// Check for percentage
	if p.current().Type == lexer.TokenPercent {
		p.advance()

		// Check for "of"
		if p.current().Type == lexer.TokenOf {
			p.advance()
			of, err := p.parseAdditive()
			if err != nil {
				return nil, err
			}
			return &PercentOfExpr{Percent: expr, Of: of}, nil
		}

		expr = &PercentExpr{Value: expr}
	}

	return expr, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	tok := p.current()

	switch tok.Type {
	case lexer.TokenNumber:
		normalized := p.normalizeNumber(tok.Literal)
		val, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", tok.Literal)
		}
		p.advance()
		return &NumberExpr{Value: val}, nil

	case lexer.TokenString:
		// String literal
		val := tok.Literal
		p.advance()
		return &StringExpr{Value: val}, nil

	case lexer.TokenCurrency:
		currency := tok.Literal
		p.advance()

		if p.current().Type != lexer.TokenNumber {
			return nil, fmt.Errorf("expected number after currency symbol")
		}

		normalized := p.normalizeNumber(p.current().Literal)
		val, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", p.current().Literal)
		}
		p.advance()

		return &CurrencyExpr{
			Value:    &NumberExpr{Value: val},
			Currency: currency,
		}, nil

	case lexer.TokenIdent:
		// Try to parse as number words first
		if val, ok := p.tryParseNumberWords(); ok {
			return &NumberExpr{Value: val}, nil
		}
		name := tok.Literal
		p.advance()

		// Check for function call
		if p.current().Type == lexer.TokenLParen {
			return p.parseFunctionCall(name)
		}

		return &IdentExpr{Name: name}, nil

	case lexer.TokenUnit:
		// Allow function names that collide with unit tokens, e.g., "min(...)"
		// Also allow unit tokens to be used as variable names
		name := tok.Literal
		p.advance()
		if p.current().Type == lexer.TokenLParen {
			return p.parseFunctionCall(name)
		}
		// Treat as a variable reference
		return &IdentExpr{Name: name}, nil

	case lexer.TokenConstant:
		// Physical constants are parsed as identifiers and resolved to their constant values during evaluation
		name := tok.Literal
		p.advance()
		return &IdentExpr{Name: name}, nil

	case lexer.TokenThree:
		// Could be "three quarters" or just "three" as a number
		if p.peek(1).Type == lexer.TokenQuarters {
			// It's "three quarters" - handle in tryParseFuzzyPhrase
			return nil, fmt.Errorf("unexpected token: %s", tok.Type)
		}
		// It's just "three" as a number word
		if val, ok := p.tryParseNumberWords(); ok {
			return &NumberExpr{Value: val}, nil
		}
		return nil, fmt.Errorf("unexpected token: %s", tok.Type)

	case lexer.TokenLParen:
		p.advance()
		// Allow conversions inside parentheses
		expr, err := p.parseConversion()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TokenRParen); err != nil {
			return nil, err
		}
		return expr, nil

	case lexer.TokenSum, lexer.TokenAverage, lexer.TokenMean, lexer.TokenTotal:
		return p.parseFunctionCall(tok.Literal)

	case lexer.TokenToday, lexer.TokenTomorrow, lexer.TokenYesterday:
		return p.parseDateKeyword()

	case lexer.TokenNext, lexer.TokenLast:
		return p.parseWeekday()

	case lexer.TokenMonday, lexer.TokenTuesday, lexer.TokenWednesday,
		lexer.TokenThursday, lexer.TokenFriday, lexer.TokenSaturday, lexer.TokenSunday:
		return p.parseWeekday()

	case lexer.TokenJanuary, lexer.TokenFebruary, lexer.TokenMarch, lexer.TokenApril,
		lexer.TokenMay, lexer.TokenJune, lexer.TokenJuly, lexer.TokenAugust,
		lexer.TokenSeptember, lexer.TokenOctober, lexer.TokenNovember, lexer.TokenDecember:
		return p.parseMonth()

	case lexer.TokenNow:
		p.advance()
		return &TimeExpr{Time: time.Now()}, nil

	case lexer.TokenTimeValue:
		// Parse time in HH:MM or HH:MM:SS format
		timeStr := tok.Literal
		parts := strings.Split(timeStr, ":")

		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("invalid time format: %s", timeStr)
		}

		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid hours in time: %s", parts[0])
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minutes in time: %s", parts[1])
		}

		seconds := 0
		if len(parts) == 3 {
			seconds, err = strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid seconds in time: %s", parts[2])
			}
		}

		// Convert to decimal hours and store as a time unit
		decimalHours := float64(hours) + float64(minutes)/60.0 + float64(seconds)/3600.0

		p.advance()
		// Return as a unit expression with "time" unit to preserve time format
		return &UnitExpr{
			Value: &NumberExpr{Value: decimalHours},
			Unit:  "time",
		}, nil

	case lexer.TokenDate:
		// Parse date in DD/MM/YYYY format (British style)
		// TODO: Add locale support for MM/DD/YYYY (US style)
		dateStr := tok.Literal
		p.advance()

		parts := strings.Split(dateStr, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid date format: %s", dateStr)
		}

		day, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid day in date: %s", parts[0])
		}

		month, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid month in date: %s", parts[1])
		}

		year, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid year in date: %s", parts[2])
		}

		// Validate ranges
		if day < 1 || day > 31 {
			return nil, fmt.Errorf("invalid day: %d", day)
		}
		if month < 1 || month > 12 {
			return nil, fmt.Errorf("invalid month: %d", month)
		}
		if year < 1000 || year > 9999 {
			return nil, fmt.Errorf("invalid year: %d", year)
		}

		// Create date (DD/MM/YYYY format)
		parsedDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		// Validate that the date is valid (e.g., not 31/02/2024)
		if parsedDate.Day() != day || parsedDate.Month() != time.Month(month) || parsedDate.Year() != year {
			return nil, fmt.Errorf("invalid date: %s (day/month/year out of range)", dateStr)
		}

		return &DateExpr{Date: parsedDate}, nil

	case lexer.TokenPrev:
		// Parse prev, prev~, prev~1, prev~5, prev#15, etc.
		literal := tok.Literal
		p.advance()

		// Default offset is 0 (just "prev"), relative mode
		offset := 0
		absolute := false

		// Check if the literal contains '~' (relative offset)
		if strings.Contains(literal, "~") {
			parts := strings.Split(literal, "~")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid prev syntax: %s (expected 'prev~N')", literal)
			}
			if parts[1] == "" {
				// "prev~" means offset 1
				offset = 1
			} else {
				// "prev~N" means offset N
				val, err := strconv.Atoi(parts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid prev offset: %s", parts[1])
				}
				if val < 0 {
					return nil, fmt.Errorf("prev offset must be non-negative, got: %d", val)
				}
				offset = val
			}
		} else if strings.Contains(literal, "#") {
			// Check if the literal contains '#' (absolute position)
			parts := strings.Split(literal, "#")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid prev syntax: %s (expected 'prev#N')", literal)
			}
			if parts[1] == "" {
				// "prev#" without a number is an error
				return nil, fmt.Errorf("prev# requires a line number (e.g., 'prev#15')")
			}
			// "prev#N" means absolute line number N
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid prev line number: %s", parts[1])
			}
			if val < 1 {
				return nil, fmt.Errorf("prev line number must be positive, got: %d", val)
			}
			offset = val
			absolute = true
		}

		return &PrevExpr{Offset: offset, Absolute: absolute}, nil

	default:
		// If it's a keyword token being used as a variable name, treat it as an identifier
		if p.isKeywordToken(tok.Type) {
			name := tok.Literal
			p.advance()
			// Check for function call
			if p.current().Type == lexer.TokenLParen {
				return p.parseFunctionCall(name)
			}
			return &IdentExpr{Name: name}, nil
		}
		return nil, fmt.Errorf("unexpected token: %s", tok.Type)
	}
}

func (p *Parser) parseFunctionCall(name string) (Expr, error) {
	p.advance() // skip function name if not already done

	if p.current().Type == lexer.TokenLParen {
		p.advance() // skip '('
	} else if p.current().Type == lexer.TokenOf {
		p.advance() // skip 'of' for natural language
	}

	var args []Expr

	for p.current().Type != lexer.TokenRParen && p.current().Type != lexer.TokenEOF {
		// Allow conversions within function arguments
		arg, err := p.parseConversion()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current().Type == lexer.TokenComma {
			p.advance()
		} else {
			break
		}
	}

	if p.current().Type == lexer.TokenRParen {
		p.advance()
	}

	return &FunctionCallExpr{
		Name: strings.ToLower(name),
		Args: args,
	}, nil
}

// getCurrencySymbol maps known currency names/codes to their canonical symbol
func (p *Parser) getCurrencySymbol(code string) string {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "$", "usd", "dollar", "dollars":
		return "$"
	case "£", "gbp":
		return "£"
	case "€", "eur", "euro", "euros":
		return "€"
	case "¥", "jpy", "yen":
		return "¥"
	// For other currencies without a unique 1-char symbol, return the uppercase code
	case "aud", "cad", "nzd", "chf", "sek", "nok", "dkk", "pln", "czk", "huf", "ron",
		"rub", "try", "aed", "sar", "ils", "cny", "hkd", "sgd", "inr", "krw", "twd",
		"thb", "myr", "idr", "php", "mxn", "brl", "zar":
		return strings.ToUpper(strings.TrimSpace(code))
	default:
		// Fallback: return as-is
		return code
	}
}

func (p *Parser) parseDateKeyword() (Expr, error) {
	tok := p.current()
	p.advance()

	var base time.Time
	switch tok.Type {
	case lexer.TokenToday:
		base = time.Now()
	case lexer.TokenTomorrow:
		base = time.Now().AddDate(0, 0, 1)
	case lexer.TokenYesterday:
		base = time.Now().AddDate(0, 0, -1)
	}

	// Normalise to start of day
	base = time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())

	expr := &DateExpr{Date: base}

	// Check for date arithmetic
	if p.current().Type == lexer.TokenPlus || p.current().Type == lexer.TokenMinus {
		op := p.current().Literal
		p.advance()

		offset, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		unit := ""
		if p.current().Type == lexer.TokenUnit || p.current().Type == lexer.TokenIdent {
			unit = p.current().Literal
			p.advance()
		}

		return &DateArithmeticExpr{
			Base:     expr,
			Operator: op,
			Offset:   offset,
			Unit:     unit,
		}, nil
	}

	return expr, nil
}

func (p *Parser) parseWeekday() (Expr, error) {
	modifier := ""

	// Check for "next" or "last"
	if p.current().Type == lexer.TokenNext || p.current().Type == lexer.TokenLast {
		modifier = strings.ToLower(p.current().Literal)
		p.advance()
	}

	// Get the weekday
	tok := p.current()
	var weekday time.Weekday

	switch tok.Type {
	case lexer.TokenMonday:
		weekday = time.Monday
	case lexer.TokenTuesday:
		weekday = time.Tuesday
	case lexer.TokenWednesday:
		weekday = time.Wednesday
	case lexer.TokenThursday:
		weekday = time.Thursday
	case lexer.TokenFriday:
		weekday = time.Friday
	case lexer.TokenSaturday:
		weekday = time.Saturday
	case lexer.TokenSunday:
		weekday = time.Sunday
	default:
		return nil, fmt.Errorf("expected weekday, got %s", tok.Type)
	}

	p.advance()

	return &WeekdayExpr{
		Weekday:  weekday,
		Modifier: modifier,
	}, nil
}

func (p *Parser) parseMonth() (Expr, error) {
	tok := p.current()
	var monthName string

	switch tok.Type {
	case lexer.TokenJanuary:
		monthName = "January"
	case lexer.TokenFebruary:
		monthName = "February"
	case lexer.TokenMarch:
		monthName = "March"
	case lexer.TokenApril:
		monthName = "April"
	case lexer.TokenMay:
		monthName = "May"
	case lexer.TokenJune:
		monthName = "June"
	case lexer.TokenJuly:
		monthName = "July"
	case lexer.TokenAugust:
		monthName = "August"
	case lexer.TokenSeptember:
		monthName = "September"
	case lexer.TokenOctober:
		monthName = "October"
	case lexer.TokenNovember:
		monthName = "November"
	case lexer.TokenDecember:
		monthName = "December"
	default:
		return nil, fmt.Errorf("expected month name, got %s", tok.Type)
	}

	p.advance()

	return &MonthExpr{
		Month: monthName,
	}, nil
}

// tryParseNumberWords attempts to parse a sequence of words as a number
// Returns the number value and true if successful
func (p *Parser) tryParseNumberWords() (float64, bool) {
	startPos := p.pos
	var words []string

	// Collect consecutive tokens that might be number words
	// This includes TokenIdent and certain keywords like TokenThree
	for {
		tok := p.current()
		var word string

		switch tok.Type {
		case lexer.TokenIdent:
			word = tok.Literal
		case lexer.TokenThree:
			// Special case: "three" is a keyword but also a number word
			// Check if next token is "quarters" - if so, don't treat as number
			if p.peek(1).Type == lexer.TokenQuarters {
				goto done
			}
			word = tok.Literal
		case lexer.TokenEOF:
			goto done
		default:
			goto done
		}

		// Check if this could be a number word (using en_GB as default)
		if !lexer.IsNumberWord(word, "en_GB") {
			break
		}
		words = append(words, word)
		p.advance()
	}

done:
	if len(words) == 0 {
		return 0, false
	}

	// Try to parse the collected words as a number
	val, ok := lexer.ParseNumberWords(words, "en_GB")
	if !ok {
		// Restore position if parsing failed
		p.pos = startPos
		return 0, false
	}

	return val, true
}
