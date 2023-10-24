package models

import (
	"fmt"
	"strings"
	"time"
)

type (
	Order struct {
		Number     string      `json:"number" db:"number"`
		User       string      `db:"user_id"`
		Status     string      `json:"status" db:"status"`
		Accrual    float64     `json:"accrual,omitempty" db:"accrual,omitempty"`
		UploadedAt RFC3339Time `json:"uploaded_at" db:"uploaded_at"`
	}

	RFC3339Time struct {
		time.Time
	}
)

func (f *RFC3339Time) Scan(value interface{}) (err error) {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cant assert %T to time.Time", value)
	}
	f.Time = t
	return
}

func (f *RFC3339Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		f.Time = time.Time{}
		return
	}
	f.Time, err = time.Parse(time.RFC3339, s)
	return
}

func (f RFC3339Time) MarshalJSON() ([]byte, error) {
	if f.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", f.Time.Format(time.RFC3339))), nil
}
