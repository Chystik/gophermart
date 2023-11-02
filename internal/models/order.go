package models

import (
	"strconv"
)

type Status string

const (
	Invalid    Status = "INVALID"
	New        Status = "NEW"
	Registered Status = "REGISTERED"
	Processing Status = "PROCESSING"
	Processed  Status = "PROCESSED"
)

type (
	Order struct {
		Number     string      `json:"number" db:"number"`
		User       string      `db:"user_id"`
		Status     Status      `json:"status" db:"status"`
		Accrual    Money       `json:"accrual,omitempty" db:"accrual,omitempty"`
		UploadedAt RFC3339Time `json:"uploaded_at" db:"uploaded_at"`
	}
)

// ValidLuhnNumber checks the order number using the Luhn algorithm
func (o Order) ValidLuhnNumber() bool {
	luhn, err := strconv.Atoi(o.Number)
	if err != nil {
		return false
	}

	checksum := func(n int) int {
		for i := 0; n > 0; i++ {
			cur := n % 10

			if i%2 == 0 {
				cur = cur * 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			luhn += cur
			n = n / 10
		}
		return luhn % 10
	}

	return (luhn%10+checksum(luhn/10))%10 == 0
}
