package cart

import (
	"time"
)

// Cart models the users' carts.
// Each user can have only a cart at a time.
type Cart struct {
	UserID    string    `json:"-" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Version   int       `json:"-" db:"version"`
	Items     []Item    `json:"items" db:"-"`
}

// Item models the item of a cart.
// A cart can have many items.
type Item struct {
	UserID    string    `json:"-" db:"user_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ItemNew models the data required to insert a
// new item on the user's cart.
type ItemNew struct {
	CourseID string `json:"course_id" db:"course_id"`
}
