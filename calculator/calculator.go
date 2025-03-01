package calculator

import (
	"math"
	"time"
)

type Calculator struct {
	Crunch          bool
	WeekendWorkTime float64
	IdealPercent    float64
}

func New(crunch bool, weekendWorkTime, ideal float64) Calculator {
	return Calculator{
		Crunch:          crunch,
		WeekendWorkTime: weekendWorkTime,
		IdealPercent:    ideal,
	}
}

func daysInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isWeekDay(d time.Weekday) bool {
	return d >= time.Monday && d <= time.Friday
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

func (c Calculator) CalculateWorkToday(gap float64) float64 {
	now := time.Now()
	weekDays, weekEndDays := daysRemainingInMonth(now.Day(), now)
	if c.Crunch {
		return math.Ceil(gap / (weekEndDays + weekDays))
	}
	if isWeekDay(now.Weekday()) {
		return math.Ceil((gap - c.WeekendWorkTime*weekDays) / weekDays)
	} else {
		return c.WeekendWorkTime
	}
}

func (c Calculator) CalculateIdeal(minutes float64, ath time.Time) float64 {
	current := time.Now()
	weekendATH, weekdayATH := daysRemainingInMonth(0, ath)
	weekendCurrent, weekdayCurrent := daysRemainingInMonth(current.Day(), current)
	weekDay := (minutes - weekendATH*c.WeekendWorkTime) / weekdayATH
	weekDay *= c.IdealPercent
	idealForMonth := (weekDay * weekdayCurrent) + (c.WeekendWorkTime * weekendCurrent)
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
