package lexer

// TokenType represents the type of token.
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError
	TokenWhitespace
	
	// Literals
	TokenNumber
	TokenIdent
	TokenString
	
	// Operators
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenPercent
	TokenEquals
	
	// Delimiters
	TokenLParen
	TokenRParen
	TokenComma
	TokenColon
	
	// Keywords
	TokenIn
	TokenOf
	TokenPer
	TokenBy
	TokenWhat
	TokenIs
	TokenIncrease
	TokenDecrease
	TokenSum
	TokenAverage
	TokenMean
	TokenTotal
	TokenHalf
	TokenDouble
	TokenTwice
	TokenQuarters
	TokenThree
	TokenAfter
	TokenBefore
	TokenFrom
	TokenAgo
	TokenNow
	TokenToday
	TokenTomorrow
	TokenYesterday
	TokenNext
	TokenLast
	TokenMonday
	TokenTuesday
	TokenWednesday
	TokenThursday
	TokenFriday
	TokenSaturday
	TokenSunday
	
	// Units (commonly recognised)
	TokenUnit
	
	// Currency symbols
	TokenCurrency
	
	// Date/Time
	TokenDate
	TokenTime
)

// Token represents a single lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenWhitespace:
		return "WHITESPACE"
	case TokenNumber:
		return "NUMBER"
	case TokenIdent:
		return "IDENT"
	case TokenString:
		return "STRING"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenPercent:
		return "%"
	case TokenEquals:
		return "="
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenComma:
		return ","
	case TokenColon:
		return ":"
	case TokenIn:
		return "in"
	case TokenOf:
		return "of"
	case TokenPer:
		return "per"
	case TokenBy:
		return "by"
	case TokenWhat:
		return "what"
	case TokenIs:
		return "is"
	case TokenIncrease:
		return "increase"
	case TokenDecrease:
		return "decrease"
	case TokenSum:
		return "sum"
	case TokenAverage:
		return "average"
	case TokenMean:
		return "mean"
	case TokenTotal:
		return "total"
	case TokenHalf:
		return "half"
	case TokenDouble:
		return "double"
	case TokenTwice:
		return "twice"
	case TokenQuarters:
		return "quarters"
	case TokenThree:
		return "three"
	case TokenAfter:
		return "after"
	case TokenBefore:
		return "before"
	case TokenFrom:
		return "from"
	case TokenAgo:
		return "ago"
	case TokenNow:
		return "now"
	case TokenToday:
		return "today"
	case TokenTomorrow:
		return "tomorrow"
	case TokenYesterday:
		return "yesterday"
	case TokenNext:
		return "next"
	case TokenLast:
		return "last"
	case TokenMonday:
		return "monday"
	case TokenTuesday:
		return "tuesday"
	case TokenWednesday:
		return "wednesday"
	case TokenThursday:
		return "thursday"
	case TokenFriday:
		return "friday"
	case TokenSaturday:
		return "saturday"
	case TokenSunday:
		return "sunday"
	case TokenUnit:
		return "UNIT"
	case TokenCurrency:
		return "CURRENCY"
	case TokenDate:
		return "DATE"
	case TokenTime:
		return "TIME"
	default:
		return "UNKNOWN"
	}
}
