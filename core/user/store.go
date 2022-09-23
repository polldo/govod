package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, user User) error {
	const q = `
	INSERT INTO users
		(id, name, email, password_hash, role, active, created_at, updated_at)
	VALUES
	(:id, :name, :email, :password_hash, :role, :active, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, user); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, user User) error {
	const q = `
	UPDATE users
	SET
		name = :name,
		email = :email,
		role = :role,
		active = :active,
		password_hash = :password_hash,
		updated_at = :updated_at
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, db, q, user); err != nil {
		return fmt.Errorf("updating user[%s]: %w", user.ID, err)
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

func FetchByEmail(ctx context.Context, db sqlx.ExtContext, email string) (User, error) {
	in := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		email = :email`

	var user User
	if err := database.NamedQueryStruct(ctx, db, q, in, &user); err != nil {
		return User{}, fmt.Errorf("selecting user email[%q]: %w", email, err)
	}

	return user, nil
}

func FetchByToken(ctx context.Context, db sqlx.ExtContext, tokenHash []byte, tokenScope string) (User, error) {
	in := struct {
		Hash  []byte    `db:"hash"`
		Scope string    `db:"scope"`
		Time  time.Time `db:"time"`
	}{
		Hash:  tokenHash,
		Scope: tokenScope,
		Time:  time.Now().UTC(),
	}

	const q = `
	SELECT
		u.*
	FROM
		users AS u
	LEFT JOIN
		tokens AS t ON t.user_id = u.id
	WHERE 
		t.hash = :hash AND t.scope = :scope  AND t.expiry > :time`

	var user User
	if err := database.NamedQueryStruct(ctx, db, q, in, &user); err != nil {
		return User{}, fmt.Errorf("selecting user by token: %w", err)
	}

	return user, nil
}
