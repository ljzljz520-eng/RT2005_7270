package service

import (
	"context"
	"drupal-scheduler/model"
	"drupal-scheduler/store"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkflowOne(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	s := New(db)
	defer s.Close()
	r := model.NewRecord("w1", "notice", "body", time.Now().Add(time.Hour))
	if e := s.Register(context.Background(), r); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Query("queued", ""); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowTwo(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	s := New(db)
	defer s.Close()
	r := model.NewRecord("w2", "notice", "body", time.Now().Add(time.Hour))
	_ = s.Register(context.Background(), r)
	if e := s.Approve(context.Background(), r.ID); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowThree(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	s := New(db)
	defer s.Close()
	r := model.NewRecord("w3", "notice", "body", time.Now().Add(time.Hour))
	_ = s.Register(context.Background(), r)
	if e := s.Reschedule(context.Background(), r.ID, time.Now().Add(2*time.Hour)); e != nil {
		t.Fatal(e)
	}
}
