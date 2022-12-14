package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/core/course"
	mock "github.com/stripe/stripe-mock/param"
)

type mockPaypal struct {
	expectedCart []course.Course
}

func (m *mockPaypal) handle() http.Handler {
	checkout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Take the payload to perform some basic checks on the items received.
		var pu struct {
			Units []paypal.PurchaseUnitRequest `json:"purchase_units"`
		}
		if err := json.NewDecoder(r.Body).Decode(&pu); err != nil {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		if len(pu.Units) != 1 {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		// Check the number of items, should equals the length of the cart.
		if len(pu.Units[0].Items) != len(m.expectedCart) {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		var tot int
		for _, c := range m.expectedCart {
			tot += c.Price
		}

		// Check the paypal amount against the total of the cart.
		if pu.Units[0].Amount.Value != strconv.Itoa(tot) {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		// Generate a random provider-id that will be used to capture this order.
		randID := fmt.Sprintf("paypal-%d", rand.Intn(300))
		ord := paypal.Order{ID: randID}
		web.Respond(context.Background(), w, ord, 200)
	})

	capture := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Happy path: the paypal payment was completed.
		ord := paypal.Order{Status: "COMPLETED"}
		web.Respond(context.Background(), w, ord, 200)
	})

	r := mux.NewRouter()
	r.Handle("/v2/checkout/orders", checkout).Methods("POST")
	r.Handle("/v2/checkout/orders/{id}/capture", capture).Methods("POST")
	return r
}

type mockStripe struct {
	expectedCart []course.Course
}

func (m *mockStripe) handle() http.Handler {
	checkout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract the stripe payload to check its fields.
		params, _ := mock.ParseParams(r)
		lines := params["line_items"].(map[string]any)

		n := 0
		tot := 0
		for _, li := range lines {
			it := li.(map[string]any)

			if it["quantity"] != "1" {
				web.Respond(context.Background(), w, nil, 400)
				return
			}

			pd := it["price_data"].(map[string]any)
			// Stripe expresses prices in cents.
			s := pd["unit_amount"].(string)
			amount, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				web.Respond(context.Background(), w, err, 400)
				return
			}

			n += 1
			tot += int(amount / 100)
		}

		// Check the number of items against the cart.
		if n != len(m.expectedCart) {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		exp := 0
		for _, c := range m.expectedCart {
			exp += c.Price
		}

		// Check the total amount against the cart.
		if tot != exp {
			web.Respond(context.Background(), w, nil, 400)
			return
		}

		// Generate a random provider-id that will be used to capture this order.
		randID := fmt.Sprintf("stripe-%d", rand.Intn(300))
		ord := map[string]any{"ID": randID, "URL": randID}
		web.Respond(context.Background(), w, ord, 201)
	})

	r := mux.NewRouter()
	r.Handle("/v1/checkout/sessions", checkout).Methods("POST")
	return r
}
