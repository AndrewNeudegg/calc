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
	// Major and regional world cities with IANA time zones
	var locations = []Location{
		// Europe
		{"London", "Europe/London", 0},
		{"Birmingham", "Europe/London", 0},
		{"Manchester", "Europe/London", 0},
		{"Dublin", "Europe/Dublin", 0},
		{"Edinburgh", "Europe/London", 0},
		{"Paris", "Europe/Paris", 1},
		{"Lyon", "Europe/Paris", 1},
		{"Marseille", "Europe/Paris", 1},
		{"Berlin", "Europe/Berlin", 1},
		{"Hamburg", "Europe/Berlin", 1},
		{"Munich", "Europe/Berlin", 1},
		{"Frankfurt", "Europe/Berlin", 1},
		{"Amsterdam", "Europe/Amsterdam", 1},
		{"Brussels", "Europe/Brussels", 1},
		{"Zurich", "Europe/Zurich", 1},
		{"Geneva", "Europe/Zurich", 1},
		{"Madrid", "Europe/Madrid", 1},
		{"Barcelona", "Europe/Madrid", 1},
		{"Rome", "Europe/Rome", 1},
		{"Milan", "Europe/Rome", 1},
		{"Athens", "Europe/Athens", 2},
		{"Istanbul", "Europe/Istanbul", 3},
		{"Oslo", "Europe/Oslo", 1},
		{"Stockholm", "Europe/Stockholm", 1},
		{"Copenhagen", "Europe/Copenhagen", 1},
		{"Warsaw", "Europe/Warsaw", 1},
		{"Prague", "Europe/Prague", 1},
		{"Vienna", "Europe/Vienna", 1},
		{"Budapest", "Europe/Budapest", 1},
		{"Lisbon", "Europe/Lisbon", 0},
		{"Reykjavik", "Atlantic/Reykjavik", 0},
		{"Moscow", "Europe/Moscow", 3},
		{"St Petersburg", "Europe/Moscow", 3},
		{"Helsinki", "Europe/Helsinki", 2},
		{"Tallinn", "Europe/Tallinn", 2},
		{"Riga", "Europe/Riga", 2},
		{"Vilnius", "Europe/Vilnius", 2},
		{"Bucharest", "Europe/Bucharest", 2},
		{"Sofia", "Europe/Sofia", 2},
		{"Belgrade", "Europe/Belgrade", 1},
		{"Zagreb", "Europe/Zagreb", 1},
		{"Sarajevo", "Europe/Sarajevo", 1},
		{"Ljubljana", "Europe/Ljubljana", 1},
		{"Kiev", "Europe/Kyiv", 2},
		{"Minsk", "Europe/Minsk", 3},

		// Middle East
		{"Istanbul", "Europe/Istanbul", 3},
		{"Ankara", "Europe/Istanbul", 3},
		{"Jerusalem", "Asia/Jerusalem", 2},
		{"Tel Aviv", "Asia/Jerusalem", 2},
		{"Riyadh", "Asia/Riyadh", 3},
		{"Jeddah", "Asia/Riyadh", 3},
		{"Dubai", "Asia/Dubai", 4},
		{"Abu Dhabi", "Asia/Dubai", 4},
		{"Doha", "Asia/Qatar", 3},
		{"Kuwait City", "Asia/Kuwait", 3},
		{"Manama", "Asia/Bahrain", 3},
		{"Muscat", "Asia/Muscat", 4},
		{"Tehran", "Asia/Tehran", 3},
		{"Baghdad", "Asia/Baghdad", 3},
		{"Beirut", "Asia/Beirut", 2},
		{"Amman", "Asia/Amman", 2},

		// Africa
		{"Cairo", "Africa/Cairo", 2},
		{"Alexandria", "Africa/Cairo", 2},
		{"Casablanca", "Africa/Casablanca", 0},
		{"Marrakesh", "Africa/Casablanca", 0},
		{"Tunis", "Africa/Tunis", 1},
		{"Algiers", "Africa/Algiers", 1},
		{"Nairobi", "Africa/Nairobi", 3},
		{"Addis Ababa", "Africa/Addis_Ababa", 3},
		{"Khartoum", "Africa/Khartoum", 2},
		{"Johannesburg", "Africa/Johannesburg", 2},
		{"Cape Town", "Africa/Johannesburg", 2},
		{"Durban", "Africa/Johannesburg", 2},
		{"Lagos", "Africa/Lagos", 1},
		{"Abuja", "Africa/Lagos", 1},
		{"Accra", "Africa/Accra", 0},
		{"Dakar", "Africa/Dakar", 0},
		{"Kampala", "Africa/Kampala", 3},
		{"Dar es Salaam", "Africa/Dar_es_Salaam", 3},
		{"Kigali", "Africa/Kigali", 2},
		{"Luanda", "Africa/Luanda", 1},

		// Asia
		{"Dubai", "Asia/Dubai", 4},
		{"Doha", "Asia/Qatar", 3},
		{"Tel Aviv", "Asia/Jerusalem", 2},
		{"Tehran", "Asia/Tehran", 3},
		{"Karachi", "Asia/Karachi", 5},
		{"Lahore", "Asia/Karachi", 5},
		{"Mumbai", "Asia/Kolkata", 5},
		{"Delhi", "Asia/Kolkata", 5},
		{"Bangalore", "Asia/Kolkata", 5},
		{"Chennai", "Asia/Kolkata", 5},
		{"Kolkata", "Asia/Kolkata", 5},
		{"Kathmandu", "Asia/Kathmandu", 5},
		{"Dhaka", "Asia/Dhaka", 6},
		{"Colombo", "Asia/Colombo", 5},
		{"Yangon", "Asia/Yangon", 6},
		{"Bangkok", "Asia/Bangkok", 7},
		{"Chiang Mai", "Asia/Bangkok", 7},
		{"Hanoi", "Asia/Bangkok", 7},
		{"Ho Chi Minh City", "Asia/Ho_Chi_Minh", 7},
		{"Phnom Penh", "Asia/Phnom_Penh", 7},
		{"Kuala Lumpur", "Asia/Kuala_Lumpur", 8},
		{"Penang", "Asia/Kuala_Lumpur", 8},
		{"Singapore", "Asia/Singapore", 8},
		{"Jakarta", "Asia/Jakarta", 7},
		{"Surabaya", "Asia/Jakarta", 7},
		{"Bali", "Asia/Makassar", 8},
		{"Manila", "Asia/Manila", 8},
		{"Hong Kong", "Asia/Hong_Kong", 8},
		{"Macau", "Asia/Macau", 8},
		{"Shanghai", "Asia/Shanghai", 8},
		{"Beijing", "Asia/Shanghai", 8},
		{"Shenzhen", "Asia/Shanghai", 8},
		{"Guangzhou", "Asia/Shanghai", 8},
		{"Seoul", "Asia/Seoul", 9},
		{"Busan", "Asia/Seoul", 9},
		{"Tokyo", "Asia/Tokyo", 9},
		{"Osaka", "Asia/Tokyo", 9},
		{"Sapporo", "Asia/Tokyo", 9},
		{"Fukuoka", "Asia/Tokyo", 9},
		{"Taipei", "Asia/Taipei", 8},
		{"Ulaanbaatar", "Asia/Ulaanbaatar", 8},
		{"Almaty", "Asia/Almaty", 6},
		{"Tashkent", "Asia/Tashkent", 5},
		{"Bishkek", "Asia/Bishkek", 6},
		{"Astana", "Asia/Almaty", 6},

		// Oceania
		{"Sydney", "Australia/Sydney", 10},
		{"Melbourne", "Australia/Melbourne", 10},
		{"Canberra", "Australia/Sydney", 10},
		{"Brisbane", "Australia/Brisbane", 10},
		{"Perth", "Australia/Perth", 8},
		{"Adelaide", "Australia/Adelaide", 9},
		{"Hobart", "Australia/Hobart", 10},
		{"Darwin", "Australia/Darwin", 9},
		{"Auckland", "Pacific/Auckland", 12},
		{"Wellington", "Pacific/Auckland", 12},
		{"Christchurch", "Pacific/Auckland", 12},
		{"Suva", "Pacific/Fiji", 12},
		{"Port Moresby", "Pacific/Port_Moresby", 10},
		{"Honolulu", "Pacific/Honolulu", -10},

		// North America
		{"New York", "America/New_York", -5},
		{"Boston", "America/New_York", -5},
		{"Philadelphia", "America/New_York", -5},
		{"Washington DC", "America/New_York", -5},
		{"Miami", "America/New_York", -5},
		{"Atlanta", "America/New_York", -5},
		{"Chicago", "America/Chicago", -6},
		{"Dallas", "America/Chicago", -6},
		{"Houston", "America/Chicago", -6},
		{"Minneapolis", "America/Chicago", -6},
		{"Denver", "America/Denver", -7},
		{"Salt Lake City", "America/Denver", -7},
		{"Phoenix", "America/Phoenix", -7},
		{"Los Angeles", "America/Los_Angeles", -8},
		{"San Francisco", "America/Los_Angeles", -8},
		{"Seattle", "America/Los_Angeles", -8},
		{"Portland (US)", "America/Los_Angeles", -8},
		{"Las Vegas", "America/Los_Angeles", -8},
		{"Vancouver", "America/Vancouver", -8},
		{"Calgary", "America/Edmonton", -7},
		{"Toronto", "America/Toronto", -5},
		{"Ottawa", "America/Toronto", -5},
		{"Montreal", "America/Toronto", -5},
		{"Quebec City", "America/Toronto", -5},
		{"Winnipeg", "America/Winnipeg", -6},
		{"Halifax", "America/Halifax", -4},
		{"St John’s", "America/St_Johns", -3},
		{"Mexico City", "America/Mexico_City", -6},
		{"Guadalajara", "America/Mexico_City", -6},
		{"Monterrey", "America/Monterrey", -6},

		// Central & South America
		{"Bogota", "America/Bogota", -5},
		{"Medellin", "America/Bogota", -5},
		{"Quito", "America/Guayaquil", -5},
		{"Lima", "America/Lima", -5},
		{"Caracas", "America/Caracas", -4},
		{"Santiago", "America/Santiago", -4},
		{"Valparaiso", "America/Santiago", -4},
		{"Buenos Aires", "America/Argentina/Buenos_Aires", -3},
		{"Cordoba", "America/Argentina/Cordoba", -3},
		{"Montevideo", "America/Montevideo", -3},
		{"Asuncion", "America/Asuncion", -4},
		{"Sao Paulo", "America/Sao_Paulo", -3},
		{"Rio de Janeiro", "America/Sao_Paulo", -3},
		{"Brasilia", "America/Sao_Paulo", -3},
		{"Salvador", "America/Bahia", -3},
		{"Recife", "America/Recife", -3},
		{"La Paz", "America/La_Paz", -4},
		{"Bogota", "America/Bogota", -5},
		{"Panama City", "America/Panama", -5},
		{"San Jose (CR)", "America/Costa_Rica", -6},
		{"Havana", "America/Havana", -5},
		{"Kingston", "America/Jamaica", -5},
		{"Santo Domingo", "America/Santo_Domingo", -4},

		// Smaller / Pacific Islands
		{"Pago Pago", "Pacific/Pago_Pago", -11},
		{"Tahiti", "Pacific/Tahiti", -10},
		{"Noumea", "Pacific/Noumea", 11},
		{"Guam", "Pacific/Guam", 10},
		{"Saipan", "Pacific/Saipan", 10},
		{"Apia", "Pacific/Apia", 13},
		{"Nuku'alofa", "Pacific/Tongatapu", 13},
		{"Papeete", "Pacific/Tahiti", -10},
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
