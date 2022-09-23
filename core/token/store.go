package token

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, token Token) error {
	const q = `
	INSERT INTO tokens
		(hash, user_id, expiry, scope)
	VALUES
		(:hash, :user_id, :expiry, :scope)`

	if err := database.NamedExecContext(ctx, db, q, token); err != nil {
		return fmt.Errorf("inserting token: %w", err)
	}

	return nil
}

func DeleteByUser(ctx context.Context, db sqlx.ExtContext, userID string, scope string) error {
	data := struct {
		UserID string `db:"user_id"`
		Scope  string `db:"scope"`
	}{
		UserID: userID,
		Scope:  scope,
	}

	const q = `
	DELETE FROM tokens
	WHERE user_id = :user_id AND scope = :scope`

	if err := database.NamedExecContext(ctx, db, q, data); err != nil {
		return fmt.Errorf("deleting token: %w", err)
	}

	return nil
}
