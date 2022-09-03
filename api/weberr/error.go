package weberr

import (
	"errors"
	"net/http"
)

// ErrorResponse contains the error message in the following form:
// `{ "error": "some error message" }`.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RequestError is used to pass an error during the request through the
// application with web specific context.
// RequestError wraps a provided error with HTTP details that can be used later on
// to build an appropriate HTTP error response.
type RequestError struct {
	Err error
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (r *RequestError) Error() string { return r.Err.Error() }

// Unwrap allows to propagate inner error behaviors.
func (e *RequestError) Unwrap() error { return e.Err }

// NewError wraps a provided error with HTTP details that can be used later on
// to build and log an appropriate HTTP error response.
// This function should be used when handlers encounter expected errors.
func NewError(err error, status int, opts ...Opt) error {
	e := &RequestError{Err: err}

	// Add the Response behavior.
	opts = append(opts, WithResponse(
		&ErrorResponse{err.Error()},
		status,
	))

	// Return err wrapped with the provided behaviors and returns it.
	return Wrap(e, opts...)
}

// NotFound returns a new `Status Not Found` request error.
func NotFound(opts ...Opt) error {
	return NewError(
		errors.New("the resource could not be found"),
		http.StatusNotFound,
		opts...,
	)
}
