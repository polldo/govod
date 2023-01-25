package sub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/database"
	"github.com/stripe/stripe-go/v74"
	stripecl "github.com/stripe/stripe-go/v74/client"
	"github.com/stripe/stripe-go/v74/webhook"
)

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
//
// TODO: We may cache the active subscriptions in users' session to reduce database hits.
// We should include both the subscription state and its expiration to make it working.
func IsActive(ctx context.Context, db *sqlx.DB, userID string) (bool, error) {
	sub, err := FetchLastByOwnerStatus(ctx, db, userID, Active)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return false, nil
		}
		return false, err
	}

	// Wait for 5 days before considering the subscription cancelled.
	// Just to give extra time to the user to pay.
	exp := sub.Expiry.Add(24 * time.Hour * 5)

	// Not expired yet, the subscription is still active.
	if time.Now().UTC().Before(exp) {
		return true, nil
	}

	// The subscription has expired.
	// TODO: before considering it cancelled we could fetch its status
	// on the provider site (paypal/stripe) to double check.
	// We may have missed a webhook.
	//
	// Let's update the subscription state to 'cancelled'.
	err = UpdateStatus(ctx, db, StatusUp{
		ProviderID: sub.ProviderID,
		Status:     Cancelled,
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		return false, fmt.Errorf("subscription[%s] to be cancelled but failed to update the db", sub.ID)
	}

	return false, nil
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
			PlanID:             plan.PaypalID,
			ApplicationContext: &paypal.ApplicationContext{ShippingPreference: "NO_SHIPPING"},
		})
		if err != nil {
			return fmt.Errorf("creating paypal order: %w", err)
		}

		now := time.Now().UTC()
		sub := Sub{
			ID:         validate.GenerateID(),
			PlanID:     plan.ID,
			UserID:     clm.UserID,
			Provider:   Paypal,
			ProviderID: ppsub.SubscriptionDetails.ID,
			Status:     Pending,
			CreatedAt:  now,
			UpdatedAt:  now,
			Expiry:     now,
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

		evt := paypal.AnyEvent{}
		if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
			return err
		}

		switch evt.EventType {
		case "PAYMENT.SALE.COMPLETED":
			sale := struct {
				ID string `json:"billing_agreement_id"`
			}{}

			if err := json.Unmarshal(evt.Resource, &sale); err != nil {
				return err
			}

			// Retrieve subscription from paypal to know the exact expiration.
			// The alternative is to fetch that info from the plan; but that
			// would not be accurate because the user may pay the subscription later.
			info, err := pp.GetSubscriptionDetails(ctx, sale.ID)
			if err != nil {
				return err
			}

			if info.SubscriptionStatus != paypal.SubscriptionStatusActive {
				return fmt.Errorf("sale completed but subscription[%s] is not active yet", sale.ID)
			}

			up := StatusUp{
				ProviderID: sale.ID,
				Status:     Active,
				Expiry:     info.BillingInfo.NextBillingTime.UTC(),
				UpdatedAt:  time.Now().UTC(),
			}

			if err := UpdateStatus(ctx, db, up); err != nil {
				return fmt.Errorf("changing status of subscription with provider_id[%s] to 'active': %w", sale.ID, err)
			}
			return web.Respond(ctx, w, nil, http.StatusOK)

		default:
			return fmt.Errorf("event type[%s] is unknown", evt.EventType)
		}
	}
}

func HandleStripeCheckout(db *sqlx.DB, strp *stripecl.API, cfg config.Stripe) web.Handler {
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

		// We don't know yet the provider id. We should override it in webhook events.
		now := time.Now().UTC()
		sub := Sub{
			ID:         validate.GenerateID(),
			PlanID:     plan.ID,
			UserID:     clm.UserID,
			Provider:   Stripe,
			ProviderID: "ToBeOverridden",
			Status:     Pending,
			CreatedAt:  now,
			UpdatedAt:  now,
			Expiry:     now,
		}

		if err := Create(ctx, db, sub); err != nil {
			return fmt.Errorf("creating subscription: %w", err)
		}

		// Pass the internal subscription id to the stripe subscription object
		// so that we can bind them during webhook events retrieval.
		// A valid alternative is to use a customer_id, storing it in the database.
		params := &stripe.CheckoutSessionParams{
			SuccessURL: stripe.String(cfg.SuccessURL),
			CancelURL:  stripe.String(cfg.CancelURL),
			Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
				Metadata: map[string]string{"internal_id": sub.ID},
			},
			LineItems: []*stripe.CheckoutSessionLineItemParams{{
				Quantity: stripe.Int64(1),
				Price:    &plan.StripeID,
			}},
		}

		// Create a new stripe checkout for subscription.
		s, err := strp.CheckoutSessions.New(params)
		if err != nil {
			return fmt.Errorf("creating stripe session: %w", err)
		}

		http.Redirect(w, r, s.URL, http.StatusSeeOther)
		return nil
	}
}

// TODO: We need to handle recurrent payments that require 3d secure auth.
// A simple solution is to pay stripe to handle that flow automatically.
// More info at (https://stripe.com/docs/billing/migration/strong-customer-authentication).
func HandleStripeWebhook(db *sqlx.DB, strp *stripecl.API, cfg config.Stripe) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("cannot read the request body: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		sig := r.Header.Get("Stripe-Signature")
		if sig == "" {
			err = errors.New("received stripe event is not signed")
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		event, err := webhook.ConstructEvent(b, sig, cfg.WebhookSecret)
		if err != nil {
			err = fmt.Errorf("cannot construct stripe event: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		// Filter all the events but the 'invoice.paid' one.
		// That is when the user pays for a subscription.
		if event.Type != "invoice.paid" {
			return web.Respond(ctx, w, nil, http.StatusOK)
		}

		var invoice stripe.Invoice
		if err = json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			err = fmt.Errorf("unable to decode stripe event: %w", err)
			return weberr.NewError(err, err.Error(), http.StatusBadRequest)
		}

		if invoice.Subscription == nil {
			return fmt.Errorf("paid invoice[%s] doesn't contain a subscription id: %w", invoice.ID, err)
		}

		sub, err := strp.Subscriptions.Get(invoice.Subscription.ID, nil)
		if err != nil {
			return fmt.Errorf("invoice paid but cannot fetch subscription[%s]: %w", invoice.Subscription.ID, err)
		}

		if sub.Status != stripe.SubscriptionStatusActive {
			return fmt.Errorf("invoice paid but subscription[%s] is not active yet", invoice.Subscription.ID)
		}

		internalID, ok := sub.Metadata["internal_id"]
		if !ok {
			return fmt.Errorf("invoice paid but subscription[%s] has no 'internal_id'", invoice.Subscription.ID)
		}

		// Could work but it's ugly. Make a single query to update the database.
		// Alternatively, we could create the subscription directly here.
		if err := UpdateProviderID(ctx, db, internalID, sub.ID); err != nil {
			return fmt.Errorf("binding subscription[%s] with provider_id[%s]: %w", internalID, sub.ID, err)
		}

		up := StatusUp{
			ProviderID: sub.ID,
			Status:     Active,
			Expiry:     time.Unix(sub.CurrentPeriodEnd, 0),
			UpdatedAt:  time.Now().UTC(),
		}

		if err := UpdateStatus(ctx, db, up); err != nil {
			return fmt.Errorf("changing status of subscription with provider_id[%s] to 'active': %w", sub.ID, err)
		}

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}
