package sub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/database"
)

// listen to webhooks to update subscriptions state

func HandleListPlans(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		plans, err := FetchAllPlans(ctx, db)
		if err != nil {
			return weberr.InternalError(err)
		}

		return web.Respond(ctx, w, plans, http.StatusOK)
	}
}

// IsActive checks whether a user has an active subscription.
func IsActive(ctx context.Context, db *sqlx.DB, userID string) (bool, error) {
	if _, err := FetchActiveByOwner(ctx, db, userID); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func HandlePaypalCheckout(db *sqlx.DB, pp *paypal.Client) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		actv, err := IsActive(ctx, db, clm.UserID)
		if err != nil {
			return weberr.InternalError(err)
		}

		if actv {
			err = errors.New("user is already subscribed")
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		var snew SubNew
		if err := web.Decode(r, &snew); err != nil {
			err = fmt.Errorf("unable to decode payload: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if err := validate.CheckID(snew.PlanID); err != nil {
			err = fmt.Errorf("passed plan id is not valid: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		plan, err := FetchPlan(ctx, db, snew.PlanID)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				err = errors.New("unknown plan_id")
				return weberr.NewError(err, err.Error(), http.StatusBadRequest)
			}
			return err
		}

		ppsub, err := pp.CreateSubscription(ctx, paypal.SubscriptionBase{
			PlanID: plan.PaypalID,
		})
		if err != nil {
			return fmt.Errorf("creating paypal order: %w", err)
		}

		now := time.Now().UTC()
		sub := Sub{
			ID:        ppsub.SubscriptionDetails.ID,
			PlanID:    plan.ID,
			UserID:    clm.UserID,
			Provider:  Paypal,
			Status:    Pending,
			CreatedAt: now,
			UpdatedAt: now,
			Expiry:    now,
		}

		if err := Create(ctx, db, sub); err != nil {
			return fmt.Errorf("creating subscription: %w", err)
		}

		return web.Respond(ctx, w, sub, http.StatusOK)
	}
}

func HandlePaypalWebhook(db *sqlx.DB, pp *paypal.Client, webhookID string) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		rsp, err := pp.VerifyWebhookSignature(ctx, r, webhookID)
		if err != nil {
			return err
		}

		if rsp.VerificationStatus != "SUCCESS" {
			return err
		}

		evt := paypal.Event{}
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			return err
		}

		fmt.Printf("Webhook: %+v\n", evt)

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}
