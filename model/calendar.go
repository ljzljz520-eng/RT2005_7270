package model

import "time"

type Window struct {
	Open, Close time.Time
	Capacity    int
}

func (w Window) Valid() bool               { return w.Close.After(w.Open) && w.Capacity > 0 }
func (w Window) Contains(t time.Time) bool { return !t.Before(w.Open) && t.Before(w.Close) }
func (w Window) Slots(step time.Duration) []time.Time {
	var out []time.Time
	if !w.Valid() || step <= 0 {
		return out
	}
	for t := w.Open; t.Before(w.Close); t = t.Add(step) {
		out = append(out, t)
	}
	return out
}
func WeekdayName(d time.Weekday) string {
	switch d {
	case time.Monday:
		return "Monday"
	case time.Tuesday:
		return "Tuesday"
	case time.Wednesday:
		return "Wednesday"
	case time.Thursday:
		return "Thursday"
	case time.Friday:
		return "Friday"
	case time.Saturday:
		return "Saturday"
	default:
		return "Sunday"
	}
}
