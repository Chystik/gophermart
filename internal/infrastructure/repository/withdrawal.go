package repository

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
)

type withdrawalRepository struct {
	*sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewWithdrawalRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) *withdrawalRepository {
	return &withdrawalRepository{
		DB:     db,
		getter: getter,
	}
}

func (wr *withdrawalRepository) Create(ctx context.Context, w models.Withdrawal) error {
	query := `
			INSERT INTO	praktikum.withdrawal (order_id, sum, processed_at)
			VALUES (:order_id, :sum, :processed_at)`

	_, err := sqlx.NamedExecContext(ctx, wr.getter.DefaultTrOrDB(ctx, wr.DB), query, w)
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

	err := sqlx.SelectContext(ctx, wr.getter.DefaultTrOrDB(ctx, wr.DB), &w, query)
	if err != nil {
		return nil, err
	}

	return w, nil
}
