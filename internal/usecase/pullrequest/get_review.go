package pullrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) GetReview(c context.Context, input GetReviewInput) (GetReviewOutput, error) {
	var out GetReviewOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		// достаём pull request по ID
		prs, err := u.prRepo.GetByUserID(ctx, input.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		out = GetReviewOutput{
			UserID:       input.UserID,
			PullRequests: fromDomainPRtoShortPRSlice(prs),
		}

		return nil
	})

	return out, err
}
