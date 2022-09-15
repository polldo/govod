package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/polldo/govod/core/user"
)

type userTest struct {
	srv *httptest.Server
}

func TestUser(t *testing.T) {
	env, err := NewTestEnv(t, "user_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	ut := &userTest{srv: env.Server}

	usr, err := Signup(ut.srv, user.UserNew{
		Name:            "Paolo Calao",
		Email:           "polldo@test.com",
		Password:        "pass",
		PasswordConfirm: "pass",
		Role:            "USER",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := Login(ut.srv, "polldo@test.com", "pass"); err != nil {
		t.Fatal(err)
	}

	ut.getUserOK(t, usr.ID)
}

func (ut *userTest) getUserOK(t *testing.T, id string) user.User {
	r, err := http.NewRequest(http.MethodGet, ut.srv.URL+"/users/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ut.srv.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}

	var got user.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal created user: %v", err)
	}

	exp := got
	exp.Name = "Paolo Calao"
	exp.Email = "polldo@test.com"
	exp.Role = "USER"

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong user payload. Diff: \n%s", diff)
	}

	return got
}
