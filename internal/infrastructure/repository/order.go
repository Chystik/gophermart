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

	_, err := or.ExecContext(ctx, query, order.Number, order.User, order.Status, order.Accrual, order.UploadedAt.Time)
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

func (or *orderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	var claimsKey models.ClaimsKey = "props"

	claims, ok := ctx.Value(claimsKey).(*models.AuthClaims)
	if !ok {
		return orders, errors.New("bad login")
	}

	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE user_id = $1
			ORDER BY uploaded_at ASC`

	err := or.SelectContext(ctx, &orders, query, claims.Login)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *orderRepository) Update(ctx context.Context, order models.Order) error {
	query := `
			UPDATE praktikum.order 
			SET status = $1, accrual = $2
			WHERE number = $3`

	_, err := or.ExecContext(ctx, query, order.Status, order.Accrual, order.Number)

	return err
}

func (or *orderRepository) GetUnprocessed(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE status = $1
			OR status = $2
			OR status = $3
			ORDER BY uploaded_at ASC`

	err := or.SelectContext(ctx, &orders, query, "PROCESSING", "NEW", "REGISTERED")
	if err != nil {
		return nil, err
	}

	return orders, nil
}
