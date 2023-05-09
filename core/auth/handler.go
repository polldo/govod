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

func HandleLogin(db *sqlx.DB, session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		email, pass, ok := r.BasicAuth()
		if !ok {
			err := errors.New("must provide email and password in Basic auth")
			return weberr.NewError(err, err.Error(), http.StatusUnauthorized)
		}

		u, err := user.FetchByEmail(ctx, db, email)
		if err != nil {
			err := fmt.Errorf("fetching user by email %s: %w", email, err)
			return weberr.NotAuthorized(err)
		}

		if !u.Active {
			err := fmt.Errorf("user %s is not active yet", u.Email)
			return weberr.NewError(err, err.Error(), http.StatusForbidden)
		}

		err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pass))
		if err != nil {
			return weberr.NotAuthorized(err)
		}

		// TODO: Save the entire user struct in the session
		// or just some info?
		session.Put(ctx, userKey, u.ID)
		session.Put(ctx, roleKey, u.Role)
		if err := session.RenewToken(ctx); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

func HandleOauthLogin(session *scs.SessionManager, provs map[string]Provider) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		p := web.Param(r, "provider")
		prov, ok := provs[p]
		if !ok {
			return weberr.NotFound(fmt.Errorf("provider %s not found", p))
		}

		state, err := random.StringSecure(32)
		if err != nil {
			return weberr.InternalError(err)
		}

		url := prov.AuthCodeURL(state)

		session.Put(ctx, oauthKey, state)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return nil
	}
}

func HandleOauthCallback(db *sqlx.DB, session *scs.SessionManager, provs map[string]Provider) web.Handler {
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
			return weberr.InternalError(fmt.Errorf("extracting info from oauth claims: %w", err))
		}

		if info.Name == "" || info.Email == "" {
			return weberr.InternalError(fmt.Errorf("name or email not found in idToken claims: %+v", info))
		}

		u, err := user.FetchByEmail(ctx, db, info.Email)
		if err != nil {
			// Just fail and return on any unexpected error.
			if !errors.Is(err, database.ErrDBNotFound) {
				err := fmt.Errorf("fetching user by email %s: %w", info.Email, err)
				return weberr.InternalError(err)
			}

			// If user not found instead, create a new user with an unguessable password.
			// The password can be recovered later on with the dedicated handler.
			now := time.Now().UTC()
			pass, err := random.StringSecure(16)
			if err != nil {
				return weberr.InternalError(err)
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
				return weberr.InternalError(err)
			}
		}

		// Create a session for the user.
		session.Put(ctx, userKey, u.ID)
		session.Put(ctx, roleKey, u.Role)
		if err := session.RenewToken(ctx); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

func HandleLogout(session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		session.Remove(ctx, userKey)
		session.Remove(ctx, roleKey)

		if err := session.RenewToken(ctx); err != nil {
			return err
		}

		return nil
	}
}

func HandleSignup(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var u user.UserSignup
		if err := web.Decode(r, &u); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
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
			Active:       false,
		}

		if err := user.Create(ctx, db, usr); err != nil {
			return err
		}

		return web.Respond(ctx, w, usr, http.StatusCreated)
	}
}
