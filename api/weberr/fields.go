package weberr

import "errors"

type fielder interface {
	Fields() map[string]interface{}
}

// Fields extracts fields to be logged together with the error, if possible.
// An error has fields if it implements the interface:
//     type fielder interface {
//          Fields() map[string]interface{}
//     }
// If the error does not implement 'Fields' behavior, it returns
// 'ok' to false and other parameters should be ignored.
func Fields(err error) (fields map[string]interface{}, ok bool) {
	var fe fielder
	if errors.As(err, &fe) {
		return fe.Fields(), true
	}
	return nil, false
}

// fieldsError wraps an error adding the 'Fields' behavior to it.
type fieldsError struct {
	error
	fields map[string]interface{}
}

func (e *fieldsError) Fields() map[string]interface{} { return e.fields }

func (e *fieldsError) Unwrap() error { return e.error }
