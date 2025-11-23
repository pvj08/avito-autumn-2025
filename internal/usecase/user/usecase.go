package user

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	UpdateSetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
}

type Usecase interface {
	SetIsActive(c context.Context, input SetIsActiveInput) (SetIsActiveOutput, error)
}

type usecase struct {
	tx       txmanager.TxManager
	userRepo UserRepository
	log      logger.Logger
}

func New(tx txmanager.TxManager, u UserRepository, log logger.Logger) *usecase {
	return &usecase{
		tx:       tx,
		userRepo: u,
		log:      log,
	}
}

func fromDomainUser(d domain.User) User {
	return User{
		UserID:   d.UserID,
		Username: d.Username,
		IsActive: d.IsActive,
		TeamName: d.TeamName,
	}
}
