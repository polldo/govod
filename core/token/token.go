package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ActivationToken = "activation"
	RecoveryToken   = "recovery"
)

// Token models tokens to be sent to users for
// activation and recovery purposes.
type Token struct {
	Hash   []byte    `json:"-" db:"hash"`
	UserID string    `json:"userId" db:"user_id"`
	Expiry time.Time `json:"expiry" db:"expiry"`
	Scope  string    `json:"scope" db:"scope"`
}

// GenToken generates a new random token for a user.
func GenToken(userID string, ttl time.Duration, scope string) (string, Token, error) {
	token := Token{
		UserID: userID,
		Expiry: time.Now().UTC().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", Token{}, err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plaintext))
	token.Hash = hash[:]

	return plaintext, token, nil
}
