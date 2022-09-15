package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
)

func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		userID := web.Param(r, "id")

		// TODO: Add auth checks.
		// claims, err := auth.GetClaims(ctx)
		// if err != nil {
		// 	return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		// }
		//
		// // If you are not an admin and looking to retrieve someone other than yourself.
		// if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		// 	return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		// }

		user, err := Fetch(ctx, db, userID)
		if err != nil {
			// TODO: Improve errors.
			// switch {
			// case errors.Is(err, user.ErrInvalidID):
			// 	return v1Web.NewRequestError(err, http.StatusBadRequest)
			// case errors.Is(err, user.ErrNotFound):
			// 	return v1Web.NewRequestError(err, http.StatusNotFound)
			// default:
			// 	return fmt.Errorf("ID[%s]: %w", userID, err)
			// }
			return fmt.Errorf("ID[%s]: %w", userID, err)
		}

		return web.Respond(ctx, w, user, http.StatusOK)
	}
}
