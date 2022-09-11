package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/polldo/govod/api"
	"github.com/polldo/govod/core/user"
	"github.com/sirupsen/logrus"
)

type userTest struct {
	api http.Handler
}

func TestUser(t *testing.T) {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	api := api.APIMux(api.APIConfig{
		Log: log,
		DB:  db,
	})

	ut := &userTest{api: api}

	u := ut.postUserOK(t)
	ut.getUserOK(t, u.ID)
}

func (ut *userTest) postUserOK(t *testing.T) user.User {
	body, err := json.Marshal(&user.UserNew{
		Name:            "Paolo Calao",
		Email:           "polldo@test.com",
		Role:            "USER",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	ut.api.ServeHTTP(w, r)

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

func (ut *userTest) getUserOK(t *testing.T, id string) user.User {
	r := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	w := httptest.NewRecorder()
	ut.api.ServeHTTP(w, r)

	fmt.Printf("%+v", w)

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
