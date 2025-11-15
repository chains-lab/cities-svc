package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/repo/pgdb"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func (r *Repo) CreateCity(ctx context.Context, m models.City) (models.City, error) {
	schema := cityModelToSchema(m)

	err := r.sql.cities.New().Insert(ctx, schema)
	if err != nil {
		return models.City{}, err
	}

	return citySchemaToModel(schema), nil
}

func (r *Repo) GetCityByID(ctx context.Context, id uuid.UUID) (models.City, error) {
	row, err := r.sql.cities.New().FilterID(id).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.City{}, nil
	case err != nil:
		return models.City{}, err
	}

	return citySchemaToModel(row), nil
}

func (r *Repo) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	row, err := r.sql.cities.New().FilterSlug(slug).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.City{}, nil
	case err != nil:
		return models.City{}, err
	}

	return citySchemaToModel(row), nil
}

func (r *Repo) GetCityByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error) {
	row, err := r.sql.cities.New().FilterWithinRadiusMeters(point, radius).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.City{}, nil
	case err != nil:
		return models.City{}, err
	}

	return citySchemaToModel(row), nil
}

func (r *Repo) FilterCities(
	ctx context.Context,
	filter city.FilterParams,
	page, size uint64,
) (models.CitiesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := r.sql.cities.New()

	if filter.CountryID != nil {
		query.FilterCountryID(*filter.CountryID)
	}
	if filter.Name != nil {
		query.FilterNameLike(*filter.Name)
	}
	if filter.Status != nil {
		query.FilterStatus(*filter.Status)
	}
	if filter.Location != nil {
		query.FilterWithinRadiusMeters(filter.Location.Point, filter.Location.RadiusM)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CitiesCollection{}, err
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.CitiesCollection{}, err
	}

	cities := make([]models.City, 0, len(rows))
	for _, r := range rows {
		cities = append(cities, citySchemaToModel(r))
	}

	return models.CitiesCollection{
		Data:  cities,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (r *Repo) UpdateCity(
	ctx context.Context,
	cityID uuid.UUID,
	params city.UpdateParams,
	updatedAt time.Time,
) error {
	query := r.sql.cities.New().FilterID(cityID)

	if params.Name != nil {
		query = query.UpdateName(*params.Name)
	}
	if params.Slug != nil {
		query = query.UpdateSlug(sql.NullString{Valid: true, String: *params.Slug})
	}
	if params.Icon != nil {
		query = query.UpdateIcon(sql.NullString{Valid: true, String: *params.Icon})
	}
	if params.Timezone != nil {
		query = query.UpdateTimezone(*params.Timezone)
	}
	if params.Point != nil {
		query = query.UpdatePoint(*params.Point)
	}

	if params == (city.UpdateParams{}) {
		return nil
	}

	err := query.Update(ctx, updatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateCityStatus(ctx context.Context, id uuid.UUID, status string, updatedAt time.Time) error {
	return r.sql.cities.New().
		FilterID(id).
		UpdateStatus(status).
		Update(ctx, updatedAt)
}

func (r *Repo) DeleteadminForCity(ctx context.Context, cityID uuid.UUID) error {
	err := r.sql.cityAdmin.New().FilterCityID(cityID).Delete(ctx)
	return err
}

func citySchemaToModel(s pgdb.City) models.City {
	res := models.City{
		ID:        s.ID,
		CountryID: s.CountryID,
		Point:     s.Point,
		Status:    s.Status,
		Name:      s.Name,
		Icon:      s.Icon,
		Timezone:  s.Timezone,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	return res
}

func cityModelToSchema(m models.City) pgdb.City {
	res := pgdb.City{
		ID:        m.ID,
		CountryID: m.CountryID,
		Point:     m.Point,
		Status:    m.Status,
		Name:      m.Name,
		Icon:      m.Icon,
		Slug:      m.Slug,
		Timezone:  m.Timezone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return res
}
