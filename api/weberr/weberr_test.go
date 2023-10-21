package weberr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWrap(t *testing.T) {
	foo := func() error {
		return Wrap(errors.New("some err"), WithFields(map[string]any{"info": "bar"}))
	}

	// Just propagate the error.
	service := func() error {
		if err := foo(); err != nil {
			return fmt.Errorf("cannot run foo: %w", err)
		}
		return nil
	}

	handler := func() {
		err := service()
		if err == nil {
			return
		}

		if ff, ok := Fields(err); ok {
			logrus.WithFields(logrus.Fields(ff)).Error(err)
		} else {
			t.Error("error should have had fields")
		}
	}
	handler()
}
