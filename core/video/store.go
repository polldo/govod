package video

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, video Video) error {
	const q = `
	INSERT INTO videos
		(video_id, course_id, index, name, description, free, created_at, updated_at)
	VALUES
	(:video_id, :course_id, :index, :name, :description, :free, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, video); err != nil {
		return fmt.Errorf("inserting video: %w", err)
	}

	return nil
}

func CreateURL(ctx context.Context, db sqlx.ExtContext, url URL) error {
	const q = `
	INSERT INTO videos_url
		(video_id, url, created_at, updated_at)
	VALUES
	(:video_id, :url, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, url); err != nil {
		return fmt.Errorf("inserting video url: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, video Video) error {
	const q = `
	UPDATE videos
	SET
		course_id = :course_id,
		index = :index,
		name = :name,
		description = :description,
		free = :free,
		updated_at = :updated_at
	WHERE
		video_id = :video_id`

	if err := database.NamedExecContext(ctx, db, q, video); err != nil {
		return fmt.Errorf("updating video[%s]: %w", video.ID, err)
	}

	return nil
}

func UpdateURL(ctx context.Context, db sqlx.ExtContext, url URL) error {
	const q = `
	UPDATE videos_url
	SET
		url = :url,
		updated_at = :updated_at
	WHERE
		video_id = :video_id`

	if err := database.NamedExecContext(ctx, db, q, url); err != nil {
		return fmt.Errorf("updating url of video[%s]: %w", url.VideoID, err)
	}

	return nil
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
		return Video{}, fmt.Errorf("updating video[%s]: %w", video.ID, err)
	}

	return video, nil
}

func FetchURL(ctx context.Context, db sqlx.ExtContext, id string) (URL, error) {
	in := struct {
		ID string `db:"video_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		*
	FROM
		videos_url
	WHERE
		video_id = :video_id`

	var url URL
	if err := database.NamedQueryStruct(ctx, db, q, in, &url); err != nil {
		return URL{}, fmt.Errorf("updating url of video[%s]: %w", url.VideoID, err)
	}

	return url, nil
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
	if err := database.NamedQuerySlice(ctx, db, q, nil, &videos); err != nil {
		return nil, fmt.Errorf("selecting videos: %w", err)
	}

	return videos, nil
}
