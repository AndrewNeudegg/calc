package commands

import (
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/settings"
)

func TestExecuteHelp(t *testing.T) {
	s := settings.Default()
	h := New(s)

	result := h.Execute("help", nil)

	if !strings.Contains(result, "Available commands") && !strings.Contains(result, "Commands") {
		t.Errorf("Help output should contain command list, got: %s", result)
	}
}

func TestExecuteSet(t *testing.T) {
	s := settings.Default()
	h := New(s)

	tests := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{"locale", "en_US"}, false},
		{[]string{"precision", "4"}, false},
		{[]string{"invalid"}, true},
		{nil, true},
	}

	for _, tt := range tests {
		result := h.Execute("set", tt.args)
		hasError := strings.Contains(result, "error") || strings.Contains(result, "usage")
		if hasError != tt.wantErr {
			t.Errorf("Execute(set, %v) error = %v, wantErr %v, result: %s", tt.args, hasError, tt.wantErr, result)
		}
	}
}

func TestExecuteTzList(t *testing.T) {
	s := settings.Default()
	h := New(s)

	result := h.Execute("tz", []string{"list"})

	expectedCities := []string{"London", "New York", "Tokyo", "Singapore"}
	for _, city := range expectedCities {
		if !strings.Contains(result, city) {
			t.Errorf("Timezone list should contain %q, got: %s", city, result)
		}
	}
}

func TestExecuteSave(t *testing.T) {
	s := settings.Default()
	s.ConfigPath = t.TempDir() + "/settings.json"
	h := New(s)

	result := h.Execute("save", []string{"test.calc"})

	if strings.Contains(result, "error") && !strings.Contains(result, "saved") {
		t.Errorf("Save should work: %s", result)
	}
}

func TestExecuteOpen(t *testing.T) {
	s := settings.Default()
	h := New(s)

	result := h.Execute("open", []string{"test.calc"})

	// Open should not crash
	if result == "" {
		t.Error("Open should return a message")
	}
}

func TestExecuteUnknown(t *testing.T) {
	s := settings.Default()
	h := New(s)

	result := h.Execute("unknown", nil)

	if !strings.Contains(result, "unknown") {
		t.Error("Execute(unknown) should return unknown command error")
	}
}

func TestSettingsIntegration(t *testing.T) {
	s := settings.Default()
	h := New(s)

	// Set locale
	result := h.Execute("set", []string{"locale", "en_US"})
	if strings.Contains(result, "error") {
		t.Fatalf("Failed to set locale: %s", result)
	}

	if s.Locale != "en_US" {
		t.Errorf("Expected locale en_US, got %s", s.Locale)
	}

	// Set precision
	result = h.Execute("set", []string{"precision", "5"})
	if strings.Contains(result, "error") {
		t.Fatalf("Failed to set precision: %s", result)
	}

	if s.Precision != 5 {
		t.Errorf("Expected precision 5, got %d", s.Precision)
	}
}
