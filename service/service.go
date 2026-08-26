package service

import (
	"context"
	"drupal-scheduler/model"
	"drupal-scheduler/scheduler"
	"drupal-scheduler/store"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	db     *store.Store
	engine *scheduler.Engine
}

func New(db *store.Store) *Service {
	e := scheduler.New(db, 2)
	e.Start()
	return &Service{db: db, engine: e}
}
func (s *Service) Close() { s.engine.Stop() }
func (s *Service) Register(ctx context.Context, r model.Record) error {
	if err := model.ValidateRecord(r); err != nil {
		return err
	}
	if r.StartAt.IsZero() {
		r.StartAt = time.Now().UTC()
	}
	r.Status = "queued"
	if err := s.db.SaveRecord(r); err != nil {
		return err
	}
	return s.db.SaveEvent(model.NewEvent(r.ID+"-register", r.ID, "registered", r.Title))
}
func (s *Service) Approve(ctx context.Context, id string) error {
	r, e := s.db.GetRecord(id)
	if e != nil {
		return e
	}
	if err := model.StatusTransition(r.Status, "processing"); err != nil {
		return err
	}
	r.Status = "processing"
	if e = s.db.SaveRecord(r); e != nil {
		return e
	}
	return s.engine.Enqueue(ctx, r)
}
func (s *Service) Process(ctx context.Context, id string) error {
	r, e := s.db.GetRecord(id)
	if e != nil {
		return e
	}
	if r.IsExpired(time.Now()) {
		r.Status = "expired"
		_ = s.db.SaveRecord(r)
		return fmt.Errorf("deadline exceeded")
	}
	return s.engine.RunNow(ctx, r)
}
func (s *Service) Archive(ctx context.Context, id string) error {
	r, e := s.db.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status != "published" && r.Status != "expired" {
		return errors.New("record not archivable")
	}
	r.Status = "archived"
	return s.db.SaveRecord(r)
}
func (s *Service) Query(status, q string) ([]model.Record, error) { return s.db.Search(status, q) }
func (s *Service) CreateProfile(p model.Profile) error {
	if e := model.ValidateProfile(p); e != nil {
		return e
	}
	return s.db.SaveProfile(p)
}
