package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type cityQ interface {
	New() dbx.CityQ

	Insert(ctx context.Context, input dbx.CityModel) error
	Update(ctx context.Context, input dbx.UpdateCityInput) error
	Get(ctx context.Context) (dbx.CityModel, error)
	Select(ctx context.Context) ([]dbx.CityModel, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CityQ
	FilterCountryID(countryID uuid.UUID) dbx.CityQ
	FilterStatus(status string) dbx.CityQ
	FilterName(name string) dbx.CityQ

	SearchName(name string) dbx.CityQ

	SortedNameAlphabet() dbx.CityQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CityQ
}

// Create methods for city

type CreateCityInput struct {
	CountryID uuid.UUID
	Name      string
	Status    string
}

func (a App) CreateCity(ctx context.Context, input CreateCityInput) (models.City, error) {
	country, err := a.GetCountryByID(ctx, input.CountryID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSuspended {
		return models.City{}, errx.RaiseCountryStatusIsNotApplicable(
			ctx,
			fmt.Errorf("country with ID '%s' is not '%s', current status: '%s'", input.CountryID, enum.CountryStatusSuspended, country.Status),
			input.CountryID,
			enum.CountryStatusSuspended,
			country.Status,
		)
	}

	ID := uuid.New()

	city := dbx.CityModel{
		ID:        ID,
		CountryID: input.CountryID,
		Name:      input.Name,
		Status:    input.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = a.citiesQ.New().Insert(ctx, city)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityModel(city), nil
}

// Read methods for city

func (a App) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	city, err := a.citiesQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseCityNotFoundByID(
				ctx,
				fmt.Errorf("city with ID '%s' not found cause: %s", ID, err),
				ID,
			)
		default:
			return models.City{}, errx.RaiseInternal(ctx, err)
		}
	}

	return cityModel(city), nil
}

func (a App) SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, pag pagination.Request) ([]models.City, pagination.Response, error) {
	_, err := a.GetCountryByID(ctx, countryID)
	if err != nil {
		return []models.City{}, pagination.Response{}, err
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	cities, err := a.citiesQ.New().
		FilterCountryID(countryID).
		SortedNameAlphabet().
		SearchName(like).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, pagination.Response{}, errx.RaiseCityNotFoundByName(
				ctx,
				fmt.Errorf("city with name '%s' not found in country with ID '%s' cause: %s", like, countryID, err),
				like,
			)
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	total, err := a.citiesQ.New().Count(context.Background())
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			total = 0
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	res, pagRes := citiesArray(cities, limit, offset, total)

	return res, pagRes, nil
}

// Update methods for city

func (a App) UpdateCityName(ctx context.Context, cityID uuid.UUID, name string) (models.City, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	cityUpdate := dbx.UpdateCityInput{
		Name:      &name,
		UpdatedAt: time.Now().UTC(),
	}

	err = a.citiesQ.New().FilterID(cityID).Update(ctx, cityUpdate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseInternal(ctx, err)
		}
	}

	return a.GetCityByID(ctx, cityID)
}

func (a App) UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	city, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	country, err := a.GetCountryByID(ctx, city.CountryID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSupported {
		return models.City{}, errx.RaiseCountryStatusIsNotApplicable(
			ctx,
			fmt.Errorf("country with ID '%s' is not %s, current status: %s", country.ID, enum.CountryStatusSupported, country.Status),
			country.ID,
			country.Status,
			enum.CountryStatusSupported,
		)
	}

	_, err = enum.ParseCityStatus(status)
	if err != nil {
		return models.City{}, errx.RaiseInvalidCityStatus(ctx, err, status)
	}

	err = a.citiesQ.New().FilterID(city.ID).Update(ctx, dbx.UpdateCityInput{
		Status: &status,
	})
	if err != nil {
		return models.City{}, errx.RaiseInternal(ctx, err)
	}

	return a.GetCityByID(ctx, cityID)
}

// updateStatusForCitiesByCountryID updates the status of all cities in a given country.
// Its internal method used to update cities status when country status is changed.
func (a App) updateStatusForCitiesByCountryID(ctx context.Context, countryID uuid.UUID, status string) error {
	_, err := a.countriesQ.New().FilterID(countryID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseCountryNotFoundByID(
				ctx,
				fmt.Errorf("country with ID '%s' not found cause: %s", countryID, err),
				countryID,
			)
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	cities, err := a.citiesQ.New().FilterCountryID(countryID).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = nil
		default:
			return errx.RaiseInternal(ctx, err)
		}
	}

	for _, city := range cities {
		err = a.citiesQ.New().FilterID(city.ID).Update(ctx, dbx.UpdateCityInput{
			Status: &status,
		})
		if err != nil {
			return errx.RaiseInternal(ctx, err)
		}
	}

	return nil
}

// internal methods  for city
func cityModel(city dbx.CityModel) models.City {
	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Name:      city.Name,
		Status:    city.Status,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}
}

func citiesArray(cities []dbx.CityModel, limit, offset, total uint64) ([]models.City, pagination.Response) {
	res := make([]models.City, 0, len(cities))
	for _, city := range cities {
		res = append(res, cityModel(city))
	}

	pag := pagination.Response{
		Page:  offset/limit + 1,
		Size:  limit,
		Total: total,
	}

	return res, pag
}
