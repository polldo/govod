package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// UserInfo includes the information of a user
// to be retrieved from external providers.
type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProviderConfig contains the information needed
// to setup an Oauth provider.
type ProviderConfig struct {
	Name        string
	Client      string
	Secret      string
	URL         string
	RedirectURL string
}

// Provider wraps an external Oauth provider.
type Provider struct {
	*oauth2.Config
	*oidc.Provider
}

// MakeProviders builds supported Oauth providers.
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
