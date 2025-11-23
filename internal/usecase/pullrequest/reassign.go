package pullrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

// Reassign — переназначить конкретного ревьювера на другого из его команды.
func (u *usecase) Reassign(c context.Context, input ReassignInput) (ReassignOutput, error) {
	var out ReassignOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		// 1. Забираем PR
		pr, err := u.prRepo.GetByID(ctx, input.PullRequestID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get pullrequest: %w", err)
		}

		// 2. Нельзя переназначать ревьюверов у MERGED PR
		if pr.Status == domain.PullRequestStatusMERGED {
			return domain.ErrMerged
		}

		// 3. Проверяем, что этот ревьювер вообще назначен
		idx := -1
		for i, id := range pr.AssignedReviewers {
			if id == input.UserID {
				idx = i
				break
			}
		}
		if idx == -1 {
			return domain.ErrNotAssigned
		}

		// 4. Берём команду ревьювера и список участников
		team, err := u.teamRepo.GetByUserID(ctx, input.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		// 5. Формируем список кандидатов: исключаем автора и текущих ревьюверов
		newReviewerID, err := reassignChooseReviewer(team, pr.AuthorID, input.UserID)
		if err != nil {
			return err
		}

		// 7. Заменяем старого ревьювера на нового
		pr.AssignedReviewers[idx] = newReviewerID

		// 8. Сохраняем изменения
		if err := u.prRepo.Save(ctx, pr); err != nil {
			return err
		}

		// 9. Маппим в DTO
		out = ReassignOutput{
			PullRequest: PullRequest{
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
