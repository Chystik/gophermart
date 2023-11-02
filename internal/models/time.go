package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type RFC3339Time struct {
	time.Time
}

// Scan - Implement the database/sql scanner interface
func (f *RFC3339Time) Scan(value interface{}) (err error) {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cant assert %T to time.Time", value)
	}
	f.Time = t
	return
}

// Value - Implementation of valuer for database/sql
func (f RFC3339Time) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	return f.Time, nil
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
