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
		// достаём pull request по ID
		pr, err := u.prRepo.GetByID(ctx, input.PullRequestID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		// уже MERGED → просто отдаем как есть, без Save
		if pr.Status == domain.PullRequestStatusMERGED {
			out = MergeOutput{
				PullRequest: fromDomainPullRequest(pr),
			}
			return nil
		}

		// вызываем метод доменной модели для смены статуса
		pr.Merge()

		// сохраняем изменения в репозитории
		if err := u.prRepo.Save(ctx, pr); err != nil {
			return err
		}

		out = MergeOutput{
			PullRequest: fromDomainPullRequest(pr),
		}

		return nil
	})

	return out, err
}
