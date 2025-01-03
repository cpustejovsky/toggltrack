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

	"github.com/cpustejovsky/toggltrack/calculator"
	"github.com/cpustejovsky/toggltrack/flags"
	"github.com/cpustejovsky/toggltrack/record"
)

var (
	compareYear         = flags.Int("compareYear", "COMPARE_YEAR", 1971, "The month you are comparing to (1-12)")
	compareMonth        = flags.Int("compareMonth", "COMPARE_MONTH", 1, "The month you are comparing to (1-12)")
	compareMonth_hour   = flags.Int("compareMonth_hour", "COMPARE_MONTH_HOUR", 1, "Total hours for your most worked month")
	compareMonth_minute = flags.Int("compareMonth_minute", "COMPARE_MONTH_MINUTE", 0, "Total minutes mod 60 for your most worked month")
	crunch              = flags.Bool("crunch", "CRUNCH", false, "Whether you want to work without weekend breaks")
	showCompareMonth    = flags.Bool("showCompareMonth", "SHOW_COMPARE_MONTH", true, "Whether you want to show the month you're comparing to")
	ideal               = flags.Float64("ideal", "IDEAL", 0.05, "The percentage you want to grow")
	showIdeal           = flags.Bool("showIdeal", "SHOW_IDEAL", true, "Whether you want to show your ideal goal")
	weekendWork         = flags.Float64("weekendWork", "WEEKEND_WORK", 0.0, "Total minutes you want to work on Saturdays and Sundays")
)

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

func main() {
	flag.Parse()
	args := flag.Args()
	date := time.Date(*compareYear, time.Month(*compareMonth), 1, 0, 0, 0, 0, time.UTC)
	compareMonth := record.New(fmt.Sprintf("%s %d", date.Month(), date.Year()), *compareMonth_hour, *compareMonth_minute)
	i := compareMonth.TotalMinutes() * 1.20
	Ideal := record.New("Ideal", int(math.Round(i/60.0)), 00)
	if len(args) < 4 {
		log.Println("please provide hours and minutes")
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	//Parse arguments
	nums := argsToInts(args...)
	now := time.Now()
	nowName := fmt.Sprintf("%s %d", now.Month(), now.Year())
	start := record.New(nowName, (nums[0]), (nums[1]))
	current := record.New(nowName, (nums[2]), (nums[3]))
	// TODO: determine whether to keep or set conditionally
	// fmt.Printf("Current time is %d:%02d\n", now.Hour(), now.Minute())
	if *showCompareMonth {
		OutputStats(w, start, current, compareMonth)
	}
	if *showIdeal {
		OutputStats(w, start, current, Ideal)
	}
}

func OutputStats(w *tabwriter.Writer, start, current, goal record.Record) {
	//Calculate minLeft
	initialMinutes := start.TotalMinutes()
	currentMinutes := current.TotalMinutes() + initialMinutes
	athMinutes := goal.TotalMinutes()
	goalpercentage := (currentMinutes / athMinutes) * 100
	//TODO: get this aligned?
	// fmt.Fprintln(w, "\tTime\tPercentage")
	fmt.Fprintf(w, "%s\t%.0fh %.0fm\t\n",
		goal.Name(), goal.TotalMinutes()/60.0, math.Mod(goal.TotalMinutes(), 60.0))
	curMin := current.TotalMinutes() + start.TotalMinutes()
	fmt.Fprintf(w, "Work Done\t%.0fh %.0fm\t%.1f%%\t\n", (curMin / 60.0), math.Mod(curMin, 60.0), goalpercentage)
	if currentMinutes > athMinutes {
		t := currentMinutes - athMinutes
		fmt.Printf("%dhr %dm (%.1f%%) extra\n",
			int(t)/60,
			int(t)%60,
			goalpercentage-100)
	} else {
		//Calculate weekdays and weekend days
		gapMin := athMinutes - initialMinutes
		work := calculator.CalculateWorkToday(gapMin, *weekendWork, *crunch)
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

		fmt.Fprintf(w, "Work Left\t%.0fh %.0fm\t%.1f%%\t\n", (gapMin-workDone)/60.0, math.Mod((gapMin-workDone), 60.0), 100-goalpercentage)
		w.Flush()
	}
}
