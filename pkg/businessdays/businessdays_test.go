package businessdays

import (
	"testing"
	"time"
)

func TestNewWeeklySchedule(t *testing.T) {
	schedule := NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday)
	
	// Test business days
	if !schedule.BusinessDays[time.Monday] {
		t.Error("Monday should be a business day")
	}
	if !schedule.BusinessDays[time.Friday] {
		t.Error("Friday should be a business day")
	}
	
	// Test non-business days
	if schedule.BusinessDays[time.Saturday] {
		t.Error("Saturday should not be a business day")
	}
	if schedule.BusinessDays[time.Sunday] {
		t.Error("Sunday should not be a business day")
	}
}

func TestIsBusinessDay(t *testing.T) {
	schedule := NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday)
	
	tests := []struct {
		date       time.Time
		isBusiness bool
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), true},  // Monday
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), true},  // Tuesday
		{time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), true},  // Friday
		{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), false}, // Saturday
		{time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), false}, // Sunday
	}
	
	for _, tt := range tests {
		result := schedule.IsBusinessDay(tt.date)
		if result != tt.isBusiness {
			t.Errorf("IsBusinessDay(%s) = %v, want %v", tt.date.Weekday(), result, tt.isBusiness)
		}
	}
}

func TestGetScheduleForLocale(t *testing.T) {
	tests := []struct {
		locale         string
		expectedMon    bool
		expectedFri    bool
		expectedSat    bool
		expectedSun    bool
	}{
		{"en_GB", true, true, false, false},  // Mon-Fri
		{"en_US", true, true, false, false},  // Mon-Fri
		{"ar_SA", true, false, false, true},  // Sun-Thu
		{"ar_AE", true, false, false, true},  // Sun-Thu
		{"he_IL", true, false, false, true},  // Sun-Thu
		{"unknown", true, true, false, false}, // Default to Mon-Fri
	}
	
	for _, tt := range tests {
		schedule := GetScheduleForLocale(tt.locale)
		if schedule.BusinessDays[time.Monday] != tt.expectedMon {
			t.Errorf("Locale %s: Monday business day = %v, want %v", tt.locale, schedule.BusinessDays[time.Monday], tt.expectedMon)
		}
		if schedule.BusinessDays[time.Friday] != tt.expectedFri {
			t.Errorf("Locale %s: Friday business day = %v, want %v", tt.locale, schedule.BusinessDays[time.Friday], tt.expectedFri)
		}
		if schedule.BusinessDays[time.Saturday] != tt.expectedSat {
			t.Errorf("Locale %s: Saturday business day = %v, want %v", tt.locale, schedule.BusinessDays[time.Saturday], tt.expectedSat)
		}
		if schedule.BusinessDays[time.Sunday] != tt.expectedSun {
			t.Errorf("Locale %s: Sunday business day = %v, want %v", tt.locale, schedule.BusinessDays[time.Sunday], tt.expectedSun)
		}
	}
}

func TestAddBusinessDays(t *testing.T) {
	schedule := NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday)
	
	tests := []struct {
		name     string
		start    time.Time
		days     int
		expected time.Time
	}{
		{
			name:     "add 1 business day from Monday",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
			days:     1,
			expected: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), // Tuesday
		},
		{
			name:     "add 5 business days from Monday",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
			days:     5,
			expected: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Next Monday
		},
		{
			name:     "add 1 business day from Friday",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Friday
			days:     1,
			expected: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday (skips weekend)
		},
		{
			name:     "add 3 business days from Friday",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Friday
			days:     3,
			expected: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), // Wednesday (skips weekend)
		},
		{
			name:     "subtract 1 business day from Tuesday",
			start:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), // Tuesday
			days:     -1,
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			name:     "subtract 1 business day from Monday",
			start:    time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			days:     -1,
			expected: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Previous Friday (skips weekend)
		},
		{
			name:     "add 0 business days",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Monday
			days:     0,
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Same day
		},
		{
			name:     "add 10 business days from Wednesday",
			start:    time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), // Wednesday
			days:     10,
			expected: time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC), // Wednesday (2 weeks later)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddBusinessDays(tt.start, tt.days, schedule)
			// Compare dates without time
			resultDate := time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, time.UTC)
			if !resultDate.Equal(tt.expected) {
				t.Errorf("AddBusinessDays(%s, %d) = %s, want %s",
					tt.start.Format("2006-01-02"), tt.days,
					resultDate.Format("2006-01-02"), tt.expected.Format("2006-01-02"))
			}
		})
	}
}

func TestAddBusinessDaysWithMiddleEasternSchedule(t *testing.T) {
	// Sunday-Thursday schedule
	schedule := NewWeeklySchedule(time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday)
	
	tests := []struct {
		name     string
		start    time.Time
		days     int
		expected time.Time
	}{
		{
			name:     "add 1 business day from Sunday",
			start:    time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
			days:     1,
			expected: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			name:     "add 1 business day from Thursday",
			start:    time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), // Thursday
			days:     1,
			expected: time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday (skips Fri-Sat weekend)
		},
		{
			name:     "add 5 business days from Sunday",
			start:    time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
			days:     5,
			expected: time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC), // Next Sunday (skips Fri-Sat)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddBusinessDays(tt.start, tt.days, schedule)
			resultDate := time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, time.UTC)
			if !resultDate.Equal(tt.expected) {
				t.Errorf("AddBusinessDays(%s, %d) = %s, want %s",
					tt.start.Format("2006-01-02"), tt.days,
					resultDate.Format("2006-01-02"), tt.expected.Format("2006-01-02"))
			}
		})
	}
}

func TestCountBusinessDays(t *testing.T) {
	schedule := NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday)
	
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "same day",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Monday to Tuesday",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Monday to Friday (same week)",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: 4,
		},
		{
			name:     "Friday to Monday (over weekend)",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			expected: 1, // Only Monday counts (Friday to Monday is 1 business day)
		},
		{
			name:     "Monday to next Monday (full week)",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			expected: 5, // 5 business days (Mon-Fri)
		},
		{
			name:     "backwards (Tuesday to Monday)",
			start:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: -1,
		},
		{
			name:     "backwards (next Monday to previous Friday over weekend)",
			start:    time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: -1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountBusinessDays(tt.start, tt.end, schedule)
			if result != tt.expected {
				t.Errorf("CountBusinessDays(%s, %s) = %d, want %d",
					tt.start.Format("2006-01-02"), tt.end.Format("2006-01-02"),
					result, tt.expected)
			}
		})
	}
}
