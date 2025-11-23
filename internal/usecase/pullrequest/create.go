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
		team, err := u.teamRepo.GetByUserID(ctx, input.AuthorID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		// 2. выбираем до 2 активных ревьюверов
		reviewers, err := chooseReviewers(team, input.AuthorID)
		if err != nil {
			return err
		}

		// 3. собираем доменную модель
		pr := domain.PullRequest{
			PullRequestID:     input.PullRequestID,
			PullRequestName:   input.PullRequestName,
			AuthorID:          input.AuthorID,
			AssignedReviewers: reviewers,
			Status:            domain.PullRequestStatusOPEN,
		}

		created, err := u.prRepo.Create(ctx, pr)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return err
			}
			return fmt.Errorf("failed to create pullrequest: %w", err)
		}

		out = CreateOutput{
			PullRequest: PullRequest{
				AssignedReviewers: created.AssignedReviewers,
				AuthorID:          created.AuthorID,
				PullRequestID:     created.PullRequestID,
				PullRequestName:   created.PullRequestName,
				Status:            PullRequestStatus(created.Status),
				CreatedAt:         created.CreatedAt,
				MergedAt:          created.MergedAt,
			},
		}

		return nil
	})

	return out, err
}
