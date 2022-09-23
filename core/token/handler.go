package token

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/user"
	"github.com/polldo/govod/database"
	"golang.org/x/crypto/bcrypt"
)

type Mailer interface {
	Send(dest string, tmplName string, token string) error
}

func HandleToken(db *sqlx.DB, mailer Mailer) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Email string `json:"email" validate:"required"`
			Scope string `json:"scope" validate:"required"`
		}

		if err := web.Decode(r, &in); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(in); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		usr, err := user.FetchByEmail(ctx, db, in.Email)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
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

		// Wrap in a transaction the following 2 db queries.

		// Delete pending tokens.
		if err := DeleteByUser(ctx, db, usr.ID, scope); err != nil {
			return err
		}

		text, token, err := GenToken(usr.ID, 6*time.Hour, scope)
		if err != nil {
			return err
		}

		if err := Create(ctx, db, token); err != nil {
			return err
		}

		// In a goroutine with multiple tries ??
		if err := mailer.Send(usr.Email, scope, text); err != nil {
			return err
		}

		return nil
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
			return fmt.Errorf("validating data: %w", err)
		}

		hash := sha256.Sum256([]byte(in.Token))

		// Probably a transaction here -> update the user only if the token gets deleted.
		usr, err := user.FetchByToken(ctx, db, hash[:], ActivationToken)
		if err != nil {
			return err
		}

		usr.Active = true
		usr.UpdatedAt = time.Now().UTC()
		if err := user.Update(ctx, db, usr); err != nil {
			return err
		}

		if err := DeleteByUser(ctx, db, usr.ID, ActivationToken); err != nil {
			return err
		}

		return nil
	}
}

func HandleRecovery(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Token           string `json:"token" validate:"required"`
			Password        string `json:"password" validate:"required"`
			PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
		}

		if err := web.Decode(r, &in); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(in); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		tokh := sha256.Sum256([]byte(in.Token))

		// Probably a transaction here -> update the user only if the token gets deleted.
		usr, err := user.FetchByToken(ctx, db, tokh[:], RecoveryToken)
		if err != nil {
			return err
		}

		passh, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		usr.PasswordHash = passh
		usr.UpdatedAt = time.Now().UTC()
		if err := user.Update(ctx, db, usr); err != nil {
			return err
		}

		if err := DeleteByUser(ctx, db, usr.ID, RecoveryToken); err != nil {
			return err
		}

		return nil
	}
}
