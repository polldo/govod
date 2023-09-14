package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, course Course) error {
	const q = `
	INSERT INTO courses
		(course_id, name, description, price, created_at, updated_at)
	VALUES
	(:course_id, :name, :description, :price, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, course); err != nil {
		return fmt.Errorf("inserting course: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, course Course) (Course, error) {
	const q = `
	UPDATE courses
	SET
		name = :name,
		description = :description,
		price = :price,
		updated_at = :updated_at,
		version = version + 1
	WHERE
		course_id = :course_id AND
		version = :version
	RETURNING version`

	v := struct {
		Version int `db:"version"`
	}{}

	if err := database.NamedQueryStruct(ctx, db, q, course, &v); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Course{}, fmt.Errorf("updating course[%s]: version conflict", course.ID)
		}
		return Course{}, fmt.Errorf("updating course[%s]: %w", course.ID, err)
	}

	course.Version = v.Version

	return course, nil
}

func Fetch(ctx context.Context, db sqlx.ExtContext, id string) (Course, error) {
	in := struct {
		ID string `db:"course_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		*
	FROM
		courses
	WHERE
		course_id = :course_id`

	var course Course
	if err := database.NamedQueryStruct(ctx, db, q, in, &course); err != nil {
		return Course{}, fmt.Errorf("selecting course[%s]: %w", id, err)
	}

	return course, nil
}

func FetchAll(ctx context.Context, db sqlx.ExtContext) ([]Course, error) {
	const q = `
	SELECT
		*
	FROM
		courses
	ORDER BY
		course_id`

	cs := []Course{}
	if err := database.NamedQuerySlice(ctx, db, q, struct{}{}, &cs); err != nil {
		return nil, fmt.Errorf("selecting all courses: %w", err)
	}

	return cs, nil
}

func FetchByOwner(ctx context.Context, db sqlx.ExtContext, userID string) ([]Course, error) {
	in := struct {
		ID     string `db:"user_id"`
		Status string `db:"status"`
	}{
		ID: userID,

		// WARNING: This magic string is tech debt.
		// TODO: Use a const instead of this magic value.
		// This seems a good reason to move all the handlers in the api package.
		// Or just create a models package with all the struct and const.
		Status: "success",
	}

	const q = `
	SELECT
		c.*
	FROM
		orders AS o
	INNER JOIN
		order_items AS i ON i.order_id = o.order_id
	INNER JOIN
		courses AS c ON i.course_id = c.course_id
	WHERE
		o.status = :status AND
		o.user_id = :user_id
	ORDER BY
		c.course_id`

	cs := []Course{}
	if err := database.NamedQuerySlice(ctx, db, q, in, &cs); err != nil {
		return nil, fmt.Errorf("selecting all courses: %w", err)
	}

	return cs, nil
}

func FetchOwned(ctx context.Context, db sqlx.ExtContext, courseID string, userID string) (Course, error) {
	in := struct {
		UserID   string `db:"user_id"`
		CourseID string `db:"course_id"`
		Status   string `db:"status"`
	}{
		UserID:   userID,
		CourseID: courseID,
		Status:   "success",
	}

	const q = `
	SELECT
		c.*
	FROM
		orders AS o
	INNER JOIN
		order_items AS i ON i.order_id = o.order_id
	INNER JOIN
		courses AS c ON i.course_id = c.course_id
	WHERE
		o.status = :status AND
		o.user_id = :user_id AND
		c.course_id = :course_id`

	var cs Course
	if err := database.NamedQueryStruct(ctx, db, q, in, &cs); err != nil {
		return Course{}, fmt.Errorf("selecting owned course: %w", err)
	}

	return cs, nil
}
