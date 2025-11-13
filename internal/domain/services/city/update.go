package city

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateParams struct {
	Point    *orb.Point
	Name     *string
	Icon     *string
	Slug     *string
	Timezone *string
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

	admins, err := s.db.GetCityAdmins(ctx, cityID)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins, cause: %w", err),
		)
	}

	err = s.event.PublishCityUpdated(ctx, city, admins.GetUserIDs())
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city updated event, cause: %w", err),
		)
	}

	return city, nil
}

func (s Service) UpdateStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	err := enum.CheckCityStatus(status)
	if err != nil {
		return models.City{}, errx.ErrorInvalidCityStatus.Raise(err)
	}

	now := time.Now().UTC()

	city, err := s.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	recipients, err := s.db.GetCityAdmins(ctx, city.ID)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", city.ID, err),
		)
	}

	err = s.db.Transaction(ctx, func(ctx context.Context) error {
		switch status {
		case enum.CityStatusSupported:
			err = s.CountryIsSupported(ctx, city.CountryID)
			if err != nil {
				return err
			}

		case enum.CityStatusSuspended:
			err = s.CountryIsSupported(ctx, city.CountryID)
			if err != nil {
				return err
			}

			err = s.db.DeleteAdminsForCity(ctx, city.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status, cause: %w", err),
				)
			}

		case enum.CityStatusUnsupported:
			err = s.db.DeleteAdminsForCity(ctx, city.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status, cause: %w", err),
				)
			}
		}

		err = s.db.UpdateCityStatus(ctx, cityID, status, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city status, cause: %w", err),
			)
		}

		city.Status = status
		city.UpdatedAt = now

		return nil
	})
	if err != nil {
		return models.City{}, err
	}

	err = s.event.PublishCityUpdatedStatus(ctx, city, status, recipients.GetUserIDs())
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city updated status event, cause: %w", err),
		)
	}

	return city, nil
}
