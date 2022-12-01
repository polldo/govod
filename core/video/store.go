package video

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, video Video) error {
	const q = `
	INSERT INTO videos
		(video_id, course_id, index, name, description, free, url, created_at, updated_at)
	VALUES
	(:video_id, :course_id, :index, :name, :description, :free, :url, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, video); err != nil {
		return fmt.Errorf("inserting video: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, video Video) (Video, error) {
	const q = `
	UPDATE videos
	SET
		course_id = :course_id,
		index = :index,
		name = :name,
		description = :description,
		free = :free,
		updated_at = :updated_at,
		version = version + 1
	WHERE
		video_id = :video_id AND
		version = :version
	RETURNING version`

	v := struct {
		Version int `db:"version"`
	}{}

	if err := database.NamedQueryStruct(ctx, db, q, video, &v); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Video{}, fmt.Errorf("updating video[%s]: version conflict", video.ID)
		}
		return Video{}, fmt.Errorf("updating video[%s]: %w", video.ID, err)
	}

	video.Version = v.Version

	return video, nil
}

func Fetch(ctx context.Context, db sqlx.ExtContext, id string) (Video, error) {
	in := struct {
		ID string `db:"video_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		*
	FROM
		videos
	WHERE
		video_id = :video_id`

	var video Video
	if err := database.NamedQueryStruct(ctx, db, q, in, &video); err != nil {
		return Video{}, fmt.Errorf("fetching video[%s]: %w", id, err)
	}

	return video, nil
}

func FetchAll(ctx context.Context, db sqlx.ExtContext) ([]Video, error) {
	const q = `
	SELECT
		*
	FROM
		videos
	ORDER BY
		video_id`

	var videos []Video
	if err := database.NamedQuerySlice(ctx, db, q, struct{}{}, &videos); err != nil {
		return nil, fmt.Errorf("selecting videos: %w", err)
	}

	return videos, nil
}
