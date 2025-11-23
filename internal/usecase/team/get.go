package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) Get(c context.Context, input GetInput) (GetOutput, error) {
	var out GetOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		team, err := u.teamRepo.GetByTeamName(ctx, input.TeamName)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to get team: %w", err)
		}

		out = GetOutput{
			Team: fromDomainTeam(team),
		}
		return nil
	})

	return out, err
}
