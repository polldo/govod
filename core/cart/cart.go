package cart

import (
	"time"
)

// Cart is the model of a user's cart.
// A cart is very similar to an order, but it's modeled a separate
// entity because they're used in different ways.
// A cart can be always updated safely, while an order is more complex to handle
// because data races during payments must be considered.
//
// Also, do we really need a cart ? we already have items! Yes because it will contain coupons and discounts.
type Cart struct {
	UserID    string    `json:"-" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Version   int       `json:"-" db:"version"`
	Items     []Item    `json:"items" db:"-"`
}

type Item struct {
	UserID    string    `json:"-" db:"user_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ItemNew struct {
	CourseID string `json:"course_id" db:"course_id"`
}
