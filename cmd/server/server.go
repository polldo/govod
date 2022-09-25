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
	"github.com/polldo/govod/api"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/email"
	"github.com/sirupsen/logrus"
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

	// Construct the mux for the API calls.
	mux := api.APIMux(api.APIConfig{
		Log:     logger,
		DB:      db,
		Session: sessionManager,
		Mailer:  mail,
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
	}
	return nil
}
