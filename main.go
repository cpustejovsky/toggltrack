package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type Record struct {
	name    string
	hours   int
	minutes int
}

func NewRecord(name string, hr, min int) Record {
	return Record{
		name:    name,
		hours:   hr,
		minutes: min,
	}
}

func (t Record) TotalMinutes() float64 { return float64(t.hours*60 + t.minutes) }

var (
	compareYear         = flag.Int("compareYear", EnvVarToInt(os.Getenv("COMPARE_YEAR"), 1971), "The month you are comparing to (1-12)")
	compareMonth        = flag.Int("compareMonth", EnvVarToInt(os.Getenv("COMPARE_MONTH"), 1), "The month you are comparing to (1-12)")
	compareMonth_hour   = flag.Int("compareMonth_hour", EnvVarToInt(os.Getenv("COMPARE_MONTH_HOUR"), 1), "Total hours for your most worked month")
	compareMonth_minute = flag.Int("compareMonth_minute", EnvVarToInt(os.Getenv("COMPARE_MONTH_MINUTE"), 0), "Total minutes mod 60 for your most worked month")
	crunch              = flag.Bool("crunch", EnvVarToBool(os.Getenv("CRUNCH"), false), "Whether you want to work without weekend breaks")
	showCompareMonth    = flag.Bool("showCompareMonth", EnvVarToBool(os.Getenv("SHOW_COMPARE_MONTH"), true), "Whether you want to show the month you're comparing to")
	showIdeal           = flag.Bool("showIdeal", EnvVarToBool(os.Getenv("SHOW_IDEAL"), true), "Whether you want to show your ideal goal")
	weekendWork         = flag.Float64("weekendWork", EnvVarToFloat64(os.Getenv("WEEKEND_WORK"), 0.0), "Total minutes you want to work on Saturdays and Sundays")
)

func EnvVarToFloat64(s string, def float64) float64 {
	v, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return def
	}
	return v
}

func EnvVarToInt(s string, def int) int {
	v, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return def
	}
	return int(v)
}

func EnvVarToBool(s string, def bool) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}

func main() {
	flag.Parse()
	args := flag.Args()
	date := time.Date(*compareYear, time.Month(*compareMonth), 1, 0, 0, 0, 0, time.UTC)
	compareMonth := NewRecord(fmt.Sprintf("%s %d", date.Month(), date.Year()), *compareMonth_hour, *compareMonth_minute)
	i := compareMonth.TotalMinutes() * 1.20
	Ideal := NewRecord("Ideal", int(math.Round(i/60.0)), 00)
	if len(args) < 4 {
		log.Println("please provide hours and minutes")
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	//Parse arguments
	nums := argsToInts(args...)
	now := time.Now()
	nowName := fmt.Sprintf("%s %d", now.Month(), now.Year())
	start := NewRecord(nowName, (nums[0]), (nums[1]))
	current := NewRecord(nowName, (nums[2]), (nums[3]))
	// TODO: determine whether to keep or set conditionally
	// fmt.Printf("Current time is %d:%02d\n", now.Hour(), now.Minute())
	curMin := current.TotalMinutes() + start.TotalMinutes()
	fmt.Printf("Work done this month: %dh %dm\n", int(curMin/60.0), int(curMin)%60)
	if *showCompareMonth {
		OutputStats(w, start, current, compareMonth)
	}
	if *showIdeal {
		OutputStats(w, start, current, Ideal)
	}
}

func OutputStats(w *tabwriter.Writer, start, current, goal Record) {
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

		fmt.Fprintf(w, "Work Left\t%dhr %dm \t%.1f%%\t\n", int(gapMin-workDone)/60, int(gapMin-workDone)%60, 100-goalpercentage)
		w.Flush()
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
