package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/polldo/govod/api/web"

	"context"
)

const (
	// RequestIDHeader is the name of the header used to transmit the request ID.
	RequestIDHeader = "X-Request-Id"

	// DefaultRequestIDLengthLimit is the default maximum length for the request ID header value.
	DefaultRequestIDLengthLimit = 128
)

// reqIDKeyCtx is the private type used to store the request id in the context.
// It is private to avoid possible collisions with keys used by other packages.
type reqIDKeyCtx int

// reqIDKey is the context key used to store the request ID value.
const reqIDKey reqIDKeyCtx = 1

// Counter used to create new request ids.
var reqID int64

// Common prefix to all newly created request ids for this process.
var reqPrefix string

// Initialize common prefix on process startup.
// Algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L44-L50 .
func init() {
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	reqPrefix = string(b64[0:10])
}

// RequestID is a middleware that injects a request ID into the context of each request.
// Retrieve it using ctx.Value(ReqIDKey). If the incoming request has a RequestIDHeader header then
// that value is used else a random value is generated.
func RequestID() web.Middleware {
	lengthLimit := DefaultRequestIDLengthLimit
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			id := r.Header.Get(RequestIDHeader)
			if id == "" {
				id = fmt.Sprintf("%s-%d", reqPrefix, atomic.AddInt64(&reqID, 1))
			} else if lengthLimit >= 0 && len(id) > lengthLimit {
				id = id[:lengthLimit]
			}
			ctx = context.WithValue(ctx, reqIDKey, id)

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

// ContextRequestID extracts the Request ID from the context.
func ContextRequestID(ctx context.Context) (reqID string) {
	id := ctx.Value(reqIDKey)
	if id != nil {
		reqID = id.(string)
	}
	return
}
