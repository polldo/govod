package middleware

import (
	"context"
	"net/http"

	"github.com/polldo/govod/api/web"
	"github.com/sirupsen/logrus"
)

type fields interface{ Fields() map[string]interface{} }

// Fields extracts fields to be logged together with the error.
// If the error does not implement the Fields behavior, it returns
// 'ok' to false and other parameters should be ignored.
func Fields(err error) (map[string]interface{}, bool) {
	if fe, ok := err.(fields); ok {
		return fe.Fields(), true
	}
	return nil, false
}

type response interface{ Response() (interface{}, int) }

// Response returns a body and status code to use as a web response.
// If the error does not implement the Response behavior, it returns
// false as third parameter and other parameters should be ignored.
func Response(err error) (interface{}, int, bool) {
	if re, ok := err.(response); ok {
		body, code := re.Response()
		return body, code, true
	}
	return nil, 0, false
}

// Errors handles errors coming out of the call chain.
// This middleware leverages a technique of opaque errors that
// allows to customize errors with behaviors without coupling them to
// a specific type.
// In this way, it's easier to create new errors compatible with
// the behavior used here.
func Errors(log logrus.FieldLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			err := handler(ctx, w, r)
			if err == nil {
				return nil
			}

			// Prepare fields to be logged.
			fields := map[string]interface{}{
				"req_id":  ContextRequestID(ctx),
				"message": err,
			}
			if f, ok := Fields(err); ok {
				for k, v := range f {
					fields[k] = v
				}
			}

			// Log the error.
			log.WithFields(logrus.Fields(fields)).Error("ERROR")

			// Try to retrieve a response from the error.
			if body, code, ok := Response(err); ok {
				return web.Respond(ctx, w, body, code)
			}

			// Unknown error, respond with Internal Server Error.
			er := struct {
				Error string `json:"error"`
			}{
				http.StatusText(http.StatusInternalServerError),
			}
			return web.Respond(ctx, w, er, http.StatusInternalServerError)
		}
		return h
	}
	return m
}
