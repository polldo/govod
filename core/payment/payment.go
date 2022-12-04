package payment

import "time"

type Payment struct {
	ID         string    `json:"id" db:"payment_id"`
	OrderID    string    `json:"order_id" db:"order_id"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	Status     string    `json:"status" db:"status"`
	Amount     float64   `json:"amount" db:"amount"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StatusUp struct {
	ID        string    `db:"payment_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}
