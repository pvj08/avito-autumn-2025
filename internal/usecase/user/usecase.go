package user

import "github.com/pvj08/avito-autumn-2025/pkg/logger"

type UserRepository interface {
	// TODO: объявить интерфейс репозитория
}

type UserUsecase struct {
	log  logger.Logger
	repo UserRepository
}

func NewUserUsecase(log logger.Logger, repo UserRepository) *UserUsecase {
	return &UserUsecase{
		log:  log,
		repo: repo,
	}
}
