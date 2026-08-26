package audit

import (
	"drupal-scheduler/model"
	"drupal-scheduler/store"
	"fmt"
	"time"
)

type Logger struct{ db *store.Store }

func New(db *store.Store) *Logger { return &Logger{db: db} }
func (l *Logger) Log(actor, action, target string, meta map[string]string) error {
	if meta == nil {
		meta = map[string]string{}
	}
	a := model.Audit{ID: fmt.Sprintf("%d-%s", time.Now().UnixNano(), target), Actor: actor, Action: action, Target: target, At: time.Now().UTC(), Metadata: meta}
	return l.db.SaveAudit(a)
}
