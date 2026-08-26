package service

import (
	"context"
	"drupal-scheduler/model"
	"fmt"
	"sync"
	"time"
)

type Notifier struct {
	mu   sync.Mutex
	sent []string
}

func (n *Notifier) Notify(ctx context.Context, p model.Profile, r model.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.sent = append(n.sent, fmt.Sprintf("%s:%s:%s", p.Email, r.ID, time.Now().UTC().Format(time.RFC3339)))
	return nil
}
func (n *Notifier) Sent() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]string(nil), n.sent...)
}
func BuildMessage(r model.Record) string { return fmt.Sprintf("Announcement %s: %s", r.Title, r.Body) }
