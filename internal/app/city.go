package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateCityParams struct {
	Name      string
	CountryID uuid.UUID
	Point     orb.Point
	Timezone  string
}

func (a App) CreateCity(ctx context.Context, params CreateCityParams) (models.City, error) {
	return a.cities.Create(ctx, entities.CreateCityParams{
		Name:      params.Name,
		CountryID: params.CountryID,
		Point:     params.Point,
		Timezone:  params.Timezone,
	})
}

func (a App) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	return a.cities.GetByID(ctx, ID)
}

type SearchCityFilters struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Point     *orb.Point
	Radius    *uint // in meters
}

// SearchCities searches for cities by name, country ID, and status with pagination and sorting.
// This method for sysadmin
func (a App) SearchCities(
	ctx context.Context,
	filters SearchCityFilters,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
	paramsToEntity := entities.SelectCityFilters{}
	if filters.Name != nil {
		paramsToEntity.Name = filters.Name
	}
	if filters.Status != nil {
		paramsToEntity.Status = filters.Status
	}
	if filters.CountryID != nil {
		paramsToEntity.CountryID = filters.CountryID
	}
	if filters.Point != nil && filters.Radius != nil {
		paramsToEntity.Point = filters.Point
		paramsToEntity.Radius = filters.Radius
	}
	return a.cities.Select(ctx, paramsToEntity, pag, sort)
}

func (a App) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	return a.cities.GetBySlug(ctx, slug)
}

type UpdateCityParams struct {
	Name     *string
	Status   *string
	Point    *orb.Point
	Icon     *string
	Slug     *string
	Timezone *string
}

func (a App) UpdateCity(ctx context.Context, cityID uuid.UUID, params UpdateCityParams) (models.City, error) {
	_, err := a.cities.GetByID(ctx, cityID)
	if errors.Is(err, errx.ErrorCityNotFound) {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			errors.New("city not found"),
		)
	}
	if err != nil {
		return models.City{}, err
	}

	update := entities.UpdateCityParams{}

	if params.Point != nil {
		if params.Point[0] < -180 || params.Point[0] > 180 || params.Point[1] < -90 || params.Point[1] > 90 {
			return models.City{}, errx.ErrorInvalidPoint.Raise(
				fmt.Errorf("invalid point coordinates: longitude %f, latitude %f", params.Point[0], params.Point[1]),
			)
		}
		update.Point = params.Point
	}
	if params.Slug != nil {
		_, err = a.cities.GetBySlug(ctx, *params.Slug)
		if err == nil {
			return models.City{}, errx.ErrorCityAlreadyExistsWithThisSlug.Raise(
				fmt.Errorf("city with slug: %s already exists", *params.Slug),
			)
		} else if !errors.Is(err, errx.ErrorCityNotFound) {
			return models.City{}, err
		}

		update.Slug = params.Slug
	}
	if params.Name != nil {
		update.Name = params.Name
	}
	if params.Timezone != nil {
		update.Timezone = params.Timezone
	}
	if params.Icon != nil {
		update.Icon = params.Icon
	}
	if params.Status != nil {
		err = constant.CheckCityStatus(*params.Status)
		if err != nil {
			return models.City{}, errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to parse city status: %w", err),
			)
		}

		update.Status = params.Status
	}

	if update == (entities.UpdateCityParams{}) {
		return a.cities.GetByID(ctx, cityID)
	}

	update.UpdatedAt = time.Now().UTC()

	err = a.cities.UpdateOne(ctx, cityID, update)
	if err != nil {
		return models.City{}, err
	}

	return a.cities.GetByID(ctx, cityID)
}

func (a App) SetCityStatusCommunity(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	_, err := a.cities.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	status := constant.CityStatusCommunity

	txErr := a.transaction(func(ctx context.Context) error {
		err = a.cities.UpdateOne(ctx, cityID, entities.UpdateCityParams{
			Status:    &status,
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		err = a.gov.DeleteMany(ctx, entities.DeleteGovsFilters{
			CityID: &cityID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.City{}, txErr
	}

	return a.GetCityByID(ctx, cityID)
}

func (a App) SetCityStatusOfficial(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	_, err := a.cities.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	status := constant.CityStatusOfficial

	err = a.cities.UpdateOne(ctx, cityID, entities.UpdateCityParams{
		Status:    &status,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return models.City{}, err
	}

	return a.GetCityByID(ctx, cityID)
}

func (a App) SetCityStatusDeprecated(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	_, err := a.cities.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	status := constant.CityStatusDeprecated

	txErr := a.transaction(func(ctx context.Context) error {
		err = a.cities.UpdateOne(ctx, cityID, entities.UpdateCityParams{
			Status:    &status,
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		err = a.gov.DeleteMany(ctx, entities.DeleteGovsFilters{
			CityID: &cityID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.City{}, txErr
	}

	return a.GetCityByID(ctx, cityID)
}
