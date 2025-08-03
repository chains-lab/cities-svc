package entities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/ape"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
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
	FilterStatus(status string) dbx.CountriesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CountriesQ
}

type Countries struct {
	queries countriesQ
}

func NewCountries(db *sql.DB) (Countries, error) {
	return Countries{
		queries: dbx.NewCountries(db),
	}, nil
}

func (c Countries) Create(ctx context.Context, name, status string) error {
	if enum.CheckCountryStatus(status) == false {
		return ape.RaiseInvalidCountryStatus(status)
	}

	country := dbx.CountryModel{
		ID:        uuid.New(),
		Name:      name,
		Status:    status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := c.queries.New().Insert(ctx, country)
	if err != nil {
		switch {
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

type UpdateCountryInput struct {
	Name      *string   `db:"name"`
	Status    *string   `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c Countries) Update(ctx context.Context, input UpdateCountryInput) error {
	if input.Status != nil && enum.CheckCountryStatus(*input.Status) == false {
		return ape.RaiseInvalidCountryStatus(*input.Status)
	}

	if input.UpdatedAt.IsZero() {
		input.UpdatedAt = time.Now().UTC()
	}

	err := c.queries.New().Update(ctx, dbx.UpdateCountryInput{
		Name:      input.Name,
		Status:    input.Status,
		UpdatedAt: input.UpdatedAt,
	})
	if err != nil {
		switch {
		default:
			return ape.RaiseInternal(err)
		}
	}

	return nil
}

func (c Countries) GetByID(ctx context.Context, id uuid.UUID) (models.CountryModel, error) {
	country, err := c.queries.New().FilterID(id).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CountryModel{}, ape.RaiseCountryNotFoundByID(err, id)
		}
		return models.CountryModel{}, ape.RaiseInternal(err)
	}

	return countryDbxToModel(country), nil
}

func (c Countries) GetByName(ctx context.Context, name string) (models.CountryModel, error) {
	country, err := c.queries.New().FilterName(name).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CountryModel{}, ape.RaiseCountryNotFoundByName(err, name)
		}
		return models.CountryModel{}, ape.RaiseInternal(err)
	}

	return countryDbxToModel(country), nil
}

func (c Countries) Select(ctx context.Context, page, limit uint64) ([]models.CountryModel, error) {
	countries, err := c.queries.New().Page(page, limit).Select(ctx)
	if err != nil {
		return nil, ape.RaiseInternal(err)
	}

	res := make([]models.CountryModel, len(countries))

	for i, country := range countries {
		res[i] = countryDbxToModel(country)
	}

	return res, nil
}

func (c Countries) SelectWithStatus(ctx context.Context, status string, page, limit uint64) ([]models.CountryModel, error) {
	if !enum.CheckCountryStatus(status) {
		return nil, ape.RaiseInvalidCountryStatus(status)
	}

	countries, err := c.queries.New().FilterStatus(status).Page(page, limit).Select(ctx)
	if err != nil {
		return nil, ape.RaiseInternal(err)
	}

	res := make([]models.CountryModel, len(countries))

	for i, country := range countries {
		res[i] = countryDbxToModel(country)
	}

	return res, nil
}

func countryDbxToModel(country dbx.CountryModel) models.CountryModel {
	return models.CountryModel{
		ID:        country.ID,
		Name:      country.Name,
		Status:    country.Status,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}
}
