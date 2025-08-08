package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	errs2 "github.com/chains-lab/cities-dir-svc/internal/errs"
	"github.com/google/uuid"
)

type countriesQ interface {
	New() dbx.CountriesQ

	Insert(ctx context.Context, input dbx.CountryModel) error
	Update(ctx context.Context, input dbx.UpdateCountryInput) error
	Get(ctx context.Context) (dbx.CountryModel, error)
	Select(ctx context.Context) ([]dbx.CountryModel, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CountriesQ
	FilterName(name string) dbx.CountriesQ
	FilterStatus(status enum.CountryStatus) dbx.CountriesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CountriesQ
}

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	country := dbx.CountryModel{
		ID:        uuid.New(),
		Name:      name,
		Status:    enum.CountryStatusUnsupported,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := a.countriesQ.New().Insert(ctx, country)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, errs2.RaiseCountryAlreadyExists(
				fmt.Errorf("country with name '%s' already existt, cause: %s", name, err),
				name,
			)
		default:
			return models.Country{}, errs2.RaiseInternal(err)
		}
	}

	return countryDbxToModel(country), nil
}

func (a App) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.countriesQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Country{}, errs2.RaiseCountryNotFoundByID(
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		}
		return models.Country{}, errs2.RaiseInternal(err)
	}

	return countryDbxToModel(country), nil
}

func (a App) UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status enum.CountryStatus) (models.Country, error) {
	//TODO in future create event to kafka about country status change

	_, ok := enum.ParseCountryStatus(string(status))
	if !ok {
		return models.Country{}, errs2.RaiseInvalidCountryStatus(
			fmt.Errorf("invalid country status: %s, met be one of %s", status, enum.GetAllCountriesStatuses()),
			status,
		)
	}

	trxErr := a.transaction(func(ctx context.Context) error {
		err := a.countriesQ.New().FilterID(ID).Update(ctx, dbx.UpdateCountryInput{
			Status: &status,
		})
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return errs2.RaiseCountryNotFoundByID(
					fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
					ID,
				)
			default:
				return errs2.RaiseInternal(err)
			}
		}

		if status == enum.CountryStatusUnsupported {
			err = a.UpdateStatusForCitiesByCountryID(ctx, ID, enum.CityStatusUnsupported)
			if err != nil {
				return err
			}
		} else if status == enum.CountryStatusSuspended {
			err = a.UpdateStatusForCitiesByCountryID(ctx, ID, enum.CityStatusSuspended)
			if err != nil {
				return err
			}
		} else if status == enum.CountryStatusSupported {
			//TODO mb in future create event to kafka about country status change
		}

		return nil
	})
	if trxErr != nil {
		return models.Country{}, trxErr
	}

	return a.GetCountryByID(ctx, ID)
}

func (a App) UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error) {
	country, err := a.countriesQ.New().FilterName(name).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// country with this name does not exist, so we can update it
		default:
			return models.Country{}, errs2.RaiseInternal(err)
		}
	}
	if err == nil && country.ID != ID {
		return models.Country{}, errs2.RaiseCountryAlreadyExists(
			fmt.Errorf("country with id '%s' already exists, cause: %s", ID, country),
			name,
		)
	}

	err = a.countriesQ.New().FilterID(ID).Update(ctx, dbx.UpdateCountryInput{
		Name: &name,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, errs2.RaiseCountryNotFoundByID(
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		default:
			return models.Country{}, errs2.RaiseInternal(err)
		}
	}

	return a.GetCountryByID(ctx, ID)
}

func (a App) DeleteCountry(ctx context.Context, ID uuid.UUID) error {
	//TODO this method is not safe, because it can delete country with cities and admins, need to add check for cities and admins
	// dose not allow to delete country with cities and admins, or delete all cities and admins before deleting country

	_, err := a.GetCountryByID(ctx, ID)
	if err != nil {
		return err
	}

	cities, err := a.citiesQ.New().FilterCountryID(ID).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			break
		default:
			return errs2.RaiseInternal(err)
		}
	}

	if err == nil || len(cities) > 0 {
		return errs2.RaiseInternal(err)
	}

	if err = a.countriesQ.New().FilterID(ID).Delete(ctx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errs2.RaiseCountryNotFoundByID(
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		default:
			return errs2.RaiseInternal(err)
		}
	}

	return nil
}

func (a App) SearchCountries(ctx context.Context, name string, status enum.CountryStatus, limit, offset uint64) ([]models.Country, error) {
	countries, err := a.countriesQ.New().
		FilterName(name).
		FilterStatus(status).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return nil, errs2.RaiseInternal(err)
	}

	return arrCountryDbxToModel(countries), nil
}

func countryDbxToModel(country dbx.CountryModel) models.Country {
	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    country.Status,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}
}

func arrCountryDbxToModel(countries []dbx.CountryModel) []models.Country {
	result := make([]models.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, countryDbxToModel(country))
	}
	return result
}
