package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/polldo/govod/core/token"
	"github.com/polldo/govod/core/user"
)

type tokenTest struct {
	*TestEnv
}

func TestToken(t *testing.T) {
	env, err := NewTestEnv(t, "token_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	tt := &tokenTest{env}

	au := tt.signupTestUser(t, "mary.lu@activation.com")
	tt.notActiveOK(t, au)
	tt.activationToken(t, au)

	ru := tt.signupTestUser(t, "mary.lu@recovery.com")
	tt.recoveryToken(t, ru)
}

func (tt *tokenTest) signupTestUser(t *testing.T, email string) user.UserSignup {
	u := user.UserSignup{
		Name:            "Mary Lu",
		Email:           email,
		Password:        "marysecret",
		PasswordConfirm: "marysecret",
	}

	if _, err := Signup(tt.Server, u); err != nil {
		t.Fatal(err)
	}

	return u
}

func (tt *tokenTest) notActiveOK(t *testing.T, u user.UserSignup) {
	if err := Login(tt.Server, u.Email, u.Password); err == nil {
		t.Fatal("user should not be active at this point")
	}
}

func (tt *tokenTest) requestToken(t *testing.T, email string, scope string) string {
	body, err := json.Marshal(&struct {
		Email string
		Scope string
	}{
		Email: email,
		Scope: scope,
	})

	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, tt.Server.URL+"/tokens", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := tt.Server.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusNoContent {
		t.Fatalf("can't activate user: status code %s", w.Status)
	}

	return tt.Mailer.token
}

func (tt *tokenTest) activationToken(t *testing.T, u user.UserSignup) {
	// First token request.
	tokst := tt.requestToken(t, u.Email, token.ActivationToken)

	// Second token request.
	toknd := tt.requestToken(t, u.Email, token.ActivationToken)

	// Request a third token but with different scope.
	_ = tt.requestToken(t, u.Email, token.RecoveryToken)

	// Check if the first token has been invalidated (deleted).
	body, err := json.Marshal(&struct{ Token string }{Token: tokst})
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, tt.Server.URL+"/tokens/activate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := tt.Server.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode == http.StatusNoContent {
		t.Fatalf("first token shouldn't have been valid")
	}

	if err := Login(tt.Server, u.Email, u.Password); err == nil {
		t.Fatal("user should not be active yet")
	}

	// Check if the second token is still valid.
	body, err = json.Marshal(&struct{ Token string }{Token: toknd})
	if err != nil {
		t.Fatal(err)
	}

	r, err = http.NewRequest(http.MethodPost, tt.Server.URL+"/tokens/activate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err = tt.Server.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusNoContent {
		t.Fatalf("second token should have been valid")
	}

	if err := Login(tt.Server, u.Email, u.Password); err != nil {
		t.Fatal("user should be active at this point")
	}
}

func (tt *tokenTest) recoveryToken(t *testing.T, u user.UserSignup) {
	newPassword := "marys-new-secret"

	// First token request.
	tokst := tt.requestToken(t, u.Email, token.RecoveryToken)

	// Second token request.
	toknd := tt.requestToken(t, u.Email, token.RecoveryToken)

	// Request a third token but with different scope.
	_ = tt.requestToken(t, u.Email, token.ActivationToken)

	if err := Activate(tt.Server, u.Email, tt.Mailer); err != nil {
		t.Fatal(err)
	}

	// Check if the first token has been invalidated (deleted).
	body, err := json.Marshal(&struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}{
		Token:           tokst,
		Password:        newPassword,
		PasswordConfirm: newPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, tt.Server.URL+"/tokens/recover", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := tt.Server.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode == http.StatusOK {
		t.Fatalf("first token shouldn't have been valid")
	}

	if err := Login(tt.Server, u.Email, newPassword); err == nil {
		t.Fatal("user should have the old password")
	}

	// Check if the second token is still valid.
	body, err = json.Marshal(&struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}{
		Token:           toknd,
		Password:        newPassword,
		PasswordConfirm: newPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	r, err = http.NewRequest(http.MethodPost, tt.Server.URL+"/tokens/recover", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err = tt.Server.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("second token should have been valid")
	}

	if err := Login(tt.Server, u.Email, newPassword); err != nil {
		t.Fatal("user should have the new password at this point")
	}
}
