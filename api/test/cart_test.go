package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/polldo/govod/core/cart"
)

type cartTest struct {
	*TestEnv
}

func TestCart(t *testing.T) {
	env, err := NewTestEnv(t, "cart_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	// Auxiliary test struct to populate the db.
	// Its advantage with respect to using a seed
	// is that it's more realistic and closer to an integration test.
	ut := &courseTest{env}
	course1 := ut.createCourseOK(t)
	course2 := ut.createCourseOK(t)

	ct := &cartTest{env}

	// The cart should be initially empty.
	ct.showCartOK(t, cart.Cart{Items: []cart.Item{}})

	// Add two items and check if they're added.
	item1 := ct.createItemOK(t, course1.ID)
	item2 := ct.createItemOK(t, course2.ID)
	ct.showCartOK(t, cart.Cart{
		Items: []cart.Item{item1, item2},
	})

	// Flush the cart and check that it's empty.
	ct.deleteCartOK(t)
	ct.showCartOK(t, cart.Cart{Items: []cart.Item{}})

	// Deletion should be idempotent.
	ct.deleteCartOK(t)

	// Add items to the cart, then delete them and check that it's empty.
	ct.createItemOK(t, course1.ID)
	ct.createItemOK(t, course2.ID)
	ct.deleteItemOK(t, item1.CourseID)
	ct.deleteItemOK(t, item2.CourseID)
	ct.showCartOK(t, cart.Cart{Items: []cart.Item{}})
}

func (ct *cartTest) createItemOK(t *testing.T, courseID string) cart.Item {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := cart.ItemNew{
		CourseID: courseID,
	}

	body, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPut, ct.URL+"/cart/items", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusCreated {
		t.Fatalf("can't create cart item: status code %s", w.Status)
	}

	var got cart.Item
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal created cart item: %v", err)
	}

	exp := got
	exp.CourseID = c.CourseID

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong course payload. Diff: \n%s", diff)
	}

	return got
}

func (ct *cartTest) deleteItemOK(t *testing.T, courseID string) {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	r, err := http.NewRequest(http.MethodDelete, ct.URL+"/cart/items/"+courseID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusNoContent {
		t.Fatalf("can't delete cart item: status code %s", w.Status)
	}
}

func (ct *cartTest) deleteCartOK(t *testing.T) {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	r, err := http.NewRequest(http.MethodDelete, ct.URL+"/cart", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusNoContent {
		t.Fatalf("can't delete cart: status code %s", w.Status)
	}
}

func (ct *cartTest) showCartOK(t *testing.T, exp cart.Cart) {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	r, err := http.NewRequest(http.MethodGet, ct.URL+"/cart", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't show cart: status code %s", w.Status)
	}

	var got cart.Cart
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal retrievede cart: %v", err)
	}

	// These are don't care fields.
	exp.UpdatedAt = got.UpdatedAt
	exp.CreatedAt = got.CreatedAt

	// Don't care about dates.
	now := got.UpdatedAt
	cmp.Transformer("", func(in cart.Item) cart.Item {
		out := in
		out.CreatedAt = now
		out.UpdatedAt = now
		return out
	})

	less := func(a, b cart.Item) bool { return a.CourseID < b.CourseID }
	if diff := cmp.Diff(got, exp, cmpopts.SortSlices(less)); diff != "" {
		t.Fatalf("wrong cart payload. Diff: \n%s", diff)
	}
}
