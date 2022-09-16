package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
)

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
