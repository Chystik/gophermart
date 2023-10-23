package repository

import (
	"time"

	"github.com/Chystik/gophermart/internal/models"
)

type dsOrder struct {
	Number     int       `db:"number"`
	Status     string    `db:"status"`
	Accrual    float64   `db:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploaded_at"`
}

func fromDomainOrder(o models.Order) dsOrder {
	return dsOrder{
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAt: o.UploadedAt,
	}
}

func toDomainOrder(o dsOrder) models.Order {
	return models.Order{
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAt: o.UploadedAt,
	}
}

func toDomainOrders(o []dsOrder) []models.Order {
	orders := make([]models.Order, 0, len(o))
	for _, order := range o {
		orders = append(orders, toDomainOrder(order))
	}
	return orders
}
