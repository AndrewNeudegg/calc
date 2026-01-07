package businessdays

import (
	"time"
)

// WeeklySchedule defines which days of the week are business days.
// Days are represented as time.Weekday (0=Sunday, 1=Monday, ..., 6=Saturday)
type WeeklySchedule struct {
	BusinessDays map[time.Weekday]bool
}

// LocaleSchedules maps locale codes to their weekly schedules
var LocaleSchedules = map[string]*WeeklySchedule{
	// Default Western schedule (Monday-Friday)
	"en_GB": NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
	"en_US": NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
	"de_DE": NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
	"fr_FR": NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
	"es_ES": NewWeeklySchedule(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
	
	// Middle Eastern schedule (Sunday-Thursday)
	"ar_SA": NewWeeklySchedule(time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday),
	"ar_AE": NewWeeklySchedule(time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday),
	
	// Israeli schedule (Sunday-Thursday)
	"he_IL": NewWeeklySchedule(time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday),
}

// NewWeeklySchedule creates a new weekly schedule with the specified business days.
func NewWeeklySchedule(days ...time.Weekday) *WeeklySchedule {
	schedule := &WeeklySchedule{
		BusinessDays: make(map[time.Weekday]bool),
	}
	for _, day := range days {
		schedule.BusinessDays[day] = true
	}
	return schedule
}

// IsBusinessDay checks if a given date is a business day according to the schedule.
func (ws *WeeklySchedule) IsBusinessDay(date time.Time) bool {
	return ws.BusinessDays[date.Weekday()]
}

// GetScheduleForLocale returns the weekly schedule for a given locale.
// If the locale is not found, it returns the default Western schedule (Monday-Friday).
func GetScheduleForLocale(locale string) *WeeklySchedule {
	if schedule, exists := LocaleSchedules[locale]; exists {
		return schedule
	}
	// Default to Western schedule
	return LocaleSchedules["en_GB"]
}

// AddBusinessDays adds a specified number of business days to a date.
// Positive values add business days forward, negative values subtract backward.
func AddBusinessDays(date time.Time, days int, schedule *WeeklySchedule) time.Time {
	if days == 0 {
		return date
	}
	
	direction := 1
	if days < 0 {
		direction = -1
		days = -days
	}
	
	current := date
	businessDaysAdded := 0
	
	for businessDaysAdded < days {
		// Move to next/previous day
		current = current.AddDate(0, 0, direction)
		
		// Check if it's a business day
		if schedule.IsBusinessDay(current) {
			businessDaysAdded++
		}
	}
	
	return current
}

// CountBusinessDays counts the number of business days between two dates (exclusive of end date).
// If start is after end, returns a negative count.
func CountBusinessDays(start, end time.Time, schedule *WeeklySchedule) int {
	// Normalize to start of day for consistent comparison
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	
	if start.Equal(end) {
		return 0
	}
	
	direction := 1
	if start.After(end) {
		// Swap and count backwards
		start, end = end, start
		direction = -1
	}
	
	count := 0
	current := start
	
	for current.Before(end) {
		if schedule.IsBusinessDay(current) {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}
	
	return count * direction
}
