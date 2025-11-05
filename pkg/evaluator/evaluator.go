package evaluator

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/andrewneudegg/calc/pkg/currency"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/timezone"
	"github.com/andrewneudegg/calc/pkg/units"
)

// Environment stores variables and state.
type Environment struct {
	variables map[string]Value
	units     *units.System
	currency  *currency.System
	timezone  *timezone.System
}

// NewEnvironment creates a new evaluation environment.
func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		units:     units.NewSystem(),
		currency:  currency.NewSystem(),
		timezone:  timezone.NewSystem(),
	}
}

// Evaluator evaluates expressions.
type Evaluator struct {
	env *Environment
}

// New creates a new evaluator.
func New(env *Environment) *Evaluator {
	return &Evaluator{env: env}
}

// Eval evaluates an expression and returns a value.
func (e *Evaluator) Eval(expr parser.Expr) Value {
	if expr == nil {
		return NewError("nil expression")
	}

	switch node := expr.(type) {
	case *parser.NumberExpr:
		return NewNumber(node.Value)

	case *parser.BinaryExpr:
		return e.evalBinary(node)

	case *parser.UnaryExpr:
		return e.evalUnary(node)

	case *parser.IdentExpr:
		return e.evalIdent(node)

	case *parser.AssignExpr:
		return e.evalAssign(node)

	case *parser.UnitExpr:
		return e.evalUnit(node)

	case *parser.ConversionExpr:
		return e.evalConversion(node)

	case *parser.CurrencyExpr:
		return e.evalCurrency(node)

	case *parser.PercentExpr:
		return e.evalPercent(node)

	case *parser.PercentOfExpr:
		return e.evalPercentOf(node)

	case *parser.PercentChangeExpr:
		return e.evalPercentChange(node)

	case *parser.WhatPercentExpr:
		return e.evalWhatPercent(node)

	case *parser.FunctionCallExpr:
		return e.evalFunctionCall(node)

	case *parser.StringExpr:
		return NewString(node.Value)

	case *parser.DateExpr:
		return NewDate(node.Date)

	case *parser.TimeExpr:
		return NewDate(node.Time)

	case *parser.DateArithmeticExpr:
		return e.evalDateArithmetic(node)

	case *parser.FuzzyExpr:
		return e.evalFuzzy(node)

	case *parser.WeekdayExpr:
		return e.evalWeekday(node)

	case *parser.MonthExpr:
		return e.evalMonth(node)

	case *parser.TimeInLocationExpr:
		return e.evalTimeInLocation(node)

	case *parser.TimeDifferenceExpr:
		return e.evalTimeDifference(node)

	case *parser.TimeConversionExpr:
		return e.evalTimeConversion(node)

	case *parser.RateExpr:
		return e.evalRate(node)

	default:
		return NewError(fmt.Sprintf("unknown expression type: %T", expr))
	}
}

func (e *Evaluator) evalBinary(node *parser.BinaryExpr) Value {
	left := e.Eval(node.Left)
	if left.IsError() {
		return left
	}

	right := e.Eval(node.Right)
	if right.IsError() {
		return right
	}

	// Handle date + unit or date - unit (date arithmetic)
	if left.Type == ValueDate && right.Type == ValueUnit && (node.Operator == "+" || node.Operator == "-") {
		// Extract offset value and unit
		offset := right.Number
		if node.Operator == "-" {
			offset = -offset
		}

		// Calculate new date based on unit
		unit := right.Unit
		var newDate time.Time

		switch strings.ToLower(unit) {
		case "day", "days", "d":
			newDate = left.Date.AddDate(0, 0, int(offset))
		case "week", "weeks", "w":
			newDate = left.Date.AddDate(0, 0, int(offset*7))
		case "month", "months", "mo":
			newDate = left.Date.AddDate(0, int(offset), 0)
		case "year", "years", "y":
			newDate = left.Date.AddDate(int(offset), 0, 0)
		case "hour", "hours", "h", "hr":
			newDate = left.Date.Add(time.Duration(offset * float64(time.Hour)))
		case "minute", "minutes", "min":
			newDate = left.Date.Add(time.Duration(offset * float64(time.Minute)))
		case "second", "seconds", "s", "sec":
			newDate = left.Date.Add(time.Duration(offset * float64(time.Second)))
		default:
			return NewError(fmt.Sprintf("cannot add unit '%s' to date", unit))
		}

		return NewDate(newDate)
	}

	// Handle date-date subtraction (returns days with unit)
	if left.Type == ValueDate && right.Type == ValueDate && node.Operator == "-" {
		duration := left.Date.Sub(right.Date)
		days := duration.Hours() / 24.0
		return NewUnit(days, "days")
	}

	// Handle currency operations
	if left.Type == ValueCurrency || right.Type == ValueCurrency {
		return e.evalCurrencyBinary(left, node.Operator, right)
	}

	// Handle unit operations
	if left.Type == ValueUnit || right.Type == ValueUnit {
		return e.evalUnitBinary(left, node.Operator, right)
	}

	// Handle percentage operations
	if right.Type == ValuePercent && node.Operator == "+" {
		// e.g., "30 + 20%" = 30 + (30 * 0.20)
		return NewNumber(left.Number + (left.Number * right.Number / 100))
	}

	if right.Type == ValuePercent && node.Operator == "-" {
		// e.g., "30 - 20%" = 30 - (30 * 0.20)
		return NewNumber(left.Number - (left.Number * right.Number / 100))
	}

	// Standard numeric operations
	switch node.Operator {
	case "+":
		return NewNumber(left.Number + right.Number)
	case "-":
		return NewNumber(left.Number - right.Number)
	case "*":
		return NewNumber(left.Number * right.Number)
	case "/":
		if right.Number == 0 {
			return NewError("division by zero")
		}
		return NewNumber(left.Number / right.Number)
	default:
		return NewError(fmt.Sprintf("unknown operator: %s", node.Operator))
	}
}

func (e *Evaluator) evalUnary(node *parser.UnaryExpr) Value {
	operand := e.Eval(node.Operand)
	if operand.IsError() {
		return operand
	}

	switch node.Operator {
	case "-":
		operand.Number = -operand.Number
		return operand
	default:
		return NewError(fmt.Sprintf("unknown unary operator: %s", node.Operator))
	}
}

func (e *Evaluator) evalIdent(node *parser.IdentExpr) Value {
	val, ok := e.env.variables[node.Name]
	if !ok {
		return NewError(fmt.Sprintf("undefined variable: %s", node.Name))
	}
	return val
}

func (e *Evaluator) evalAssign(node *parser.AssignExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	e.env.variables[node.Name] = val
	return val
}

func (e *Evaluator) evalUnit(node *parser.UnitExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	return NewUnit(val.Number, node.Unit)
}

func (e *Evaluator) evalConversion(node *parser.ConversionExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	// Handle currency conversion
	if val.Type == ValueCurrency {
		result, err := e.env.currency.Convert(val.Number, val.Currency, node.ToUnit)
		if err != nil {
			return NewError(err.Error())
		}
		return NewCurrency(result, e.env.currency.GetSymbol(node.ToUnit))
	}

	// Handle unit conversion
	if val.Type == ValueUnit {
		// Special case: currency/time rates (e.g., $/day) to other currency/time (e.g., gbp/month)
		if units.IsCompoundUnit(val.Unit) || units.IsCompoundUnit(node.ToUnit) {
			fromParts := strings.Split(val.Unit, "/")
			toParts := strings.Split(node.ToUnit, "/")

			// Handle currency rate: currency in numerator and time in denominator
			if len(fromParts) == 2 && e.env.currency.IsCurrency(strings.TrimSpace(fromParts[0])) {
				fromCur := strings.TrimSpace(fromParts[0])
				fromTime := strings.TrimSpace(fromParts[1])

				// If target is compound currency/time
				if len(toParts) == 2 && e.env.currency.IsCurrency(strings.TrimSpace(toParts[0])) {
					toCur := strings.TrimSpace(toParts[0])
					toTime := strings.TrimSpace(toParts[1])

					// Scale rate to the target time period
					// factor = (1 toTime) expressed in fromTime units
					timeFactor, err := e.env.units.Convert(1, toTime, fromTime)
					if err != nil {
						return NewError(err.Error())
					}

					perTarget := val.Number * timeFactor

					// Convert currency
					converted, err := e.env.currency.Convert(perTarget, fromCur, toCur)
					if err != nil {
						return NewError(err.Error())
					}

					// Return total amount per target period as a currency value (e.g., monthly amount)
					return NewCurrency(converted, e.env.currency.GetSymbol(toCur))
				}

				// If target is a different time unit but same currency rate
				if len(toParts) == 2 && !e.env.currency.IsCurrency(strings.TrimSpace(toParts[0])) {
					// Non-currency compound target: delegate to unit conversion if possible
					result, err := e.env.units.ConvertCompoundUnit(val.Number, val.Unit, node.ToUnit)
					if err != nil {
						return NewError(err.Error())
					}
					return NewUnit(result, node.ToUnit)
				}
			}

			// Generic compound unit conversions (non-currency)
			result, err := e.env.units.ConvertCompoundUnit(val.Number, val.Unit, node.ToUnit)
			if err != nil {
				return NewError(err.Error())
			}
			return NewUnit(result, node.ToUnit)
		}

		// Regular simple unit conversion
		result, err := e.env.units.Convert(val.Number, val.Unit, node.ToUnit)
		if err != nil {
			return NewError(err.Error())
		}
		return NewUnit(result, node.ToUnit)
	}

	// Try converting a plain number with a unit
	result, err := e.env.units.Convert(val.Number, "unknown", node.ToUnit)
	if err != nil {
		return NewError(err.Error())
	}
	return NewUnit(result, node.ToUnit)
}

func (e *Evaluator) evalCurrency(node *parser.CurrencyExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	// Normalize the currency code to a symbol for display
	symbol := e.env.currency.GetSymbol(node.Currency)
	return NewCurrency(val.Number, symbol)
}

func (e *Evaluator) evalPercent(node *parser.PercentExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	return NewPercent(val.Number)
}

func (e *Evaluator) evalPercentOf(node *parser.PercentOfExpr) Value {
	percent := e.Eval(node.Percent)
	if percent.IsError() {
		return percent
	}

	of := e.Eval(node.Of)
	if of.IsError() {
		return of
	}

	result := of.Number * (percent.Number / 100)

	// Preserve the type of the "of" value
	switch of.Type {
	case ValueCurrency:
		return NewCurrency(result, of.Currency)
	case ValueUnit:
		return NewUnit(result, of.Unit)
	default:
		return NewNumber(result)
	}
}

func (e *Evaluator) evalPercentChange(node *parser.PercentChangeExpr) Value {
	base := e.Eval(node.Base)
	if base.IsError() {
		return base
	}

	percent := e.Eval(node.Percent)
	if percent.IsError() {
		return percent
	}

	var result float64
	if node.Increase {
		result = base.Number * (1 + percent.Number/100)
	} else {
		result = base.Number * (1 - percent.Number/100)
	}

	// Preserve the type
	switch base.Type {
	case ValueCurrency:
		return NewCurrency(result, base.Currency)
	case ValueUnit:
		return NewUnit(result, base.Unit)
	default:
		return NewNumber(result)
	}
}

func (e *Evaluator) evalWhatPercent(node *parser.WhatPercentExpr) Value {
	part := e.Eval(node.Part)
	if part.IsError() {
		return part
	}

	whole := e.Eval(node.Whole)
	if whole.IsError() {
		return whole
	}

	if whole.Number == 0 {
		return NewError("division by zero")
	}

	result := (part.Number / whole.Number) * 100
	return NewPercent(result)
}

func (e *Evaluator) evalFunctionCall(node *parser.FunctionCallExpr) Value {
	switch strings.ToLower(node.Name) {
	case "sum", "total":
		return e.evalSum(node.Args)
	case "average", "mean":
		return e.evalAverage(node.Args)
	case "min":
		return e.evalMin(node.Args)
	case "max":
		return e.evalMax(node.Args)
	case "print":
		return e.evalPrint(node.Args)
	default:
		return NewError(fmt.Sprintf("unknown function: %s", node.Name))
	}
}

// evalPrint returns a string after interpolating {var} placeholders using current variables.
// It does not produce side effects; the REPL will print the returned string value.
func (e *Evaluator) evalPrint(args []parser.Expr) Value {
	if len(args) != 1 {
		return NewError("print requires exactly one argument")
	}
	val := e.Eval(args[0])
	if val.IsError() {
		return val
	}
	if val.Type != ValueString {
		return NewError("print expects a string literal")
	}
	s := val.Text
	// Find {identifier} placeholders and replace
	// Simple single-pass replacement; does not support nested braces
	var out strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '{' {
			// find closing brace
			j := i + 1
			for j < len(s) && s[j] != '}' {
				j++
			}
			if j >= len(s) {
				// unmatched '{' - leave as-is
				out.WriteString(s[i:])
				break
			}
			name := strings.TrimSpace(s[i+1 : j])
			if name == "" {
				out.WriteString(s[i : j+1])
				i = j + 1
				continue
			}
			// Look up variable
			v, ok := e.env.variables[name]
			if !ok {
				return NewError(fmt.Sprintf("undefined variable: %s", name))
			}
			out.WriteString(v.String())
			i = j + 1
		} else {
			out.WriteByte(s[i])
			i++
		}
	}
	return NewString(out.String())
}

func (e *Evaluator) evalSum(args []parser.Expr) Value {
	var sum float64
	for _, arg := range args {
		val := e.Eval(arg)
		if val.IsError() {
			return val
		}
		sum += val.Number
	}
	return NewNumber(sum)
}

func (e *Evaluator) evalAverage(args []parser.Expr) Value {
	if len(args) == 0 {
		return NewError("average requires at least one argument")
	}

	sumVal := e.evalSum(args)
	if sumVal.IsError() {
		return sumVal
	}

	return NewNumber(sumVal.Number / float64(len(args)))
}

func (e *Evaluator) evalMin(args []parser.Expr) Value {
	if len(args) == 0 {
		return NewError("min requires at least one argument")
	}

	// Evaluate first argument to initialize
	first := e.Eval(args[0])
	if first.IsError() {
		return first
	}
	minVal := first.Number

	for i := 1; i < len(args); i++ {
		v := e.Eval(args[i])
		if v.IsError() {
			return v
		}
		if v.Number < minVal {
			minVal = v.Number
		}
	}

	return NewNumber(minVal)
}

func (e *Evaluator) evalMax(args []parser.Expr) Value {
	if len(args) == 0 {
		return NewError("max requires at least one argument")
	}

	// Evaluate first argument to initialize
	first := e.Eval(args[0])
	if first.IsError() {
		return first
	}
	maxVal := first.Number

	for i := 1; i < len(args); i++ {
		v := e.Eval(args[i])
		if v.IsError() {
			return v
		}
		if v.Number > maxVal {
			maxVal = v.Number
		}
	}

	return NewNumber(maxVal)
}

func (e *Evaluator) evalDateArithmetic(node *parser.DateArithmeticExpr) Value {
	base := e.Eval(node.Base)
	if base.IsError() {
		return base
	}

	offset := e.Eval(node.Offset)
	if offset.IsError() {
		return offset
	}

	offsetVal := int(offset.Number)
	if node.Operator == "-" {
		offsetVal = -offsetVal
	}

	var result time.Time
	unit := strings.ToLower(node.Unit)

	switch unit {
	case "day", "days":
		result = base.Date.AddDate(0, 0, offsetVal)
	case "week", "weeks":
		result = base.Date.AddDate(0, 0, offsetVal*7)
	case "month", "months":
		result = base.Date.AddDate(0, offsetVal, 0)
	case "year", "years":
		result = base.Date.AddDate(offsetVal, 0, 0)
	case "hour", "hours", "h", "hr":
		result = base.Date.Add(time.Duration(offsetVal) * time.Hour)
	case "minute", "minutes", "min":
		result = base.Date.Add(time.Duration(offsetVal) * time.Minute)
	case "second", "seconds", "s", "sec":
		result = base.Date.Add(time.Duration(offsetVal) * time.Second)
	default:
		return NewError(fmt.Sprintf("unknown time unit: %s", node.Unit))
	}

	return NewDate(result)
}

func (e *Evaluator) evalFuzzy(node *parser.FuzzyExpr) Value {
	val := e.Eval(node.Value)
	if val.IsError() {
		return val
	}

	pattern := strings.ToLower(node.Pattern)
	var result float64

	switch pattern {
	case "half":
		result = val.Number * 0.5
	case "double", "twice":
		result = val.Number * 2
	case "three quarters":
		result = val.Number * 0.75
	default:
		return NewError(fmt.Sprintf("unknown fuzzy pattern: %s", node.Pattern))
	}

	// Preserve type
	switch val.Type {
	case ValueCurrency:
		return NewCurrency(result, val.Currency)
	case ValueUnit:
		return NewUnit(result, val.Unit)
	default:
		return NewNumber(result)
	}
}

func (e *Evaluator) evalCurrencyBinary(left Value, op string, right Value) Value {
	// Convert both to the same currency if needed
	if left.Type == ValueCurrency && right.Type == ValueCurrency {
		if left.Currency != right.Currency {
			// Convert right to left's currency
			converted, err := e.env.currency.Convert(right.Number, right.Currency, left.Currency)
			if err != nil {
				return NewError(err.Error())
			}
			right.Number = converted
			right.Currency = left.Currency
		}
	}

	switch op {
	case "+":
		return NewCurrency(left.Number+right.Number, left.Currency)
	case "-":
		return NewCurrency(left.Number-right.Number, left.Currency)
	case "*":
		// Allow: currency * number OR number * currency
		// Reject: currency * currency
		if left.Type == ValueCurrency && right.Type == ValueCurrency {
			return NewError("cannot multiply two currencies")
		}
		// Determine which side has the currency
		if left.Type == ValueCurrency {
			return NewCurrency(left.Number*right.Number, left.Currency)
		} else {
			// right must be currency (since we're in evalCurrencyBinary)
			return NewCurrency(left.Number*right.Number, right.Currency)
		}
	case "/":
		if right.Number == 0 {
			return NewError("division by zero")
		}
		if right.Type == ValueCurrency {
			return NewNumber(left.Number / right.Number)
		}
		return NewCurrency(left.Number/right.Number, left.Currency)
	default:
		return NewError(fmt.Sprintf("unknown operator: %s", op))
	}
}

func (e *Evaluator) evalUnitBinary(left Value, op string, right Value) Value {
	switch op {
	case "+", "-":
		// For addition/subtraction, units must be compatible
		if left.Type == ValueUnit && right.Type == ValueUnit {
			if left.Unit != right.Unit {
				// Try to convert right to left's unit
				converted, err := e.env.units.Convert(right.Number, right.Unit, left.Unit)
				if err != nil {
					return NewError(err.Error())
				}
				right.Number = converted
				right.Unit = left.Unit
			}
		}

		if op == "+" {
			return NewUnit(left.Number+right.Number, left.Unit)
		}
		return NewUnit(left.Number-right.Number, left.Unit)

	case "*":
		if right.Type == ValueUnit {
			// If left is a plain number (not a unit), this is scalar multiplication
			// Result should be in the right's unit
			if left.Type != ValueUnit {
				return NewUnit(left.Number*right.Number, right.Unit)
			}
			// Both are units - creating compound unit
			return NewUnit(left.Number*right.Number, left.Unit+"·"+right.Unit)
		}
		return NewUnit(left.Number*right.Number, left.Unit)

	case "/":
		if right.Number == 0 {
			return NewError("division by zero")
		}
		if right.Type == ValueUnit {
			// For division, try to convert if possible
			if left.Unit != right.Unit {
				converted, err := e.env.units.Convert(right.Number, right.Unit, left.Unit)
				if err == nil {
					// Units are compatible, convert and divide
					right.Number = converted
					right.Unit = left.Unit
				}
				// If conversion fails, units are incompatible - we'll create a rate unit below
			}

			result := left.Number / right.Number
			// If units are the same (after conversion), return dimensionless number
			if left.Unit == right.Unit {
				return NewNumber(result)
			}
			// Otherwise, create rate unit for incompatible units
			rateUnit := left.Unit + "/" + right.Unit
			return NewUnit(result, rateUnit)
		}
		return NewUnit(left.Number/right.Number, left.Unit)

	default:
		return NewError(fmt.Sprintf("unknown operator: %s", op))
	}
}

// GetVariable retrieves a variable from the environment.
func (e *Evaluator) GetVariable(name string) (Value, bool) {
	val, ok := e.env.variables[name]
	return val, ok
}

// SetVariable sets a variable in the environment.
func (e *Evaluator) SetVariable(name string, val Value) {
	e.env.variables[name] = val
}

// Round rounds a value to the specified number of decimal places.
func Round(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}

func (e *Evaluator) evalWeekday(node *parser.WeekdayExpr) Value {
	now := time.Now()
	currentWeekday := now.Weekday()
	targetWeekday := node.Weekday

	// Calculate days until target weekday
	daysUntil := int(targetWeekday - currentWeekday)
	if daysUntil < 0 {
		daysUntil += 7
	}

	var result time.Time

	switch node.Modifier {
	case "next":
		// Next occurrence (at least 1 day away)
		if daysUntil == 0 {
			daysUntil = 7
		}
		result = now.AddDate(0, 0, daysUntil)
	case "last":
		// Last occurrence
		daysAgo := int(currentWeekday - targetWeekday)
		if daysAgo <= 0 {
			daysAgo += 7
		}
		result = now.AddDate(0, 0, -daysAgo)
	default:
		// This week (could be today or in the future this week)
		result = now.AddDate(0, 0, daysUntil)
	}

	// Normalise to start of day
	result = time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, result.Location())

	return NewDate(result)
}

func (e *Evaluator) evalMonth(node *parser.MonthExpr) Value {
	// Return the number of days in the specified month
	// We'll use the current year, or next year if we're past that month
	now := time.Now()

	// Map month name to month number
	monthMap := map[string]time.Month{
		"January": time.January, "February": time.February, "March": time.March,
		"April": time.April, "May": time.May, "June": time.June,
		"July": time.July, "August": time.August, "September": time.September,
		"October": time.October, "November": time.November, "December": time.December,
	}

	month, ok := monthMap[node.Month]
	if !ok {
		return NewError(fmt.Sprintf("unknown month: %s", node.Month))
	}

	// Use current year for the month
	year := now.Year()

	// Get the number of days in this month
	// Create date for first day of next month, then subtract one day
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	daysInMonth := float64(lastOfMonth.Day())

	return NewUnit(daysInMonth, "days")
}

func (e *Evaluator) evalTimeInLocation(node *parser.TimeInLocationExpr) Value {
	// Get current time in the specified location
	loc, err := e.env.timezone.GetLocation(node.Location)
	if err != nil {
		return NewError(err.Error())
	}

	// Get current UTC time and convert to target location
	now := time.Now().UTC()
	targetTime := now.Add(time.Duration(loc.Offset) * time.Hour)

	return NewDate(targetTime)
}

func (e *Evaluator) evalTimeDifference(node *parser.TimeDifferenceExpr) Value {
	offset, err := e.env.timezone.GetOffset(node.From, node.To)
	if err != nil {
		return NewError(err.Error())
	}

	// Convert to target unit if specified
	hours := float64(offset)
	if node.TargetUnit != "" {
		unit := strings.ToLower(node.TargetUnit)
		switch unit {
		case "day", "days", "d":
			hours = hours / 24
			return NewUnit(hours, "days")
		case "hour", "hours", "h", "hr":
			return NewUnit(hours, "hours")
		case "minute", "minutes", "min", "m":
			hours = hours * 60
			return NewUnit(hours, "minutes")
		case "second", "seconds", "sec", "s":
			hours = hours * 3600
			return NewUnit(hours, "seconds")
		default:
			return NewError(fmt.Sprintf("unsupported time unit: %s", node.TargetUnit))
		}
	}

	// Default: Return as hours
	return NewUnit(hours, "hours")
}

func (e *Evaluator) evalTimeConversion(node *parser.TimeConversionExpr) Value {
	// Start with current UTC time
	var baseTime time.Time
	if node.Time != nil {
		timeVal := e.Eval(node.Time)
		if timeVal.IsError() {
			return timeVal
		}
		baseTime = timeVal.Date
	} else {
		// Get current time in source location (as UTC + offset)
		fromLoc, err := e.env.timezone.GetLocation(node.From)
		if err != nil {
			return NewError(err.Error())
		}
		// Current time in the source location
		baseTime = time.Now().UTC().Add(time.Duration(fromLoc.Offset) * time.Hour)
	}

	// Apply offset if provided
	if node.Offset != nil {
		offsetVal := e.Eval(node.Offset)
		if offsetVal.IsError() {
			return offsetVal
		}

		// Convert offset to duration
		var offsetDuration time.Duration
		if offsetVal.Type == ValueUnit {
			// Handle unit-based offset (e.g., "3 hours")
			switch strings.ToLower(offsetVal.Unit) {
			case "hour", "hours", "h", "hr":
				offsetDuration = time.Duration(offsetVal.Number) * time.Hour
			case "minute", "minutes", "min":
				offsetDuration = time.Duration(offsetVal.Number) * time.Minute
			case "second", "seconds", "sec", "s":
				offsetDuration = time.Duration(offsetVal.Number) * time.Second
			default:
				return NewError(fmt.Sprintf("unsupported time unit: %s", offsetVal.Unit))
			}
		} else {
			// Assume hours if no unit specified
			offsetDuration = time.Duration(offsetVal.Number) * time.Hour
		}

		if node.Operator == "-" {
			offsetDuration = -offsetDuration
		}
		baseTime = baseTime.Add(offsetDuration)
	}

	// Convert to target location
	// baseTime is now in the source timezone (UTC + source offset + any offset applied)
	// To convert to target timezone, we need to: remove source offset, add target offset
	toLoc, err := e.env.timezone.GetLocation(node.To)
	if err != nil {
		return NewError(err.Error())
	}

	fromLoc, err := e.env.timezone.GetLocation(node.From)
	if err != nil {
		return NewError(err.Error())
	}

	// Remove source timezone offset to get back to UTC, then add target offset
	targetTime := baseTime.Add(time.Duration(toLoc.Offset-fromLoc.Offset) * time.Hour)

	return NewDate(targetTime)
}

func (e *Evaluator) evalRate(node *parser.RateExpr) Value {
	// Evaluate numerator and denominator
	num := e.Eval(node.Numerator)
	if num.IsError() {
		return num
	}

	den := e.Eval(node.Denominator)
	if den.IsError() {
		return den
	}

	// Both must have units
	if num.Type != ValueUnit {
		return NewError("rate numerator must have a unit")
	}

	if den.Type != ValueUnit {
		return NewError("rate denominator must have a unit")
	}

	// Calculate the rate value
	if den.Number == 0 {
		return NewError("division by zero in rate")
	}

	rateValue := num.Number / den.Number

	// Create compound unit string
	compoundUnit := num.Unit + "/" + den.Unit

	return NewUnit(rateValue, compoundUnit)
}
