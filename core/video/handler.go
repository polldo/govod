package video

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/validate"
)

func HandleCreate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var v VideoNew
		if err := web.Decode(r, &v); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(v); err != nil {
			err = fmt.Errorf("validating data: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		now := time.Now().UTC()

		video := Video{
			ID:          validate.GenerateID(),
			CourseID:    v.CourseID,
			Index:       v.Index,
			Name:        v.Name,
			Description: v.Description,
			Free:        v.Free,
			URL:         v.URL,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := Create(ctx, db, video); err != nil {
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, video, http.StatusCreated)
	}
}

func HandleUpdate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		if err := validate.CheckID(videoID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		var vup VideoUp
		if err := web.Decode(r, &vup); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.Check(vup); err != nil {
			err = fmt.Errorf("validating data: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		if vup.CourseID != nil {
			video.CourseID = *vup.CourseID
		}
		if vup.Index != nil {
			video.Index = *vup.Index
		}
		if vup.Name != nil {
			video.Name = *vup.Name
		}
		if vup.Description != nil {
			video.Description = *vup.Description
		}
		if vup.Free != nil {
			video.Free = *vup.Free
		}
		if vup.URL != nil {
			video.URL = *vup.URL
		}
		video.UpdatedAt = time.Now().UTC()

		if video, err = Update(ctx, db, video); err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, video, http.StatusOK)
	}
}

func HandleList(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videos, err := FetchAll(ctx, db)
		if err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, videos, http.StatusOK)
	}
}

func HandleListByCourse(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "course_id")

		if err := validate.CheckID(courseID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		videos, err := FetchAllByCourse(ctx, db, courseID)
		if err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, videos, http.StatusOK)
	}
}

func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		if err := validate.CheckID(videoID); err != nil {
			err = fmt.Errorf("passed id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, video, http.StatusOK)
	}
}
