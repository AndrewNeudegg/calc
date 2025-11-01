package currency

import (
	"testing"
)

func TestNewSystem(t *testing.T) {
	s := NewSystem()
	
	if s == nil {
		t.Fatal("NewSystem should not return nil")
	}
	
	if len(s.rates) == 0 {
		t.Error("NewSystem should initialize default rates")
	}
}

func TestConvert(t *testing.T) {
	s := NewSystem()
	
	tests := []struct {
		amount   float64
		from     string
		to       string
		expected float64
		wantErr  bool
	}{
		{100, "USD", "USD", 100, false},
		{100, "$", "£", 78.74, false}, // 100 / 1.27
		{127, "£", "$", 161.29, false},  // 127 * 1.27
		{100, "USD", "EUR", 90.91, false}, // 100 / 1.10
		{100, "InvalidCur", "USD", 0, true},
		{100, "USD", "InvalidCur", 0, true},
	}
	
	for _, tt := range tests {
		result, err := s.Convert(tt.amount, tt.from, tt.to)
		
		if (err != nil) != tt.wantErr {
			t.Errorf("Convert(%f, %q, %q) error = %v, wantErr %v", 
				tt.amount, tt.from, tt.to, err, tt.wantErr)
			continue
		}
		
		if !tt.wantErr {
			// Allow small floating point differences
			diff := result - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.1 {
				t.Errorf("Convert(%f, %q, %q) = %f, want ~%f", 
					tt.amount, tt.from, tt.to, result, tt.expected)
			}
		}
	}
}

func TestIsCurrency(t *testing.T) {
	s := NewSystem()
	
	tests := []struct {
		symbol   string
		expected bool
	}{
		{"$", true},
		{"USD", true},
		{"£", true},
		{"GBP", true},
		{"€", true},
		{"EUR", true},
		{"¥", true},
		{"JPY", true},
		{"INVALID", false},
		{"XXX", false},
	}
	
	for _, tt := range tests {
		result := s.IsCurrency(tt.symbol)
		if result != tt.expected {
			t.Errorf("IsCurrency(%q) = %v, want %v", tt.symbol, result, tt.expected)
		}
	}
}

func TestSetRate(t *testing.T) {
	s := NewSystem()
	
	// Set a custom rate
	err := s.SetRate("USD", "GBP", 0.75) // 1 USD = 0.75 GBP
	if err != nil {
		t.Fatalf("SetRate failed: %v", err)
	}
	
	// Verify the conversion works with new rate
	result, err := s.Convert(100, "USD", "GBP")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	
	// Should be 75 GBP
	if result < 74 || result > 76 {
		t.Errorf("After SetRate, Convert(100, USD, GBP) = %f, want ~75", result)
	}
}

func TestSetRateInvalidCurrency(t *testing.T) {
	s := NewSystem()
	
	err := s.SetRate("INVALID", "USD", 1.0)
	if err == nil {
		t.Error("SetRate with invalid currency should return error")
	}
	
	err = s.SetRate("USD", "INVALID", 1.0)
	if err == nil {
		t.Error("SetRate with invalid target currency should return error")
	}
}
