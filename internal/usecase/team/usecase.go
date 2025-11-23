package team

import (
	"context"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/txmanager"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type TeamRepository interface {
	// TODO: объявить интерфейс репозитория
	GetByUserID(ctx context.Context, userID string) (domain.Team, error)
	GetByTeamName(ctx context.Context, teamName string) (domain.Team, error)
	Create(ctx context.Context, team domain.Team) (domain.Team, error)
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

// toDomainTeam маппит AddInput (из delivery/oapi) в доменную модель.
func toDomainTeam(input AddInput) domain.Team {
	members := make([]domain.TeamMember, 0, len(input.Members))
	for _, m := range input.Members {
		members = append(members, domain.TeamMember{
			IsActive: m.IsActive,
			UserID:   m.UserID,
			Username: m.Username,
		})
	}

	return domain.Team{
		TeamName: input.TeamName,
		Members:  members,
	}
}

// fromDomainTeam маппит domain.Team в DTO-структуру, которую возвращаем наружу.
// Здесь предполагается, что в AddOutput/GetOutput есть поле Team типа Team (DTO).
func fromDomainTeam(t domain.Team) Team {
	members := make([]TeamMember, 0, len(t.Members))
	for _, m := range t.Members {
		members = append(members, TeamMember{
			IsActive: m.IsActive,
			UserID:   m.UserID,
			Username: m.Username,
		})
	}

	return Team{
		TeamName: t.TeamName,
		Members:  members,
	}
}
