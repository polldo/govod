package order

import "time"

type Status string

const (
	Pending Status = "pending"
	Success Status = "success"
	Expired Status = "expired"
)

type Order struct {
	ID         string    `json:"id" db:"order_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	Status     Status    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StatusUp struct {
	ID        string    `db:"order_id"`
	Status    Status    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Item struct {
	OrderID   string    `json:"order_id" db:"order_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Price     int       `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
