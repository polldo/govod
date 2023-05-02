package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf/v3"
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api"
	"github.com/polldo/govod/api/background"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/core/auth"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/email"
	"github.com/sirupsen/logrus"
	stripecl "github.com/stripe/stripe-go/v74/client"
)

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	if err := Run(log); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// Run starts the API server.
func Run(logger *logrus.Logger) error {
	logger.Infof("starting server")
	defer logger.Info("shutdown complete")

	// Fetch and parse the configuration.
	const prefix = "GOVOD"
	var cfg config.Config
	if _, err := conf.Parse(prefix, &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	// Build a stdlib logger for the http server.
	lw := logger.Writer()
	defer lw.Close()
	errLog := log.New(lw, "", 0)

	// Open the database connection.
	db, err := database.Open(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}

	// Init the session manager.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	// Build a mailer.
	mail := email.New(cfg.Email.Address, cfg.Email.Password, cfg.Email.Host, cfg.Email.Port)

	// Init a background manager to safely spawn go-routines.
	bg := background.New(logger)

	// Build the paypal client to allow payments.
	pp, err := paypal.NewClient(
		cfg.Paypal.ClientID,
		cfg.Paypal.Secret,
		cfg.Paypal.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to build the paypal client: %w", err)
	}

	// The paypal token must be retrieved manually only the first time.
	if _, err = pp.GetAccessToken(context.TODO()); err != nil {
		return fmt.Errorf("failed to get the first paypal access token: %w", err)
	}

	// Build the stripe client to allow payments.
	strp := &stripecl.API{}
	strp.Init(cfg.Stripe.APISecret, nil)

	// Instantiate known oauth providers.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Oauth.DiscoveryTimeout)
	defer cancel()
	google := cfg.Oauth.Google
	oauthProvs, err := auth.MakeProviders(ctx, []auth.ProviderConfig{
		{Name: "google", Client: google.Client, Secret: google.Secret, URL: google.URL, RedirectURL: google.RedirectURL},
	})
	if err != nil {
		return fmt.Errorf("failed to discover oauth providers: %w", err)
	}

	// Construct the mux for the API calls.
	mux := api.APIMux(api.APIConfig{
		CorsOrigin: cfg.Cors.Origin,
		Log:        logger,
		DB:         db,
		Session:    sessionManager,
		Mailer:     mail,
		Background: bg,
		Paypal:     pp,
		Stripe:     strp,
		StripeCfg:  cfg.Stripe,
		Providers:  oauthProvs,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Handler:      mux,
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     errLog,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	go func() {
		logger.Infof("starting api router at %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infof("shutting down: signal %s", sig)

		// Wait some time to complete pending requests.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		if err := bg.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not complete all background tasks: %w", err)
		}
	}
	return nil
}
