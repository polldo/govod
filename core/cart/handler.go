package cart

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/core/course"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/validate"
)

// HandleShow returns the cart of the user.
// Returns an empty cart if the user has no cart.
func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		// Return an empty cart if it doesn't exist yet.
		cart, err := Fetch(ctx, db, clm.UserID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return web.Respond(ctx, w, Cart{Items: []Item{}}, http.StatusOK)
			}
			return fmt.Errorf("fetching user[%s] cart: %w", clm.UserID, err)
		}

		cart.Items, err = FetchItems(ctx, db, clm.UserID)
		if err != nil {
			return fmt.Errorf("fetching user[%s] cart items: %w", clm.UserID, err)
		}

		return web.Respond(ctx, w, cart, http.StatusOK)
	}
}

// HandleDelete flushes the user's cart. It also drops
// all related items in cascade.
func HandleDelete(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		// Just delete the cart. Its items will be deleted in cascade.
		if err := Delete(ctx, db, clm.UserID); err != nil {
			return fmt.Errorf("deleting user[%s] cart: %w", clm.UserID, err)
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

// HandleCreateItem adds a new item in the user's cart.
func HandleCreateItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var itnew ItemNew
		if err := web.Decode(w, r, &itnew); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		owned, err := course.FetchByOwner(ctx, db, clm.UserID)
		if err != nil {
			return fmt.Errorf("checking if course[%s] is already owned by user[%s]: %w",
				itnew.CourseID,
				clm.UserID,
				err,
			)
		}

		for _, o := range owned {
			if itnew.CourseID == o.ID {
				err := errors.New("course already owned")
				return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
			}
		}

		if _, err := Upsert(ctx, db, clm.UserID); err != nil {
			return fmt.Errorf("upserting user[%s] cart: %w", clm.UserID, err)
		}

		now := time.Now().UTC()
		item := Item{
			UserID:    clm.UserID,
			CourseID:  itnew.CourseID,
			UpdatedAt: now,
			CreatedAt: now,
		}

		if err := CreateItem(ctx, db, item); err != nil {
			return fmt.Errorf("creating cart item[%s] for user[%s]: %w", item.CourseID, clm.UserID, err)
		}

		return web.Respond(ctx, w, item, http.StatusCreated)
	}
}

// HandleDeleteItem deletes an item from the user's cart.
func HandleDeleteItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "course_id")

		if err := validate.CheckID(courseID); err != nil {
			return weberr.BadRequest(fmt.Errorf("passed id is not valid: %w", err))
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		if _, err := Upsert(ctx, db, clm.UserID); err != nil {
			return fmt.Errorf("upserting user[%s] cart: %w", clm.UserID, err)
		}

		if err := DeleteItem(ctx, db, clm.UserID, courseID); err != nil {
			return fmt.Errorf("deleting user[%s] cart item: %w", clm.UserID, err)
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}
