package order

import "time"

type Order struct {
	ID        string    `json:"id" db:"order_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Item struct {
	// ID        string    `json:"item_id" db:"item_id"`
	OrderID   string    `json:"order_id" db:"order_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
