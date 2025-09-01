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

func (a App) SearchSupportedCitiesInCountry(
	ctx context.Context,
	Name string,
	CountryID uuid.UUID,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
	paramToEntity := entities.SelectCityFilters{
		Name:      &Name,
		CountryID: &CountryID,
		Status:    []string{constant.CityStatusCommunity, constant.CityStatusOfficial},
	}
	res, pagination, err := a.cities.SelectCities(ctx, paramToEntity, pag, sort)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	return res, pagination, nil
}

type SelectCityFilters struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Point     *orb.Point
}

// SelectCities searches for cities by name, country ID, and status with pagination and sorting.
// This method for sysadmin
func (a App) SelectCities(
	ctx context.Context,
	filters SelectCityFilters,
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
	if filters.Point != nil {
		paramsToEntity.Point = filters.Point
	}

	return a.cities.SelectCities(ctx, paramsToEntity, pag, sort)
}

func (a App) GetNearbyCity(ctx context.Context, point orb.Point) (models.City, error) {
	return a.cities.GetByRadius(ctx, point, 15000) // 15 km
}

func (a App) UpdateCitySlug(ctx context.Context, cityID uuid.UUID, slug string) (models.City, error) {
	c, err := a.cities.GetBySlug(ctx, slug)
	switch {
	case err == nil && c.ID != cityID:
		return models.City{}, errx.ErrorCityAlreadyExists.Raise(fmt.Errorf("city with slug: %s already exists", slug))
	case err != nil && !errors.Is(err, errx.ErrorCityNotFound):
		return models.City{}, err
	}

	_, err = a.cities.GetBySlug(ctx, slug)
	if !errors.Is(err, errx.ErrorCityNotFound) && err != nil {
		return models.City{}, err
	}

	err = a.cities.UpdateOne(ctx, cityID, entities.UpdateCityParams{
		Slug:      &slug,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return models.City{}, err
	}

	return a.cities.GetByID(ctx, cityID)
}

type UpdateCityParams struct {
	Point    *orb.Point
	Name     *string
	Icon     *string
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

	var point *orb.Point
	if params.Point != nil {
		point = params.Point
	}

	err = a.cities.UpdateOne(ctx, cityID, entities.UpdateCityParams{
		Point:     point,
		Name:      params.Name,
		Icon:      params.Icon,
		Timezone:  params.Timezone,
		UpdatedAt: time.Now().UTC(),
	})
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

		err = a.gov.UpdateMany(ctx, entities.UpdateGovsFilters{
			CityID: &cityID,
		}, entities.UpdateGovsParams{
			Active: func(b bool) *bool { return &b }(false),
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

		err = a.gov.UpdateMany(ctx, entities.UpdateGovsFilters{
			CityID: &cityID,
		}, entities.UpdateGovsParams{
			Active: func(b bool) *bool { return &b }(false),
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
