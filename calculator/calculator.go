package calculator

import (
	"math"
	"time"
)

func daysInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isWeekDay(d time.Weekday) bool {
	return d >= time.Monday && d <= time.Friday
}

func CalculateWorkToday(gap, weekendWork float64, crunch bool) float64 {
	now := time.Now()
	weekDays, weekEndDays := daysRemainingInMonth(now.Day(), now)
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
func daysRemainingInMonth(start int, date time.Time) (float64, float64) {
	var weekDays, weekEndDays float64
	end := daysInMonth(date)
	for i := start + 1; i <= end; i++ {
		next := time.Date(date.Year(), date.Month(), i, 0, 0, 0, 0, time.UTC)
		if isWeekDay(next.Weekday()) {
			weekDays++
		} else {
			weekEndDays++
		}
	}
	return weekDays, weekEndDays
}
func CalculateIdeal(minutes, weekendWork, ideal float64, ath, current time.Time) float64 {
	weekendATH, weekdayATH := daysRemainingInMonth(0, ath)
	weekendCurrent, weekdayCurrent := daysRemainingInMonth(current.Day(), current)
	weekDay := (minutes - weekendATH*weekendWork) / weekdayATH
	weekDay *= ideal
	idealForMonth := (weekDay * weekdayCurrent) + (weekendWork * weekendCurrent)
	return idealForMonth
}

func CalculateWorkWeekDay(gap, weekendWork float64, crunch bool) float64 {
	now := time.Now()
	weekDays, weekEndDays := daysRemainingInMonth(now.Day(), now)
	if crunch {
		return math.Ceil(gap / (weekEndDays + weekDays))
	}
	weekendWork *= weekEndDays
	return math.Ceil((gap - weekendWork) / weekDays)

}
