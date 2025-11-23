package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/pvj08/avito-autumn-2025/internal/domain"
)

func (u *usecase) SetIsActive(c context.Context, input SetIsActiveInput) (SetIsActiveOutput, error) {
	var out SetIsActiveOutput

	err := u.tx.Do(c, func(ctx context.Context) error {
		user, err := u.userRepo.UpdateSetIsActive(ctx, input.UserID, input.IsActive)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to update user: %w", err)
		}

		out = SetIsActiveOutput{
			User: fromDomainUser(user),
		}

		return nil
	})
	return out, err
}
