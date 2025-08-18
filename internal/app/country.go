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

type countryQ interface {
	New() dbx.CountryQ

	Insert(ctx context.Context, input dbx.Country) error
	Update(ctx context.Context, input dbx.UpdateCountryInput) error
	Get(ctx context.Context) (dbx.Country, error)
	Select(ctx context.Context) ([]dbx.Country, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CountryQ
	FilterName(name string) dbx.CountryQ
	FilterStatus(status string) dbx.CountryQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CountryQ
}

// Create methods for countries

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	country := dbx.Country{
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
			return models.Country{}, errx.RaiseCountryAlreadyExists(
				ctx,
				fmt.Errorf("country with name '%s' already existt, cause: %s", name, err),
				name,
			)
		default:
			return models.Country{}, errx.RaiseInternal(ctx, err)
		}
	}

	return countryModel(country), nil
}

// Read methods for countries

func (a App) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.countriesQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Country{}, errx.RaiseCountryNotFoundByID(
				ctx,
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		}
		return models.Country{}, errx.RaiseInternal(ctx, err)
	}

	return countryModel(country), nil
}

func (a App) SearchCountries(ctx context.Context, name string, status string, pag pagination.Request) ([]models.Country, pagination.Response, error) {
	limit, offset := pagination.CalculateLimitOffset(pag)

	countries, err := a.countriesQ.New().
		FilterName(name).
		FilterStatus(status).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	total, err := a.countriesQ.New().Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	res, pagRes := countriesArray(countries, limit, offset, total)

	return res, pagRes, nil
}

// Update methods for countries

func (a App) UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status string) (models.Country, error) {
	//TODO in future create event to kafka about country status change

	_, err := enum.ParseCountryStatus(status)
	if err != nil {
		return models.Country{}, errx.RaiseInvalidCountryStatus(
			ctx,
			fmt.Errorf("invalid country status '%s', cause: %s", status, err),
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
				return errx.RaiseCountryNotFoundByID(
					ctx,
					fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
					ID,
				)
			default:
				return errx.RaiseInternal(ctx, err)
			}
		}

		if status == enum.CountryStatusUnsupported {
			err = a.updateStatusForCitiesByCountryID(ctx, ID, enum.CountryStatusUnsupported)
			if err != nil {
				return err
			}
		} else if status == enum.CountryStatusSuspended {
			err = a.updateStatusForCitiesByCountryID(ctx, ID, enum.CountryStatusSuspended)
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
			return models.Country{}, errx.RaiseInternal(ctx, err)
		}
	}
	if err == nil && country.ID != ID {
		return models.Country{}, errx.RaiseCountryAlreadyExists(
			ctx,
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
			return models.Country{}, errx.RaiseCountryNotFoundByID(
				ctx,
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		default:
			return models.Country{}, errx.RaiseInternal(ctx, err)
		}
	}

	return a.GetCountryByID(ctx, ID)
}

// Helper functions for countries

func countryModel(country dbx.Country) models.Country {
	return models.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    country.Status,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}
}

func countriesArray(countries []dbx.Country, limit, offset, total uint64) ([]models.Country, pagination.Response) {
	result := make([]models.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, countryModel(country))
	}

	pag := pagination.Response{
		Page:  offset/limit + 1,
		Size:  limit,
		Total: total,
	}

	return result, pag
}
