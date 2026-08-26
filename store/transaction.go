package store

import (
	"context"
	"drupal-scheduler/model"
)

type Batch struct {
	s       *Store
	records []model.Record
}

func (s *Store) BeginBatch() *Batch { return &Batch{s: s} }
func (b *Batch) Add(r model.Record) { b.records = append(b.records, r) }
func (b *Batch) Commit(ctx context.Context) error {
	for _, r := range b.records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := b.s.SaveRecord(r); err != nil {
			return err
		}
	}
	return nil
}
func (b *Batch) Size() int { return len(b.records) }
