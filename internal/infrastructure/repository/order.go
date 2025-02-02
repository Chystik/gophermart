package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/Chystik/gophermart/pkg/transaction"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type orderRepository struct {
	*sqlx.DB
	getter transaction.CtxGetter
	logger logger.AppLogger
}

func NewOrderRepository(db *sqlx.DB, getter transaction.CtxGetter, logger logger.AppLogger) *orderRepository {
	return &orderRepository{
		DB:     db,
		getter: getter,
		logger: logger,
	}
}

func (or *orderRepository) Create(ctx context.Context, order models.Order) error {
	query := `
			INSERT INTO praktikum.order (number, user_id, status, accrual, uploaded_at)
			VALUES (:number, :user_id, :status, :accrual, :uploaded_at)
			ON CONFLICT (number) DO NOTHING`

	or.logger.Debug("OrderRepository.Create", zap.String("query", query))

	_, err := sqlx.NamedExecContext(ctx, or.getter.GetTrxOrDB(ctx, or.DB), query, order)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		} else if pgErr.Code == "23505" { // login exists: duplicate key value violates unique constraint
			return &models.AppError{Op: "orderRepository.Create", Code: models.ErrExists, Message: "order already exists"}
		}
	}

	return nil
}

func (or *orderRepository) Get(ctx context.Context, order models.Order) (models.Order, error) {
	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE number = $1`

	or.logger.Debug("OrderRepository.Get", zap.String("query", query))

	err := sqlx.GetContext(ctx, or.getter.GetTrxOrDB(ctx, or.DB), &order, query, order.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, &models.AppError{Op: "orderRepository.Get", Code: models.ErrNotFound, Message: "order not found"}
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

	or.logger.Debug("OrderRepository.GetList", zap.String("query", query))

	rows, err := sqlx.NamedQueryContext(ctx, or.getter.GetTrxOrDB(ctx, or.DB), query, login)
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

	or.logger.Debug("OrderRepository.Update", zap.String("query", query))

	_, err := sqlx.NamedExecContext(ctx, or.getter.GetTrxOrDB(ctx, or.DB), query, order)

	return err
}

func (or *orderRepository) GetUnprocessed(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	stat := []string{"PROCESSING", "NEW", "REGISTERED"} // ignore "PROCESSED", "INVALID"

	query := `
			SELECT number, user_id, status, accrual, uploaded_at
			FROM praktikum.order
			WHERE status = ANY ($1)`

	or.logger.Debug("OrderRepository.GetUnprocessed", zap.String("query", query))

	err := sqlx.SelectContext(ctx, or.getter.GetTrxOrDB(ctx, or.DB), &orders, query, stat)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
