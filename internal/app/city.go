package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/domain/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
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
	country, err := a.GetCountryByID(ctx, params.CountryID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSupported {
		return models.City{}, errx.ErrorCountryNotSupported.Raise(
			fmt.Errorf("country status is not 'supported'"),
		)
	}

	return a.cities.Create(ctx, city.CreateCityParams{
		Name:      params.Name,
		CountryID: params.CountryID,
		Point:     params.Point,
		Timezone:  params.Timezone,
	})
}

func (a App) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	return a.cities.GetByID(ctx, ID)
}

func (a App) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	return a.cities.GetBySlug(ctx, slug)
}

type FilterListCitiesParams struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Location  *FilterListCityDistance
}

type FilterListCityDistance struct {
	Point   orb.Point
	RadiusM uint64
}

// ListCities searches for cities by name, country ID, and status with pagination and sorting.
// This method for sysadmin
func (a App) ListCities(
	ctx context.Context,
	filters FilterListCitiesParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
	paramsToEntity := city.FilterListParams{}
	if filters.Name != nil {
		paramsToEntity.Name = filters.Name
	}
	if filters.Status != nil {
		paramsToEntity.Status = filters.Status
	}
	if filters.CountryID != nil {
		paramsToEntity.CountryID = filters.CountryID
	}
	if filters.Location != nil {
		paramsToEntity.Location = &city.FilterListDistance{
			Point:   filters.Location.Point,
			RadiusM: filters.Location.RadiusM,
		}
	}
	return a.cities.List(ctx, paramsToEntity, pag, sort)
}

type UpdateCityParams struct {
	Name     *string
	Point    *orb.Point
	Icon     *string
	Slug     *string
	Timezone *string
}

func (a App) UpdateCity(ctx context.Context, cityID, initiatorID uuid.UUID, role string, params UpdateCityParams) (models.City, error) {
	if role == roles.User {
		gov, err := a.GetInitiatorGov(ctx, initiatorID)
		if err != nil {
			return models.City{}, err
		}
		if gov.CityID != cityID {
			return models.City{}, errx.ErrorInitiatorIsNotThisCityGov.Raise(
				fmt.Errorf("initiator %s is not the city %s", initiatorID.String(), cityID.String()),
			)
		}
		if gov.Role == enum.CityGovRoleMayor || gov.Role == enum.CityGovRoleModerator || gov.Role == enum.CityGovRoleAdvisor {
			return models.City{}, errx.ErrorInitiatorGovRoleHaveNotEnoughRights.Raise(
				fmt.Errorf("initiator %s have not enough rights", initiatorID.String()),
			)
		}
	}

	update := city.UpdateCityParams{}

	if params.Point != nil {
		update.Point = params.Point
	}
	if params.Slug != nil {
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

	return a.cities.UpdateOne(ctx, cityID, update)
}
