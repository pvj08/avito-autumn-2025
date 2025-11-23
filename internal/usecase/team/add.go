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

		created, err := u.repo.Create(ctx, domainTeam)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return err
			}
			return fmt.Errorf("failed to create team: %w", err)
		}

		// domain -> DTO
		out = AddOutput{
			Team: fromDomainTeam(created),
		}

		return nil
	})

	return out, err
}
