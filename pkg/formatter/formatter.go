package formatter

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/andrewneudegg/calc/pkg/evaluator"
	"github.com/andrewneudegg/calc/pkg/settings"
)

// Formatter formats values according to settings.
type Formatter struct {
	settings *settings.Settings
}

// New creates a new formatter.
func New(s *settings.Settings) *Formatter {
	return &Formatter{settings: s}
}

// Format formats a value according to settings.
func (f *Formatter) Format(val evaluator.Value) string {
	if val.IsError() {
		return fmt.Sprintf("Error: %s", val.Error)
	}

	switch val.Type {
	case evaluator.ValueNumber:
		return f.formatNumber(val.Number)
	case evaluator.ValueUnit:
		return fmt.Sprintf("%s %s", f.formatNumber(val.Number), val.Unit)
	case evaluator.ValueCurrency:
		return fmt.Sprintf("%s%s", val.Currency, f.formatNumber(val.Number))
	case evaluator.ValuePercent:
		return fmt.Sprintf("%s%%", f.formatNumber(val.Number))
	case evaluator.ValueDate:
		return f.formatDate(val.Date)
	default:
		return "unknown"
	}
}

func (f *Formatter) formatDate(d time.Time) string {
	// If the time has a non-zero time component (hours, minutes, seconds),
	// show the time as well as the date
	if d.Hour() != 0 || d.Minute() != 0 || d.Second() != 0 {
		// Use a format that includes time
		return d.Format("2 Jan 2006 15:04:05 MST")
	}
	// Otherwise just show the date
	return d.Format(f.settings.DateFormat)
}

func (f *Formatter) formatNumber(n float64) string {
	// Round to precision
	rounded := f.round(n, f.settings.Precision)

	// Format with thousand separators for UK locale
	if f.settings.Locale == "en_GB" || f.settings.Locale == "en_UK" {
		return f.formatWithCommas(rounded, f.settings.Precision)
	}

	// Default format
	format := fmt.Sprintf("%%.%df", f.settings.Precision)
	return fmt.Sprintf(format, rounded)
}

func (f *Formatter) round(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}

func (f *Formatter) formatWithCommas(n float64, decimals int) string {
	// Split into integer and decimal parts
	integer := int64(math.Abs(n))
	decimal := n - float64(int64(n))

	// Format integer part with commas
	intStr := fmt.Sprintf("%d", integer)
	var parts []string

	for i := len(intStr); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{intStr[start:i]}, parts...)
	}

	result := strings.Join(parts, ",")

	// Add negative sign if needed
	if n < 0 {
		result = "-" + result
	}

	// Add decimal part if needed
	if decimals > 0 {
		decStr := fmt.Sprintf("%.*f", decimals, math.Abs(decimal))
		// Remove leading "0"
		if len(decStr) > 2 {
			result += decStr[1:]
		} else {
			result += ".00"
		}
	}

	return result
}
