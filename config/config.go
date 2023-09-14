package config

import (
	"time"
)

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

type Cors struct {
	Origin string `conf:"default="`
}

type Web struct {
	Address         string        `conf:"default:0.0.0.0:8000"`
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:120s"`
}

type DB struct {
	User         string `conf:"default:postgres"`
	Password     string `conf:"default:postgres,mask"`
	Host         string `conf:"default:localhost"`
	Name         string `conf:"default:postgres"`
	MaxIdleConns int    `conf:"default:0"`
	MaxOpenConns int    `conf:"default:0"`
	DisableTLS   bool   `conf:"default:true"`
}

type Email struct {
	Host     string
	Port     string
	Address  string
	Password string
}

type Stripe struct {
	APISecret     string
	WebhookSecret string
	SuccessURL    string `conf:"default:http://mylocal.com:3000/dashboard"`
	CancelURL     string `conf:"default:http://mylocal.com:3000/cart"`
}

type Paypal struct {
	ClientID string
	Secret   string
	URL      string `conf:"default:https://api.sandbox.paypal.com"`
}

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

type Auth struct {
	ActivationRequired bool `conf:"default:false"`
}
