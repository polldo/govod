package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/core/user"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/random"
	"github.com/polldo/govod/validate"
	"golang.org/x/crypto/bcrypt"
)

const oauthKey = "oauthstate"

// HandleLogin makes a session for the user if the passed credentials
// are correct.
func HandleLogin(db *sqlx.DB, session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		email, pass, ok := r.BasicAuth()
		if !ok {
			return weberr.BadRequest(errors.New("must provide email and password in Basic auth"))
		}

		u, err := user.FetchByEmail(ctx, db, email)
		if err != nil {
			return weberr.NotAuthorized(fmt.Errorf("fetching user by email %s: %w", email, err))
		}

		err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
		if err != nil {
			return weberr.NotAuthorized(err)
		}

		if !u.Active {
			err := fmt.Errorf("user %s is not active yet", u.Email)
			return weberr.NewError(err, err.Error(), http.StatusLocked)
		}

		session.Put(ctx, userKey, u.ID)
		session.Put(ctx, roleKey, u.Role)
		if err := session.RenewToken(ctx); err != nil {
			return fmt.Errorf("renewing token: %w", err)
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

// HandleOauthLogin starts the Oauth flow to authenticate the user.
// It returns the URL to complete the authentication on the specified external provider.
func HandleOauthLogin(session *scs.SessionManager, provs map[string]Provider) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		p := web.Param(r, "provider")
		prov, ok := provs[p]
		if !ok {
			return weberr.NotFound(fmt.Errorf("provider %s not found", p))
		}

		state, err := random.StringSecure(32)
		if err != nil {
			return fmt.Errorf("generating random secure string: %w", err)
		}

		url := prov.AuthCodeURL(state)

		session.Put(ctx, oauthKey, state)
		return web.Respond(ctx, w, url, http.StatusOK)
	}
}

// HandleOauthLogin completes the Oauth flow for the user and creates a new authenticated session.
func HandleOauthCallback(db *sqlx.DB, session *scs.SessionManager, provs map[string]Provider, redirect string) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		p := web.Param(r, "provider")
		prov, ok := provs[p]
		if !ok {
			return weberr.NotFound(fmt.Errorf("provider %s not found", p))
		}

		scstate, ok := session.Get(ctx, oauthKey).(string)
		if !ok {
			return weberr.NotAuthorized(fmt.Errorf("invalid state found in session: %+v", scstate))
		}

		state, code := r.FormValue("state"), r.FormValue("code")
		if scstate != state {
			return weberr.NotAuthorized(errors.New("wrong state"))
		}

		tok, err := prov.Exchange(ctx, code)
		if err != nil {
			return weberr.NotAuthorized(err)
		}

		rawIDTok, ok := tok.Extra("id_token").(string)
		if err != nil {
			return weberr.NotAuthorized(errors.New("id token not present"))
		}

		verifier := prov.Verifier(&oidc.Config{ClientID: prov.ClientID})
		idTok, err := verifier.Verify(ctx, rawIDTok)
		if err != nil {
			return weberr.NotAuthorized(errors.New("id token not valid"))
		}

		info := UserInfo{}
		if err := idTok.Claims(&info); err != nil {
			return fmt.Errorf("extracting info from oauth claims: %w", err)
		}

		if info.Name == "" || info.Email == "" {
			return fmt.Errorf("name or email not found in idToken claims: %+v", info)
		}

		u, err := user.FetchByEmail(ctx, db, info.Email)
		if err != nil {
			// Just fail and return on any unexpected error.
			if !errors.Is(err, database.ErrDBNotFound) {
				return fmt.Errorf("fetching user by email %s: %w", info.Email, err)
			}

			// If user not found instead, create a new user with an unguessable password.
			// The password can be recovered later on with the dedicated handler.
			now := time.Now().UTC()
			pass, err := random.StringSecure(16)
			if err != nil {
				return fmt.Errorf("generating random secure string: %w", err)
			}

			u = user.User{
				ID:           validate.GenerateID(),
				Name:         info.Name,
				Email:        info.Email,
				Role:         claims.RoleUser,
				PasswordHash: []byte(pass),
				CreatedAt:    now,
				UpdatedAt:    now,
				Active:       true,
			}

			if err := user.Create(ctx, db, u); err != nil {
				return err
			}
		}

		// Create a session for the user.
		session.Put(ctx, userKey, u.ID)
		session.Put(ctx, roleKey, u.Role)
		if err := session.RenewToken(ctx); err != nil {
			return fmt.Errorf("renewing token: %w", err)
		}

		http.Redirect(w, r, redirect, http.StatusFound)
		return nil
	}
}

// HandleLogout cancels the user's session.
func HandleLogout(session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if err := session.Destroy(ctx); err != nil {
			return fmt.Errorf("destroying session: %w", err)
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

// HandleSignup tries to register the user with the passed information.
// If activationRequired is true, users need to confirm the registration
// via email.
func HandleSignup(db *sqlx.DB, session *scs.SessionManager, activationRequired bool) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var u user.UserSignup
		if err := web.Decode(w, r, &u); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(u); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		now := time.Now().UTC()

		usr := user.User{
			ID:           validate.GenerateID(),
			Name:         u.Name,
			Email:        u.Email,
			Role:         claims.RoleUser,
			PasswordHash: hash,
			CreatedAt:    now,
			UpdatedAt:    now,
			Active:       !activationRequired,
		}

		if err := user.Create(ctx, db, usr); err != nil {
			if errors.Is(err, user.ErrUniqueEmail) {
				return weberr.NewError(err, "email already registered", http.StatusConflict)
			}
			return fmt.Errorf("creating user[%s]: %w", usr.Email, err)
		}

		if !activationRequired {
			session.Put(ctx, userKey, usr.ID)
			session.Put(ctx, roleKey, usr.Role)
			if err := session.RenewToken(ctx); err != nil {
				return fmt.Errorf("renewing token: %w", err)
			}
		}

		return web.Respond(ctx, w, usr, http.StatusCreated)
	}
}
