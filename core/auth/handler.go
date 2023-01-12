package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/core/user"
	"github.com/polldo/govod/database"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}

func HandleOauthLogin(db *sqlx.DB, session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		prov := web.Param(r, "provider")
		if prov != "google" {
			return weberr.NotFound(fmt.Errorf("provider %s not found", prov))
		}

		conf := &oauth2.Config{
			RedirectURL:  "http://mylocal.com:8000/auth/oauth-callback",
			ClientID:     "785050419234-c7ao87rji0crqpkfsu4sr8m77asp4umu.apps.googleusercontent.com",
			ClientSecret: "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",
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
		}

		state := "private-state-here"
		url := conf.AuthCodeURL(state)

		http.Redirect(w, r, url, http.StatusSeeOther)
		return nil
	}
}

func HandleOauthCallback(db *sqlx.DB, session *scs.SessionManager) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		state, code := r.FormValue("state"), r.FormValue("code")
		if state != "private-state-here" {
			return weberr.NotAuthorized(errors.New("wrong state"))
		}

		conf := &oauth2.Config{
			RedirectURL:  "http://mylocal.com:8000/auth/oauth-callback",
			ClientID:     "785050419234-c7ao87rji0crqpkfsu4sr8m77asp4umu.apps.googleusercontent.com",
			ClientSecret: "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",
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
		}

		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			return weberr.NotAuthorized(err)
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
		if err != nil {
			return weberr.NotAuthorized(err)
		}
		defer resp.Body.Close()

		var info struct {
			Name          string `json:"name"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"verified_email"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return weberr.NotAuthorized(err)
		}

		u, err := user.FetchByEmail(ctx, db, info.Email)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				fmt.Println("must register user")
			}
			err := fmt.Errorf("fetching user by email %s: %w", info.Email, err)
			return weberr.NotAuthorized(err)
		}

		// // TODO: Save the entire user struct in the session
		// // or just some info?
		session.Put(ctx, userKey, u.ID)
		session.Put(ctx, roleKey, u.Role)
		if err := session.RenewToken(ctx); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}

func HandleLogout(db *sqlx.DB, session *scs.SessionManager) web.Handler {
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
			return fmt.Errorf("validating data: %w", err)
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
