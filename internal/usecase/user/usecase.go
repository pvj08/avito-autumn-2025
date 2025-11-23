package user

import (
	"context"

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
	log  logger.Logger
	repo UserRepository
}

func New(log logger.Logger, repo UserRepository) *usecase {
	return &usecase{
		log:  log,
		repo: repo,
	}
}
