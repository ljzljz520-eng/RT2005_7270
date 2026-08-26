package audit

import (
	"drupal-scheduler/model"
	"errors"
	"strings"
)

func Validate(a model.Audit) error {
	if a.ID == "" || a.Actor == "" || a.Action == "" || a.Target == "" {
		return errors.New("audit fields required")
	}
	return nil
}
func Redact(a model.Audit) model.Audit {
	for k := range a.Metadata {
		if strings.Contains(strings.ToLower(k), "token") {
			a.Metadata[k] = "[redacted]"
		}
	}
	return a
}
func Describe(a model.Audit) string { return a.Actor + " " + a.Action + " " + a.Target }
