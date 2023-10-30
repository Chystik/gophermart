package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUploadedByUser        = errors.New("object already uploaded by user")
	ErrUploadedByAnotherUser = errors.New("object already uploaded by another user")
)

type orderRepository struct {
	*sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{db}
}

func (or *orderRepository) Create(ctx context.Context, order models.Order) error {
	query := `
			INSERT INTO praktikum.order (number, user_id, status, accrual, uploaded_at)
			VALUES (:number, :user_id, :status, :accrual, :uploaded_at)
			ON CONFLICT (number) DO NOTHING`

	_, err := sqlx.NamedExecContext(ctx, or, query, order)
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

	err := sqlx.GetContext(ctx, or, &order, query, order.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, ErrNotFound
		} else {
			return order, err
		}
	}

	return order, nil
}

func (or *orderRepository) GetList(ctx context.Context, login models.User) ([]models.Order, error) {
	var orders []models.Order

	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE user_id = :login
			ORDER BY uploaded_at ASC`

	rows, err := sqlx.NamedQueryContext(ctx, or, query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.StructScan(&order)
		if err != nil {
			return orders, rows.Err()
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (or *orderRepository) Update(ctx context.Context, order models.Order) error {
	query := `
			UPDATE praktikum.order 
			SET status = :status, accrual = :accrual
			WHERE number = :number`

	_, err := sqlx.NamedExecContext(ctx, or, query, order)

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

	err := sqlx.SelectContext(ctx, or, &orders, query, "PROCESSING", "NEW", "REGISTERED")
	if err != nil {
		return nil, err
	}

	return orders, nil
}
