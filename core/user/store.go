package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Create(ctx context.Context, db sqlx.ExtContext, user User) error {
	const q = `
	INSERT INTO users
		(id, name, email, password, role, created_at, updated_at)
	VALUES
		(:id, :name, :email, :password, :role, :created_at, :updated_at)`

	if _, err := sqlx.NamedExecContext(ctx, db, q, user); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}
