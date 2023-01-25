package sub

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/database"
)

func CreatePlan(ctx context.Context, db sqlx.ExtContext, plan Plan) error {
	const q = `
	INSERT INTO subscription_plans
		(plan_id, stripe_id, paypal_id, name, price, months_recurrence, created_at, updated_at)
	VALUES
		(:plan_id, :stripe_id, :paypal_id, :name, :price, :months_recurrence, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, db, q, plan); err != nil {
		return fmt.Errorf("inserting subscription plan: %w", err)
	}

	return nil
}

func FetchPlan(ctx context.Context, db sqlx.ExtContext, id string) (Plan, error) {
	in := struct {
		ID string `db:"plan_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT 
		*
	FROM
		subscription_plans
	WHERE
		plan_id = :plan_id`

	var plan Plan
	if err := database.NamedQueryStruct(ctx, db, q, in, &plan); err != nil {
		return Plan{}, fmt.Errorf("selecting plan[%s]: %w", id, err)
	}

	return plan, nil
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
		(subscription_id, plan_id, user_id, provider, provider_id, status, expiry, created_at, updated_at)
	VALUES
		(:subscription_id, :plan_id, :user_id, :provider, :provider_id, :status, :expiry, :created_at, :updated_at)`

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
		expiry = :expiry,
		updated_at = :updated_at
	WHERE
		provider_id = :provider_id`

	if err := database.NamedExecContext(ctx, db, q, up); err != nil {
		return fmt.Errorf("updating state of subscription with providerID[%s]: %w", up.ProviderID, err)
	}

	return nil
}

func UpdateProviderID(ctx context.Context, db sqlx.ExtContext, subID string, providerID string) error {
	in := struct {
		ID     string `db:"subscription_id"`
		provID string `db:"provider_id"`
	}{
		ID:     subID,
		provID: providerID,
	}

	const q = `
	UPDATE subscriptions
	SET
		provider_id = :provider_id
	WHERE
		subscription_id = :subscription_id`

	if err := database.NamedExecContext(ctx, db, q, in); err != nil {
		return fmt.Errorf("updating provider id of subscription[%s]: %w", in.ID, err)
	}

	return nil
}

func Fetch(ctx context.Context, db sqlx.ExtContext, id string) (Sub, error) {
	in := struct {
		ID string `db:"subscription_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT
		*
	FROM
		subscriptions
	WHERE
		subscription_id = :subscription_id`

	var sub Sub
	if err := database.NamedQueryStruct(ctx, db, q, in, &sub); err != nil {
		return Sub{}, fmt.Errorf("selecting subscription[%s]: %w", id, err)
	}

	return sub, nil
}

func FetchLastByOwnerStatus(ctx context.Context, db sqlx.ExtContext, userID string, status Status) (Sub, error) {
	in := struct {
		ID     string `db:"user_id"`
		Status string `db:"status"`
	}{
		ID:     userID,
		Status: string(status),
	}

	const q = `
	SELECT
		*
	FROM
		subscriptions
	WHERE 
		user_id = :user_id AND
		status = :status
	ORDER BY created_at DESC
	LIMIT 1`

	var sub Sub
	if err := database.NamedQueryStruct(ctx, db, q, in, &sub); err != nil {
		return Sub{}, fmt.Errorf("selecting last subscription[%s] for user[%s]: %w", status, userID, err)
	}

	return sub, nil
}

func FetchLastByOwner(ctx context.Context, db sqlx.ExtContext, userID string) (Sub, error) {
	in := struct {
		ID string `db:"user_id"`
	}{
		ID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		subscriptions
	WHERE
		user_id = :user_id
	ORDER BY created_at DESC
	LIMIT 1`

	var sub Sub
	if err := database.NamedQueryStruct(ctx, db, q, in, &sub); err != nil {
		return Sub{}, fmt.Errorf("selecting last subscription for user[%s]: %w", userID, err)
	}

	return sub, nil
}
