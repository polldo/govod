package payment

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func Create(ctx context.Context, db sqlx.ExtContext, pay Payment) error {
	const q = `
	INSERT INTO payments
		(payment_id, order_id, provider_id, status, amount, created_at, updated_at)
	VALUES
	(:payment_id, :order_id, :provider_id, :status, :amount, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, pay); err != nil {
		return fmt.Errorf("inserting payment: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status and the date.
func UpdateStatus(ctx context.Context, db sqlx.ExtContext, up StatusUp) error {
	const q = `
	UPDATE payments
	SET
		status = :status,
		updated_at = :updated_at
	WHERE
		payment_id = :payment_id`

	if err := database.NamedExecContext(ctx, db, q, up); err != nil {
		return fmt.Errorf("updating payment[%s]: %w", up.ID, err)
	}

	return nil
}
