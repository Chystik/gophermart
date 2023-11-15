package repository

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/Chystik/gophermart/pkg/transaction"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type withdrawalRepository struct {
	*sqlx.DB
	getter transaction.CtxGetter
	logger logger.AppLogger
}

func NewWithdrawalRepository(db *sqlx.DB, getter transaction.CtxGetter, logger logger.AppLogger) *withdrawalRepository {
	return &withdrawalRepository{
		DB:     db,
		getter: getter,
		logger: logger,
	}
}

func (wr *withdrawalRepository) Create(ctx context.Context, w models.Withdrawal) error {
	query := `
			INSERT INTO	praktikum.withdrawal (order_id, sum, processed_at)
			VALUES (:order_id, :sum, :processed_at)`

	wr.logger.Debug("WithdrawalRepository.Create", zap.String("query", query))

	_, err := sqlx.NamedExecContext(ctx, wr.getter.GetTrxOrDB(ctx, wr.DB), query, w)
	if err != nil {
		return err
	}

	return nil
}

func (wr *withdrawalRepository) GetAll(ctx context.Context) ([]models.Withdrawal, error) {
	var w []models.Withdrawal

	query := `
			SELECT order_id, sum, processed_at
			FROM praktikum.withdrawal
			ORDER BY processed_at ASC`

	wr.logger.Debug("WithdrawalRepository.GetAll", zap.String("query", query))

	err := sqlx.SelectContext(ctx, wr.getter.GetTrxOrDB(ctx, wr.DB), &w, query)
	if err != nil {
		return nil, err
	}

	return w, nil
}
