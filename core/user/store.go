package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, user User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, role, active, created_at, updated_at)
	VALUES
	(:user_id, :name, :email, :password_hash, :role, :active, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, user); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, user User) (User, error) {
	const q = `
	UPDATE users
	SET
		name = :name,
		email = :email,
		role = :role,
		active = :active,
		password_hash = :password_hash,
		updated_at = :updated_at,
		version = version + 1
	WHERE
		user_id = :user_id AND
		version = :version
	RETURNING version`

	v := struct {
		Version int `db:"version"`
	}{}

	if err := database.NamedQueryStruct(ctx, db, q, user, &v); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return User{}, fmt.Errorf("updating user[%s]: version conflict", user.ID)
		}
		return User{}, fmt.Errorf("updating course[%s]: %w", user.ID, err)
	}

	user.Version = v.Version

	return user, nil
}

func Fetch(ctx context.Context, db sqlx.ExtContext, id string) (User, error) {
	in := struct {
		ID string `db:"user_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

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
		tokens AS t ON t.user_id = u.user_id
	WHERE 
		t.hash = :hash AND t.scope = :scope  AND t.expiry > :time`

	var user User
	if err := database.NamedQueryStruct(ctx, db, q, in, &user); err != nil {
		return User{}, fmt.Errorf("selecting user by token: %w", err)
	}

	return user, nil
}
