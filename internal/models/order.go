package models

import "time"

type Order struct {
	Number     int
	User       string
	Status     string
	Accrual    float64
	UploadedAt time.Time
}
