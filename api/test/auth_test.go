package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/polldo/govod/core/user"
)

func Signup(srv *httptest.Server, usr user.UserNew) (user.User, error) {
	body, err := json.Marshal(&usr)
	if err != nil {
		return user.User{}, err
	}

	r, err := http.NewRequest(http.MethodPost, srv.URL+"/signup", bytes.NewBuffer(body))
	if err != nil {
		return user.User{}, err
	}

	w, err := srv.Client().Do(r)
	if err != nil {
		return user.User{}, err
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusCreated {
		return user.User{}, fmt.Errorf("can't signup: status code %s", w.Status)
	}

	var got user.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		return user.User{}, fmt.Errorf("cannot unmarshal created user: %v", err)
	}

	return got, nil
}

func Login(srv *httptest.Server, email string, pass string) error {
	r, err := http.NewRequest(http.MethodPost, srv.URL+"/login", nil)
	if err != nil {
		return err
	}

	r.SetBasicAuth(email, pass)

	w, err := srv.Client().Do(r)
	if err != nil {
		return err
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		return fmt.Errorf("can't login: status code %s", w.Status)
	}

	return nil
}

type authTest struct {
	*TestEnv
}

func TestAuth(t *testing.T) {
	env, err := NewTestEnv(t, "auth_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	at := &authTest{env}

	at.signupOK(t)
	at.signupUnauth(t)
	at.signupAlreadyExistent(t)
	at.loginOK(t)
	at.loginWrongPass(t)
}

func (at *authTest) signupOK(t *testing.T) {
	usr := user.UserNew{
		Name:            "Paolo Calao",
		Email:           "polldo@test.com",
		Role:            "USER",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	got, err := Signup(at.Server, usr)
	if err != nil {
		t.Fatal(err)
	}

	exp := got
	exp.Name = "Paolo Calao"
	exp.Email = "polldo@test.com"
	exp.Role = "USER"

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong user payload. Diff: \n%s", diff)
	}
}

func (at *authTest) signupUnauth(t *testing.T) {
	usr := user.UserNew{
		Name:            "Paolo Calao",
		Email:           "polldo@test.com",
		Role:            "ADMIN",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	_, err := Signup(at.Server, usr)
	if err == nil {
		t.Fatal("simple users cannot create admin")
	}
}

func (at *authTest) signupAlreadyExistent(t *testing.T) {
	usr := user.UserNew{
		Name:            "Paolo Calao",
		Email:           at.UserEmail,
		Role:            "ADMIN",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	_, err := Signup(at.Server, usr)
	if err == nil {
		t.Fatal("cannot create already existing user")
	}
}

func (at *authTest) loginOK(t *testing.T) {
	if err := Login(at.Server, "polldo@test.com", "testpass"); err != nil {
		t.Fatalf("login failed: %v", err)
	}
}

func (at *authTest) loginWrongPass(t *testing.T) {
	if err := Login(at.Server, "polldo@test.com", "wrongpass"); err == nil {
		t.Fatal("login should have failed")
	}
}
