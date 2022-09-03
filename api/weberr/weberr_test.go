package weberr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWrap(t *testing.T) {
	// The error is not important, let's make it quiet.
	auth := func() error {
		return Wrap(errors.New("token expired"), WithQuiet(true))
	}

	// Just propagate the error.
	service := func() error {
		if err := auth(); err != nil {
			return fmt.Errorf("cannot run someFunc: %w", err)
		}
		return nil
	}

	handler := func() {
		err := service()
		if err == nil {
			return
		}

		if IsQuiet(err) {
			logrus.WithField("error", err).Info("quiet error")
		} else {
			logrus.WithField("error", err).Error("error")
			t.Error("error should have been quiet")
		}
	}
	handler()
}
