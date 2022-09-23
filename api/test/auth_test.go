package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/polldo/govod/core/token"
	"github.com/polldo/govod/core/user"
)

func Signup(srv *httptest.Server, usr user.UserSignup) (user.User, error) {
	body, err := json.Marshal(&usr)
	if err != nil {
		return user.User{}, err
	}

	r, err := http.NewRequest(http.MethodPost, srv.URL+"/auth/signup", bytes.NewBuffer(body))
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

func Activate(srv *httptest.Server, email string, mailer *mockMailer) error {
	body, err := json.Marshal(&struct {
		Email string
		Scope string
	}{
		Email: email,
		Scope: token.ActivationToken,
	})
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, srv.URL+"/tokens", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	w, err := srv.Client().Do(r)
	if err != nil {
		return err
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		return fmt.Errorf("can't activate user: status code %s", w.Status)
	}

	token := mailer.token

	body, err = json.Marshal(&struct{ Token string }{Token: token})
	if err != nil {
		return err
	}

	r, err = http.NewRequest(http.MethodPost, srv.URL+"/tokens/activate", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	w, err = srv.Client().Do(r)
	if err != nil {
		return err
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		return fmt.Errorf("can't activate user: status code %s", w.Status)
	}

	return nil
}

func Login(srv *httptest.Server, email string, pass string) error {
	r, err := http.NewRequest(http.MethodPost, srv.URL+"/auth/login", nil)
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
	at.signupAlreadyExistent(t)
	at.signupNoPasswordConfirm(t)
	at.loginOK(t)
	at.loginWrongPass(t)
	at.loginNotActive(t)
}

func (at *authTest) signupOK(t *testing.T) {
	usr := user.UserSignup{
		Name:            "Paolo Calao",
		Email:           "polldo@test.com",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	got, err := Signup(at.Server, usr)
	if err != nil {
		t.Fatal(err)
	}

	if err := Activate(at.Server, usr.Email, at.Mailer); err != nil {
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

func (at *authTest) signupAlreadyExistent(t *testing.T) {
	usr := user.UserSignup{
		Name:            "Paolo Calao",
		Email:           at.UserEmail,
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	_, err := Signup(at.Server, usr)
	if err == nil {
		t.Fatal("cannot create already existing user")
	}
}

func (at *authTest) signupNoPasswordConfirm(t *testing.T) {
	usr := user.UserSignup{
		Name:     "Rose Nopass",
		Email:    "rose@nopass.it",
		Password: "testpass",
	}

	_, err := Signup(at.Server, usr)
	if err == nil {
		t.Fatal("cannot create user without password confirm field")
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

func (at *authTest) loginNotActive(t *testing.T) {
	usr := user.UserSignup{
		Name:            "Mr Smith",
		Email:           "inactive@test.com",
		Password:        "testpass",
		PasswordConfirm: "testpass",
	}

	got, err := Signup(at.Server, usr)
	if err != nil {
		t.Fatal(err)
	}

	exp := got
	exp.Name = "Mr Smith"
	exp.Email = "inactive@test.com"
	exp.Role = "USER"
	exp.Active = false

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong user payload. Diff: \n%s", diff)
	}

	if err := Login(at.Server, usr.Email, usr.Password); err == nil {
		t.Fatal("inactive users cannot login")
	}
}
