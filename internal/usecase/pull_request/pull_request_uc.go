package usecase

import "github.com/pvj08/avito-autumn-2025/pkg/logger"

type pullRequestUC struct {
	log  logger.Logger
	repo PullRequestRepository
}

type PullRequestRepository interface {
	// TODO: объявить интерфейс репозитория
}

func NewPullRequestUC(log logger.Logger, repo PullRequestRepository) *pullRequestUC {
	return &pullRequestUC{
		log:  log,
		repo: repo,
	}
}

func (uc *pullRequestUC) Create() {
	// TODO: Реализовать методы UC
}
