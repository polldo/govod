package claims

import (
	"context"
	"errors"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// Claims represents the authorization claims stored in the session.
type Claims struct {
	UserID string
	Role   string
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// claimsKey is used to store/retrieve a Claims value from a context.Context.
const claimsKey ctxKey = 1

// Set stores the claims in the context.
func Set(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// Get returns the claims from the context.
func Get(ctx context.Context) (Claims, error) {
	v, ok := ctx.Value(claimsKey).(Claims)
	if !ok {
		return Claims{}, errors.New("claim value missing from context")
	}
	return v, nil
}

func IsAdmin(ctx context.Context) bool {
	c, err := Get(ctx)
	if err != nil {
		return false
	}

	return c.Role == RoleAdmin
}

func IsUser(ctx context.Context, id string) bool {
	c, err := Get(ctx)
	if err != nil {
		return false
	}

	return c.UserID == id
}
