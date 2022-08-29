package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/polldo/govod/api"
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
func Run(log logrus.FieldLogger) error {
	log.Infof("starting server")
	defer log.Info("shutdown complete")

	// Construct the mux for the API calls.
	mux := api.APIMux(api.APIConfig{
		Log: log,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:    "localhost:3333",
		Handler: mux,
		// TODO: add these params to the config.
		// Addr:    cfg.Server.Address,
		// ReadTimeout:  cfg.Web.ReadTimeout,
		// WriteTimeout: cfg.Web.WriteTimeout,
		// IdleTimeout:  cfg.Web.IdleTimeout,
		// ErrorLog:     log,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	go func() {
		log.Infof("starting api router at %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infof("shutting down: signal %s", sig)

		// Wait some time to complete pending requests.
		// TODO: add this timeout to the config.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
