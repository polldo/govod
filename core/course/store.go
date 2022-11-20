package course

import (
	"context"
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

func Update(ctx context.Context, db sqlx.ExtContext, course Course) error {
	const q = `
	UPDATE courses
	SET
		name = :name,
		description = :description,
		price = :price,
		updated_at = :updated_at
	WHERE
		course_id = :course_id`

	if err := database.NamedExecContext(ctx, db, q, course); err != nil {
		return fmt.Errorf("updating course[%s]: %w", course.ID, err)
	}

	return nil
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
		return Course{}, fmt.Errorf("selecting course[%s]: %w", course.ID, err)
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

	var cs []Course
	if err := database.NamedQuerySlice(ctx, db, q, nil, &cs); err != nil {
		return nil, fmt.Errorf("selecting all courses: %w", err)
	}

	return cs, nil
}
