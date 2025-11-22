package pullrequest

import "github.com/pvj08/avito-autumn-2025/pkg/logger"

type PullRequestRepository interface {
	// TODO: объявить интерфейс репозитория
}

type PullRequestUsecase struct {
	log  logger.Logger
	repo PullRequestRepository
}

func NewPullRequestUsecase(log logger.Logger, repo PullRequestRepository) *PullRequestUsecase {
	return &PullRequestUsecase{
		log:  log,
		repo: repo,
	}
}
