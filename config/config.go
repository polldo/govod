package config

import (
	"time"
)

// Config contains all the config parameters useful
// to setup the whole server components.
type Config struct {
	Cors   Cors
	Web    Web
	DB     DB
	Email  Email
	Paypal Paypal
	Stripe Stripe
	Oauth  Oauth
	Auth   Auth
}

// Cors includes parameters for CORS setup.
type Cors struct {
	Origin string `conf:"default="`
}

// Web contains all the parameters related to the http listener.
type Web struct {
	Address         string        `conf:"default:0.0.0.0:8000"`
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:120s"`
}

// DB contains the details of the PostgreSQL to use.
type DB struct {
	User         string `conf:"default:postgres"`
	Password     string `conf:"default:postgres,mask"`
	Host         string `conf:"default:localhost"`
	Name         string `conf:"default:postgres"`
	MaxIdleConns int    `conf:"default:0"`
	MaxOpenConns int    `conf:"default:0"`
	DisableTLS   bool   `conf:"default:true"`
}

// Email includes both SMTP information and more business related
// details which regards the sending of emails.
type Email struct {
	Host          string
	Port          string
	Address       string
	Password      string
	RecoveryURL   string        `conf:"default:http://mylocal.com:3000/password/confirm?token="`
	ActivationURL string        `conf:"default:http://mylocal.com:3000/activate/confirm?token="`
	TokenTimeout  time.Duration `conf:"default:10s"`
}

// Stripe contains parameters to setup the Stripe dependency.
type Stripe struct {
	APISecret     string
	WebhookSecret string
	SuccessURL    string `conf:"default:http://mylocal.com:3000/dashboard"`
	CancelURL     string `conf:"default:http://mylocal.com:3000/cart"`
}

// Paypal contains parameters to setup the Paypal dependency.
type Paypal struct {
	ClientID string
	Secret   string
	URL      string `conf:"default:https://api.sandbox.paypal.com"`
}

// Oauth includes all details needed to setup Oauth authentication.
type Oauth struct {
	DiscoveryTimeout time.Duration `conf:"default:30s"`
	LoginRedirectURL string        `conf:"default:http://mylocal.com:3000/dashboard"`
	Google           struct {
		Client      string
		Secret      string
		URL         string `conf:"default:https://accounts.google.com"`
		RedirectURL string `conf:"default:http://mylocal.com:8000/auth/oauth-callback/google"`
	}
}

// Auth configures authentication options.
type Auth struct {
	ActivationRequired bool `conf:"default:false"`
}
