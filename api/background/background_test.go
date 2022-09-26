package background

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestBackground(t *testing.T) {
	log := logrus.New()
	genfn := func(sleep time.Duration, cnt chan int, quit chan struct{}) func() error {
		return func() error {
			select {
			case <-time.After(sleep):
				cnt <- 1
			case <-quit:
				time.Sleep(sleep)
			}
			return nil
		}
	}

	unit := time.Millisecond

	// To make this test working, each 'durations' value should not be too close to 'timeout'.
	// Otherwise the corresponding task could select the 'time.After' event instead of the 'quit'
	// channel event, resulting in a write to 'cnt' channel that should be closed at that point
	// - causing a panic.
	tests := []struct {
		name      string
		durations []time.Duration
		timeout   time.Duration
		completed int
	}{
		{
			name:      "One task completed in time",
			durations: []time.Duration{5 * unit},
			timeout:   10 * unit,
			completed: 1,
		},

		{
			name:      "Three tasks completed in time",
			durations: []time.Duration{2 * unit, 3 * unit, 6 * unit},
			timeout:   10 * unit,
			completed: 3,
		},

		{
			name:      "Three tasks completed in time - three task not completed",
			durations: []time.Duration{2 * unit, 3 * unit, 6 * unit, 20 * unit, 22 * unit, 24 * unit},
			timeout:   10 * unit,
			completed: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := New(log)

			cnt := make(chan int, len(tt.durations))
			quit := make(chan struct{})
			for _, d := range tt.durations {
				bg.Add(genfn(d, cnt, quit))
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			startAt := time.Now()

			err := bg.Shutdown(ctx)
			if tt.completed < len(tt.durations) && err == nil {
				t.Fatal("some tasks should not have been completed: an error was expected")
			}
			if tt.completed >= len(tt.durations) && err != nil {
				t.Fatal("all tasks should have been completed: an error was not expected")
			}

			close(quit)
			close(cnt)

			if time.Now().After(startAt.Add(tt.timeout + 1*unit)) {
				t.Fatalf("shutdown has taken more than %s", tt.timeout)
			}

			var c int
			for range cnt {
				c += 1
			}

			if c != tt.completed {
				t.Fatalf("expected %d completed tasks, got %d", tt.completed, c)
			}
		})
	}
}

func TestBackgroundPanic(t *testing.T) {
	log := logrus.New()
	bg := New(log)
	bg.Add(func() error {
		panic("now what?")
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	if err := bg.Shutdown(ctx); err != nil {
		t.Fatalf("panic should not result in an error: %v", err)
	}
}
