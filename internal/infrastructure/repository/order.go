package repository

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/postgres"

	"github.com/jackc/pgx/v5/pgconn"
)

type orderRepository struct {
	*postgres.PgClient
}

func NewOrderRepository(db *postgres.PgClient) *orderRepository {
	return &orderRepository{db}
}

func (or *orderRepository) Create(ctx context.Context, order models.Order) error {
	o := fromDomainOrder(order)

	query := `
			INSERT INTO praktikum.order (number, status, accrual, uploaded_at)
			VALUES ($1, $2, $3, $4)`

	_, err := or.ExecContext(ctx, query, o.Number, o.Status, o.Accrual, o.UploadedAt)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		} else if pgErr.Code == "23505" { // login exists: duplicate key value violates unique constraint
			return ErrExists
		}
		return err
	}

	return nil
}

func (or *orderRepository) GetList(ctx context.Context) ([]models.Order, error) {
	var orders []dsOrder

	query := `
			SELECT number, status, accrual, uploaded_at
			FROM praktikum.orders`

	err := or.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return toDomainOrders(orders), nil
}
