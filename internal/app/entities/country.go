package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-proto/gen/go/common/pagination"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/pagi"

	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/google/uuid"
)

type countryQ interface {
	New() dbx.CountryQ

	Insert(ctx context.Context, input dbx.Country) error
	Update(ctx context.Context, input map[string]any) error
	Get(ctx context.Context) (dbx.Country, error)
	Select(ctx context.Context) ([]dbx.Country, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CountryQ
	FilterName(name string) dbx.CountryQ
	FilterStatus(status string) dbx.CountryQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CountryQ
}

type Country struct {
	country countryQ
}

func NewCountryService(cfg config.Config) (Country, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return Country{}, err
	}

	return Country{
		country: dbx.NewCountryQ(pg),
	}, nil
}

// Create methods for countries

func (a Country) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	country := dbx.Country{
		ID:        uuid.New(),
		Name:      name,
		Status:    enum.CountryStatusUnsupported,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := a.country.New().Insert(ctx, country)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, problems.RaiseCountryAlreadyExists(
				ctx,
				fmt.Errorf("country with name '%s' already existt, cause: %s", name, err),
				name,
			)
		default:
			return models.Country{}, problems.RaiseInternal(ctx, err)
		}
	}

	return countryModel(country), nil
}

// Read methods for countries

func (a Country) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.country.New().FilterID(ID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Country{}, problems.RaiseCountryNotFoundByID(
				ctx,
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		}
		return models.Country{}, problems.RaiseInternal(ctx, err)
	}

	return countryModel(country), nil
}

func (a Country) SearchCountries(ctx context.Context, name string, status string, pag pagi.Request) ([]models.Country, pagi.Response, error) {
	limit, offset := pagi.CalculateLimitOffset(pag)

	countries, err := a.country.New().
		FilterName(name).
		FilterStatus(status).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(ctx, err)
	}

	total, err := a.country.New().Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(ctx, err)
	}

	res, pagRes := countriesArray(countries, limit, pag.Page, total)

	return res, pagRes, nil
}

// Update methods for countries

func (a Country) UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status string) (models.Country, error) {
	//TODO in future create event to kafka about country status change

	_, err := enum.ParseCountryStatus(status)
	if err != nil {
		return models.Country{}, problems.RaiseInvalidCountryStatus(
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
				return problems.RaiseCountryNotFoundByID(
					ctx,
					fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
					ID,
				)
			default:
				return problems.RaiseInternal(ctx, err)
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

func (a Country) UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error) {
	country, err := a.countriesQ.New().FilterName(name).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// country with this name does not exist, so we can update it
		default:
			return models.Country{}, problems.RaiseInternal(ctx, err)
		}
	}
	if err == nil && country.ID != ID {
		return models.Country{}, problems.RaiseCountryAlreadyExists(
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
			return models.Country{}, problems.RaiseCountryNotFoundByID(
				ctx,
				fmt.Errorf("country with id '%s' does not exist, cause: %s", ID, err),
				ID,
			)
		default:
			return models.Country{}, problems.RaiseInternal(ctx, err)
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

func countriesArray(countries []dbx.Country, limit, page, total uint64) ([]models.Country, pagi.Response) {
	result := make([]models.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, countryModel(country))
	}

	pag := pagi.Response{
		Page:  page,
		Size:  limit,
		Total: total,
	}

	return result, pag
}
