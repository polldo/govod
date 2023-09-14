package cart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Fetch(ctx context.Context, db sqlx.ExtContext, userID string) (Cart, error) {
	in := struct {
		ID string `db:"user_id"`
	}{
		ID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		carts
	WHERE
		user_id = :user_id`

	var c Cart
	if err := database.NamedQueryStruct(ctx, db, q, in, &c); err != nil {
		return Cart{}, fmt.Errorf("selecting cart of user[%s]: %w", userID, err)
	}

	return c, nil
}

func Create(ctx context.Context, db sqlx.ExtContext, cart Cart) error {
	const q = `
	INSERT INTO carts
		(user_id, created_at, updated_at)
	VALUES
	(:user_id, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, cart); err != nil {
		return fmt.Errorf("inserting cart: %w", err)
	}

	return nil
}

func Delete(ctx context.Context, db sqlx.ExtContext, userID string) error {
	in := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	DELETE FROM 
		carts
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, db, q, in); err != nil {
		return fmt.Errorf("deleting cart: %w", err)
	}

	return nil
}

func Update(ctx context.Context, db sqlx.ExtContext, cart Cart) (Cart, error) {
	const q = `
	UPDATE carts
	SET
		updated_at = :updated_at,
		version = version + 1
	WHERE
		user_id = :user_id AND
		version = :version
	RETURNING version`

	v := struct {
		Version int `db:"version"`
	}{}

	if err := database.NamedQueryStruct(ctx, db, q, cart, &v); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Cart{}, fmt.Errorf("updating cart of user[%s]: version conflict", cart.UserID)
		}
		return Cart{}, fmt.Errorf("updating cart of user[%s]: %w", cart.UserID, err)
	}

	cart.Version = v.Version

	return cart, nil
}

func Upsert(ctx context.Context, db sqlx.ExtContext, userID string) (Cart, error) {
	cart, err := Fetch(ctx, db, userID)
	if err != nil {
		// Just abort in case of unexpected errors.
		if !errors.Is(err, database.ErrDBNotFound) {
			return Cart{}, err
		}

		// Create the cart if it doesn't exist.
		now := time.Now().UTC()
		cart := Cart{
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := Create(ctx, db, cart)
		return cart, err
	}

	// Update the cart if it already exists.
	cart.UpdatedAt = time.Now().UTC()
	return Update(ctx, db, cart)
}

func FetchItems(ctx context.Context, db sqlx.ExtContext, userID string) ([]Item, error) {
	in := struct {
		ID string `db:"user_id"`
	}{
		ID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		cart_items
	WHERE
		user_id = :user_id
	ORDER BY
		course_id`

	ci := []Item{}
	if err := database.NamedQuerySlice(ctx, db, q, in, &ci); err != nil {
		return nil, fmt.Errorf("selecting cart items of user[%s]: %w", userID, err)
	}

	return ci, nil
}

func CreateItem(ctx context.Context, db sqlx.ExtContext, item Item) error {
	const q = `
	INSERT INTO cart_items
		(user_id, course_id, created_at, updated_at)
	VALUES
	(:user_id, :course_id, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, item); err != nil {
		return fmt.Errorf("inserting cart item: %w", err)
	}

	return nil
}

func DeleteItem(ctx context.Context, db sqlx.ExtContext, userID string, courseID string) error {
	in := struct {
		UserID   string `db:"user_id"`
		CourseID string `db:"course_id"`
	}{
		UserID:   userID,
		CourseID: courseID,
	}

	const q = `
	DELETE FROM 
		cart_items
	WHERE
		user_id = :user_id AND course_id = :course_id`

	if err := database.NamedExecContext(ctx, db, q, in); err != nil {
		return fmt.Errorf("deleting cart item: %w", err)
	}

	return nil
}
