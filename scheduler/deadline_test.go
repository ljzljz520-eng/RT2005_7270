package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestDeadlineContext(t *testing.T) {
	ctx, c := DeadlineContext(context.Background(), time.Now().Add(time.Millisecond))
	defer c()
	time.Sleep(3 * time.Millisecond)
	if !ShouldAbort(ctx) {
		t.Fatal("not aborted")
	}
}
