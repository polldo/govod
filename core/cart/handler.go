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

// TODO: Evaluate whether a denormalization could be useful to reduce HTTP calls.
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
			return err
		}

		cart.Items, err = FetchItems(ctx, db, clm.UserID)
		if err != nil {
			return err
		}

		return web.Respond(ctx, w, cart, http.StatusOK)
	}
}

func HandleDelete(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		// Just delete the cart. Its items will be deleted in cascade.
		if err := Delete(ctx, db, clm.UserID); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

func HandleCreateItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var itnew ItemNew
		if err := web.Decode(r, &itnew); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		owned, err := course.FetchByOwner(ctx, db, clm.UserID)
		if err != nil {
			return fmt.Errorf("checking if course is already owned by user: %w", err)
		}

		for _, o := range owned {
			if itnew.CourseID == o.ID {
				err := errors.New("course already owned")
				return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
			}
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

func HandleDeleteItem(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "course_id")

		if err := validate.CheckID(courseID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
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
