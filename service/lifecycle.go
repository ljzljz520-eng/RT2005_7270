package service

import (
	"context"
	"drupal-scheduler/model"
	"time"
)

func (s *Service) ExpireDue(ctx context.Context, now time.Time) (int, error) {
	rows, e := s.db.ListRecords()
	if e != nil {
		return 0, e
	}
	n := 0
	for _, r := range rows {
		if r.Status != "archived" && r.IsExpired(now) {
			r.Status = "expired"
			if e = s.db.SaveRecord(r); e != nil {
				return n, e
			}
			n++
		}
	}
	return n, nil
}
func (s *Service) Publish(ctx context.Context, id string) error {
	r, e := s.db.GetRecord(id)
	if e != nil {
		return e
	}
	if !r.CanPublish(time.Now()) {
		return model.StatusTransition(r.Status, "published")
	}
	r.Status = "published"
	return s.db.SaveRecord(r)
}
func (s *Service) Reschedule(ctx context.Context, id string, deadline time.Time) error {
	r, e := s.db.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status == "published" || r.Status == "archived" {
		return context.Canceled
	}
	r.Deadline = deadline.UTC()
	r.Status = "queued"
	return s.db.SaveRecord(r)
}
