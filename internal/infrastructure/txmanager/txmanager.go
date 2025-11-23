package txmanager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type sqlxTxManager struct {
	db *sqlx.DB
}

func NewSqlx(db *sqlx.DB) TxManager {
	return &sqlxTxManager{db: db}
}

func (m *sqlxTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	// 1. Если транзакция уже есть в контексте — вложенный сценарий.
	if existingTx := TxFromContext(ctx); existingTx != nil {
		return fn(ctx)
	}

	// 2. Транзакции нет — создаём новую.
	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Оборачиваем ctx с транзакцией
	ctxWithTx := withTx(ctx, tx)

	// 3. Коммит/роллбек + обработка паник
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // best effort
			panic(p)
		}

		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	// 4. Выполняем бизнес-логику
	err = fn(ctxWithTx)
	return
}
