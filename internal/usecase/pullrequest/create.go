package pullrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) Create(c context.Context, input CreateInput) (CreateOutput, error) {
	var out CreateOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		// 1. достаём команду автора
		team, err := u.teamReader.GetByUserID(ctx, input.AuthorID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		// 2. выбираем до 2 активных ревьюверов
		reviewers := chooseReviewers(team, input.AuthorID)

		// 3. собираем доменную модель
		pr := fromAddInputToDomainPR(input, reviewers)

		created, err := u.prRepo.Create(ctx, pr)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return domain.ErrPrExists
			}
			return fmt.Errorf("failed to create pullrequest: %w", err)
		}

		out = CreateOutput{
			PullRequest: fromDomainPullRequest(created),
		}

		return nil
	})

	return out, err
}
