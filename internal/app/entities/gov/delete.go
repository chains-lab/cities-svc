package gov

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

type DeleteGovsFilters struct {
	UserID    *uuid.UUID
	CityID    *uuid.UUID
	CountryID *uuid.UUID
	Role      *string
}

func (g Gov) DeleteMany(ctx context.Context, filters DeleteGovsFilters) error {
	query := g.govQ.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Role != nil {
		err := constant.CheckCityGovRole(*filters.Role)
		if err != nil {
			return errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}

	err := query.Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city govs, cause: %w", err),
		)
	}

	return nil
}
