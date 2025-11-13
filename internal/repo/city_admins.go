package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/repo/pgdb"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (r *Repo) CreateCityAdmin(ctx context.Context, cityMod models.CityAdmin) error {
	return r.sql.cityAdmin.New().Insert(ctx, CityAdminModelToSchema(cityMod))
}

//func (d *Repo) GetCityAdmin(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error) {
//	query := d.sql.cityAdmin.New()
//
//	if filters.UserID != nil {
//		query = query.FilterUserID(*filters.UserID)
//	}
//	if filters.CityID != nil {
//		query = query.FilterCityID(*filters.CityID)
//	}
//	if filters.Role != nil {
//		query = query.FilterRole(*filters.Role)
//	}
//
//	row, err := query.Get(ctx)
//	switch {
//	case errors.Is(err, sql.ErrNoRows):
//		return models.CityAdmin{}, nil
//	case err != nil:
//		return models.CityAdmin{}, err
//	}
//
//	return CityAdminSchemaToModel(row), nil
//}
//
//func (d *Repo) GetCityAdminByUserAndCityID(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error) {
//	return d.GetCityAdmin(ctx, admin.GetFilters{
//		UserID: &userID,
//		CityID: &cityID,
//	})
//}

func (r *Repo) GetCityAdminWithFilter(
	ctx context.Context,
	userID, cityID *uuid.UUID,
	role *string,
) (models.CityAdmin, error) {
	query := r.sql.cityAdmin.New()

	if userID != nil {
		query = query.FilterUserID(*userID)
	}
	if cityID != nil {
		query = query.FilterCityID(*cityID)
	}
	if role != nil {
		query = query.FilterRole(*role)
	}

	row, err := query.Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.CityAdmin{}, nil
	case err != nil:
		return models.CityAdmin{}, err
	}

	return CityAdminSchemaToModel(row), nil
}

func (r *Repo) GetCityAdminByUserID(ctx context.Context, userID uuid.UUID) (models.CityAdmin, error) {
	schemas, err := r.sql.cityAdmin.New().FilterUserID(userID).Get(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return models.CityAdmin{}, nil
	}

	if err != nil {
		return models.CityAdmin{}, err
	}
	return CityAdminSchemaToModel(schemas), nil
}

func (r *Repo) GetCityAdmins(ctx context.Context, cityID uuid.UUID, roles ...string) (models.CityAdminsCollection, error) {
	query := r.sql.cityAdmin.New().FilterCityID(cityID)
	if len(roles) > 0 {
		query = query.FilterRole(roles...)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return models.CityAdminsCollection{}, err
	}

	res := make([]models.CityAdmin, len(rows))
	for i, r := range rows {
		res[i] = CityAdminSchemaToModel(r)
	}

	return models.CityAdminsCollection{
		Data:  res,
		Page:  1,
		Size:  uint64(len(res)),
		Total: uint64(len(res)),
	}, nil
}

func (r *Repo) FilterCityAdmins(
	ctx context.Context,
	filter admin.FilterParams,
	page, size uint64,
) (models.CityAdminsCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := r.sql.cityAdmin.New()

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
		res[i] = CityAdminSchemaToModel(r)
	}

	return models.CityAdminsCollection{
		Data:  res,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (r *Repo) UpdateCityAdmin(
	ctx context.Context,
	userID uuid.UUID,
	params admin.UpdateParams,
	updatedAt time.Time,
) error {
	q := r.sql.cityAdmin.New().FilterUserID(userID)

	if params.Label != nil {
		switch *params.Label {
		case "":
			q.UpdateLabel(sql.NullString{Valid: false})
		default:
			q.UpdateLabel(sql.NullString{String: *params.Label, Valid: true})
		}
	}
	if params.Position != nil {
		switch *params.Position {
		case "":
			q.UpdatePosition(sql.NullString{Valid: false})
		default:
			q.UpdatePosition(sql.NullString{String: *params.Position, Valid: true})
		}
	}

	return q.Update(ctx, updatedAt)
}

func (r *Repo) DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error {
	return r.sql.cityAdmin.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
}

func (r *Repo) DeleteAdminsForCity(ctx context.Context, cityID uuid.UUID) error {
	return r.sql.cityAdmin.New().FilterCityID(cityID).Delete(ctx)
}

func CityAdminSchemaToModel(s pgdb.CityAdmin) models.CityAdmin {
	res := models.CityAdmin{
		UserID:    s.UserID,
		CityID:    s.CityID,
		Role:      s.Role,
		Position:  s.Position,
		Label:     s.Label,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	return res
}

func CityAdminModelToSchema(m models.CityAdmin) pgdb.CityAdmin {
	s := pgdb.CityAdmin{
		UserID:    m.UserID,
		CityID:    m.CityID,
		Role:      m.Role,
		Position:  m.Position,
		Label:     m.Label,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return s
}
