package weberr

import "errors"

type quiet interface {
	Quiet() bool
}

// IsQuiet indicates whether the error implements the 'Quiet' behavior,
// that is the interface:
//     type quiet interface {
//          Quiet() bool
//     }
//
// If the error implements the 'Quiet' behavior it should not be logged as an error.
// This is useful to deal with physiological errors - like token expirations - that
// are not interesting to log as errors (perhaps to avoid triggering any alarm) but
// should be returned in the response anyway.
//
// If the error does not implement the Quiet behavior, it returns false.
func IsQuiet(err error) bool {
	var qe quiet
	if errors.As(err, &qe) {
		return qe.Quiet()
	}
	return false
}

// quietError wraps an error adding the 'Quiet' behavior to it.
type quietError struct {
	error
	quiet bool
}

func (e *quietError) Quiet() bool { return e.quiet }

func (e *quietError) Unwrap() error { return e.error }
