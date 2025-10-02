package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/data/pgdb"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateCityMod(ctx context.Context, cityMod models.CityModer) error {
	return d.sql.cityMod.New().Insert(ctx, cityModModelToSchema(cityMod))
}

func (d *Database) GetCityModer(ctx context.Context, filters citymod.GetFilters) (models.CityModer, error) {
	query := d.sql.cityMod.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil {
		query = query.FilterRole(*filters.Role)
	}

	row, err := query.Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.CityModer{}, nil
	case err != nil:
		return models.CityModer{}, err
	}

	return cityModSchemaToModel(row), nil
}

func (d *Database) FilterCityModers(
	ctx context.Context,
	filter citymod.FilterParams,
	page, size uint64,
) (models.CityModersCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.cityMod.New()

	if filter.CityID != nil {
		query.FilterCityID(*filter.CityID)
	}
	if filter.Roles != nil {
		query.FilterRole(filter.Roles...)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.CityModersCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CityModersCollection{}, err
	}

	res := make([]models.CityModer, len(rows))
	for i, r := range rows {
		res[i] = cityModSchemaToModel(r)
	}

	return models.CityModersCollection{
		Data:  res,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateCityModer(
	ctx context.Context,
	userID uuid.UUID,
	params citymod.UpdateCityModerParams,
	updatedAt time.Time,
) error {
	q := d.sql.cityMod.New().FilterUserID(userID)

	if params.Label != nil {
		switch *params.Label {
		case "":
			q.UpdateLabel(sql.NullString{Valid: false})
		default:
			q.UpdateLabel(sql.NullString{String: *params.Label, Valid: true})
		}
	}

	return q.Update(ctx, updatedAt)
}

func (d *Database) DeleteCityModer(ctx context.Context, userID, cityID uuid.UUID) error {
	return d.sql.cityMod.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
}

func cityModSchemaToModel(s pgdb.CityMod) models.CityModer {
	res := models.CityModer{
		UserID:    s.UserID,
		CityID:    s.CityID,
		Role:      s.Role,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if !s.Label.Valid {
		res.Label = &s.Label.String
	}

	return res
}

func cityModModelToSchema(m models.CityModer) pgdb.CityMod {
	s := pgdb.CityMod{
		UserID:    m.UserID,
		CityID:    m.CityID,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Label != nil {
		s.Label = sql.NullString{String: *m.Label, Valid: true}
	}

	return s
}
