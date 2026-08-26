package scheduler

import (
	"context"
	"drupal-scheduler/model"
	"sort"
	"time"
)

type Plan struct {
	Records []model.Record
	Created time.Time
}

func BuildPlan(records []model.Record, now time.Time) Plan {
	p := Plan{Created: now}
	for _, r := range records {
		if r.Status == "queued" && !r.IsExpired(now) {
			p.Records = append(p.Records, r)
		}
	}
	sort.SliceStable(p.Records, func(i, j int) bool { return p.Records[i].Deadline.Before(p.Records[j].Deadline) })
	return p
}
func (p Plan) Next() (model.Record, bool) {
	if len(p.Records) == 0 {
		return model.Record{}, false
	}
	return p.Records[0], true
}
func (p Plan) Execute(ctx context.Context, e *Engine) error {
	for _, r := range p.Records {
		if err := e.Enqueue(ctx, r); err != nil {
			return err
		}
	}
	return nil
}
