package evaluator

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/andrewneudegg/calc/pkg/businessdays"
	"github.com/andrewneudegg/calc/pkg/constants"
	"github.com/andrewneudegg/calc/pkg/currency"
	"github.com/andrewneudegg/calc/pkg/parser"
	"github.com/andrewneudegg/calc/pkg/timezone"
	"github.com/andrewneudegg/calc/pkg/units"
)

// Environment stores variables and state.
type Environment struct {
	variables           map[string]Value
	units               *units.System
	currency            *currency.System
	timezone            *timezone.System
	constants           *constants.System
	locale              string // Current locale for business days and number formatting
	historyFunc         func(offset int) (Value, error)   // Function to get previous results by relative offset
	absoluteHistoryFunc func(lineID int) (Value, error)   // Function to get result by absolute line ID
}

// NewEnvironment creates a new evaluation environment.
func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		units:     units.NewSystem(),
		currency:  currency.NewSystem(),
		timezone:  timezone.NewSystem(),
		constants: constants.NewSystem(),
		locale:    "en_GB", // Default locale
	}
}

// SetHistoryFunc sets the function to retrieve previous results.
func (e *Environment) SetHistoryFunc(f func(offset int) (Value, error)) {
	e.historyFunc = f
}

// SetAbsoluteHistoryFunc sets the function to retrieve results by absolute line ID.
func (e *Environment) SetAbsoluteHistoryFunc(f func(lineID int) (Value, error)) {
	e.absoluteHistoryFunc = f
}

// SetVariable sets a variable in the environment.
func (e *Environment) SetVariable(name string, value Value) {
	e.variables[name] = value
}

// GetVariableNames returns a list of all variable names in the environment.
func (e *Environment) GetVariableNames() []string {
	names := make([]string, 0, len(e.variables))
	for name := range e.variables {
		names = append(names, name)
	}
	return names
}

// Units returns the units system.
func (e *Environment) Units() *units.System {
	return e.units
}

// Currency returns the currency system.
func (e *Environment) Currency() *currency.System {
	return e.currency
}

// Constants returns the constants system.
func (e *Environment) Constants() *constants.System {
	return e.constants
}

// SetLocale sets the locale for the environment (for business days and number formatting).
func (e *Environment) SetLocale(locale string) {
	e.locale = locale
}

// GetLocale returns the current locale.
func (e *Environment) GetLocale() string {
	return e.locale
}

// Eval evaluates an expression using this environment.
func (e *Environment) Eval(expr parser.Expr) Value {
	evaluator := New(e)
	return evaluator.Eval(expr)
}

// Evaluator evaluates expressions.
type Evaluator struct {
	env *Environment
}

// New creates a new evaluator.
func New(env *Environment) *Evaluator {
	return &Evaluator{env: env}
}

// isBusinessDayUnit checks if the given unit string represents business days.
// It handles both singular and plural forms, case-insensitively.
func isBusinessDayUnit(unit string) bool {
	lower := strings.ToLower(unit)
	return lower == "business day" || lower == "business days"
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

	case *parser.PrevExpr:
		return e.evalPrev(node)

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

		// Check for business days (exact match)
		lowerUnit := strings.ToLower(unit)
		if lowerUnit == "business day" || lowerUnit == "business days" {
			// Get the business day schedule for the current locale
			schedule := businessdays.GetScheduleForLocale(e.env.GetLocale())
			newDate = businessdays.AddBusinessDays(left.Date, int(offset), schedule)
		} else {
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
		}

		return NewDate(newDate)
	}

	// Handle time unit - date/now (for "time until" functionality)
	// Convert time unit to a TimeExpr on today's date and then subtract
	if left.Type == ValueUnit && left.Unit == "time" && right.Type == ValueDate && node.Operator == "-" {
		// Convert time unit (decimal hours) to time on today
		now := time.Now()
		hours := int(left.Number)
		minutes := int((left.Number - float64(hours)) * 60)
		targetTime := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())
		
		// Calculate time difference
		duration := targetTime.Sub(right.Date)
		hours64 := duration.Hours()
		
		// If the result is negative and we're within 24 hours, assume next day
		if hours64 < 0 && hours64 > -24 {
			hours64 += 24
		}
		
		// Format as HH:MM
		totalMinutes := int(hours64 * 60)
		h := totalMinutes / 60
		m := totalMinutes % 60
		
		// Return as a formatted time string
		timeStr := fmt.Sprintf("%d:%02d", h, m)
		return NewString(timeStr)
	}

	// Handle date-date subtraction (returns days with unit and stores date range for business day conversion)
	if left.Type == ValueDate && right.Type == ValueDate && node.Operator == "-" {
		// Check if both times are on the same day (time-only arithmetic)
		ly, lm, ld := left.Date.Date()
		ry, rm, rd := right.Date.Date()
		
		// Check if at least one has a non-midnight time component
		leftHasTime := left.Date.Hour() != 0 || left.Date.Minute() != 0 || left.Date.Second() != 0
		rightHasTime := right.Date.Hour() != 0 || right.Date.Minute() != 0 || right.Date.Second() != 0
		
		// If both dates are on the same day AND at least one has a time component, treat as time-only arithmetic
		if ly == ry && lm == rm && ld == rd && (leftHasTime || rightHasTime) {
			// Calculate time difference in hours
			duration := left.Date.Sub(right.Date)
			hours := duration.Hours()
			
			// If the result is negative and we're within 24 hours, assume next day
			if hours < 0 && hours > -24 {
				hours += 24
			}
			
			// Format as HH:MM
			totalMinutes := int(hours * 60)
			h := totalMinutes / 60
			m := totalMinutes % 60
			
			// Return as a formatted time string
			timeStr := fmt.Sprintf("%d:%02d", h, m)
			return NewString(timeStr)
		}
		
		// Otherwise, normalize both dates to UTC midnight to avoid DST-related fractional day counts
		leftMidnight := time.Date(ly, lm, ld, 0, 0, 0, 0, time.UTC)
		rightMidnight := time.Date(ry, rm, rd, 0, 0, 0, 0, time.UTC)
		duration := leftMidnight.Sub(rightMidnight)
		days := duration.Hours() / 24.0
		return NewDateDifference(days, right.Date, left.Date)
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
		// Check if it's a physical constant
		if e.env.constants != nil && e.env.constants.IsConstant(node.Name) {
			c, err := e.env.constants.GetConstant(node.Name)
			if err == nil {
				// Return constant as a unit value
				return NewUnit(c.Value, c.Unit)
			}
		}
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
		// Check if this is a date difference being converted to business days
		if isBusinessDayUnit(node.ToUnit) && val.StartDate != nil && val.EndDate != nil {
			// Use the stored date range to count business days
			schedule := businessdays.GetScheduleForLocale(e.env.GetLocale())
			businessDays := businessdays.CountBusinessDays(*val.StartDate, *val.EndDate, schedule)
			return NewUnit(float64(businessDays), node.ToUnit)
		}
		
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

	// Handle date-date subtraction (e.g., "today - 19/09/2025")
	// This occurs when no unit is specified and offset is a date
	if base.Type == ValueDate && offset.Type == ValueDate && node.Operator == "-" && node.Unit == "" {
		// Normalize both dates to UTC midnight to avoid DST-related fractional day counts
		by, bm, bd := base.Date.Date()
		oy, om, od := offset.Date.Date()
		baseMidnight := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
		offsetMidnight := time.Date(oy, om, od, 0, 0, 0, 0, time.UTC)
		duration := baseMidnight.Sub(offsetMidnight)
		days := duration.Hours() / 24.0
		return NewDateDifference(days, offset.Date, base.Date)
	}

	offsetVal := int(offset.Number)
	if node.Operator == "-" {
		offsetVal = -offsetVal
	}

	var result time.Time
	unit := strings.ToLower(node.Unit)

	// Check for business days (exact match)
	if unit == "business day" || unit == "business days" {
		schedule := businessdays.GetScheduleForLocale(e.env.GetLocale())
		result = businessdays.AddBusinessDays(base.Date, offsetVal, schedule)
	} else {
		switch unit {
		case "day", "days", "d":
			result = base.Date.AddDate(0, 0, offsetVal)
		case "week", "weeks", "w":
			result = base.Date.AddDate(0, 0, offsetVal*7)
		case "month", "months", "mo":
			result = base.Date.AddDate(0, offsetVal, 0)
		case "year", "years", "y":
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
			// Both are units - multiply and simplify
			resultNum := left.Number * right.Number
			simplifiedUnit := e.simplifyUnits(left.Unit, right.Unit, "*")
			
			// If simplification resulted in a dimensionless unit, return a plain number
			if simplifiedUnit == "1" || simplifiedUnit == "" {
				return NewNumber(resultNum)
			}
			
			return NewUnit(resultNum, simplifiedUnit)
		}
		return NewUnit(left.Number*right.Number, left.Unit)

	case "/":
		if right.Number == 0 {
			return NewError("division by zero")
		}
		if right.Type == ValueUnit {
			// If left is a plain number (not a unit), result is in 1/right.Unit
			if left.Type != ValueUnit {
				return NewUnit(left.Number/right.Number, "1/"+right.Unit)
			}
			
			// Both are units - divide and simplify
			resultNum := left.Number / right.Number
			
			// For division, try to convert if units are compatible first
			if left.Unit != right.Unit {
				converted, err := e.env.units.Convert(right.Number, right.Unit, left.Unit)
				if err == nil {
					// Units are compatible, convert and divide - result is dimensionless
					return NewNumber(left.Number / converted)
				}
			}
			
			// If units are the same, return dimensionless number
			if left.Unit == right.Unit {
				return NewNumber(resultNum)
			}
			
			// Units are incompatible - create and simplify compound unit
			simplifiedUnit := e.simplifyUnits(left.Unit, right.Unit, "/")
			
			// If simplification resulted in a dimensionless unit, return a plain number
			if simplifiedUnit == "1" || simplifiedUnit == "" {
				return NewNumber(resultNum)
			}
			
			return NewUnit(resultNum, simplifiedUnit)
		}
		return NewUnit(left.Number/right.Number, left.Unit)

	default:
		return NewError(fmt.Sprintf("unknown operator: %s", op))
	}
}

// simplifyUnits simplifies unit expressions by canceling matching units in numerator and denominator.
// For multiplication: $/hr * hours -> $ (hr cancels with hours)
// For division: km / hours -> km/hours (creates compound), then can cancel if matched
func (e *Evaluator) simplifyUnits(leftUnit, rightUnit, op string) string {
	if op == "*" {
		return e.simplifyMultiplication(leftUnit, rightUnit)
	}
	if op == "/" {
		return e.simplifyDivision(leftUnit, rightUnit)
	}
	return leftUnit + "/" + rightUnit
}

// simplifyMultiplication simplifies unit multiplication, handling cancellation
// Examples:
// - $/hr * hours -> $
// - $/hour * hours -> $
// - km/h * h -> km
// - m/s * s -> m
func (e *Evaluator) simplifyMultiplication(leftUnit, rightUnit string) string {
	// Parse left and right units to extract numerator and denominator
	leftParts := e.parseCompoundUnit(leftUnit)
	rightParts := e.parseCompoundUnit(rightUnit)
	
	// Combine numerators and denominators
	numerators := append(leftParts.numerators, rightParts.numerators...)
	denominators := append(leftParts.denominators, rightParts.denominators...)
	
	// Cancel matching units
	numerators, denominators = e.cancelUnits(numerators, denominators)
	
	// Build result string
	return e.buildUnitString(numerators, denominators)
}

// simplifyDivision simplifies unit division, handling cancellation
// Examples:
// - km / hour -> km/hour
// - (km/hour) / km -> 1/hour
// - m / m -> 1 (dimensionless)
func (e *Evaluator) simplifyDivision(leftUnit, rightUnit string) string {
	// Parse left and right units to extract numerator and denominator
	leftParts := e.parseCompoundUnit(leftUnit)
	rightParts := e.parseCompoundUnit(rightUnit)
	
	// For division: left/right means left numerators over (left denominators + right numerators)
	// and left denominators become (left denominators + right numerators)
	// But we need to think of it as: (leftNum/leftDen) / (rightNum/rightDen) = (leftNum * rightDen) / (leftDen * rightNum)
	numerators := append(leftParts.numerators, rightParts.denominators...)
	denominators := append(leftParts.denominators, rightParts.numerators...)
	
	// Cancel matching units
	numerators, denominators = e.cancelUnits(numerators, denominators)
	
	// Build result string
	return e.buildUnitString(numerators, denominators)
}

// unitParts holds the parsed parts of a compound unit
type unitParts struct {
	numerators   []string
	denominators []string
}

// parseCompoundUnit parses a unit string into numerator and denominator parts
// Examples:
// - "$/hr" -> numerators: ["$"], denominators: ["hr"]
// - "m" -> numerators: ["m"], denominators: []
// - "m/s" -> numerators: ["m"], denominators: ["s"]
// - "$/hr·hours" -> numerators: ["$", "hours"], denominators: ["hr"]
// - "$/hr*hours" -> numerators: ["$", "hours"], denominators: ["hr"]
func (e *Evaluator) parseCompoundUnit(unit string) unitParts {
	parts := unitParts{
		numerators:   []string{},
		denominators: []string{},
	}
	
	if unit == "" || unit == "1" {
		return parts
	}
	
	// First, split by · (middle dot) or * to handle products
	// Support both characters for better usability
	products := []string{}
	if strings.Contains(unit, "·") {
		products = strings.Split(unit, "·")
	} else if strings.Contains(unit, "*") {
		products = strings.Split(unit, "*")
	} else {
		products = []string{unit}
	}
	
	for _, product := range products {
		// Each product might be a simple unit or a ratio (num/den)
		if strings.Contains(product, "/") {
			// It's a ratio
			ratio := strings.Split(product, "/")
			if len(ratio) == 2 {
				num := strings.TrimSpace(ratio[0])
				den := strings.TrimSpace(ratio[1])
				if num != "" && num != "1" {
					parts.numerators = append(parts.numerators, num)
				}
				if den != "" && den != "1" {
					parts.denominators = append(parts.denominators, den)
				}
			}
		} else {
			// Simple unit - goes to numerator
			simple := strings.TrimSpace(product)
			if simple != "" && simple != "1" {
				parts.numerators = append(parts.numerators, simple)
			}
		}
	}
	
	return parts
}

// cancelUnits cancels matching units between numerators and denominators
// Returns the simplified numerators and denominators
func (e *Evaluator) cancelUnits(numerators, denominators []string) ([]string, []string) {
	// Create copies to avoid modifying the originals
	nums := make([]string, len(numerators))
	copy(nums, numerators)
	dens := make([]string, len(denominators))
	copy(dens, denominators)
	
	// Try to cancel each numerator with each denominator
	for i := 0; i < len(nums); i++ {
		if nums[i] == "" {
			continue
		}
		for j := 0; j < len(dens); j++ {
			if dens[j] == "" {
				continue
			}
			// Check if units match (accounting for variations)
			if e.unitsMatch(nums[i], dens[j]) {
				// Cancel them out
				nums[i] = ""
				dens[j] = ""
				break
			}
		}
	}
	
	// Filter out empty strings
	resultNums := []string{}
	for _, n := range nums {
		if n != "" {
			resultNums = append(resultNums, n)
		}
	}
	
	resultDens := []string{}
	for _, d := range dens {
		if d != "" {
			resultDens = append(resultDens, d)
		}
	}
	
	return resultNums, resultDens
}

// unitsMatch checks if two unit strings represent the same unit
// Handles variations like "hr" vs "hours", "h" vs "hour", etc.
// Note: This performs dimension lookups for each check, but since unit cancellations
// are relatively infrequent and unit lists are small, the performance impact is minimal.
// If needed in the future, consider caching unit equivalence results.
func (e *Evaluator) unitsMatch(unit1, unit2 string) bool {
	// Exact match
	if unit1 == unit2 {
		return true
	}
	
	// Normalize and compare
	u1 := strings.ToLower(strings.TrimSpace(unit1))
	u2 := strings.ToLower(strings.TrimSpace(unit2))
	
	// Check if they're the same unit in the units system
	// Try to get the dimension for each unit
	if e.env.units.IsUnit(u1) && e.env.units.IsUnit(u2) {
		dim1, err1 := e.env.units.GetDimension(u1)
		dim2, err2 := e.env.units.GetDimension(u2)
		
		// If both are recognized units with the same dimension
		if err1 == nil && err2 == nil && dim1 == dim2 {
			// Check if they convert 1:1 (i.e., they're the same unit with different names)
			// e.g., "meter" and "metre"
			val, err := e.env.units.Convert(1.0, u1, u2)
			if err == nil && val == 1.0 {
				return true
			}
		}
	}
	
	return false
}

// buildUnitString builds a unit string from numerators and denominators
func (e *Evaluator) buildUnitString(numerators, denominators []string) string {
	if len(numerators) == 0 && len(denominators) == 0 {
		return "1" // Dimensionless
	}
	
	numStr := ""
	if len(numerators) == 0 {
		numStr = "1"
	} else {
		numStr = strings.Join(numerators, "·")
	}
	
	if len(denominators) == 0 {
		return numStr
	}
	
	denStr := strings.Join(denominators, "·")
	return numStr + "/" + denStr
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

func (e *Evaluator) evalPrev(node *parser.PrevExpr) Value {
	if node.Absolute {
		// Absolute position: prev#N
		if e.env.absoluteHistoryFunc == nil {
			return NewError("prev is only available in REPL mode")
		}
		
		val, err := e.env.absoluteHistoryFunc(node.Offset)
		if err != nil {
			return NewError(err.Error())
		}
		
		return val
	} else {
		// Relative offset: prev, prev~N
		if e.env.historyFunc == nil {
			return NewError("prev is only available in REPL mode")
		}
		
		val, err := e.env.historyFunc(node.Offset)
		if err != nil {
			return NewError(err.Error())
		}
		
		return val
	}
}
