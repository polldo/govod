package order

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

// Create order during checkout -> create also items with the current price.
// Let's assume that orders cannot be modified at the moment. This simplifies a lot. Otherwise
// we need to be extremely careful when a user checkouts but then go back and modifies the order.
// Orders are then fetched to assess if a user owns a course -> join order item + order + payment.

func Create(ctx context.Context, db sqlx.ExtContext, order Order) error {
	const q = `
	INSERT INTO orders
		(order_id, user_id, created_at)
	VALUES
		(:order_id, :user_id, :created_at)`

	if err := database.NamedExecContext(ctx, db, q, order); err != nil {
		return fmt.Errorf("inserting order: %w", err)
	}

	return nil
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
