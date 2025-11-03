package currency

import (
	"fmt"
	"strings"
)

// System manages currency conversions.
type System struct {
	rates map[string]float64 // rates relative to USD
}

// NewSystem creates a new currency system with default rates.
func NewSystem() *System {
	s := &System{
		rates: make(map[string]float64),
	}
	s.initDefaultRates()
	return s
}

func (s *System) initDefaultRates() {
	// Base: USD = 1.0
	s.rates["$"] = 1.0
	s.rates["USD"] = 1.0
	s.rates["usd"] = 1.0

	// British Pound
	s.rates["£"] = 1.27 // 1 GBP = 1.27 USD
	s.rates["GBP"] = 1.27
	s.rates["gbp"] = 1.27

	// Euro
	s.rates["€"] = 1.10 // 1 EUR = 1.10 USD
	s.rates["EUR"] = 1.10
	s.rates["eur"] = 1.10

	// Japanese Yen
	s.rates["¥"] = 0.0067 // 1 JPY = 0.0067 USD
	s.rates["JPY"] = 0.0067
	s.rates["jpy"] = 0.0067
}

// SetRate sets a custom exchange rate.
func (s *System) SetRate(from, to string, rate float64) error {
	from = s.normaliseCurrency(from)
	to = s.normaliseCurrency(to)

	// Convert both to their USD equivalents
	fromRate, ok := s.rates[from]
	if !ok {
		return fmt.Errorf("unknown currency: %s", from)
	}

	_, ok = s.rates[to]
	if !ok {
		return fmt.Errorf("unknown currency: %s", to)
	}

	// Update the conversion rate
	// If 1 USD = X GBP, then we need to update GBP's rate relative to USD
	s.rates[to] = fromRate / rate

	return nil
}

// Convert converts an amount from one currency to another.
func (s *System) Convert(amount float64, from, to string) (float64, error) {
	from = s.normaliseCurrency(from)
	to = s.normaliseCurrency(to)

	fromRate, ok := s.rates[from]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", from)
	}

	toRate, ok := s.rates[to]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", to)
	}

	// Convert to USD, then to target currency
	usd := amount * fromRate
	result := usd / toRate

	return result, nil
}

// IsCurrency checks if a string is a known currency.
func (s *System) IsCurrency(symbol string) bool {
	_, ok := s.rates[s.normaliseCurrency(symbol)]
	return ok
}

func (s *System) normaliseCurrency(cur string) string {
	cur = strings.TrimSpace(cur)

	// Map symbols to codes
	switch cur {
	case "$":
		return "USD"
	case "£":
		return "GBP"
	case "€":
		return "EUR"
	case "¥":
		return "JPY"
	default:
		// Handle currency names
		upper := strings.ToUpper(cur)
		switch upper {
		case "DOLLAR", "DOLLARS":
			return "USD"
		case "EURO", "EUROS":
			return "EUR"
		case "YEN":
			return "JPY"
		// "POUND" and "POUNDS" are ambiguous (weight vs currency)
		// So we don't map them here - users should use "gbp" or "£"
		default:
			return upper
		}
	}
}

// GetSymbol returns the symbol for a currency code.
func (s *System) GetSymbol(code string) string {
	// First normalize the currency name/code
	normalized := s.normaliseCurrency(code)

	switch normalized {
	case "USD":
		return "$"
	case "GBP":
		return "£"
	case "EUR":
		return "€"
	case "JPY":
		return "¥"
	default:
		return normalized
	}
}
