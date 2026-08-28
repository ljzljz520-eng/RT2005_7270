package model

import (
	"errors"
	"strings"
	"time"
)

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title required")
	}
	if r.Deadline.IsZero() {
		return errors.New("deadline required")
	}
	return nil
}
func ValidateProfile(p Profile) error {
	if p.ID == "" || p.Name == "" {
		return errors.New("profile identity required")
	}
	if !strings.Contains(p.Email, "@") {
		return errors.New("valid email required")
	}
	return nil
}
func NormalizeDeadline(t time.Time, now time.Time) time.Time {
	if t.IsZero() {
		return now.Add(24 * time.Hour)
	}
	return t.UTC()
}
func StatusTransition(from, to string) error {
	allowed := map[string][]string{"draft": {"queued"}, "queued": {"processing", "cancelled"}, "processing": {"approved", "failed", "expired"}, "approved": {"published", "expired"}, "published": {}, "failed": {"queued"}, "expired": {}}
	for _, v := range allowed[from] {
		if v == to {
			return nil
		}
	}
	return errors.New("invalid status transition")
}
