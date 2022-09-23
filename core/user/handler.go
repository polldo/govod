package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/database"
	"golang.org/x/crypto/bcrypt"
)

func HandleCreate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var u UserNew
		if err := web.Decode(r, &u); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if !claims.IsAdmin(ctx) {
			return weberr.NotAuthorized(errors.New("only admin can create other admins"))
		}

		if err := validate.Check(u); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		now := time.Now().UTC()

		usr := User{
			ID:           validate.GenerateID(),
			Name:         u.Name,
			Email:        u.Email,
			Role:         u.Role,
			PasswordHash: hash,
			CreatedAt:    now,
			UpdatedAt:    now,
			Active:       true,
		}

		if err := Create(ctx, db, usr); err != nil {
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return err
		}

		return web.Respond(ctx, w, usr, http.StatusCreated)
	}
}

func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		userID := web.Param(r, "id")
		if !claims.IsUser(ctx, userID) && !claims.IsAdmin(ctx) {
			return weberr.NotAuthorized(errors.New("user trying to fetch another user"))
		}

		user, err := Fetch(ctx, db, userID)
		if err != nil {
			return fmt.Errorf("ID[%s]: %w", userID, err)
		}

		return web.Respond(ctx, w, user, http.StatusOK)
	}
}
