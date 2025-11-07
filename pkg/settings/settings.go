package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Settings holds user preferences.
type Settings struct {
	Precision    int    `json:"precision"`
	DateFormat   string `json:"date_format"`
	Currency     string `json:"currency"`
	Locale       string `json:"locale"`
	FuzzyMode    bool   `json:"fuzzy_mode"`
	Autocomplete bool   `json:"autocomplete"`
	ConfigPath   string `json:"-"`
}

// Default returns default settings.
func Default() *Settings {
	return &Settings{
		Precision:    2,
		DateFormat:   "2 Jan 2006",
		Currency:     "GBP",
		Locale:       "en_GB",
		FuzzyMode:    true,
		Autocomplete: true,
	}
}

// Load loads settings from a file.
func Load(path string) (*Settings, error) {
	s := Default()
	s.ConfigPath = path

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	s.ConfigPath = path
	return s, nil
}

// Save saves settings to a file.
func (s *Settings) Save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(s.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.ConfigPath, data, 0644)
}

// Set updates a setting by name.
func (s *Settings) Set(name, value string) error {
	switch name {
	case "precision":
		var p int
		if _, err := fmt.Sscanf(value, "%d", &p); err != nil {
			return err
		}
		s.Precision = p
	case "dateformat", "date_format":
		s.DateFormat = value
	case "currency":
		s.Currency = value
	case "locale":
		s.Locale = value
	case "fuzzy", "fuzzy_mode":
		s.FuzzyMode = value == "on" || value == "true" || value == "1"
	case "autocomplete":
		s.Autocomplete = value == "on" || value == "true" || value == "1"
	default:
		return fmt.Errorf("unknown setting: %s", name)
	}
	return nil
}
