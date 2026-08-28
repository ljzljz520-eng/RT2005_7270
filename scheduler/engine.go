package scheduler

import (
	"context"
	"drupal-scheduler/model"
	"drupal-scheduler/store"
	"errors"
	"sync"
	"time"
)

type Engine struct {
	db      *store.Store
	workers int
	queue   chan model.Record
	wg      sync.WaitGroup
	stop    chan struct{}
}

func New(db *store.Store, workers int) *Engine {
	if workers < 1 {
		workers = 1
	}
	return &Engine{db: db, workers: workers, queue: make(chan model.Record, 32), stop: make(chan struct{})}
}
func (e *Engine) Start() {
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.worker()
	}
}
func (e *Engine) worker() {
	defer e.wg.Done()
	for {
		select {
		case r := <-e.queue:
			e.process(context.Background(), r)
		case <-e.stop:
			return
		}
	}
}
func (e *Engine) Stop() { close(e.stop); e.wg.Wait() }
func (e *Engine) Enqueue(ctx context.Context, r model.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e.queue <- r:
		return nil
	}
}
func (e *Engine) process(ctx context.Context, r model.Record) error {
	if r.IsExpired(time.Now()) {
		r.Status = "expired"
		return e.db.SaveRecord(r)
	}
	if err := e.db.SaveEvent(model.NewEvent(r.ID+"-process", r.ID, "processing", "worker started")); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	r.Status = "published"
	return e.db.SaveRecord(r)
}
func (e *Engine) RunNow(ctx context.Context, r model.Record) error {
	if r.IsExpired(time.Now()) {
		r.Status = "expired"
		_ = e.db.SaveRecord(r)
		return errors.New("deadline exceeded")
	}
	return e.process(ctx, r)
}
