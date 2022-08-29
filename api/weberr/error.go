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
// 'Status' indicates the status code of the response to be built.
type RequestError struct {
	Err    error
	Status int
}

// ErrOpt defines the type for RequestError options.
type ErrOpt func(*RequestError)

// Quiet sets the error as quiet (not very relevant to be logged).
func Quiet(q bool) ErrOpt {
	return func(err *RequestError) {
		err.Err = &quietError{error: err.Err, quiet: q}
	}
}

// Fields sets custom fields to be added to the error log.
func Fields(f map[string]interface{}) ErrOpt {
	return func(err *RequestError) {
		err.Err = &fieldsError{error: err.Err, fields: f}
	}
}

// NewError wraps a provided error with HTTP details that can be used later on
// to build and log an appropriate HTTP error response.
// This function should be used when handlers encounter expected errors.
func NewError(err error, status int, opts ...ErrOpt) error {
	e := &RequestError{Err: err, Status: status}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// NotFound returns a new `Status Not Found` request error.
func NotFound(opts ...ErrOpt) error {
	e := &RequestError{
		Err:    errors.New("the resource could not be found"),
		Status: http.StatusNotFound,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (r *RequestError) Error() string { return r.Err.Error() }

// Unwrap allows to propagate inner error behaviors.
func (e *RequestError) Unwrap() error { return e.Err }

// Response converts and returns the error in a body and status code
// to be written as response.
func (r *RequestError) Response() (body interface{}, code int) {
	return &ErrorResponse{
		Error: r.Err.Error(),
	}, r.Status
}

// quietError wraps an error adding the 'Quiet' behavior to it.
type quietError struct {
	error
	quiet bool
}

// Quiet indicates whether the error is not very relevant to be logged.
func (e *quietError) Quiet() bool { return e.quiet }

func (e *quietError) Unwrap() error { return e.error }

// fieldsError wraps an error adding the 'Fields' behavior to it.
type fieldsError struct {
	error
	fields map[string]interface{}
}

// Fields returns the fields to be logged together with the error.
func (e *fieldsError) Fields() map[string]interface{} { return e.fields }

func (e *fieldsError) Unwrap() error { return e.error }
