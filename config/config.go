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
}

type Cors struct {
	Origin string `conf:"default="`
}

type Web struct {
	Address         string        `conf:"default:0.0.0.0:3000"`
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
	SuccessURL    string
	CancelURL     string
}

type Paypal struct {
	ClientID string
	Secret   string
	URL      string `conf:"default:https://api.sandbox.paypal.com"`
}

type Oauth struct {
	DiscoveryTimeout time.Duration `conf:"default:30s"`
	Google           struct {
		Client      string
		Secret      string
		URL         string `conf:"default:https://accounts.google.com"`
		RedirectURL string `conf:"default:http://mylocal.com:8000/auth/oauth-callback/google"`
	}
}
