package reporting

import (
	"drupal-scheduler/model"
	"encoding/csv"
	"io"
	"strconv"
)

func WriteCSV(w io.Writer, rows []model.Record) error {
	c := csv.NewWriter(w)
	if e := c.Write([]string{"id", "title", "status", "deadline"}); e != nil {
		return e
	}
	for _, r := range rows {
		if e := c.Write([]string{r.ID, r.Title, r.Status, r.Deadline.UTC().Format("2006-01-02T15:04:05Z07:00")}); e != nil {
			return e
		}
	}
	c.Flush()
	return c.Error()
}
func Percent(part, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(part) * 100 / float64(total)
}
func CountStatus(rows []model.Record, status string) int {
	n := 0
	for _, r := range rows {
		if r.Status == status {
			n++
		}
	}
	return n
}
func FormatCount(n int) string { return strconv.Itoa(n) }
