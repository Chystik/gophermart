package repository

import (
	"context"

	"github.com/Chystik/gophermart/internal/models"

	"github.com/jmoiron/sqlx"
)

type withdrawalRepository struct {
	*sqlx.DB
}

func NewWithdrawalRepository(db *sqlx.DB) *withdrawalRepository {
	return &withdrawalRepository{db}
}

func (wr *withdrawalRepository) Create(ctx context.Context, w models.Withdrawal) error {
	query := `
			INSERT INTO	praktikum.withdrawal (order_id, sum, processed_at)
			VALUES (:order_id, :sum, :processed_at)`

	_, err := sqlx.NamedExecContext(ctx, wr, query, w)
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

	err := sqlx.SelectContext(ctx, wr, &w, query)
	if err != nil {
		return nil, err
	}

	return w, nil
}
