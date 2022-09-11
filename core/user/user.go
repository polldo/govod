package user

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	PasswordHash []byte    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserNew struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Role            string `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

type UserUp struct {
	Name            *string `json:"name"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Role            *string `json:"role"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}
