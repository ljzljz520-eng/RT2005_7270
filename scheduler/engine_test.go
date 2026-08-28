package scheduler

import (
	"context"
	"drupal-scheduler/model"
	"drupal-scheduler/store"
	"path/filepath"
	"testing"
	"time"
)

func TestEngineRun(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	e := New(s, 1)
	r := model.NewRecord("x", "x", "", time.Now().Add(time.Hour))
	r.Status = "processing"
	if err := e.RunNow(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetRecord("x")
	if got.Status != "published" {
		t.Fatal(got.Status)
	}
}
