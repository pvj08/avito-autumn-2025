package pullrequest

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type PullRequestRepository interface {
	// TODO: объявить интерфейс репозитория
}

type Usecase interface {
	Create(c context.Context, input CreateInput) (CreateOutput, error)

	// Пометить PR как MERGED (идемпотентная операция)
	Merge(c context.Context, input MergeInput) (MergeOutput, error)

	// Переназначить конкретного ревьювера на другого из его команды
	// Тут дохера ошибок проверить надо
	Reassign(c context.Context, input ReassignInput) (ReassignOutput, error)
}

type usecase struct {
	log  logger.Logger
	repo PullRequestRepository
}

func New(log logger.Logger, repo PullRequestRepository) *usecase {
	return &usecase{
		log:  log,
		repo: repo,
	}
}
