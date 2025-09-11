package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) DeleteGov(ctx context.Context, initiatorID, userID uuid.UUID) (models.Gov, error) {
	initiator, err := a.GetInitiatorGov(ctx, initiatorID)
	if err != nil {
		return models.Gov{}, err
	}

	g, err := a.gov.Get(ctx, gov.GetGovFilters{
		UserID: &userID,
	})
	if err != nil {
		return models.Gov{}, err
	}

	access, err := constant.CompareCityGovRoles(initiator.Role, g.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}
	if access != 1 {
		return models.Gov{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("no access to delete gov %s by initiator %s", userID, initiatorID),
		)
	}

	if g.Role == constant.CityGovRoleMayor {
		return models.Gov{}, errx.ErrorCannotRefuseMayor.Raise(
			fmt.Errorf("mayor %s cannot refuse self gov", userID),
		)
	}

	err = a.gov.DeleteOne(ctx, g.UserID)
	if err != nil {
		return models.Gov{}, fmt.Errorf("refuse own gov: %w", err)
	}

	return g, nil
}
