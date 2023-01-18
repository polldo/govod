package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	*oauth2.Config
	*oidc.Provider
}

type ProviderConfig struct {
	Name        string
	Client      string
	Secret      string
	URL         string
	RedirectURL string
}

func MakeProviders(ctx context.Context, cfg []ProviderConfig) (map[string]Provider, error) {
	provs := make(map[string]Provider)

	for _, c := range cfg {
		p, err := oidc.NewProvider(ctx, c.URL)
		if err != nil {
			return nil, fmt.Errorf("loading provider for [%s]: %w", c.Name, err)
		}

		provs[c.Name] = Provider{
			Provider: p,
			Config: &oauth2.Config{
				ClientID:     c.Client,
				ClientSecret: c.Secret,
				Endpoint:     p.Endpoint(),
				RedirectURL:  c.RedirectURL,
				Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			},
		}
	}

	return provs, nil
}
