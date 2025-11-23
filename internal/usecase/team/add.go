package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) Add(c context.Context, input AddInput) (AddOutput, error) {
	var out AddOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		// DTO -> domain
		domainTeam := toDomainTeam(input)

		created, err := u.teamRepo.Create(ctx, domainTeam)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return err
			}
			return fmt.Errorf("failed to create team: %w", err)
		}

		for _, member := range input.Members {
			user := toDomainUser(member, input.TeamName)

			err = u.userCreator.Create(ctx, user)
			if err != nil && !errors.Is(err, domain.ErrAlreadyExists) {
				return fmt.Errorf("failed to add member to team: %w", err)
			}
		}

		// domain -> DTO
		out = AddOutput{
			Team: fromDomainTeam(created),
		}

		return nil
	})

	return out, err
}

func toDomainUser(member TeamMember, tname string) domain.User {
	return domain.User{
		IsActive: member.IsActive,
		UserID:   member.UserID,
		Username: member.Username,
		TeamName: tname,
	}
}
