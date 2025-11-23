package user

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type UserRepository interface {
	// TODO: объявить интерфейс репозитория
}

type Usecase interface {
	GetReview(c context.Context, input GetReviewInput) (GetReviewOutput, error)
	SetIsActive(c context.Context, input SetIsActiveInput) (SetIsActiveOutput, error)
}

type usecase struct {
	tx   txmanager.TxManager
	repo UserRepository
	log  logger.Logger
}

func New(tx txmanager.TxManager, repo UserRepository, log logger.Logger) *usecase {
	return &usecase{
		tx:   tx,
		log:  log,
		repo: repo,
	}
}
