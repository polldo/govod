package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, user User) error {
	const q = `
	INSERT INTO users
		(id, name, email, password, role, created_at, updated_at)
	VALUES
		(:id, :name, :email, :password, :role, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, user); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

func Fetch(ctx context.Context, db sqlx.ExtContext, id string) (User, error) {
	in := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		id = :id`

	var user User
	if err := database.NamedQueryStruct(ctx, db, q, in, &user); err != nil {
		return User{}, fmt.Errorf("selecting user id[%q]: %w", id, err)
	}

	return user, nil
}
