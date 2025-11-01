package timezone

import (
	"testing"
	"time"
)

func TestGetLocation(t *testing.T) {
	s := NewSystem()
	
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"London", false},
		{"Singapore", false},
		{"Tokyo", false},
		{"New York", false},
		{"Unknown City", true},
	}
	
	for _, tt := range tests {
		_, err := s.GetLocation(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetLocation(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestGetOffset(t *testing.T) {
	s := NewSystem()
	
	tests := []struct {
		from     string
		to       string
		expected int
	}{
		{"London", "Singapore", 8},
		{"Singapore", "London", -8},
		{"New York", "London", 5},
		{"Tokyo", "London", -9},
	}
	
	for _, tt := range tests {
		offset, err := s.GetOffset(tt.from, tt.to)
		if err != nil {
			t.Errorf("GetOffset(%q, %q) error = %v", tt.from, tt.to, err)
			continue
		}
		
		if offset != tt.expected {
			t.Errorf("GetOffset(%q, %q) = %d, want %d", tt.from, tt.to, offset, tt.expected)
		}
	}
}

func TestConvertTime(t *testing.T) {
	s := NewSystem()
	
	// 10:00 in London
	londonTime := time.Date(2025, 11, 1, 10, 0, 0, 0, time.UTC)
	
	result, err := s.ConvertTime(londonTime, "London", "Singapore")
	if err != nil {
		t.Fatalf("ConvertTime failed: %v", err)
	}
	
	// Should be 18:00 in Singapore (8 hours ahead)
	expected := time.Date(2025, 11, 1, 18, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("ConvertTime got %v, want %v", result, expected)
	}
}

func TestParseTimeString(t *testing.T) {
	tests := []struct {
		input       string
		expectedHour int
		expectedMin  int
		wantErr     bool
	}{
		{"10:00", 10, 0, false},
		{"14:30", 14, 30, false},
		{"9:45", 9, 45, false},
		{"invalid", 0, 0, true},
	}
	
	for _, tt := range tests {
		result, err := ParseTimeString(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseTimeString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		
		if !tt.wantErr {
			if result.Hour() != tt.expectedHour || result.Minute() != tt.expectedMin {
				t.Errorf("ParseTimeString(%q) = %02d:%02d, want %02d:%02d",
					tt.input, result.Hour(), result.Minute(), tt.expectedHour, tt.expectedMin)
			}
		}
	}
}
