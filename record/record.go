package record

type Record struct {
	name    string
	hours   int
	minutes int
}

func New(name string, hr, min int) Record {
	return Record{
		name:    name,
		hours:   hr,
		minutes: min,
	}
}

func (r Record) Name() string { return r.name }

func (r Record) TotalMinutes() float64 { return float64(r.hours*60 + r.minutes) }
