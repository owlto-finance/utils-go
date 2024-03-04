package timer

import (
	"context"
	"time"
)

func PeriodicTask(ctx context.Context, task func(), waitSecond time.Duration) {
	for {
		task()
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitSecond):
		}
	}
}
