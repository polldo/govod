package course

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
)

// Admin should be able to create and updated courses.
func HandleCreate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var c CourseNew
		if err := web.Decode(r, &c); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(c); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		now := time.Now().UTC()

		course := Course{
			ID:          validate.GenerateID(),
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := Create(ctx, db, course); err != nil {
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return err
		}

		return web.Respond(ctx, w, course, http.StatusCreated)
	}
}

func HandleUpdate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "id")

		if err := validate.CheckID(courseID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		var cup CourseUp
		if err := web.Decode(r, &cup); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if !claims.IsAdmin(ctx) {
			return weberr.NotAuthorized(errors.New("only admin can update courses"))
		}

		if err := validate.Check(cup); err != nil {
			return fmt.Errorf("validating data: %w", err)
		}

		course, err := Fetch(ctx, db, courseID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		if cup.Name != nil {
			course.Name = *cup.Name
		}
		if cup.Description != nil {
			course.Description = *cup.Description
		}
		if cup.Price != nil {
			course.Price = *cup.Price
		}
		course.UpdatedAt = time.Now().UTC()

		if err := Update(ctx, db, course); err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, course, http.StatusOK)
	}
}

// Users should be able to list all the courses and to fetch specific ones.
func HandleList(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courses, err := FetchAll(ctx, db)
		if err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, courses, http.StatusOK)
	}
}

func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "id")

		if err := validate.CheckID(courseID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		course, err := Fetch(ctx, db, courseID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, course, http.StatusOK)
	}
}
