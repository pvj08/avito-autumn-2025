package txmanager

import (
	"context"
)

type TxManager interface {
	// Do запускает fn в транзакции.
	// Если транзакция уже есть в ctx — просто переиспользует её (вложенный сценарий).
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
