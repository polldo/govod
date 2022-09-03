package weberr

import "errors"

type responder interface {
	Response() (body interface{}, status int)
}

// Response extracts a web response body and a status code from the error, if possible.
//
// An error has a response if it satisfies the interface:
//    type responder interface {
//        Response() (interface{}, int)
//    }
// If the error does not have the Response behavior, this function returns
// 'ok' to false and other return parameters should be ignored.
func Response(err error) (body interface{}, status int, ok bool) {
	var re responder
	if errors.As(err, &re) {
		body, code := re.Response()
		return body, code, true
	}
	return nil, 0, false
}

// responseError wraps an error adding the 'Response' behavior to it.
type responseError struct {
	error
	body   interface{}
	status int
}

func (e *responseError) Response() (interface{}, int) {
	return e.body, e.status
}

func (e *responseError) Unwrap() error {
	return e.error
}
