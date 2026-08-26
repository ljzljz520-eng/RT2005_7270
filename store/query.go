package store

import (
	"drupal-scheduler/model"
	"sort"
	"strings"
)

func FilterRecords(records []model.Record, status, query string) []model.Record {
	out := make([]model.Record, 0)
	for _, r := range records {
		if status != "" && r.Status != status {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(r.Title), strings.ToLower(query)) {
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartAt.Before(out[j].StartAt) })
	return out
}
func (s *Store) Search(status, query string) ([]model.Record, error) {
	r, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	return FilterRecords(r, status, query), nil
}
