package user

import (
	"time"
)

// User models users. Email address is a unique field.
type User struct {
	ID           string    `json:"id" db:"user_id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	Active       bool      `json:"active" db:"active"`
	PasswordHash []byte    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	Version      int       `json:"-" db:"version"`
}

// UserNew includes the information an administrator needs
// to create a new user.
type UserNew struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Role            string `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"eqfield=Password"`
}

// UserSignup contains information needed to register a user.
type UserSignup struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,gte=8,lte=50"`
	PasswordConfirm string `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}

// UserUp specifies information of a user which can be updated.
type UserUp struct {
	Name            *string `json:"name"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Role            *string `json:"role"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}
