package store

import (
	"drupal-scheduler/model"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"sync"
	"time"
)

var buckets = []string{"records", "profiles", "events", "audits"}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, n := range buckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(n)); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string      { return s.path }
func encode(v any) ([]byte, error) { return json.Marshal(v) }
func (s *Store) put(bucket, key string, v any) error {
	b, err := encode(v)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), b) })
}
func (s *Store) get(bucket, key string, out any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if b == nil {
			return fmt.Errorf("%s %s not found", bucket, key)
		}
		return json.Unmarshal(b, out)
	})
}
func (s *Store) SaveRecord(r model.Record) error { return s.put("records", r.ID, r) }
func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	err := s.get("records", id, &r)
	return r, err
}
func (s *Store) SaveProfile(p model.Profile) error { return s.put("profiles", p.ID, p) }
func (s *Store) GetProfile(id string) (model.Profile, error) {
	var p model.Profile
	err := s.get("profiles", id, &p)
	return p, err
}
func (s *Store) SaveEvent(e model.Event) error {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	return s.put("events", e.ID, e)
}
func (s *Store) SaveAudit(a model.Audit) error { return s.put("audits", a.ID, a) }
func (s *Store) ListRecords() ([]model.Record, error) {
	var out []model.Record
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r model.Record
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			out = append(out, r)
			return nil
		})
	})
	return out, err
}
