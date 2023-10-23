package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/postgres"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUploadedByUser        = errors.New("object already uploaded by user")
	ErrUploadedByAnotherUser = errors.New("object already uploaded by another user")
)

type orderRepository struct {
	*postgres.PgClient
}

func NewOrderRepository(db *postgres.PgClient) *orderRepository {
	return &orderRepository{db}
}

func (or *orderRepository) Create(ctx context.Context, order models.Order) error {
	query := `
			INSERT INTO praktikum.order (number, user_id, status, accrual, uploaded_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (number) DO NOTHING`

	_, err := or.ExecContext(ctx, query, order.Number, order.User, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		} else if pgErr.Code == "23505" { // login exists: duplicate key value violates unique constraint
			return ErrExists
		}
	}

	return nil
}

func (or *orderRepository) Get(ctx context.Context, order models.Order) (models.Order, error) {
	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE number = $1`

	err := or.GetContext(ctx, &order, query, order.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, ErrNotFound
		} else {
			return order, err
		}
	}

	return order, nil
}

func (or *orderRepository) GetList(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			ORDER BY uploaded_at ASC`

	err := or.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
