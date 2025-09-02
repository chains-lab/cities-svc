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
	Role   string
}

func (g Gov) CreateGov(ctx context.Context, params CreateGovParams) (models.Gov, error) {
	err := constant.CheckCityGovRole(params.Role)
	if err != nil {
		return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	stmt := dbx.Gov{
		UserID:    params.UserID,
		CityID:    params.CityID,
		Role:      params.Role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.Gov{
		UserID:    params.UserID,
		CityID:    params.CityID,
		Role:      params.Role,
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
	UserID *uuid.UUID
	CityID *uuid.UUID
	Role   *string
}

func (g Gov) Get(ctx context.Context, filters GetGovFilters) (models.Gov, error) {
	query := g.govQ.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil {
		err := constant.CheckCityGovRole(*filters.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
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
	CityID *uuid.UUID
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

	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.Role != nil && len(filters.Role) > 0 {
		for _, r := range filters.Role {
			err := constant.CheckCityGovRole(r)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidGovRole.Raise(
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
	Role      *string
	Label     *string
	UpdatedAt time.Time
}

func (g Gov) UpdateOne(ctx context.Context, userID uuid.UUID, params UpdateGovParams) (models.Gov, error) {
	if (params.Role == nil) && (params.Label == nil) {
		return models.Gov{}, nil
	}

	stmt := dbx.UpdateCityGovParams{}

	if params.Role != nil {
		err := constant.CheckCityGovRole(*params.Role)
		if err != nil {
			return models.Gov{}, errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		stmt.Role = params.Role
	}

	if params.Label != nil {
		if *params.Label != "" {
			stmt.Label.String = *params.Label
		} else {
			stmt.Label.Valid = false
		}
	}

	stmt.UpdatedAt = &params.UpdatedAt

	err := g.govQ.New().FilterUserID(userID).Update(ctx, stmt)
	if err != nil {
		return models.Gov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city gov, cause: %w", err),
		)
	}

	return g.Get(ctx, GetGovFilters{UserID: &userID})
}

func (g Gov) DeleteOne(ctx context.Context, userID uuid.UUID) error {
	err := g.govQ.New().FilterUserID(userID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city gov, cause: %w", err),
		)
	}

	return nil
}

type DeleteGovsFilters struct {
	UserID    *uuid.UUID
	CityID    *uuid.UUID
	CountryID *uuid.UUID
	Role      *string
}

func (g Gov) DeleteMany(ctx context.Context, filters DeleteGovsFilters) error {
	query := g.govQ.New()

	if filters.UserID != nil {
		query = query.FilterUserID(*filters.UserID)
	}
	if filters.CityID != nil {
		query = query.FilterCityID(*filters.CityID)
	}
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Role != nil {
		err := constant.CheckCityGovRole(*filters.Role)
		if err != nil {
			return errx.ErrorInvalidGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
		query = query.FilterRole(*filters.Role)
	}

	err := query.Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city govs, cause: %w", err),
		)
	}

	return nil
}

func govFromDb(g dbx.Gov) models.Gov {
	res := models.Gov{
		UserID:    g.UserID,
		CityID:    g.CityID,
		Role:      g.Role,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.Label.Valid {
		res.Label = &g.Label.String
	}

	return res
}
