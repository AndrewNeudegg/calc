package evaluator

import (
	"fmt"
	"time"
)

// ValueType represents the type of a value.
type ValueType int

const (
	ValueNumber ValueType = iota
	ValueUnit
	ValueCurrency
	ValuePercent
	ValueDate
	ValueString
	ValueError
)

// Value represents an evaluated value.
type Value struct {
	Type     ValueType
	Number   float64
	Unit     string
	Currency string
	Date     time.Time
	Text     string
	Error    string
}

// NewNumber creates a new number value.
func NewNumber(n float64) Value {
	return Value{Type: ValueNumber, Number: n}
}

// NewUnit creates a new unit value.
func NewUnit(n float64, unit string) Value {
	return Value{Type: ValueUnit, Number: n, Unit: unit}
}

// NewCurrency creates a new currency value.
func NewCurrency(n float64, currency string) Value {
	return Value{Type: ValueCurrency, Number: n, Currency: currency}
}

// NewPercent creates a new percentage value.
func NewPercent(n float64) Value {
	return Value{Type: ValuePercent, Number: n}
}

// NewDate creates a new date value.
func NewDate(d time.Time) Value {
	return Value{Type: ValueDate, Date: d}
}

// NewString creates a new string value.
func NewString(s string) Value {
	return Value{Type: ValueString, Text: s}
}

// NewError creates a new error value.
func NewError(msg string) Value {
	return Value{Type: ValueError, Error: msg}
}

// IsError returns true if the value is an error.
func (v Value) IsError() bool {
	return v.Type == ValueError
}

// String returns a string representation of the value.
func (v Value) String() string {
	switch v.Type {
	case ValueNumber:
		return fmt.Sprintf("%.2f", v.Number)
	case ValueUnit:
		return fmt.Sprintf("%.2f %s", v.Number, v.Unit)
	case ValueCurrency:
		return fmt.Sprintf("%s%.2f", v.Currency, v.Number)
	case ValuePercent:
		return fmt.Sprintf("%.2f%%", v.Number)
	case ValueDate:
		// If time-of-day is non-zero, include time and timezone like formatter.formatDate
		if v.Date.Hour() != 0 || v.Date.Minute() != 0 || v.Date.Second() != 0 {
			return v.Date.Format("2 Jan 2006 15:04:05 MST")
		}
		return v.Date.Format("2 Jan 2006")
	case ValueString:
		return v.Text
	case ValueError:
		return fmt.Sprintf("Error: %s", v.Error)
	default:
		return "unknown"
	}
}
