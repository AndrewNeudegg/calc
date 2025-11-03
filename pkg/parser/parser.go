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
}

// New creates a new parser from tokens.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
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

func (p *Parser) parseExpression() (Expr, error) {
	// Check for command
	if p.current().Type == lexer.TokenColon {
		return p.parseCommand()
	}

	// Check for assignment
	if p.current().Type == lexer.TokenIdent && p.peek(1).Type == lexer.TokenEquals {
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

	if p.current().Type != lexer.TokenIdent {
		return nil, fmt.Errorf("expected command name")
	}

	command := p.current().Literal
	p.advance()

	var args []string
	for p.current().Type != lexer.TokenEOF {
		args = append(args, p.current().Literal)
		p.advance()
	}

	return &CommandExpr{
		Command: command,
		Args:    args,
	}, nil
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
			return &PercentChangeExpr{Base: base, Percent: percent, Increase: true}, true
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
			return &PercentChangeExpr{Base: base, Percent: percent, Increase: false}, true
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
			return &WhatPercentExpr{Part: part, Whole: whole}, true
		}
	}

	return nil, false
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
				strings.ToLower(p.peek(i).Literal) != strings.ToLower(part) {
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
	expr, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	// Check for "in" conversion
	if p.current().Type == lexer.TokenIn {
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

		return &ConversionExpr{Value: expr, ToUnit: toUnit}, nil
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
		val, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", tok.Literal)
		}
		p.advance()
		return &NumberExpr{Value: val}, nil

	case lexer.TokenCurrency:
		currency := tok.Literal
		p.advance()

		if p.current().Type != lexer.TokenNumber {
			return nil, fmt.Errorf("expected number after currency symbol")
		}

		val, err := strconv.ParseFloat(p.current().Literal, 64)
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
		expr, err := p.parseAdditive()
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

	default:
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
		arg, err := p.parseAdditive()
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
