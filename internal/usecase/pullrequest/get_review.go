package pullrequest

import (
	"context"
	"fmt"
)

// Не проверяет есть ли user в таблице вообще.
// Если не находит назначенных PR этому юзеру, вернут пустой слайс PRShort
func (u *usecase) GetReview(c context.Context, input GetReviewInput) (GetReviewOutput, error) {
	var out GetReviewOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		// достаём pull request по ID
		prs, err := u.prRepo.GetByUserID(ctx, input.UserID)
		if err != nil {
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
