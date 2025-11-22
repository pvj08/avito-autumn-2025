package team

import "github.com/pvj08/avito-autumn-2025/pkg/logger"

type TeamRepository interface {
	// TODO: объявить интерфейс репозитория
}

type TeamUsecase struct {
	log  logger.Logger
	repo TeamRepository
}

func NewTeamUsecase(log logger.Logger, repo TeamRepository) *TeamUsecase {
	return &TeamUsecase{
		log:  log,
		repo: repo,
	}
}
