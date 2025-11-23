package pullrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) Merge(c context.Context, input MergeInput) (MergeOutput, error) {
	var out MergeOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		pr, err := u.prRepo.GetByID(ctx, input.PullRequestID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		pr.Merge()

		if err := u.prRepo.Save(ctx, pr); err != nil {
			return err
		}

		out = MergeOutput{
			PullRequest{
				AssignedReviewers: pr.AssignedReviewers,
				AuthorID:          pr.AuthorID,
				PullRequestID:     pr.PullRequestID,
				PullRequestName:   pr.PullRequestName,
				Status:            PullRequestStatus(pr.Status),
				CreatedAt:         pr.CreatedAt,
				MergedAt:          pr.MergedAt,
			},
		}

		return nil
	})

	return out, err
}
