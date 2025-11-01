package timezone

import (
	"fmt"
	"strings"
	"time"
)

// Location represents a timezone location.
type Location struct {
	Name     string
	IanaName string
	Offset   int // offset in hours from UTC
}

// System manages timezone operations.
type System struct {
	locations map[string]*Location
}

// NewSystem creates a new timezone system.
func NewSystem() *System {
	s := &System{
		locations: make(map[string]*Location),
	}
	s.initLocations()
	return s
}

func (s *System) initLocations() {
	// Major cities with their IANA timezone names
	locations := []Location{
		{"London", "Europe/London", 0},
		{"Paris", "Europe/Paris", 1},
		{"Berlin", "Europe/Berlin", 1},
		{"Rome", "Europe/Rome", 1},
		{"Madrid", "Europe/Madrid", 1},
		{"New York", "America/New_York", -5},
		{"Los Angeles", "America/Los_Angeles", -8},
		{"Chicago", "America/Chicago", -6},
		{"Denver", "America/Denver", -7},
		{"Toronto", "America/Toronto", -5},
		{"Singapore", "Asia/Singapore", 8},
		{"Tokyo", "Asia/Tokyo", 9},
		{"Hong Kong", "Asia/Hong_Kong", 8},
		{"Shanghai", "Asia/Shanghai", 8},
		{"Sydney", "Australia/Sydney", 10},
		{"Melbourne", "Australia/Melbourne", 10},
		{"Dubai", "Asia/Dubai", 4},
		{"Moscow", "Europe/Moscow", 3},
		{"Mumbai", "Asia/Kolkata", 5},
		{"Bangkok", "Asia/Bangkok", 7},
	}
	
	for _, loc := range locations {
		key := strings.ToLower(loc.Name)
		s.locations[key] = &Location{
			Name:     loc.Name,
			IanaName: loc.IanaName,
			Offset:   loc.Offset,
		}
	}
}

// GetLocation retrieves a location by name.
func (s *System) GetLocation(name string) (*Location, error) {
	loc, ok := s.locations[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown timezone: %s", name)
	}
	return loc, nil
}

// GetOffset returns the time difference between two locations in hours.
func (s *System) GetOffset(from, to string) (int, error) {
	fromLoc, err := s.GetLocation(from)
	if err != nil {
		return 0, err
	}
	
	toLoc, err := s.GetLocation(to)
	if err != nil {
		return 0, err
	}
	
	return toLoc.Offset - fromLoc.Offset, nil
}

// ConvertTime converts a time from one timezone to another.
func (s *System) ConvertTime(t time.Time, from, to string) (time.Time, error) {
	offset, err := s.GetOffset(from, to)
	if err != nil {
		return time.Time{}, err
	}
	
	return t.Add(time.Duration(offset) * time.Hour), nil
}

// ListLocations returns all available timezone locations.
func (s *System) ListLocations() []string {
	var names []string
	for _, loc := range s.locations {
		names = append(names, loc.Name)
	}
	return names
}

// ParseTimeString parses a time string like "10:00", "14:30", etc.
func ParseTimeString(s string) (time.Time, error) {
	// Try various time formats
	formats := []string{
		"15:04",
		"3:04pm",
		"3:04 pm",
		"3pm",
		"15:04:05",
	}
	
	now := time.Now()
	
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			// Combine parsed time with today's date
			return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local), nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}
