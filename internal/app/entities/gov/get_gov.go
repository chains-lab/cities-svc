package gov

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type GetGovFilters struct {
	UserID *uuid.UUID
	CityID *uuid.UUID
	Role   *string
}

func (g Gov) GetGov(ctx context.Context, filters GetGovFilters) (models.Gov, error) {
	query := g.govQ.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil {
		err := enum.CheckCityGovRole(*filters.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}

	gov, err := query.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Gov{}, errx.ErrorCityGovNotFound.Raise(
				fmt.Errorf("city gov not found, cause: %w", err),
			)
		default:
			return models.Gov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov, cause: %w", err),
			)
		}
	}

	return govFromDb(gov), nil
}

func (g Gov) GetInitiatorGov(ctx context.Context, initiatorID uuid.UUID) (models.Gov, error) {
	initiator, err := g.GetGov(ctx, GetGovFilters{
		UserID: &initiatorID,
	})
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			return models.Gov{}, errx.ErrorInitiatorIsNotActiveCityGov.Raise(
				fmt.Errorf("initiator %s is not an active city gov", initiatorID),
			)
		default:
			return models.Gov{}, err
		}
	}

	return initiator, nil
}
