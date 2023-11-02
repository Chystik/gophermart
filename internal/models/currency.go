package models

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Money struct {
	Amount Amount
}

type Amount uint64

// Scan - Implement the database/sql scanner interface
func (m *Money) Scan(value interface{}) (err error) {
	f, ok := value.(float64)
	if !ok {
		return fmt.Errorf("cant assert %T to float64", value)
	}
	m.FromFloat(f)
	return
}

// Value - Implementation of valuer for database/sql
func (m Money) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	return m.ToFloat(), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" || s == "0" {
		m.Amount = 0
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	m.FromFloat(f)
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	if m.Amount == 0 {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(m.ToFloat(), 'f', -1, 64)), nil
}

func (m *Money) ToFloat() float64 {
	return float64(m.Amount) / 100
}

func (m *Money) FromFloat(amount float64) {
	currencyDecimals := math.Pow10(2)
	val := uint64(amount * currencyDecimals)

	largeUnit := val / 100
	smallUnit := val % 100

	m.Amount = Amount((largeUnit * 100) + smallUnit)
}

func (m Money) LargeUnit() uint64 {
	return uint64(m.Amount % 100)
}

func (m Money) SmallUnit() uint64 {
	return uint64(m.Amount / 100)
}

func (m Money) String() string {
	var builder strings.Builder

	builder.WriteString(strconv.FormatUint(uint64(m.Amount/100), 10))
	builder.WriteByte('.')
	smallUnit := strconv.FormatUint(uint64(m.Amount%100), 10)
	if len(smallUnit) == 1 {
		builder.WriteByte('0')
	}
	builder.WriteString(smallUnit)

	return builder.String()
}

func (m *Money) Add(money Money) {
	m.Amount += money.Amount
}

func (m *Money) Substract(money Money) {
	m.Amount -= money.Amount
}

func (m Money) LessThan(money Money) bool {
	return m.Amount < money.Amount
}
