package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrewneudegg/calc/pkg/settings"
)

func TestSaveWritesSettingsAndMessage(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "settings.json")
	s := settings.Default()
	s.Precision = 5
	s.Currency = "EUR"
	s.ConfigPath = cfg

	h := New(s)
	msg := h.Execute("save", []string{"session.log"})
	if msg != "saved to session.log" {
		t.Fatalf("unexpected message: %q", msg)
	}

	// File should exist with JSON reflecting our settings
	b, err := os.ReadFile(cfg)
	if err != nil {
		t.Fatalf("reading saved settings: %v", err)
	}
	var got settings.Settings
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal settings: %v", err)
	}
	if got.Precision != 5 || got.Currency != "EUR" || got.DateFormat != s.DateFormat || got.Locale != s.Locale || got.FuzzyMode != s.FuzzyMode {
		t.Fatalf("settings not saved correctly: %+v", got)
	}
}

func TestSaveUsageNoArgs(t *testing.T) {
	s := settings.Default()
	s.ConfigPath = filepath.Join(t.TempDir(), "settings.json")
	h := New(s)
	msg := h.Execute("save", nil)
	if msg != "usage: :save <filename>" {
		t.Fatalf("unexpected usage: %q", msg)
	}
}

func TestOpenMessageAndUsage(t *testing.T) {
	s := settings.Default()
	s.ConfigPath = filepath.Join(t.TempDir(), "settings.json")
	h := New(s)

	msg := h.Execute("open", []string{"workspace.calc"})
	if msg != "loaded workspace.calc" {
		t.Fatalf("unexpected open message: %q", msg)
	}

	msg = h.Execute("open", nil)
	if msg != "usage: :open <filename>" {
		t.Fatalf("unexpected open usage: %q", msg)
	}
}
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

func TestExecuteClearInvokesCallbackAndReturnsAnsi(t *testing.T) {
	s := settings.Default()
	h := New(s)

	called := false
	h.ClearWorkspace = func() error { called = true; return nil }

	out := h.Execute("clear", nil)
	if !called {
		t.Fatalf(":clear did not invoke ClearWorkspace callback")
	}
	if !strings.Contains(out, "\x1b[2J") {
		t.Fatalf("expected ANSI clear sequence in output, got: %q", out)
	}
}

func TestExecuteQuietTogglesAndSets(t *testing.T) {
	s := settings.Default()
	h := New(s)

	state := false
	h.GetQuiet = func() bool { return state }
	h.SetQuiet = func(b bool) { state = b }
	h.ToggleQuiet = func() bool { state = !state; return state }

	// Toggle with no args
	out := h.Execute("quiet", nil)
	if !state || !strings.Contains(out, "quiet: on") {
		t.Fatalf(":quiet should toggle on, got state=%v, out=%q", state, out)
	}

	// Explicit off
	out = h.Execute("quiet", []string{"off"})
	if state || !strings.Contains(out, "quiet: off") {
		t.Fatalf(":quiet off should set off, got state=%v, out=%q", state, out)
	}

	// Explicit on
	out = h.Execute("quiet", []string{"on"})
	if !state || !strings.Contains(out, "quiet: on") {
		t.Fatalf(":quiet on should set on, got state=%v, out=%q", state, out)
	}

	// Bad arg
	out = h.Execute("quiet", []string{"maybe"})
	if !strings.Contains(out, "usage") {
		t.Fatalf(":quiet with bad arg should show usage, got %q", out)
	}
}
