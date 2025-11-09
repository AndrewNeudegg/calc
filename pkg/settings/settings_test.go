package settings

import (
	"path/filepath"
	"testing"
)

func TestDefaultSettings(t *testing.T) {
	s := Default()

	if s.Locale != "en_GB" {
		t.Errorf("Expected default locale en_GB, got %s", s.Locale)
	}

	if s.Precision != 2 {
		t.Errorf("Expected default precision 2, got %d", s.Precision)
	}

	if s.Currency != "GBP" {
		t.Errorf("Expected default currency GBP, got %s", s.Currency)
	}
}

func TestSaveLoad(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test_settings.json")

	// Create settings
	s := Default()
	s.ConfigPath = testPath
	s.Locale = "en_US"
	s.Precision = 4
	s.Currency = "USD"

	// Save
	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save settings: %v", err)
	}

	// Load
	loaded, err := Load(testPath)
	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}

	// Verify
	if loaded.Locale != "en_US" {
		t.Errorf("Expected locale en_US, got %s", loaded.Locale)
	}

	if loaded.Precision != 4 {
		t.Errorf("Expected precision 4, got %d", loaded.Precision)
	}

	if loaded.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", loaded.Currency)
	}
}

func TestLoadNonExistent(t *testing.T) {
	s, err := Load("/nonexistent/path/settings.json")

	// Should not error on missing file, return defaults
	if err != nil {
		t.Errorf("Unexpected error loading non-existent file: %v", err)
	}

	if s == nil {
		t.Error("Expected default settings, got nil")
	}
}

func TestSet(t *testing.T) {
	s := Default()

	tests := []struct {
		name    string
		value   string
		check   func(*Settings) bool
		wantErr bool
	}{
		{
			name:  "precision",
			value: "4",
			check: func(s *Settings) bool { return s.Precision == 4 },
		},
		{
			name:  "locale",
			value: "en_US",
			check: func(s *Settings) bool { return s.Locale == "en_US" },
		},
		{
			name:  "currency",
			value: "USD",
			check: func(s *Settings) bool { return s.Currency == "USD" },
		},
		{
			name:    "unknown",
			value:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Set(tt.name, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set(%q, %q) error = %v, wantErr %v", tt.name, tt.value, err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil && !tt.check(s) {
				t.Errorf("Set(%q, %q) failed check", tt.name, tt.value)
			}
		})
	}
}

func TestSaveInvalidPath(t *testing.T) {
	s := Default()
	s.ConfigPath = "/invalid/nonexistent/path/settings.json"
	err := s.Save()

	if err == nil {
		t.Error("Expected error saving to invalid path")
	}
}
