package weberr

import (
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
func NewError(err error, msg string, status int, opts ...Opt) error {
	e := &RequestError{Err: err}

	// Add the Response behavior.
	opts = append(opts, WithResponse(
		&ErrorResponse{msg},
		status,
	))

	// Return err wrapped with the provided behaviors and returns it.
	return Wrap(e, opts...)
}

// NotFound returns a new `Status Not Found` request error.
func NotFound(err error, opts ...Opt) error {
	return NewError(
		err,
		"the resource could not be found",
		http.StatusNotFound,
		opts...,
	)
}

// NotAuthorized returns a new `Status Not Authorized` request error.
func NotAuthorized(err error, opts ...Opt) error {
	return NewError(
		err,
		"not authorized to access resource",
		http.StatusUnauthorized,
		opts...,
	)
}

// InternalError returns a new `Internal Server Error` request error.
func InternalError(err error, opts ...Opt) error {
	return NewError(
		err,
		"the server encountered a problem and could not process your request",
		http.StatusInternalServerError,
		opts...,
	)
}

// BadRequest returns a new `Bad Request` request error.
func BadRequest(err error, opts ...Opt) error {
	return NewError(
		err,
		"bad request",
		http.StatusBadRequest,
		opts...,
	)
}
