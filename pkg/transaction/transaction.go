package transaction

import (
	"context"
	"database/sql"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn TxFn) error
}

type TxFn func(ctx context.Context) error

type manager struct {
	*sqlx.DB
	logger logger.AppLogger
}

func NewManager(db *sqlx.DB, l logger.AppLogger) *manager {
	return &manager{
		DB:     db,
		logger: l,
	}
}

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`.
func (m *manager) WithTransaction(ctx context.Context, fn TxFn) error {
	//logger := ctxval.GetLogger(ctx)
	tx, err := m.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		m.logger.Error(err.Error())
		return &models.AppError{Op: "transaction.managerWithTransaction", Message: "cannot begin transaction", Err: err}
	}

	m.logger.Debug("Starting database transaction")
	err = fn(injectTx(ctx, tx))
	if err != nil {
		m.logger.Debug("Rolling database transaction back")
		tx.Rollback()
		return err
	}

	m.logger.Debug("Committing database transaction")
	return tx.Commit()
}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
