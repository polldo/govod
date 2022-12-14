package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/core/course"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

type orderTest struct {
	*TestEnv
}

func TestOrder(t *testing.T) {
	env, err := NewTestEnv(t, "order_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	ot := &orderTest{env}
	ct := &courseTest{env}
	rt := &cartTest{env}

	// Prepping the env.
	c1 := ct.createCourseOK(t)
	c2 := ct.createCourseOK(t)
	_ = ct.createCourseOK(t)
	_ = ct.createCourseOK(t)
	c3 := ct.createCourseOK(t)
	c4 := ct.createCourseOK(t)

	// Initially the user doesn't own any course.
	ct.listCoursesOwnedOK(t, nil)

	// Add courses to the cart.
	rt.createItemOK(t, c1.ID)
	rt.createItemOK(t, c2.ID)

	// Perform a paypal payment.
	ot.Paypal.expectedCart = []course.Course{c1, c2}
	ot.testPaypal(t)

	// Check if the paypal payment has been correctly fulfilled.
	ct.listCoursesOwnedOK(t, []course.Course{c1, c2})

	// Add new courses to the cart.
	rt.createItemOK(t, c3.ID)
	rt.createItemOK(t, c4.ID)

	// Perform a stripe payment.
	ot.Stripe.expectedCart = []course.Course{c3, c4}
	ot.testStripe(t)

	// Check if the stripe payment has been correctly fulfilled.
	ct.listCoursesOwnedOK(t, []course.Course{c1, c2, c3, c4})
}

func (ot *orderTest) testPaypal(t *testing.T) {
	if err := Login(ot.Server, ot.UserEmail, ot.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ot.Server)

	// Checkout the order via paypal.
	r, err := http.NewRequest(http.MethodPost, ot.URL+"/orders/paypal", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ot.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't create paypal order: status code %s", w.Status)
	}

	var ord paypal.Order
	if err := json.NewDecoder(w.Body).Decode(&ord); err != nil {
		t.Fatalf("cannot unmarshal paypal order: %v", err)
	}

	// Capture the paypal order.
	// We are using a mocked paypal server that returns OK on captures
	// to simulate the user payment happy path.
	r, err = http.NewRequest(http.MethodPost, ot.URL+"/orders/paypal/"+ord.ID+"/capture", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err = ot.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't capture paypal order: status code %s", w.Status)
	}
}

func (ot *orderTest) testStripe(t *testing.T) {
	if err := Login(ot.Server, ot.UserEmail, ot.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ot.Server)

	// Checkout the order via stripe.
	r, err := http.NewRequest(http.MethodPost, ot.URL+"/orders/stripe", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ot.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusSeeOther {
		t.Fatalf("can't create stripe order: status code %s", w.Status)
	}

	// Now simulate the payment by triggering a stripe webhook.
	//
	// Extract the checkout session id from the location returned by the mock.
	u, err := w.Location()
	if err != nil {
		t.Fatal(err)
	}

	// Generate the webhook payload.
	//
	// Set the same checkout id previously obtained.
	obj := map[string]any{
		// Mocked stripe returns the id in the URL.
		"id": path.Base(u.Path),
	}

	raw, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	// evt is the complete payload for the webhook.
	evt := stripe.Event{
		// Required by stripe-go 74.2.0 .
		APIVersion: "2022-11-15",
		Type:       "checkout.session.completed",
		Data: &stripe.EventData{
			Raw: json.RawMessage(raw),
		},
	}

	b, err := json.Marshal(evt)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the payload with the appropriate secret.
	signed := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload:   b,
		Secret:    ot.WebhookSecret,
		Timestamp: time.Now(),
	})

	// Finally trigger the webhook.
	r, err = http.NewRequest(http.MethodPost, ot.URL+"/orders/stripe/capture", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Stripe-Signature", signed.Header)

	w, err = ot.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't trigger stripe webhook: status code %s", w.Status)
	}
}
