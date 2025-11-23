package team

import (
	"context"

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
	log  logger.Logger
	repo TeamRepository
}

func New(log logger.Logger, repo TeamRepository) *usecase {
	return &usecase{
		log:  log,
		repo: repo,
	}
}
