package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"golang.org/x/crypto/bcrypt"
)

func HandleCreate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var u UserNew
		if err := web.Decode(r, &u); err != nil {
			return fmt.Errorf("unable to decode payload: %w", err)
			// TODO: Use significative request errors.
			// return weberr.NewError(err, http.StatusInternalServerError, weberr.WithMsg("we couldn't decode your payload!"))
		}

		if err := validate.Check(u); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}

		now := time.Now().UTC()

		user := User{
			ID:           validate.GenerateID(),
			Name:         u.Name,
			Email:        u.Email,
			Role:         u.Role,
			PasswordHash: hash,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		return Create(ctx, db, user)
	}
}
