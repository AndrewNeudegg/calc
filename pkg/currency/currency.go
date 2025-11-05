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

	// Additional common currencies (approximate static rates relative to USD)
	// Oceania
	s.rates["AUD"] = 0.654 // 1 AUD ≈ 0.654 USD
	s.rates["aud"] = 0.654
	s.rates["NZD"] = 0.595 // 1 NZD ≈ 0.595 USD
	s.rates["nzd"] = 0.595

	// Americas
	s.rates["CAD"] = 0.730 // 1 CAD ≈ 0.73 USD
	s.rates["cad"] = 0.730
	s.rates["MXN"] = 0.057 // 1 MXN ≈ 0.057 USD
	s.rates["mxn"] = 0.057
	s.rates["BRL"] = 0.196 // 1 BRL ≈ 0.196 USD
	s.rates["brl"] = 0.196

	// Europe (non-EUR)
	s.rates["CHF"] = 1.110 // 1 CHF ≈ 1.11 USD
	s.rates["chf"] = 1.110
	s.rates["SEK"] = 0.091 // 1 SEK ≈ 0.091 USD
	s.rates["sek"] = 0.091
	s.rates["NOK"] = 0.091 // 1 NOK ≈ 0.091 USD
	s.rates["nok"] = 0.091
	s.rates["DKK"] = 0.143 // 1 DKK ≈ 0.143 USD
	s.rates["dkk"] = 0.143
	s.rates["PLN"] = 0.250 // 1 PLN ≈ 0.25 USD
	s.rates["pln"] = 0.250
	s.rates["CZK"] = 0.043 // 1 CZK ≈ 0.043 USD
	s.rates["czk"] = 0.043
	s.rates["HUF"] = 0.0028 // 1 HUF ≈ 0.0028 USD
	s.rates["huf"] = 0.0028
	s.rates["RON"] = 0.217 // 1 RON ≈ 0.217 USD
	s.rates["ron"] = 0.217
	s.rates["RUB"] = 0.010 // 1 RUB ≈ 0.01 USD
	s.rates["rub"] = 0.010
	s.rates["TRY"] = 0.036 // 1 TRY ≈ 0.036 USD
	s.rates["try"] = 0.036

	// Middle East
	s.rates["AED"] = 0.272 // 1 AED ≈ 0.272 USD
	s.rates["aed"] = 0.272
	s.rates["SAR"] = 0.267 // 1 SAR ≈ 0.267 USD
	s.rates["sar"] = 0.267
	s.rates["ILS"] = 0.263 // 1 ILS ≈ 0.263 USD
	s.rates["ils"] = 0.263

	// Asia
	s.rates["CNY"] = 0.137 // 1 CNY ≈ 0.137 USD
	s.rates["cny"] = 0.137
	s.rates["HKD"] = 0.128 // 1 HKD ≈ 0.128 USD
	s.rates["hkd"] = 0.128
	s.rates["SGD"] = 0.735 // 1 SGD ≈ 0.735 USD
	s.rates["sgd"] = 0.735
	s.rates["INR"] = 0.012 // 1 INR ≈ 0.012 USD
	s.rates["inr"] = 0.012
	s.rates["KRW"] = 0.00074 // 1 KRW ≈ 0.00074 USD
	s.rates["krw"] = 0.00074
	s.rates["TWD"] = 0.031 // 1 TWD ≈ 0.031 USD
	s.rates["twd"] = 0.031
	s.rates["THB"] = 0.028 // 1 THB ≈ 0.028 USD
	s.rates["thb"] = 0.028
	s.rates["MYR"] = 0.213 // 1 MYR ≈ 0.213 USD
	s.rates["myr"] = 0.213
	s.rates["IDR"] = 0.000064 // 1 IDR ≈ 0.000064 USD
	s.rates["idr"] = 0.000064
	s.rates["PHP"] = 0.018 // 1 PHP ≈ 0.018 USD
	s.rates["php"] = 0.018

	// Africa
	s.rates["ZAR"] = 0.054 // 1 ZAR ≈ 0.054 USD
	s.rates["zar"] = 0.054
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
