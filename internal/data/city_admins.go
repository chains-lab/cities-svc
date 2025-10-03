package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/data/pgdb"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateCityAdmin(ctx context.Context, cityMod models.CityAdmin) error {
	return d.sql.cityMod.New().Insert(ctx, cityAdminModelToSchema(cityMod))
}

func (d *Database) GetCityAdmin(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error) {
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
		return models.CityAdmin{}, nil
	case err != nil:
		return models.CityAdmin{}, err
	}

	return cityAdminSchemaToModel(row), nil
}

func (d *Database) GetCityAdminByUserAndCityID(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error) {
	return d.GetCityAdmin(ctx, admin.GetFilters{
		UserID: &userID,
		CityID: &cityID,
	})
}

func (d *Database) GetCityAdminByUserID(ctx context.Context, userID uuid.UUID) (models.CityAdmin, error) {
	return d.GetCityAdmin(ctx, admin.GetFilters{
		UserID: &userID,
	})
}

func (d *Database) FilterCityAdmins(
	ctx context.Context,
	filter admin.FilterParams,
	page, size uint64,
) (models.CityAdminsCollection, error) {
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
		return models.CityAdminsCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CityAdminsCollection{}, err
	}

	res := make([]models.CityAdmin, len(rows))
	for i, r := range rows {
		res[i] = cityAdminSchemaToModel(r)
	}

	return models.CityAdminsCollection{
		Data:  res,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateCityAdmin(
	ctx context.Context,
	userID uuid.UUID,
	params admin.UpdateParams,
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

func (d *Database) DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error {
	return d.sql.cityMod.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
}

func cityAdminSchemaToModel(s pgdb.CityAdmin) models.CityAdmin {
	res := models.CityAdmin{
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

func cityAdminModelToSchema(m models.CityAdmin) pgdb.CityAdmin {
	s := pgdb.CityAdmin{
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
