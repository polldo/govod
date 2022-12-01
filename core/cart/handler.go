package cart

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
)

// Add item to cart -> if cart doesn't exist, create it.
// A user can have at most one cart.
// Delete a cart instead to flush it, when an order gets payed.
// User can add and remove items to the cart.

// Only authenticated users.
//
// TODO: Decide if a denormalization could be useful to reduce HTTP calls.
func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return err
		}

		cart, err := Fetch(ctx, db, clm.UserID)
		if err != nil {
			return err
		}

		cart.Items, err = FetchItems(ctx, db, clm.UserID)
		if err != nil {
			return err
		}

		return web.Respond(ctx, w, cart, http.StatusOK)
	}
}

// Only authenticated users.
func HandleCreateItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var itnew ItemNew
		if err := web.Decode(r, &itnew); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return err
		}

		if _, err := Upsert(ctx, db, clm.UserID); err != nil {
			return err
		}

		now := time.Now().UTC()
		item := Item{
			UserID:    clm.UserID,
			CourseID:  itnew.CourseID,
			UpdatedAt: now,
			CreatedAt: now,
		}

		if err := CreateItem(ctx, db, item); err != nil {
			return err
		}

		return web.Respond(ctx, w, item, http.StatusCreated)
	}
}

// Only authenticated users.
func HandleDeleteItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "id")

		if err := validate.CheckID(courseID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return err
		}

		if _, err := Upsert(ctx, db, clm.UserID); err != nil {
			return err
		}

		if err := DeleteItem(ctx, db, clm.UserID, courseID); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}
