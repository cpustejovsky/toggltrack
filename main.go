package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type TogglRecord struct {
	name    string
	hours   int
	minutes int
}

func NewTogglRecord(name string, hr, min int) TogglRecord {
	return TogglRecord{
		name:    name,
		hours:   hr,
		minutes: min,
	}
}

func (t TogglRecord) TotalMinutes() float64 { return float64(t.hours*60 + t.minutes) }

var (
	//TODO: Set flags for date to compare to
	ATH_hour    = flag.Int("ATH_hour", 1, "Total hours for your most worked month")
	ATH_minute  = flag.Int("ATH_minute", 0, "Total minutes mod 60 for your most worked month")
	weekendWork = flag.Float64("weekendWork", 0.0, "Total minutes you want to work on Saturdays and Sundays")
	crunch      = flag.Bool("crunch", false, "Whether you want to work without weekend breaks")
)

func main() {
	flag.Parse()
	args := flag.Args()
	date := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	ATH := NewTogglRecord(fmt.Sprintf("%s %d", date.Month(), date.Year()), *ATH_hour, *ATH_minute)
	i := ATH.TotalMinutes() * 1.20
	Ideal := NewTogglRecord("Ideal", int(math.Round(i/60.0)), 00)
	if len(args) < 4 {
		log.Println("please provide hours and minutes")
		os.Exit(1)
	}

	//Parse arguments
	nums := argsToInts(args...)
	now := time.Now()
	nowName := fmt.Sprintf("%s %d", now.Month(), now.Year())
	start := NewTogglRecord(nowName, (nums[0]), (nums[1]))
	current := NewTogglRecord(nowName, (nums[2]), (nums[3]))
	// TODO: determine whether to keep or set conditionally
	// fmt.Printf("Current time is %d:%02d\n", now.Hour(), now.Minute())
	curMin := current.TotalMinutes() + start.TotalMinutes()
	fmt.Printf("Work done this month: %dh %dm\n", int(curMin/60.0), int(curMin)%60)
	OutputStats(start, current, ATH)
	OutputStats(start, current, Ideal)
}

func OutputStats(start, current, goal TogglRecord) {
	//Calculate minLeft
	initialMinutes := start.TotalMinutes()
	currentMinutes := current.TotalMinutes() + initialMinutes
	athMinutes := goal.TotalMinutes()
	goalpercentage := (currentMinutes / athMinutes) * 100
	fmt.Println("=====================================")
	fmt.Printf("For %s (%dhr %dm)\n",
		goal.name, goal.hours, goal.minutes)
	if currentMinutes > athMinutes {
		t := currentMinutes - athMinutes
		fmt.Printf("%dhr %dm (%.1f%%) extra\n",
			int(t)/60,
			int(t)%60,
			goalpercentage-100)
	} else {
		//Calculate weekdays and weekend days
		gapMin := athMinutes - initialMinutes
		work := workCalculator(gapMin, *weekendWork)
		workDone := (currentMinutes - initialMinutes)

		var minLeft float64
		if workDone < work {
			minLeft = (work - workDone)
			fmt.Printf("Do %dhr %dm more work (%dmin)\n",
				int(minLeft)/60, int(minLeft)%60, int(minLeft))
			fmt.Printf("That's x<=%d pomodoros\n",
				int(math.Ceil(minLeft/25.0)))
		} else if workDone == work {
			fmt.Println("you've done all the work you needed to do today!")
		} else {
			minLeft = (workDone - work)
			fmt.Printf("you've done %dhr %dm of extra work!\n",
				int(minLeft)/60, int(minLeft)%60)
		}
		fmt.Println()

		fmt.Printf("Work left to reach %s:\n%dhr %dm (%.1f%%)\n", goal.name, int(gapMin-workDone)/60, int(gapMin-workDone)%60, goalpercentage)
		fmt.Println()
	}
}
func dayInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isWeekDay(d time.Weekday) bool {
	return d >= time.Monday && d <= time.Friday
}

func argsToInts(args ...string) []int {
	nums := make([]int, len(args))
	for i, arg := range args {
		x, err := strconv.Atoi(arg)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		nums[i] = x
	}
	return nums
}

func workCalculator(gap, weekendWork float64) float64 {
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
	if *crunch {
		return math.Ceil(gap / (weekEndDays + weekDays))
	}
	if isWeekDay(now.Weekday()) {
		weekendWork *= weekEndDays
		return math.Ceil((gap - weekendWork) / weekDays)
	} else {
		return weekendWork
	}
}
