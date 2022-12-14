package order

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

// Create order during checkout -> create also items with the current price.
// Let's assume that order items cannot be modified at the moment. This simplifies a lot. Otherwise
// we need to be extremely careful when a user checkouts but then go back and modifies the order.
// Orders are then fetched to assess if a user owns a course -> join order item + order + payment.

func Create(ctx context.Context, db sqlx.ExtContext, order Order) error {
	const q = `
	INSERT INTO orders
		(order_id, user_id, provider_id, status, created_at, updated_at)
	VALUES
		(:order_id, :user_id, :provider_id, :status, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, order); err != nil {
		return fmt.Errorf("inserting order: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status and the date of an order.
func UpdateStatus(ctx context.Context, db sqlx.ExtContext, up StatusUp) error {
	const q = `
	UPDATE orders
	SET
		status = :status,
		updated_at = :updated_at
	WHERE
		order_id = :order_id`

	if err := database.NamedExecContext(ctx, db, q, up); err != nil {
		return fmt.Errorf("updating state of order[%s]: %w", up.ID, err)
	}

	return nil
}

func FetchByProviderID(ctx context.Context, db sqlx.ExtContext, provID string) (Order, error) {
	in := struct {
		ProviderID string `db:"provider_id"`
	}{
		ProviderID: provID,
	}

	const q = `
	SELECT
		*
	FROM
		orders
	WHERE 
		provider_id = :provider_id`

	var order Order
	if err := database.NamedQueryStruct(ctx, db, q, in, &order); err != nil {
		return Order{}, fmt.Errorf("selecting order by provider_id[%s]: %w", provID, err)
	}

	return order, nil
}

func CreateItem(ctx context.Context, db sqlx.ExtContext, item Item) error {
	const q = `
	INSERT INTO order_items
		(order_id, course_id, price, created_at)
	VALUES
	(:order_id, :course_id, :price, :created_at)`

	if err := database.NamedExecContext(ctx, db, q, item); err != nil {
		return fmt.Errorf("inserting order item: %w", err)
	}

	return nil
}
