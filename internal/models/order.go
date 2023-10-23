package models

import "time"

type Order struct {
	Number     string    `json:"number" db:"number"`
	User       string    `db:"user_id"`
	Status     string    `json:"status" db:"status"`
	Accrual    float64   `json:"accrual,omitempty" db:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
