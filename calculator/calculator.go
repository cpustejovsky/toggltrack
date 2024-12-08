package calculator

import (
	"math"
	"time"
)

func dayInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isWeekDay(d time.Weekday) bool {
	return d >= time.Monday && d <= time.Friday
}

func CalculateWorkToday(gap, weekendWork float64, crunch bool) float64 {
	now := time.Now()
	var weekDays, weekEndDays float64
	year := now.Year()
	month := now.Month()
	daysInMonth := dayInMonth(month, year)
	for i := now.Day(); i <= daysInMonth; i++ {
		next := time.Date(year, month, i, 0, 0, 0, 0, time.UTC)
		if isWeekDay(next.Weekday()) {
			weekDays++
		} else {
			weekEndDays++
		}
	}
	if crunch {
		return math.Ceil(gap / (weekEndDays + weekDays))
	}
	if isWeekDay(now.Weekday()) {
		weekendWork *= weekEndDays
		return math.Ceil((gap - weekendWork) / weekDays)
	} else {
		return weekendWork
	}
}
