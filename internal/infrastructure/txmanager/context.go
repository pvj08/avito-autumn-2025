package txmanager

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// приватный тип ключа, чтобы никто снаружи не мог случайно пересечься по ключу
type txCtxKey struct{}

// кладём *sqlx.Tx в контекст
func withTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

// достаём *sqlx.Tx из контекста (если есть)
func TxFromContext(ctx context.Context) *sqlx.Tx {
	val := ctx.Value(txCtxKey{})
	if tx, ok := val.(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
