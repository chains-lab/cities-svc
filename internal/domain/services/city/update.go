package city

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateParams struct {
	CountryID *uuid.UUID
	Point     *orb.Point
	Name      *string
	Icon      *string
	Slug      *string
	Timezone  *string
}

func (s Service) Update(ctx context.Context, cityID uuid.UUID, params UpdateParams) (models.City, error) {
	city, err := s.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	if params.Point != nil {
		err = validatePoint(*params.Point)
		if err != nil {
			return models.City{}, err
		}

		city.Point = *params.Point
	}

	if params.Name != nil {
		err = validateName(*params.Name)
		if err != nil {
			return models.City{}, err
		}

		city.Name = *params.Name
	}

	if params.Slug != nil {
		err = validateSlug(*params.Slug)

		_, err = s.GetBySlug(ctx, *params.Slug)
		if err != nil && !errors.Is(err, errx.ErrorCityNotFound) {
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by slug, cause: %w", err),
			)
		} else if err == nil {
			return models.City{}, errx.ErrorCityAlreadyExistsWithThisSlug.Raise(
				fmt.Errorf("city with slug: %s already exists", *params.Slug),
			)
		}

		city.Slug = params.Slug
	}

	if params.Timezone != nil {
		err = validateTimezone(*params.Timezone)
		if err != nil {
			return models.City{}, err
		}
		city.Timezone = *params.Timezone
	}

	now := time.Now().UTC()

	err = s.db.UpdateCity(ctx, cityID, params, now)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city, cause: %w", err),
		)
	}

	return city, nil
}
