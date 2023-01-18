package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// 		conf := &oauth2.Config{
// 			RedirectURL:  "http://mylocal.com:8000/auth/oauth-callback/google",
// 			ClientID:     "785050419234-c7ao87rji0crqpkfsu4sr8m77asp4umu.apps.googleusercontent.com",
// 			ClientSecret: "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",
// 			Endpoint:     google.Endpoint,
// 			Scopes: []string{
// 				"https://www.googleapis.com/auth/userinfo.profile",
// 				"https://www.googleapis.com/auth/userinfo.email",
// 			},
// 		}

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ProviderConfig struct {
	Name        string
	Client      string
	Secret      string
	URL         string
	RedirectURL string
}

type Provider struct {
	*oauth2.Config
	*oidc.Provider
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
				Scopes:       []string{oidc.ScopeOpenID},
				// Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			},
		}
	}

	return provs, nil
}
