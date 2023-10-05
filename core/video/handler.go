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
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/core/course"
	"github.com/polldo/govod/database"
	"github.com/polldo/govod/validate"
)

func HandleCreate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var v VideoNew
		if err := web.Decode(w, r, &v); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(v); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
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
			ImageURL:    v.ImageURL,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := Create(ctx, db, video); err != nil {
			err := fmt.Errorf("creating video: %w", err)
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return weberr.BadRequest(err)
			}
			return err
		}

		return web.Respond(ctx, w, video, http.StatusCreated)
	}
}

func HandleUpdate(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		if err := validate.CheckID(videoID); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		var vup VideoUp
		if err := web.Decode(w, r, &vup); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(vup); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			err := fmt.Errorf("fetching video[%s]: %w", videoID, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NotFound(err)
			}
			return err
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
		if vup.ImageURL != nil {
			video.ImageURL = *vup.ImageURL
		}
		video.UpdatedAt = time.Now().UTC()

		if video, err = Update(ctx, db, video); err != nil {
			return fmt.Errorf("updating video[%s]: %w", videoID, err)
		}

		return web.Respond(ctx, w, video, http.StatusOK)
	}
}

func HandleList(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videos, err := FetchAll(ctx, db)
		if err != nil {
			return fmt.Errorf("fetching all videos: %w", err)
		}

		return web.Respond(ctx, w, videos, http.StatusOK)
	}
}

func HandleListByCourse(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "course_id")

		if err := validate.CheckID(courseID); err != nil {
			return weberr.BadRequest(fmt.Errorf("passed id is not valid: %w", err))
		}

		videos, err := FetchAllByCourse(ctx, db, courseID)
		if err != nil {
			return fmt.Errorf("fetching all videos by course[%s]: %w", courseID, err)
		}

		return web.Respond(ctx, w, videos, http.StatusOK)
	}
}

func HandleShow(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		if err := validate.CheckID(videoID); err != nil {
			return weberr.BadRequest(fmt.Errorf("passed id is not valid: %w", err))
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			err := fmt.Errorf("fetching video[%s]: %w", videoID, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NotFound(err)
			}
			return err
		}

		return web.Respond(ctx, w, video, http.StatusOK)
	}
}

func HandleShowFull(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		if err := validate.CheckID(videoID); err != nil {
			return weberr.BadRequest(fmt.Errorf("passed id is not valid: %w", err))
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			err := fmt.Errorf("fetching video[%s]: %w", videoID, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NotFound(err)
			}
			return err
		}

		var crs course.Course
		if video.Free {
			crs, err = course.Fetch(ctx, db, video.CourseID)
			if err != nil {
				return fmt.Errorf("fetching course of free video[%s]: %w", video.ID, err)
			}
		} else {
			crs, err = course.FetchOwned(ctx, db, video.CourseID, clm.UserID)
			if err != nil {
				err := fmt.Errorf("fetching course[%s] owned by user[%s]: %w", video.CourseID, clm.UserID, err)
				if errors.Is(err, database.ErrDBNotFound) {
					return weberr.NewError(err, "access forbidden", http.StatusForbidden)
				}
				return err
			}
		}

		videos, err := FetchAllByCourse(ctx, db, video.CourseID)
		if err != nil {
			err := fmt.Errorf("fetching all videos of course[%s]: %w", video.CourseID, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NotFound(err)
			}
			return err
		}

		progress, err := FetchUserProgressByCourse(ctx, db, clm.UserID, video.CourseID)
		if err != nil {
			return fmt.Errorf("fetching user[%s] progress by course[%s]: %w", clm.UserID, video.CourseID, err)
		}

		fullVideo := struct {
			Course      course.Course `json:"course"`
			Video       Video         `json:"video"`
			AllVideos   []Video       `json:"all_videos"`
			AllProgress []Progress    `json:"all_progress"`
			URL         string        `json:"url"`
		}{
			Course:      crs,
			Video:       video,
			AllVideos:   videos,
			AllProgress: progress,
			URL:         video.URL,
		}

		return web.Respond(ctx, w, fullVideo, http.StatusOK)
	}
}

func HandleShowFree(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		if err := validate.CheckID(videoID); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		video, err := Fetch(ctx, db, videoID)
		if err != nil {
			err := fmt.Errorf("fetching video[%s]: %w", videoID, err)
			if errors.Is(err, database.ErrDBNotFound) {
				return weberr.NotFound(err)
			}
			return err
		}

		if !video.Free {
			return weberr.NewError(err, "access forbidden", http.StatusForbidden)
		}

		crs, err := course.Fetch(ctx, db, video.CourseID)
		if err != nil {
			return fmt.Errorf("fetching course[%s]: %w", video.CourseID, err)
		}

		freeVideo := struct {
			Course course.Course `json:"course"`
			Video  Video         `json:"video"`
			URL    string        `json:"url"`
		}{
			Course: crs,
			Video:  video,
			URL:    video.URL,
		}

		return web.Respond(ctx, w, freeVideo, http.StatusOK)
	}
}

func HandleUpdateProgress(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID := web.Param(r, "id")

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		var up ProgressUp
		if err := web.Decode(w, r, &up); err != nil {
			return weberr.BadRequest(fmt.Errorf("unable to decode payload: %w", err))
		}

		if err := validate.Check(up); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		if err := UpdateProgress(ctx, db, clm.UserID, videoID, up.Progress); err != nil {
			return fmt.Errorf("updating video[%s] progress for user[%s]: %w", videoID, clm.UserID, err)
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}

func HandleListProgressByCourse(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		courseID := web.Param(r, "course_id")

		if err := validate.CheckID(courseID); err != nil {
			return weberr.NewError(err, err.Error(), http.StatusUnprocessableEntity)
		}

		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		progress, err := FetchUserProgressByCourse(ctx, db, clm.UserID, courseID)
		if err != nil {
			return fmt.Errorf("fetching user[%s] progress by course[%s]: %w", clm.UserID, courseID, err)
		}

		return web.Respond(ctx, w, progress, http.StatusOK)
	}
}
