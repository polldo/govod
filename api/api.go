// This file has been adapted from the excellent ardanlabs repo:
// https://github.com/ardanlabs/service ardanlabs .
package api

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/middleware"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/core/auth"
	"github.com/polldo/govod/core/user"
	"github.com/sirupsen/logrus"
)

// APIConfig contains all the mandatory dependencies required by handlers.
type APIConfig struct {
	Log     logrus.FieldLogger
	DB      *sqlx.DB
	Session *scs.SessionManager
}

// api represents our server api.
type api struct {
	*mux.Router
	mw  []web.Middleware
	log logrus.FieldLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIConfig) http.Handler {
	a := &api{
		Router: mux.NewRouter(),
		log:    cfg.Log,
	}

	// Setup the middleware common to each handler.
	a.mw = append(a.mw, auth.LoadAndSave(cfg.Session))
	a.mw = append(a.mw, middleware.RequestID())
	a.mw = append(a.mw, middleware.Logger(cfg.Log))
	a.mw = append(a.mw, middleware.Errors(cfg.Log))
	a.mw = append(a.mw, middleware.Panics())

	authen := auth.Authenticate(cfg.Session)

	// Setup the handlers.
	a.Handle(http.MethodPost, "/auth/signup", auth.HandleSignup(cfg.DB))
	a.Handle(http.MethodPost, "/auth/login", auth.HandleLogin(cfg.DB, cfg.Session))

	a.Handle(http.MethodGet, "/users/{id}", user.HandleShow(cfg.DB), authen)
	a.Handle(http.MethodPost, "/users", user.HandleCreate(cfg.DB), authen)

	return a.Router
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application router.
func (a *api) Handle(method string, path string, handler web.Handler, mw ...web.Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = web.WrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = web.WrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {

			// Some bad and unrecoverable error happened.
			a.log.WithFields(logrus.Fields{
				"req_id":  middleware.ContextRequestID(ctx),
				"message": err,
			}).Error("ERROR")
		}
	})

	a.Router.Handle(path, h).Methods(method)
}
