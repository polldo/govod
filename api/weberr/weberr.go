// Package weberr allows to add behaviors to errors.
// The idea is to decorate errors with behaviors without
// the needs of creating custom error types that directly implement them.
//
// You can decorate also custom errors.
// The advantage of adding behaviors to custom type in this way
// - rather than making them implement such behaviors directly -
// is that behaviors of wrapped errors are implicitly propagated.
package weberr

type Opt func(error) error

// Wrap allows to assign behaviors to an error.
// It leverages functional options for the
// selection of the required behaviors.
func Wrap(err error, opts ...Opt) error {
	for _, opt := range opts {
		err = opt(err)
	}
	return err
}

// WithResponse returns a functional option that
// adds the 'Response' behavior to the error.
func WithResponse(body interface{}, status int) Opt {
	return func(err error) error {
		return &responseError{error: err, body: body, status: status}
	}
}

// WithFields returns a functional option that
// adds the 'Fields' behavior to the error.
func WithFields(fields map[string]interface{}) Opt {
	return func(err error) error {
		return &fieldsError{error: err, fields: fields}
	}
}
