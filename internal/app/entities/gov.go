package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type Gov struct {
	govQ dbx.GovQ
}

func NewGov(db *sql.DB) Gov {
	return Gov{
		govQ: dbx.NewCityGovQ(db),
	}
}

type CreateGovParams struct {
	UserID uuid.UUID
	CityID uuid.UUID

	Role  string
	Label string
}

func (g Gov) CreateGov(ctx context.Context, params CreateGovParams) (models.Gov, error) {
	err := constant.CheckCityGovRole(params.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}

	now := time.Now().UTC()
	ID := uuid.New()

	stmt := dbx.Gov{
		ID:        ID,
		UserID:    params.UserID,
		CityID:    params.CityID,
		Status:    constant.GovStatusActive,
		Role:      params.Role,
		Label:     params.Label,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.Gov{
		ID:        ID,
		UserID:    params.UserID,
		CityID:    params.CityID,
		Status:    constant.GovStatusActive,
		Role:      params.Role,
		Label:     params.Label,
		UpdatedAt: now,
		CreatedAt: now,
	}

	err = g.govQ.New().Insert(ctx, stmt)
	if err != nil {
		return models.Gov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city gov: %w", err),
		)
	}

	return resp, nil
}

type GetGovFilters struct {
	ID     *uuid.UUID
	UserID *uuid.UUID
	CityID *uuid.UUID
	Status *string
	Role   *string
}

func (g Gov) Get(ctx context.Context, filters GetGovFilters) (models.Gov, error) {
	query := g.govQ.New()
	if filters.ID != nil {
		query = query.FilterID(*filters.ID)
	}
	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil {
		err := constant.CheckCityGovRole(*filters.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidCityGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}
	if filters.Status != nil {
		err := constant.CheckGovStatus(*filters.Status)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidGovStatus.Raise(
				fmt.Errorf("invalid city gov status, cause: %w", err),
			)
		}
		query = query.FilterStatus(*filters.Status)
	}

	gov, err := query.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Gov{}, errx.ErrorCityGovNotFound.Raise(
				fmt.Errorf("city gov not found, cause: %w", err),
			)
		default:
			return models.Gov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov, cause: %w", err),
			)
		}
	}

	return govFromDb(gov), nil
}

type SelectGovsFilters struct {
	UserID *uuid.UUID
	CityID *uuid.UUID
	Status []string
	Role   []string
}

func (g Gov) SelectGovs(
	ctx context.Context,
	filters SelectGovsFilters,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Gov, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1 // +1 чтобы определить наличие next
	offset := (pag.Page - 1) * pag.Size

	query := g.govQ.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Status != nil && len(filters.Status) > 0 {
		for _, s := range filters.Status {
			err := constant.CheckGovStatus(s)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidGovStatus.Raise(
					fmt.Errorf("invalid city gov status, cause: %w", err),
				)
			}
		}

		query = query.FilterStatus(filters.Status...)
	}
	if filters.Role != nil && len(filters.Role) > 0 {
		for _, r := range filters.Role {
			err := constant.CheckCityGovRole(r)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidCityGovRole.Raise(
					fmt.Errorf("invalid city gov role, cause: %w", err),
				)
			}
		}
		query = query.FilterRole(filters.Role...)
	}

	for _, s := range sort {
		switch s.Field {
		case "role":
			query = query.OrderByRole(s.Ascend)
		case "created_at":
			query = query.OrderByCreatedAt(s.Ascend)
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count city govs, cause: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to select city govs, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	res := make([]models.Gov, 0, len(rows))
	for _, cg := range rows {
		res = append(res, govFromDb(cg))
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

type UpdateGovParams struct {
	Status    *string
	Role      *string
	Label     *string
	UpdatedAt time.Time
}

func (g Gov) UpdateOne(ctx context.Context, ID uuid.UUID, params UpdateGovParams) error {
	if (params.Role == nil) && (params.Label == nil) && (params.Status == nil) {
		return nil
	}

	stmt := dbx.UpdateCityGovParams{}

	if params.Status != nil {
		err := constant.CheckGovStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidGovStatus.Raise(
				fmt.Errorf("invalid city gov status, cause: %w", err),
			)
		}
		stmt.Status = params.Status
	}

	if params.Role != nil {
		err := constant.CheckCityGovRole(*params.Role)
		if err != nil {
			return errx.ErrorInvalidCityGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		stmt.Role = params.Role
	}

	if params.Label != nil {
		stmt.Label = params.Label
	}

	if params.Status != nil && *params.Status == constant.GovStatusInactive {
		now := sql.NullTime{Time: time.Now().UTC(), Valid: true}
		stmt.DeactivatedAt = &now
	}

	stmt.UpdatedAt = &params.UpdatedAt

	err := g.govQ.New().FilterID(ID).Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city gov, cause: %w", err),
		)
	}

	return nil
}

type UpdateGovsFilters struct {
	CityID    *uuid.UUID
	UserID    *uuid.UUID
	CountryID *uuid.UUID
	Status    *string
	Role      *string
}

type UpdateGovsParams struct {
	Status *string
	Role   *string
}

func (g Gov) UpdateMany(ctx context.Context, filters UpdateGovsFilters, params UpdateGovsParams) error {
	if (params.Status == nil) && (params.Role == nil) {
		return nil
	}

	stmt := dbx.UpdateCityGovParams{}
	if params.Status != nil {
		err := constant.CheckGovStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidGovStatus.Raise(
				fmt.Errorf("invalid city gov status, cause: %w", err),
			)
		}

		stmt.Status = params.Status
		if *params.Status == constant.GovStatusInactive {
			now := sql.NullTime{Time: time.Now().UTC(), Valid: true}
			stmt.DeactivatedAt = &now
		}
	}

	if params.Role != nil {
		err := constant.CheckCityGovRole(*params.Role)
		if err != nil {
			return errx.ErrorInvalidCityGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}

		stmt.Role = params.Role
	}

	now := time.Now().UTC()
	stmt.UpdatedAt = &now

	query := g.govQ.New()

	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Status != nil {
		query = query.FilterStatus(*filters.Status)
	}
	if filters.Role != nil {
		err := constant.CheckCityGovRole(*filters.Role)
		if err != nil {
			return errx.ErrorInvalidCityGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}

	err := query.Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city govs, cause: %w", err),
		)
	}

	return nil
}

func govFromDb(g dbx.Gov) models.Gov {
	res := models.Gov{
		ID:        g.ID,
		UserID:    g.UserID,
		CityID:    g.CityID,
		Status:    g.Status,
		Role:      g.Role,
		Label:     g.Label,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.DeactivatedAt.Valid {
		res.DeactivatedAt = &g.DeactivatedAt.Time
	}

	return res
}
