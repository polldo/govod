package order

import "time"

type Order struct {
	ID         string    `json:"id" db:"order_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	Status     string    `json:"status" db:"status"`
	Amount     float64   `json:"amount" db:"amount"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StatusUp struct {
	ID        string    `db:"order_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Item struct {
	OrderID   string    `json:"order_id" db:"order_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
