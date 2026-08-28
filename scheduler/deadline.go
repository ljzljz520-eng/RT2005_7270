package scheduler

import (
	"context"
	"time"
)

func DeadlineContext(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	if deadline.IsZero() {
		return context.WithCancel(parent)
	}
	return context.WithDeadline(parent, deadline)
}
func Remaining(deadline, now time.Time) time.Duration {
	if deadline.IsZero() {
		return 0
	}
	d := time.Until(deadline)
	if now.IsZero() {
		return d
	}
	return deadline.Sub(now)
}
func ShouldAbort(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
