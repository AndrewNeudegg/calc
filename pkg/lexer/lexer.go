package lexer

import (
	"strings"
	"unicode"
)

// Lexer tokenises input text.
type Lexer struct {
	input    string
	pos      int
	line     int
	column   int
	keywords map[string]TokenType
}

// New creates a new lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		keywords: map[string]TokenType{
			"in":        TokenIn,
			"of":        TokenOf,
			"per":       TokenPer,
			"by":        TokenBy,
			"what":      TokenWhat,
			"is":        TokenIs,
			"increase":  TokenIncrease,
			"decrease":  TokenDecrease,
			"sum":       TokenSum,
			"average":   TokenAverage,
			"mean":      TokenMean,
			"total":     TokenTotal,
			"half":      TokenHalf,
			"double":    TokenDouble,
			"twice":     TokenTwice,
			"quarters":  TokenQuarters,
			"three":     TokenThree,
			"after":     TokenAfter,
			"before":    TokenBefore,
			"from":      TokenFrom,
			"ago":       TokenAgo,
			"now":       TokenNow,
			"today":     TokenToday,
			"tomorrow":  TokenTomorrow,
			"yesterday": TokenYesterday,
			"next":      TokenNext,
			"last":      TokenLast,
			"time":      TokenTime,
			"monday":    TokenMonday,
			"tuesday":   TokenTuesday,
			"wednesday": TokenWednesday,
			"thursday":  TokenThursday,
			"friday":    TokenFriday,
			"saturday":  TokenSaturday,
			"sunday":    TokenSunday,
			"january":   TokenJanuary,
			"february":  TokenFebruary,
			"march":     TokenMarch,
			"april":     TokenApril,
			"may":       TokenMay,
			"june":      TokenJune,
			"july":      TokenJuly,
			"august":    TokenAugust,
			"september": TokenSeptember,
			"october":   TokenOctober,
			"november":  TokenNovember,
			"december":  TokenDecember,
		},
	}
	return l
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipIgnored()

	if l.pos >= len(l.input) {
		return l.makeToken(TokenEOF, "")
	}

	ch := l.input[l.pos]

	// Single-character tokens
	switch ch {
	case '+':
		return l.advance(TokenPlus)
	case '-':
		return l.advance(TokenMinus)
	case '*':
		return l.advance(TokenMultiply)
	case '/':
		return l.advance(TokenDivide)
	case '%':
		return l.advance(TokenPercent)
	case '=':
		return l.advance(TokenEquals)
	case '(':
		return l.advance(TokenLParen)
	case ')':
		return l.advance(TokenRParen)
	case ',':
		return l.advance(TokenComma)
	case ':':
		return l.advance(TokenColon)
	case '$':
		return l.scanCurrency()
	}

	// Check for multi-byte UTF-8 currency symbols
	if l.pos+1 < len(l.input) {
		// Check for £ (C2 A3), € (E2 82 AC), ¥ (C2 A5)
		if ch == 0xC2 && l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			if next == 0xA3 || next == 0xA5 { // £ or ¥
				return l.scanCurrency()
			}
		}
		if ch == 0xE2 && l.pos+2 < len(l.input) {
			if l.input[l.pos+1] == 0x82 && l.input[l.pos+2] == 0xAC { // €
				return l.scanCurrency()
			}
		}
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) {
		return l.scanNumber()
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.scanIdentifier()
	}

	// Unknown character
	return l.advance(TokenError)
}

func (l *Lexer) advance(typ TokenType) Token {
	tok := l.makeToken(typ, string(l.input[l.pos]))
	l.pos++
	l.column++
	return tok
}

func (l *Lexer) makeToken(typ TokenType, literal string) Token {
	return Token{
		Type:    typ,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

// skipIgnored skips whitespace and '//' line comments repeatedly.
func (l *Lexer) skipIgnored() {
	for {
		// Skip any whitespace first
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			return
		}
		// Skip '//' comments
		if l.input[l.pos] == '/' && l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
			// Advance until end of line or input
			l.pos += 2
			l.column += 2
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
				l.column++
			}
			// If at newline, consume it and move to next line
			if l.pos < len(l.input) && l.input[l.pos] == '\n' {
				l.pos++
				l.line++
				l.column = 1
			}
			// Loop to skip following whitespace/comments
			continue
		}
		return
	}
}

func (l *Lexer) scanNumber() Token {
	start := l.pos
	startCol := l.column

	// Scan integer part
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.pos++
		l.column++
	}

	// Check for date format (DD/MM/YYYY or MM/DD/YYYY)
	if l.pos < len(l.input) && l.input[l.pos] == '/' {
		// Look ahead to see if this could be a date
		savedPos := l.pos
		savedCol := l.column

		l.pos++ // consume first '/'
		l.column++

		// Scan second number (month or day)
		secondStart := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.pos++
			l.column++
		}

		// Check for second '/'
		if l.pos < len(l.input) && l.input[l.pos] == '/' && l.pos > secondStart {
			l.pos++ // consume second '/'
			l.column++

			// Scan year
			yearStart := l.pos
			for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
				l.pos++
				l.column++
			}

			// If we have year digits, this is a date
			if l.pos > yearStart {
				literal := l.input[start:l.pos]
				return Token{
					Type:    TokenDate,
					Literal: literal,
					Line:    l.line,
					Column:  startCol,
				}
			}
		}

		// Not a date, reset
		l.pos = savedPos
		l.column = savedCol
	}

	// Check for time format (HH:MM or H:MM)
	if l.pos < len(l.input) && l.input[l.pos] == ':' {
		// Look ahead to see if this could be a time (colon followed by digits)
		if l.pos+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.pos+1])) {
			// This looks like a time - scan it
			l.pos++ // consume ':'
			l.column++

			// Scan minutes
			minuteStart := l.pos
			for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
				l.pos++
				l.column++
			}

			// Check for seconds (optional)
			if l.pos < len(l.input) && l.input[l.pos] == ':' {
				l.pos++ // consume second ':'
				l.column++

				// Scan seconds
				for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
					l.pos++
					l.column++
				}
			}

			// Make sure we have at least 2 digits for minutes
			if l.pos-minuteStart >= 2 || l.pos-minuteStart == 1 {
				literal := l.input[start:l.pos]
				return Token{
					Type:    TokenTimeValue,
					Literal: literal,
					Line:    l.line,
					Column:  startCol,
				}
			}
		}
	}

	// Check for decimal point
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		l.column++

		// Scan fractional part
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.pos++
			l.column++
		}
	}

	// Check for comma separators (UK/European format)
	literal := l.input[start:l.pos]

	return Token{
		Type:    TokenNumber,
		Literal: literal,
		Line:    l.line,
		Column:  startCol,
	}
}

func (l *Lexer) scanIdentifier() Token {
	start := l.pos
	startCol := l.column

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if !unicode.IsLetter(rune(ch)) && !unicode.IsDigit(rune(ch)) && ch != '_' {
			break
		}
		l.pos++
		l.column++
	}

	literal := l.input[start:l.pos]
	lowerLiteral := strings.ToLower(literal)

	// Check if it's a keyword
	if typ, ok := l.keywords[lowerLiteral]; ok {
		return Token{
			Type:    typ,
			Literal: literal,
			Line:    l.line,
			Column:  startCol,
		}
	}

	// Check if it's a known unit
	if l.isKnownUnit(literal) {
		return Token{
			Type:    TokenUnit,
			Literal: literal,
			Line:    l.line,
			Column:  startCol,
		}
	}

	return Token{
		Type:    TokenIdent,
		Literal: literal,
		Line:    l.line,
		Column:  startCol,
	}
}

func (l *Lexer) scanCurrency() Token {
	start := l.pos
	startCol := l.column

	// Handle ASCII currency ($)
	if l.input[l.pos] == '$' {
		l.pos++
		l.column++
	} else {
		// Handle UTF-8 multi-byte currency symbols
		// £ = C2 A3 (2 bytes)
		// € = E2 82 AC (3 bytes)
		// ¥ = C2 A5 (2 bytes)
		firstByte := l.input[l.pos]
		if firstByte == 0xC2 {
			l.pos += 2
			l.column++
		} else if firstByte == 0xE2 {
			l.pos += 3
			l.column++
		}
	}

	return Token{
		Type:    TokenCurrency,
		Literal: l.input[start:l.pos],
		Line:    l.line,
		Column:  startCol,
	}
}

func (l *Lexer) isCurrencySymbol(ch byte) bool {
	// Not used anymore, kept for compatibility
	return ch == '$'
}

func (l *Lexer) isKnownUnit(s string) bool {
	knownUnits := map[string]bool{
		// Length
		"m": true, "cm": true, "mm": true, "km": true,
		"ft": true, "in": true, "yd": true, "mi": true,
		"mile": true, "miles": true, "metre": true, "metres": true,
		"meter": true, "meters": true, "foot": true, "feet": true,
		"inch": true, "inches": true, "yard": true, "yards": true,

		// Mass
		"g": true, "kg": true, "mg": true, "µg": true, "ug": true,
		"lb": true, "lbs": true, "oz": true,
		"gram": true, "grams": true, "milligram": true, "milligrams": true,
		"microgram": true, "micrograms": true,
		"kilogram": true, "kilograms": true,
		"pound": true, "pounds": true, "ounce": true, "ounces": true,
		"stone": true, "stones": true, "st": true,
		"carat": true, "carats": true, "ct": true,
		"troyounce": true, "troyounces": true, "troyoz": true, "ozt": true,
		"tonne": true, "tonnes": true, "ton": true, "tons": true,

		// Time
		"ns": true, "nanosecond": true, "nanoseconds": true,
		"µs": true, "us": true, "microsecond": true, "microseconds": true,
		"ms": true, "millisecond": true, "milliseconds": true,
		"s": true, "sec": true, "second": true, "seconds": true,
		"min": true, "minute": true, "minutes": true,
		"h": true, "hr": true, "hour": true, "hours": true,
		"day": true, "days": true, "week": true, "weeks": true,
		"fortnight": true, "fortnights": true,
		"month": true, "months": true,
		"quarter": true, "quarters": true,
		"semester": true, "semesters": true,
		"year": true, "years": true,

		// Volume
		"l": true, "ml": true, "litre": true, "litres": true, "liter": true, "liters": true,
		"millilitre": true, "millilitres": true, "milliliter": true, "milliliters": true,
		"cl": true, "centilitre": true, "centilitres": true, "centiliter": true, "centiliters": true,
		"dl": true, "decilitre": true, "decilitres": true, "deciliter": true, "deciliters": true,
		"m3": true, "m³": true, "cm3": true, "cm³": true, "mm3": true, "mm³": true,
		"ft3": true, "ft³": true, "in3": true, "in³": true, "cc": true,
		"gal": true, "gallon": true, "gallons": true,
		"usgal": true, "usgallon": true, "usgallons": true,
		"ukgal": true, "ukgallon": true, "ukgallons": true,
		"impgal": true, "imperialgallon": true,
		"quart": true, "quarts": true, "qt": true,
		"usquart": true, "usquarts": true, "ukquart": true, "ukquarts": true,
		"pint": true, "pints": true, "pt": true,
		"uspint": true, "uspints": true, "ukpint": true, "ukpints": true,
		"imppint": true, "imperialpint": true,
		"cup": true, "cups": true,
		"floz": true, "fluidounce": true, "fluidounces": true,
		"tbsp": true, "tablespoon": true, "tablespoons": true,
		"tsp": true, "teaspoon": true, "teaspoons": true,

		// Area
		"sqm": true, "m2": true, "m²": true,
		"sqmm": true, "mm2": true, "mm²": true,
		"sqcm": true, "cm2": true, "cm²": true,
		"sqkm": true, "km2": true, "km²": true,
		"sqft": true, "ft2": true, "ft²": true,
		"sqin": true, "in2": true, "in²": true,
		"sqyd": true, "yd2": true, "yd²": true,
		"sqmi": true, "mi2": true, "mi²": true,
		"squaremetre": true, "squaremetres": true, "squaremeter": true, "squaremeters": true,
		"squarefoot": true, "squarefeet": true,
		"squareinch": true, "squareinches": true,
		"squareyard": true, "squareyards": true,
		"squaremile": true, "squaremiles": true,
		"squarekilometre": true, "squarekilometres": true, "squarekilometer": true, "squarekilometers": true,
		"acre": true, "acres": true,
		"hectare": true, "hectares": true, "ha": true,
		"are": true, "ares": true, "decare": true, "decares": true,

		// Temperature
		"c": true, "f": true, "celsius": true, "fahrenheit": true,
		"k": true, "kelvin": true,
		"r": true, "rankine": true, "°r": true,

		// Speed
		"mps": true, "kph": true, "kmh": true, "mph": true,
		"fps": true, "knot": true, "knots": true, "kn": true,

		// Pressure
		"pa": true, "pascal": true, "pascals": true,
		"kpa": true, "kilopascal": true, "kilopascals": true,
		"mpa": true, "megapascal": true, "megapascals": true,
		"bar": true, "bars": true, "mbar": true, "millibar": true, "millibars": true,
		"atm": true, "atmosphere": true, "atmospheres": true,
		"psi": true, "torr": true, "mmhg": true, "inhg": true,

		// Force
		"n": true, "newton": true, "newtons": true,
		"kilonewton": true, "kilonewtons": true,
		"mn": true, "meganewton": true, "meganewtons": true,
		"lbf": true, "poundforce": true, "poundsforce": true,
		"kgf": true, "kilogramforce": true,
		"dyne": true, "dynes": true,

		// Angle
		"deg": true, "degree": true, "degrees": true, "°": true,
		"rad": true, "radian": true, "radians": true,
		"grad": true, "gradian": true, "gradians": true, "gon": true,
		"turn": true, "turns": true, "revolution": true, "revolutions": true,

		// Frequency
		"hz": true, "hertz": true,
		"khz": true, "kilohertz": true,
		"mhz": true, "megahertz": true,
		"ghz": true, "gigahertz": true,
		"thz": true, "terahertz": true,
		"rpm": true,

		// Digital storage (bytes)
		"b": true, "byte": true, "bytes": true,
		"kb": true, "kilobyte": true, "kilobytes": true,
		"mb": true, "megabyte": true, "megabytes": true,
		"gb": true, "gigabyte": true, "gigabytes": true,
		"tb": true, "terabyte": true, "terabytes": true,
		"pb": true, "petabyte": true, "petabytes": true,

		// Digital storage (bits)
		"bit": true, "bits": true,
		"kbit": true, "kilobit": true, "kilobits": true,
		"mbit": true, "megabit": true, "megabits": true,
		"gbit": true, "gigabit": true, "gigabits": true,
		"tbit": true, "terabit": true, "terabits": true,
		"pbit": true, "petabit": true, "petabits": true,

		// Data rate (bytes per second)
		"bps": true, "kbps": true, "mbps": true, "gbps": true, "tbps": true,
		"Bps": true, "KBps": true, "MBps": true, "GBps": true, "TBps": true,

		// Data rate (bits per second)
		"bitps": true, "kbitps": true, "mbitps": true, "gbitps": true, "tbitps": true,

		// Currency codes and names
		"usd": true, "dollar": true, "dollars": true,
		"gbp": true, // "pound" and "pounds" already exist as mass units
		"eur": true, "euro": true, "euros": true,
		"jpy": true, "yen": true,
		// Expanded currency codes
		"aud": true, "cad": true, "nzd": true,
		"chf": true, "sek": true, "nok": true, "dkk": true,
		"pln": true, "czk": true, "huf": true, "ron": true,
		"rub": true, "try": true,
		"aed": true, "sar": true, "ils": true,
		"cny": true, "hkd": true, "sgd": true, "inr": true, "krw": true, "twd": true, "thb": true, "myr": true, "idr": true, "php": true,
		"mxn": true, "brl": true, "zar": true,
	}

	return knownUnits[strings.ToLower(s)]
}

// AllTokens returns all tokens from the input as a slice.
func (l *Lexer) AllTokens() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}
