package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/polldo/govod/api/web"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web/mutil"
)

// Logger writes some information about the request to the logs.
// Influenced by https://github.com/zenazn/goji/blob/master/web/middleware/logger.go
// and https://github.com/ardanlabs/service/blob/master/business/web/v1/mid/logger.go
func Logger(log logrus.FieldLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Logs the request id if it's found in context.
			if rid := ContextRequestID(ctx); rid != "" {
				log = log.WithField("req_id", rid)
			}
			log = log.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"remoteaddr": r.RemoteAddr,
			})
			log.Info("started")
			startTime := time.Now().UTC()

			// Wrap the ResponseWriter to fetch its status code later on.
			lw := mutil.WrapWriter(w)
			err := handler(ctx, lw, r)

			log = log.WithFields(logrus.Fields{
				"statuscode": lw.Status(),
				"bytes":      lw.BytesWritten(),
				"since":      time.Since(startTime).Nanoseconds(),
			})
			log.Info("completed")
			return err
		}
		return h
	}
	return m
}
