package task

import (
	"context"
	"log"
	"runtime/debug"
	"time"
)

func RunTask(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v, stack: %s", r, string(debug.Stack()))
			}
		}()
		fn()
	}()
}

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
