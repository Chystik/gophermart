package transaction

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type ctxGetter struct {
	txKey txKey
}

func NewCtxGetter() *ctxGetter {
	return &ctxGetter{txKey: txKey{}}
}

// GetTrxOrDB gets and returns sqlx.Tx from the ctx context.Context,
// or returns sqlx.ExtContext on failure.
func (g *ctxGetter) GetTrxOrDB(ctx context.Context, db sqlx.ExtContext) sqlx.ExtContext {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}

	return db
}

type CtxGetter interface {
	GetTrxOrDB(ctx context.Context, db sqlx.ExtContext) sqlx.ExtContext
}
