package app

import (
	"context"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/entities"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	_, err := a.country.GetByName(ctx, name)
	if err == nil {
		return models.Country{}, errx.ErrorCountryAlreadyExists.Raise(err)
	}

	country, err := a.country.Create(ctx, name, constant.CountryStatusUnsupported)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}

func (a App) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.country.GetByID(ctx, ID)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}

func (a App) GetCountryByName(ctx context.Context, name string) (models.Country, error) {
	country, err := a.country.GetByName(ctx, name)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}

type SearchCountriesFilters struct {
	Name     string
	Statuses []string
}

func (a App) SearchCountries(
	ctx context.Context,
	filters SearchCountriesFilters,
	pagination pagi.Request,
	sort []pagi.SortField,
) ([]models.Country, pagi.Response, error) {
	return a.country.Select(ctx, entities.SelectCountriesFilters{
		Name:     filters.Name,
		Statuses: filters.Statuses,
	}, pagination, sort)
}

func (a App) SetCountryStatusSupported(ctx context.Context, countryID uuid.UUID) (models.Country, error) {
	country, err := a.country.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	updatedAt := time.Now().UTC()
	countryStatus := constant.CountryStatusSupported

	if err = a.country.Update(ctx, country.ID, entities.UpdateCountryParams{
		Status:    &countryStatus,
		UpdatedAt: updatedAt,
	}); err != nil {
		return models.Country{}, err
	}

	//TODO in future kafka event about country status change

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    countryStatus,
		CreatedAt: country.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (a App) SetCountryStatusDeprecated(ctx context.Context, countryID uuid.UUID) (models.Country, error) {
	country, err := a.country.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	updatedAt := time.Now().UTC()
	countryStatus := constant.CountryStatusDeprecated
	cityStatus := constant.CountryStatusDeprecated

	txErr := a.transaction(func(ctx context.Context) error {
		if err = a.country.Update(ctx, country.ID, entities.UpdateCountryParams{
			Status:    &countryStatus,
			UpdatedAt: updatedAt,
		}); err != nil {
			return err
		}

		err = a.cities.UpdateMany(ctx, entities.UpdateCitiesFilters{
			CountryID: &country.ID,
			Status:    []string{constant.CityStatusOfficial, constant.CityStatusCommunity},
		}, entities.UpdateCitiesParams{
			Status:    &cityStatus,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}

		err = a.gov.UpdateMany(ctx, entities.UpdateGovsFilters{
			CountryID: &country.ID,
		}, entities.UpdateGovsParams{
			Active: func(b bool) *bool { return &b }(false),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.Country{}, txErr
	}

	//TODO in future kafka event about country status change

	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    countryStatus,
		CreatedAt: country.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}

type UpdateCountryParams struct {
	Name *string
}

func (a App) UpdateCountry(ctx context.Context, countryID uuid.UUID, params UpdateCountryParams) (models.Country, error) {
	country, err := a.country.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	updatedAt := time.Now().UTC()
	update := entities.UpdateCountryParams{
		UpdatedAt: updatedAt,
	}

	if params.Name != nil {
		_, err = a.country.GetByName(ctx, *params.Name)
		if err == nil {
			return models.Country{}, errx.ErrorCountryAlreadyExists.Raise(
				err,
			)
		} else if !errors.Is(err, errx.ErrorCountryNotFound) {
			return models.Country{}, err
		}
		update.Name = params.Name
		country.Name = *params.Name
	}

	if err = a.country.Update(ctx, country.ID, update); err != nil {
		return models.Country{}, err
	}

	return country, nil
}
