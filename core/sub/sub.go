package sub

import "time"

type Status string

const (
	Pending   Status = "pending"
	Active    Status = "active"
	Cancelled Status = "cancelled"
	Expired   Status = "expired"
)

type Provider string

const (
	Paypal Status = "paypal"
	Stripe Status = "stripe"
)

type Plan struct {
	ID               string    `json:"id" db:"plan_id"`
	Name             string    `json:"name" db:"name"`
	Price            int       `json:"price" db:"price"`
	MonthsRecurrence int       `json:"months_recurrence" db:"months_recurrence"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type Sub struct {
	ID        string    `json:"id" db:"subscription_id"`
	PlanID    string    `json:"plan_id" db:"plan_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Provider  Provider  `json:"provider" db:"provider"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Expiry    time.Time `json:"expiry" db:"expiry"`
}

type StatusUp struct {
	ID        string    `db:"subscription_id"`
	Status    Status    `db:"status"`
	Expiry    time.Time `db:"expiry"`
	UpdatedAt time.Time `db:"updated_at"`
}
