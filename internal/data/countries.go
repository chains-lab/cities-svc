package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/data/pgdb"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateCountry(ctx context.Context, country models.Country) (models.Country, error) {
	schema := pgdb.Country{
		ID:        country.ID,
		Name:      country.Name,
		Status:    country.Status,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}

	err := d.sql.countries.New().Insert(ctx, schema)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}

func (d *Database) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	c, err := d.sql.countries.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Country{}, nil
	case err != nil:
		return models.Country{}, err
	}

	return countrySchemaToModel(c), err
}

func (d *Database) GetCountryByName(ctx context.Context, name string) (models.Country, error) {
	c, err := d.sql.countries.New().FilterName(name).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Country{}, nil
	case err != nil:
		return models.Country{}, err
	}

	return countrySchemaToModel(c), err
}

func (d *Database) FilterCountries(
	ctx context.Context,
	filters country.FilterParams,
	page, size uint64,
) (models.CountriesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.countries.New()

	if filters.Name != nil {
		query.FilterNameLike(*filters.Name)
	}
	if filters.Statuses != nil {
		query.FilterStatus(filters.Statuses...)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CountriesCollection{}, err
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.CountriesCollection{}, err
	}

	countries := make([]models.Country, 0, len(rows))
	for _, row := range rows {
		countries = append(countries, countrySchemaToModel(row))
	}

	return models.CountriesCollection{
		Data:  countries,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateCountry(
	ctx context.Context,
	countryID uuid.UUID,
	params country.UpdateParams,
	updatedAt time.Time,
) error {
	q := d.sql.countries.New().FilterID(countryID)

	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}

	return q.Update(ctx, updatedAt)
}

func (d *Database) UpdateCountryStatus(ctx context.Context, countryID uuid.UUID, status string, updatedAt time.Time) error {
	err := d.sql.countries.New().FilterID(countryID).UpdateStatus(status).Update(ctx, updatedAt)
	return err
}

func (d *Database) DeleteCitiesAndHerModersForCountry(ctx context.Context, countryID uuid.UUID) error {
	err := d.sql.cities.New().FilterCountryID(countryID).Delete(ctx)
	if err != nil {
		return err
	}

	err = d.sql.cityMod.New().FilterCountryID(countryID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateStatusForAllCountryCities(ctx context.Context, countryID uuid.UUID, status string, updatedAt time.Time) error {
	err := d.sql.cities.New().FilterCountryID(countryID).UpdateStatus(status).Update(ctx, updatedAt)
	if err != nil {
		return err
	}

	return nil
}

func countrySchemaToModel(row pgdb.Country) models.Country {
	return models.Country{
		ID:        row.ID,
		Name:      row.Name,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
