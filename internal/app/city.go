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
	"github.com/google/uuid"
)

type cityQ interface {
	New() dbx.CityQ

	Insert(ctx context.Context, input dbx.CityModels) error
	Update(ctx context.Context, input dbx.UpdateCityInput) error
	Get(ctx context.Context) (dbx.CityModels, error)
	Select(ctx context.Context) ([]dbx.CityModels, error)
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

type CreateCityInput struct {
	CountryID uuid.UUID
	Name      string
	Status    string
}

// its internal method for update city status, careful with it
func (a App) updateCityStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	//TODO in future realize confirmation of status change to email or something like that
	//TODO add kafka event for city status change

	err := a.citiesQ.New().FilterID(cityID).Update(ctx, dbx.UpdateCityInput{
		Status: &status,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseCityNotFoundByID(err, cityID)
		default:
			return models.City{}, errx.RaiseInternal(err)
		}
	}

	return a.GetCityByID(ctx, cityID)
}

func (a App) CreateCity(ctx context.Context, input CreateCityInput) (models.City, error) {
	country, err := a.GetCountryByID(ctx, input.CountryID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSuspended {
		return models.City{}, errx.RaiseCountryStatusIsNotApplicable(
			fmt.Errorf("country with ID '%s' is not '%s', current status: '%s'", input.CountryID, enum.CountryStatusSuspended, country.Status),
			input.CountryID,
			enum.CountryStatusSuspended,
			country.Status,
		)
	}

	ID := uuid.New()

	city := dbx.CityModels{
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
			return models.City{}, errx.RaiseInternal(err)
		}
	}

	return cityDbxToModel(city), nil
}

func (a App) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	city, err := a.citiesQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseCityNotFoundByID(
				fmt.Errorf("city with ID '%s' not found cause: %s", ID, err),
				ID,
			)
		default:
			return models.City{}, errx.RaiseInternal(err)
		}
	}

	return cityDbxToModel(city), nil
}

func (a App) DeleteCity(ctx context.Context, ID uuid.UUID) error {
	err := a.citiesQ.New().FilterID(ID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseCityNotFoundByID(
				fmt.Errorf("city with ID '%s' not found cause: %s", ID, err),
				ID,
			)
		case err != nil:
			return errx.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, page, limit uint64) ([]models.City, error) {
	cities, err := a.citiesQ.New().
		FilterCountryID(countryID).
		SortedNameAlphabet().
		SearchName(like).
		Page(page, limit).
		Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errx.RaiseCityNotFoundByName(
				fmt.Errorf("city with name '%s' not found in country with ID '%s' cause: %s", like, countryID, err),
				like,
			)
		default:
			return nil, errx.RaiseInternal(err)
		}
	}

	return citiesArrayDbxToModel(cities), nil
}

func (a App) UpdateCitiesStatusByOwner(ctx context.Context, initiatorID, cityID uuid.UUID, status string) (models.City, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	country, err := a.GetCountryByID(ctx, initiatorID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSuspended {
		return models.City{}, errx.RaiseCountryStatusIsNotApplicable(
			fmt.Errorf("country with ID '%s' is not %s, current status: %s", country.ID, enum.CountryStatusSuspended, country.Status),
			country.ID,
			country.Status,
			enum.CountryStatusSuspended,
		)
	}

	initiator, err := a.GetInitiatorCityAdmin(ctx, cityID, initiatorID)
	if err != nil {
		return models.City{}, errx.RaiseInternal(err)
	}

	if initiator.Role != enum.CityAdminRoleOwner {
		return models.City{}, errx.RaiseCityAdminHaveNotEnoughRights(
			fmt.Errorf("initiator: '%s', is not owner of city: '%s'",
				initiatorID,
				cityID,
			),
			cityID,
			initiatorID,
		)
	}

	_, err = enum.ParseCityStatus(string(status))
	if err != nil {
		return models.City{}, errx.RaiseInvalidCityStatus(err, status)
	}

	return a.updateCityStatus(ctx, cityID, status)
}

func (a App) UpdateCitiesStatusBySysAdmin(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	_, err = enum.ParseCityStatus(status)
	if err != nil {
		return models.City{}, errx.RaiseInvalidCityStatus(err, status)
	}

	return a.updateCityStatus(ctx, cityID, status)
}

func (a App) UpdateStatusForCitiesByCountryID(ctx context.Context, countryID uuid.UUID, status string) error {
	_, err := a.countriesQ.New().FilterID(countryID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.RaiseCountryNotFoundByID(
				fmt.Errorf("country with ID '%s' not found cause: %s", countryID, err),
				countryID,
			)
		default:
			return errx.RaiseInternal(err)
		}
	}

	cities, err := a.citiesQ.New().FilterCountryID(countryID).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = nil
		default:
			return errx.RaiseInternal(err)
		}
	}

	for _, city := range cities {
		_, err := a.updateCityStatus(ctx, city.ID, status)
		if err != nil {
			return errx.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) UpdateCityName(ctx context.Context, initiatorID, cityID uuid.UUID, name string) (models.City, error) {
	_, err := a.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	initiator, err := a.GetCityAdminForCity(ctx, initiatorID, cityID)
	if err != nil {
		return models.City{}, err
	}

	if initiator.Role != enum.CityAdminRoleOwner {
		return models.City{}, errx.RaiseCityAdminHaveNotEnoughRights(
			fmt.Errorf("initiator: '%s', is not owner of city: '%s'",
				initiatorID,
				cityID,
			),
			cityID,
			initiatorID,
		)
	}

	cityUpdate := dbx.UpdateCityInput{
		Name:      &name,
		UpdatedAt: time.Now().UTC(),
	}

	err = a.citiesQ.New().FilterID(cityID).Update(ctx, cityUpdate)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.RaiseInternal(err)
		}
	}

	return a.GetCityByID(ctx, cityID)
}

func cityDbxToModel(city dbx.CityModels) models.City {
	return models.City{
		ID:        city.ID,
		CountryID: city.CountryID,
		Name:      city.Name,
		Status:    city.Status,
		CreatedAt: city.CreatedAt,
		UpdatedAt: city.UpdatedAt,
	}
}

func citiesArrayDbxToModel(cities []dbx.CityModels) []models.City {
	res := make([]models.City, 0, len(cities))
	for _, city := range cities {
		res = append(res, cityDbxToModel(city))
	}
	return res
}
