package api

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api/background"
	"github.com/polldo/govod/api/middleware"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/core/auth"
	"github.com/polldo/govod/core/cart"
	"github.com/polldo/govod/core/course"
	"github.com/polldo/govod/core/order"
	"github.com/polldo/govod/core/token"
	"github.com/polldo/govod/core/user"
	"github.com/polldo/govod/core/video"
	"github.com/sirupsen/logrus"
	stripecl "github.com/stripe/stripe-go/v74/client"
)

// APIConfig contains all the mandatory dependencies required by handlers.
type APIConfig struct {
	CorsOrigin         string
	Log                logrus.FieldLogger
	DB                 *sqlx.DB
	Session            *scs.SessionManager
	Mailer             token.Mailer
	TokenTimeout       time.Duration
	Background         *background.Background
	Paypal             *paypal.Client
	Stripe             *stripecl.API
	StripeCfg          config.Stripe
	Providers          map[string]auth.Provider
	LoginRedirectURL   string
	ActivationRequired bool
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

	if cfg.CorsOrigin != "" {
		a.mw = append(a.mw, middleware.Cors(cfg.CorsOrigin))

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}

		// Allow any OPTIONS method when CORS is set.
		a.Handle(http.MethodOptions, "/{path:.*}", h)
	}

	authen := auth.Authenticate(cfg.Session)
	admin := auth.Admin(cfg.Session)

	// Setup the handlers.
	a.Handle(http.MethodPost, "/auth/signup", auth.HandleSignup(cfg.DB, cfg.Session, cfg.ActivationRequired))
	a.Handle(http.MethodPost, "/auth/login", auth.HandleLogin(cfg.DB, cfg.Session))
	a.Handle(http.MethodPost, "/auth/logout", auth.HandleLogout(cfg.Session))
	a.Handle(http.MethodGet, "/auth/oauth-login/{provider}", auth.HandleOauthLogin(cfg.Session, cfg.Providers))
	a.Handle(http.MethodGet, "/auth/oauth-callback/{provider}", auth.HandleOauthCallback(cfg.DB, cfg.Session, cfg.Providers, cfg.LoginRedirectURL))

	a.Handle(http.MethodPost, "/tokens", token.HandleToken(cfg.DB, cfg.Mailer, cfg.TokenTimeout, cfg.Background))
	a.Handle(http.MethodPost, "/tokens/activate", token.HandleActivation(cfg.DB, cfg.Session))
	a.Handle(http.MethodPost, "/tokens/recover", token.HandleRecovery(cfg.DB))

	a.Handle(http.MethodGet, "/users/current", user.HandleShowCurrent(cfg.DB), authen)
	a.Handle(http.MethodGet, "/users/{id}", user.HandleShow(cfg.DB), authen)
	a.Handle(http.MethodPost, "/users", user.HandleCreate(cfg.DB), authen)

	a.Handle(http.MethodGet, "/courses/owned", course.HandleListOwned(cfg.DB), authen)
	a.Handle(http.MethodGet, "/courses/{course_id}/videos", video.HandleListByCourse(cfg.DB))
	a.Handle(http.MethodGet, "/courses/{course_id}/progress", video.HandleListProgressByCourse(cfg.DB), authen)
	a.Handle(http.MethodGet, "/courses/{id}", course.HandleShow(cfg.DB))
	a.Handle(http.MethodGet, "/courses", course.HandleList(cfg.DB))
	a.Handle(http.MethodPost, "/courses", course.HandleCreate(cfg.DB), admin)
	a.Handle(http.MethodPut, "/courses/{id}", course.HandleUpdate(cfg.DB), admin)

	a.Handle(http.MethodGet, "/videos/{id}/full", video.HandleShowFull(cfg.DB), authen)
	a.Handle(http.MethodGet, "/videos/{id}/free", video.HandleShowFree(cfg.DB))
	a.Handle(http.MethodGet, "/videos/{id}", video.HandleShow(cfg.DB))
	a.Handle(http.MethodGet, "/videos", video.HandleList(cfg.DB))
	a.Handle(http.MethodPost, "/videos", video.HandleCreate(cfg.DB), admin)
	a.Handle(http.MethodPut, "/videos/{id}/progress", video.HandleUpdateProgress(cfg.DB), authen)
	a.Handle(http.MethodPut, "/videos/{id}", video.HandleUpdate(cfg.DB), admin)

	a.Handle(http.MethodGet, "/cart", cart.HandleShow(cfg.DB), authen)
	a.Handle(http.MethodDelete, "/cart", cart.HandleDelete(cfg.DB), authen)
	a.Handle(http.MethodPut, "/cart/items", cart.HandleCreateItem(cfg.DB), authen)
	a.Handle(http.MethodDelete, "/cart/items/{course_id}", cart.HandleDeleteItem(cfg.DB), authen)

	a.Handle(http.MethodPost, "/orders/paypal", order.HandlePaypalCheckout(cfg.DB, cfg.Paypal), authen)
	a.Handle(http.MethodPost, "/orders/paypal/{id}/capture", order.HandlePaypalCapture(cfg.DB, cfg.Paypal), authen)
	a.Handle(http.MethodPost, "/orders/stripe", order.HandleStripeCheckout(cfg.DB, cfg.Stripe, cfg.StripeCfg), authen)
	a.Handle(http.MethodPost, "/orders/stripe/capture", order.HandleStripeCapture(cfg.DB, cfg.StripeCfg))

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
