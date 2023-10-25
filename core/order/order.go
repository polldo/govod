package order

import "time"

// Status models the possible states of an order.
type Status string

const (
	Pending Status = "pending"
	Success Status = "success"
	Expired Status = "expired"
)

// Order models orders.
// Orders have a one-to-many relationship with items.
type Order struct {
	ID         string    `json:"id" db:"order_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	Status     Status    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// StatusUp contains the information needed to update an order.
type StatusUp struct {
	ID        string    `db:"order_id"`
	Status    Status    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Item models the item of an order.
// An item can only belong to one order.
// An order can have many items.
type Item struct {
	OrderID   string    `json:"order_id" db:"order_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Price     int       `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
