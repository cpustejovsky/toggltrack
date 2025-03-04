package calculator_test

import (
	"testing"
	"time"

	"github.com/cpustejovsky/toggltrack/calculator"
)

func TestDaysRemainingInMonth(t *testing.T) {
	expectWeekdays := 21.0
	expectWeekendDays := 10.0
	march := time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)
	gotWeekdays := calculator.WeekdaysRemaining(1, march)
	gotWeekendDays := calculator.WeekendDaysRemaining(1, march)
	if gotWeekendDays != expectWeekendDays {
		t.Logf("expected %v; got %v", expectWeekendDays, gotWeekendDays)
	}
	if gotWeekdays != expectWeekdays {
		t.Logf("expected %v; got %v", expectWeekdays, gotWeekdays)
	}
}

func TestCalculateIdeal(t *testing.T) {
	crunch := false
	weekend := 120.0
	ideal := 1.11
	calc := calculator.New(crunch, weekend, ideal)
	februaryMin := (86.0 * 60) + 57.0
	february := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	got := calc.CalculateIdeal(februaryMin, february)
	expect := (februaryMin - (8 * 120.0)) / 20.0
	expect *= 1.11
	expect = (expect * 21) + (10 * 120.0)
	if expect != got {
		t.Logf("expected %v; got %v", expect, got)
	}
}

func TestCalculateWorkToday(t *testing.T) {
}
