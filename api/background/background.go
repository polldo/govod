package background

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

type Background struct {
	wg  sync.WaitGroup
	log logrus.FieldLogger
}

func New(log logrus.FieldLogger) *Background {
	return &Background{
		log: log,
	}
}

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
