package city

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateCityParams struct {
	CountryID *uuid.UUID
	Point     *orb.Point
	Status    *string
	Name      *string
	Icon      *string
	Slug      *string
	Timezone  *string
}

func (c City) UpdateOne(ctx context.Context, cityID uuid.UUID, params UpdateCityParams) (models.City, error) {
	city, err := c.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	if params.CountryID == nil && params.Point == nil && params.Status == nil &&
		params.Name == nil && params.Icon == nil && params.Slug == nil && params.Timezone == nil {
		return models.City{}, nil
	}

	stmt := dbx.UpdateCityParams{}
	if params.CountryID != nil {
		city.CountryID = *params.CountryID
		stmt.CountryID = params.CountryID
	}

	if params.Point != nil {
		err = c.validatePoint(*params.Point)
		if err != nil {
			return models.City{}, err
		}

		city.Point = *params.Point
		stmt.Point = params.Point
	}

	if params.Status != nil {
		err = constant.CheckCityStatus(*params.Status)
		if err != nil {
			return models.City{}, errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status, cause: %s", err),
			)
		}

		city.Status = *params.Status
		stmt.Status = params.Status
	}

	if params.Name != nil {
		err = c.validateName(*params.Name)
		if err != nil {
			return models.City{}, err
		}

		city.Name = *params.Name
		stmt.Name = params.Name
	}

	if params.Icon != nil && *params.Icon != "" {

		city.Icon = params.Icon
		stmt.Icon = &sql.NullString{String: *params.Icon, Valid: true}
	} else if params.Icon != nil && *params.Icon == "" {

		city.Icon = nil
		stmt.Icon = &sql.NullString{String: "", Valid: false}
	}

	if params.Slug != nil && *params.Slug != "" {
		err = c.validateSlug(*params.Slug)

		_, err = c.GetBySlug(ctx, *params.Slug)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by slug: %w", err),
			)
		} else if err == nil {
			return models.City{}, errx.ErrorCityAlreadyExistsWithThisSlug.Raise(
				fmt.Errorf("city with slug: %s already exists", *params.Slug),
			)
		}

		city.Slug = params.Slug
		stmt.Slug = &sql.NullString{String: *params.Slug, Valid: true}
	} else if params.Slug != nil && *params.Slug == "" {

		city.Slug = nil
		stmt.Slug = &sql.NullString{String: "", Valid: false}
	}

	if params.Timezone != nil {
		err = c.validateTimezone(*params.Timezone)
		if err != nil {
			return models.City{}, err
		}
		stmt.Timezone = params.Timezone
	}

	stmt.UpdatedAt = time.Now().UTC()
	city.UpdatedAt = stmt.UpdatedAt

	err = c.citiesQ.New().FilterID(cityID).Update(ctx, stmt)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city, cause: %w", err),
		)
	}

	return city, nil
}

type UpdateCitiesFilters struct {
	CountryID *uuid.UUID
	Status    []string
}

type UpdateCitiesParams struct {
	Status   *string
	Timezone *string

	UpdatedAt time.Time
}

func (c City) UpdateMany(ctx context.Context, filters UpdateCitiesFilters, params UpdateCitiesParams) error {
	query := c.citiesQ.New()
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Status != nil {
		for _, s := range filters.Status {
			err := constant.CheckCityStatus(s)
			if err != nil {
				return errx.ErrorInvalidCityStatus.Raise(
					fmt.Errorf("failed to invalid city status: %s, cause: %w", s, err),
				)
			}
		}
		query = query.FilterStatus(filters.Status...)
	}

	if params.Status == nil && params.Timezone == nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("0 filters provided for update cities"),
		)
	}

	stmt := dbx.UpdateCityParams{
		UpdatedAt: params.UpdatedAt,
	}

	if params.Status != nil {
		err := constant.CheckCityStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status: %s, cause: %w", *params.Status, err),
			)
		}
		stmt.Status = params.Status
	}
	if params.Timezone != nil {
		err := c.validateTimezone(*params.Timezone)
		if err != nil {
			return err
		}
		stmt.Timezone = params.Timezone
	}

	err := query.Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update cities by country_id, cause: %w", err),
		)
	}

	return nil
}
