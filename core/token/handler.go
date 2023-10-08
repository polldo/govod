package token

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/background"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/user"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/rate"
	"github.com/polldo/govod/validate"
	"golang.org/x/crypto/bcrypt"
)

// Mailer should be able to send emails to users
// for handling their activation and their password recovery.
type Mailer interface {
	SendActivationToken(token string, to string) error
	SendRecoveryToken(token string, to string) error
}

// HandleToken is used to send specific tokens to users via email.
// A valid scope must be provided by users, together with their email.
// It doesn't require a user to be logged in, because users who need
// tokens will probably not be able to login at all.
// This function leverages a rate limiter to avoid too many emails.
func HandleToken(db *sqlx.DB, mailer Mailer, timeout time.Duration, bg *background.Background) web.Handler {
	rps := rate.Every(timeout)
	limiter := rate.NewLimiter(1, 10, float64(rps))

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Email string `json:"email" validate:"required,email"`
			Scope string `json:"scope" validate:"required"`
		}

		if err := web.Decode(w, r, &in); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		if !limiter.Check(in.Email) {
			err := errors.New("too many requests")
			return weberr.NewError(err, err.Error(), http.StatusTooManyRequests)
		}

		usr, err := user.FetchByEmail(ctx, db, in.Email)
		if err != nil {
			err := fmt.Errorf("fetching token by user[%s]: %w", in.Email, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, "Email is not registered", http.StatusUnprocessableEntity)
			}
			return err
		}

		scope := in.Scope
		switch scope {
		case ActivationToken:
			if usr.Active {
				return weberr.BadRequest(fmt.Errorf("user %s is already active", usr.Email))
			}
		case RecoveryToken:
		default:
			return weberr.BadRequest(fmt.Errorf("scope %s is not supported", scope))
		}

		text, token, err := GenToken(usr.ID, 6*time.Hour, scope)
		if err != nil {
			return fmt.Errorf("generating random token: %w", err)
		}

		// Delete pending tokens only if the new token is actually stored.
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, scope); err != nil {
				return fmt.Errorf("deleting token by user[%s]: %w", usr.ID, err)
			}

			if err := Create(ctx, tx, token); err != nil {
				return fmt.Errorf("creating new token for user[%s]: %w", usr.ID, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		// Send the email in background.
		bg.Add(func() error {
			switch scope {
			case ActivationToken:
				if err := mailer.SendActivationToken(text, usr.Email); err != nil {
					return fmt.Errorf("failed to send activation token %s to %s: %w", scope, usr.Email, err)
				}
			case RecoveryToken:
				if err := mailer.SendRecoveryToken(text, usr.Email); err != nil {
					return fmt.Errorf("failed to send recovery token %s to %s: %w", scope, usr.Email, err)
				}
			default:
				return fmt.Errorf("scope %s is not supported", scope)
			}
			return nil
		})

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

// HandleActivation validates the passed token and, if correct,
// it activates the user.
func HandleActivation(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Token string `json:"token" validate:"required"`
		}

		if err := web.Decode(w, r, &in); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		hash := sha256.Sum256([]byte(in.Token))

		usr, err := user.FetchByToken(ctx, db, hash[:], ActivationToken)
		if err != nil {
			err := fmt.Errorf("fetching user by token: %w", err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.BadRequest(err)
			}
			return err
		}

		// Delete the token only if the user gets updated correctly (and viceversa).
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, ActivationToken); err != nil {
				return fmt.Errorf("deleting token by user[%s]: %w", usr.ID, err)
			}

			usr.Active = true
			usr.UpdatedAt = time.Now().UTC()
			if _, err := user.Update(ctx, tx, usr); err != nil {
				return fmt.Errorf("activating user[%s]: %w", usr.ID, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	}
}

// HandleRecovery validates the passed token and, if correct,
// changes the user's password with the one provided.
func HandleRecovery(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Token           string `json:"token" validate:"required"`
			Password        string `json:"password" validate:"required,gte=8,lte=50"`
			PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
		}

		if err := web.Decode(w, r, &in); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		tokh := sha256.Sum256([]byte(in.Token))

		usr, err := user.FetchByToken(ctx, db, tokh[:], RecoveryToken)
		if err != nil {
			err := fmt.Errorf("fetch user by token: %w", err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.BadRequest(err)
			}
			return err
		}

		passh, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		// Delete the token only if the user gets updated correctly (and viceversa).
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, RecoveryToken); err != nil {
				return fmt.Errorf("deleting token by user[%s]: %w", usr.ID, err)
			}

			usr.PasswordHash = passh
			usr.UpdatedAt = time.Now().UTC()
			if _, err := user.Update(ctx, tx, usr); err != nil {
				return fmt.Errorf("recoverying user[%s]: %w", usr.ID, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	}
}
