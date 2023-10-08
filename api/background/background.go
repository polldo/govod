package background

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

// Background is a container of concurrently executed tasks.
// It handles tasks's errors by logging them.
// It also recovers in case of panics.
type Background struct {
	wg  sync.WaitGroup
	log logrus.FieldLogger
}

// New constructs and returns a new Background.
func New(log logrus.FieldLogger) *Background {
	return &Background{
		log: log,
	}
}

// Add inserts a new task, which will be executed in background.
func (bg *Background) Add(f func() error) {
	bg.wg.Add(1)

	go func() {
		defer bg.wg.Done()

		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				bg.log.WithField("message", err).Error("PANIC")
			}
		}()

		if err := f(); err != nil {
			bg.log.WithField("message", err).Error("ERROR")
		}
	}()
}

// Shutdown waits for tasks to complete. If the passed context
// expires then it returns an error indicating that some task
// didn't terminate in time.
func (bg *Background) Shutdown(ctx context.Context) error {
	quit := make(chan struct{})
	go func() {
		bg.wg.Wait()
		quit <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-quit:
		return nil
	}
}
