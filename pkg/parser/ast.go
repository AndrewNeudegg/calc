package parser

import "time"

// Node represents an AST node.
type Node interface {
	node()
}

// Expr represents an expression node.
type Expr interface {
	Node
	expr()
}

// NumberExpr represents a numeric literal.
type NumberExpr struct {
	Value float64
}

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Left     Expr
	Operator string
	Right    Expr
}

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Operator string
	Operand  Expr
}

// IdentExpr represents a variable reference.
type IdentExpr struct {
	Name string
}

// AssignExpr represents a variable assignment.
type AssignExpr struct {
	Name  string
	Value Expr
}

// UnitExpr represents a value with a unit.
type UnitExpr struct {
	Value Expr
	Unit  string
}

// ConversionExpr represents a unit conversion.
type ConversionExpr struct {
	Value  Expr
	ToUnit string
}

// CurrencyExpr represents a currency value.
type CurrencyExpr struct {
	Value    Expr
	Currency string
}

// PercentExpr represents a percentage.
type PercentExpr struct {
	Value Expr
}

// PercentOfExpr represents "X% of Y".
type PercentOfExpr struct {
	Percent Expr
	Of      Expr
}

// PercentChangeExpr represents "increase/decrease X by Y%".
type PercentChangeExpr struct {
	Base     Expr
	Percent  Expr
	Increase bool // true for increase, false for decrease
}

// WhatPercentExpr represents "X is what % of Y".
type WhatPercentExpr struct {
	Part  Expr
	Whole Expr
}

// FunctionCallExpr represents a function call like sum(), average(), etc.
type FunctionCallExpr struct {
	Name string
	Args []Expr
}

// DateExpr represents a date value.
type DateExpr struct {
	Date time.Time
}

// TimeExpr represents a time value.
type TimeExpr struct {
	Time time.Time
}

// DateArithmeticExpr represents date arithmetic like "today + 3 days".
type DateArithmeticExpr struct {
	Base     Expr
	Operator string
	Offset   Expr
	Unit     string // days, weeks, months, years
}

// FuzzyExpr represents fuzzy phrases like "half of X", "double X".
type FuzzyExpr struct {
	Pattern string // "half", "double", "twice", etc.
	Value   Expr
}

// CommandExpr represents a command like ":save file.txt".
type CommandExpr struct {
	Command string
	Args    []string
}

// RateExpr represents a rate like "100 km / 2 hours".
type RateExpr struct {
	Numerator   Expr
	Denominator Expr
}

// WeekdayExpr represents "next monday", "last friday", etc.
type WeekdayExpr struct {
	Weekday  time.Weekday
	Modifier string // "next", "last", or empty for "this week"
}

// Implement node() for all types
func (*NumberExpr) node()         {}
func (*BinaryExpr) node()         {}
func (*UnaryExpr) node()          {}
func (*IdentExpr) node()          {}
func (*AssignExpr) node()         {}
func (*UnitExpr) node()           {}
func (*ConversionExpr) node()     {}
func (*CurrencyExpr) node()       {}
func (*PercentExpr) node()        {}
func (*PercentOfExpr) node()      {}
func (*PercentChangeExpr) node()  {}
func (*WhatPercentExpr) node()    {}
func (*FunctionCallExpr) node()   {}
func (*DateExpr) node()           {}
func (*TimeExpr) node()           {}
func (*DateArithmeticExpr) node() {}
func (*FuzzyExpr) node()          {}
func (*CommandExpr) node()        {}
func (*RateExpr) node()           {}
func (*WeekdayExpr) node()        {}

// Implement expr() for expression types
func (*NumberExpr) expr()         {}
func (*BinaryExpr) expr()         {}
func (*UnaryExpr) expr()          {}
func (*IdentExpr) expr()          {}
func (*AssignExpr) expr()         {}
func (*UnitExpr) expr()           {}
func (*ConversionExpr) expr()     {}
func (*CurrencyExpr) expr()       {}
func (*PercentExpr) expr()        {}
func (*PercentOfExpr) expr()      {}
func (*PercentChangeExpr) expr()  {}
func (*WhatPercentExpr) expr()    {}
func (*FunctionCallExpr) expr()   {}
func (*DateExpr) expr()           {}
func (*TimeExpr) expr()           {}
func (*DateArithmeticExpr) expr() {}
func (*FuzzyExpr) expr()          {}
func (*CommandExpr) expr()        {}
func (*RateExpr) expr()           {}
func (*WeekdayExpr) expr()        {}
