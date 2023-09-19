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
	"github.com/polldo/govod/validate"
	"golang.org/x/crypto/bcrypt"
)

type Mailer interface {
	SendActivationToken(token string, to string) error
	SendResetToken(token string, to string) error
}

// TODO: Rate-limit or introduce a timeout between several requests.
// TODO: If there exists a not expired token, return it instead of creating a new one.
// Or, in alternative, deletes all tokens only when the user is actually activated.
func HandleToken(db *sqlx.DB, mailer Mailer, bg *background.Background) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Email string `json:"email" validate:"required,email"`
			Scope string `json:"scope" validate:"required"`
		}

		if err := web.Decode(r, &in); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		usr, err := user.FetchByEmail(ctx, db, in.Email)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, "Email is not registered", http.StatusUnprocessableEntity)
			}
			return err
		}

		scope := in.Scope
		switch scope {

		case ActivationToken:
			fmt.Println("Activation!")
			if usr.Active {
				err := fmt.Errorf("user %s is already active", usr.Email)
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}

		case RecoveryToken:
			fmt.Println("Reset!")

		default:
			err := fmt.Errorf("scope %s is not supported", scope)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		text, token, err := GenToken(usr.ID, 6*time.Hour, scope)
		if err != nil {
			return err
		}

		// Delete pending tokens only if the new token is actually stored.
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, scope); err != nil {
				return err
			}

			if err := Create(ctx, tx, token); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return weberr.InternalError(err)
		}

		bg.Add(func() error {
			switch scope {
			case ActivationToken:
				if err := mailer.SendActivationToken(text, usr.Email); err != nil {
					return fmt.Errorf("failed to send activation token %s to %s: %w", scope, usr.Email, err)
				}
			case RecoveryToken:
				if err := mailer.SendResetToken(text, usr.Email); err != nil {
					return fmt.Errorf("failed to send reset token %s to %s: %w", scope, usr.Email, err)
				}
			default:
				return fmt.Errorf("scope %s is not supported", scope)
			}
			return nil
		})

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

func HandleActivation(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Token string `json:"token" validate:"required"`
		}

		if err := web.Decode(r, &in); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, "not valid", http.StatusUnprocessableEntity)
		}

		hash := sha256.Sum256([]byte(in.Token))

		usr, err := user.FetchByToken(ctx, db, hash[:], ActivationToken)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		// Delete the token only if the user gets updated correctly (and viceversa).
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, ActivationToken); err != nil {
				return err
			}

			usr.Active = true
			usr.UpdatedAt = time.Now().UTC()
			if _, err := user.Update(ctx, tx, usr); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return weberr.InternalError(err)
		}

		return nil
	}
}

func HandleRecovery(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Token           string `json:"token" validate:"required"`
			Password        string `json:"password" validate:"required,gte=8,lte=50"`
			PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
		}

		if err := web.Decode(r, &in); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(in); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		tokh := sha256.Sum256([]byte(in.Token))

		usr, err := user.FetchByToken(ctx, db, tokh[:], RecoveryToken)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		passh, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		// Delete the token only if the user gets updated correctly (and viceversa).
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := DeleteByUser(ctx, tx, usr.ID, RecoveryToken); err != nil {
				return err
			}

			usr.PasswordHash = passh
			usr.UpdatedAt = time.Now().UTC()
			if _, err := user.Update(ctx, tx, usr); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return weberr.InternalError(err)
		}

		return nil
	}
}
