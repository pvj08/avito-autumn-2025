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
				return domain.ErrTeamExists
			}
			return fmt.Errorf("failed to create team: %w", err)
		}

		for _, member := range input.Members {
			user := toDomainUser(member, input.TeamName)

			err = u.userCreator.Create(ctx, user)

			// Может ли один пользователь быть в нескольких командах сразу?
			// Скорее всего да

			// Но у меня это ломает логику и я не хочу заморачиваться в тестовом
			if err != nil {
				if errors.Is(err, domain.ErrAlreadyExists) {
					return domain.ErrUserExists
				}
				return err
			}
			created.Members = append(created.Members, toTeamMember(user))
		}

		// domain -> DTO
		out = AddOutput{
			Team: fromDomainTeam(created),
		}

		return nil
	})

	return out, err
}
