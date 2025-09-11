package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) RefuseOwnGov(ctx context.Context, userID uuid.UUID) error {
	gov, err := a.GetInitiatorGov(ctx, userID)
	if err != nil {
		return err
	}

	if gov.Role == constant.CityGovRoleMayor {
		return errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self gov", userID),
		)
	}

	err = a.gov.DeleteOne(ctx, gov.UserID)
	if err != nil {
		return fmt.Errorf("refuse own gov: %w", err)
	}

	return nil
}
