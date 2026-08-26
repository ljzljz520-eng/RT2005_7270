package reporting

import (
	"drupal-scheduler/model"
	"strings"
	"time"
)

type Summary struct {
	Total, Queued, Published, Expired int
	GeneratedAt                       time.Time
}

func Summarize(records []model.Record) Summary {
	s := Summary{GeneratedAt: time.Now().UTC()}
	for _, r := range records {
		s.Total++
		switch r.Status {
		case "queued":
			s.Queued++
		case "published":
			s.Published++
		case "expired":
			s.Expired++
		}
	}
	return s
}
func Render(s Summary) string {
	return strings.Join([]string{"total=" + itoa(s.Total), "queued=" + itoa(s.Queued), "published=" + itoa(s.Published), "expired=" + itoa(s.Expired)}, " ")
}
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	digits := ""
	for v > 0 {
		digits = string(rune('0'+v%10)) + digits
		v /= 10
	}
	return digits
}
