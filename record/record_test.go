package record_test

import (
	"testing"

	"github.com/cpustejovsky/toggltracker/record"
)

func TestRecord(t *testing.T) {
	var r record.Record
	name := "Test Record"
	hours := 420
	minutes := 42
	t.Run("New Record", func(t *testing.T) {
		r = record.New(name, hours, minutes)
	})
	t.Run("Record Name", func(t *testing.T) {
		got := r.Name()
		if got != name {
			t.Fatalf("got %s; expected %s\n", got, name)
		}
	})
	t.Run("Record TotalMinutes", func(t *testing.T) {
		got := r.TotalMinutes()
		expect := float64(hours*60 + minutes)
		if got != expect {
			t.Fatalf("got %.1f; expectecd %.1f\n", got, expect)
		}
	})
}
