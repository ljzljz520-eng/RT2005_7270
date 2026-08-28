package service

import (
	"context"
	"drupal-scheduler/model"
	"drupal-scheduler/store"
	"path/filepath"
	"testing"
	"time"
)

func TestBusinessChain09(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	s := New(db)
	defer s.Close()
	r := model.NewRecord("late", "expired notice", "body", time.Now().Add(time.Hour))
	if e := s.Register(context.Background(), r); e != nil {
		t.Fatal(e)
	}
	ctx, c := context.WithCancel(context.Background())
	defer c()
	c()
	if e := s.Process(ctx, r.ID); e == nil {
		t.Fatal("expected deadline")
	}
	got, _ := db.GetRecord(r.ID)
	if got.Status == "published" {
		t.Fatal("expired notice published")
	}
}
