package team

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type TeamRepository interface {
	// TODO: объявить интерфейс репозитория
}

type Usecase interface {
	// Создать команду с участниками (создаёт/обновляет пользователей)
	Add(c context.Context, input AddInput) (AddOutput, error)

	Get(c context.Context, input GetInput) (GetOutput, error)
}

type usecase struct {
	tx   txmanager.TxManager
	repo TeamRepository
	log  logger.Logger
}

func New(tx txmanager.TxManager, repo TeamRepository, log logger.Logger) *usecase {
	return &usecase{
		tx:   tx,
		repo: repo,
		log:  log,
	}
}
