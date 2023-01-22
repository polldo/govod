package sub

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func CreatePlan(ctx context.Context, db sqlx.ExtContext, plan Plan) error {
	const q = `
	INSERT INTO subscription_plans
		(plan_id, name, price, months_recurrence, created_at, updated_at)
	VALUES
		(:plan_id, :name, :price, :months_recurrence, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, plan); err != nil {
		return fmt.Errorf("inserting subscription plan: %w", err)
	}

	return nil
}

func FetchAllPlans(ctx context.Context, db sqlx.ExtContext) ([]Plan, error) {
	const q = `
	SELECT
		*
	FROM
		subscription_plans
	ORDER BY
		plan_id`

	var ps []Plan
	if err := database.NamedQuerySlice(ctx, db, q, struct{}{}, &ps); err != nil {
		return nil, fmt.Errorf("selecting all susbcription plans: %w", err)
	}

	return ps, nil
}

func Create(ctx context.Context, db sqlx.ExtContext, sub Sub) error {
	const q = `
	INSERT INTO subscriptions
		(subscription_id, plan_id, user_id, provider, status, expiry, created_at, updated_at)
	VALUES
		(:subscription_id, :plan_id, :user_id, :provider, :status, :expiry, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, sub); err != nil {
		return fmt.Errorf("inserting subscription: %w", err)
	}

	return nil
}

func UpdateStatus(ctx context.Context, db sqlx.ExtContext, up StatusUp) error {
	const q = `
	UPDATE subscriptions
	SET
		status = :status,
		expiry = :expiry
		updated_at = :updated_at
	WHERE
		subscription_id = :subscription_id`

	if err := database.NamedExecContext(ctx, db, q, up); err != nil {
		return fmt.Errorf("updating state of subscription[%s]: %w", up.ID, err)
	}

	return nil
}

func FetchActiveByOwner(ctx context.Context, db sqlx.ExtContext, userID string) (Sub, error) {
	in := struct {
		ID     string    `db:"user_id"`
		Status string    `db:"status"`
		Time   time.Time `db:"time"`
	}{
		ID:     userID,
		Status: string(Active),
		Time:   time.Now().UTC(),
	}

	const q = `
	SELECT
		*
	FROM
		subscriptions
	WHERE 
		user_id = :user_id AND
		status = :status AND
		expiry > :time`

	var sub Sub
	if err := database.NamedQueryStruct(ctx, db, q, in, &sub); err != nil {
		return Sub{}, fmt.Errorf("selecting active subscription for user[%s]: %w", userID, err)
	}

	return sub, nil
}
