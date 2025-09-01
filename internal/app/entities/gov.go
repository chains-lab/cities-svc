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
	govQ dbx.CityGovQ
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
	Label  *string
}

func (g Gov) CreateGov(ctx context.Context, params CreateGovParams) (models.CityGov, error) {
	err := constant.ParseCityGovRole(params.Role)
	if err != nil {
		return models.CityGov{}, errx.ErrorInvalidCityGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}

	now := time.Now().UTC()
	ID := uuid.New()

	stmt := dbx.CityGov{
		ID:        ID,
		UserID:    params.UserID,
		CityID:    params.CityID,
		Active:    true,
		Role:      params.Role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.CityGov{
		ID:        ID,
		UserID:    params.UserID,
		CityID:    params.CityID,
		Active:    true,
		Role:      params.Role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	if params.Label != nil {
		stmt.Label = sql.NullString{String: *params.Label, Valid: true}
		resp.Label = params.Label
	}

	err = g.govQ.New().Insert(ctx, stmt)
	if err != nil {
		return models.CityGov{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city gov: %w", err),
		)
	}

	return resp, nil
}

// GetByID methods for citygov

func (g Gov) GetForUserAndCity(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	gov, err := g.govQ.New().FilterCityID(cityID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorCityGovNotFound.Raise(
				fmt.Errorf("city gov not found by city_id: %s, user_id: %s, cause: %w", cityID, userID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov by city_id: %s, user_id: %s, cause: %w", cityID, userID, err),
			)
		}
	}

	return govFromDb(gov), nil
}

func (g Gov) Get(ctx context.Context, ID uuid.UUID) (models.CityGov, error) {
	gov, err := g.govQ.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorCityGovNotFound.Raise(
				fmt.Errorf("city gov not found by id: %s, cause: %w", ID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov by id: %s, cause: %w", ID, err),
			)
		}
	}

	return govFromDb(gov), nil
}

func (g Gov) GetInitiatorForCity(ctx context.Context, cityID, initiatorID uuid.UUID) (models.CityGov, error) {
	initiator, err := g.govQ.New().FilterUserID(initiatorID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorInitiatorIsNotCityGov.Raise(
				fmt.Errorf("user_id: %s, not city gov for city: %s, cause: %w", initiatorID, cityID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov for user_id %s city_id %s, cause: %w", initiatorID, cityID, err),
			)
		}
	}

	return govFromDb(initiator), nil
}

func (g Gov) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityGov, error) {
	initiator, err := g.govQ.New().FilterID(initiatorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorInitiatorIsNotCityGov.Raise(
				fmt.Errorf("id: %s, not city gov, cause: %w", initiatorID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov for id: %s, cause: %w", initiatorID, err),
			)
		}
	}

	return govFromDb(initiator), nil
}

type SelectGovsParams struct {
	UserID *uuid.UUID
	CityID *uuid.UUID
	Role   []string
	Label  *string
	Active *bool
}

func (g Gov) SelectGovs(
	ctx context.Context,
	params SelectGovsParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.CityGov, pagi.Response, error) {
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

	if params.UserID != nil {
		query = query.FilterUserID(*params.UserID)
	}
	if params.CityID != nil {
		query = query.FilterCityID(*params.CityID)
	}
	if params.Label != nil {
		query = query.FilterLabelLike(*params.Label)
	}
	if params.Active != nil {
		query = query.FilterActive(*params.Active)
	}
	if params.Role != nil && len(params.Role) > 0 {
		for _, r := range params.Role {
			err := constant.ParseCityGovRole(r)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidCityGovRole.Raise(
					fmt.Errorf("invalid city gov role, cause: %w", err),
				)
			}
		}
		query = query.FilterRole(params.Role...)
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

	res := make([]models.CityGov, 0, len(rows))
	for _, cg := range rows {
		res = append(res, govFromDb(cg))
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

// Delete methods for citygov

type UpdateGovParams struct {
	Active    *bool
	Role      *string
	Label     *string
	UpdatedAt time.Time
}

func (g Gov) UpdateOne(ctx context.Context, ID uuid.UUID, params UpdateGovParams) error {
	if (params.Active == nil) && (params.Role == nil) && (params.Label == nil) {
		return nil
	}

	stmt := dbx.UpdateCityGovParams{}
	if params.Active != nil {
		stmt.Active = params.Active
	}
	if params.Role != nil {
		err := constant.ParseCityGovRole(*params.Role)
		if err != nil {
			return errx.ErrorInvalidCityGovRole.Raise(
				fmt.Errorf("invalid city gov role, cause: %w", err),
			)
		}
	}
	if params.Label != nil {
		if *params.Label == "" {
			stmt.Label = &sql.NullString{Valid: false}
		} else {
			stmt.Label = &sql.NullString{String: *params.Label, Valid: true}
		}
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
	Role      *string
	Active    *bool
}

type UpdateGovsParams struct {
	Role   *string
	Active *bool
}

func (g Gov) UpdateMany(ctx context.Context, filters UpdateGovsFilters, params UpdateGovsParams) error {
	if (params.Active == nil) && (params.Role == nil) {
		return nil
	}

	stmt := dbx.UpdateCityGovParams{}
	if params.Active != nil {
		stmt.Active = params.Active
	}

	if params.Role != nil {
		err := constant.ParseCityGovRole(*params.Role)
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
	if filters.Active != nil {
		query = query.FilterActive(*filters.Active)
	}
	if filters.Role != nil {
		err := constant.ParseCityGovRole(*filters.Role)
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

func govFromDb(g dbx.CityGov) models.CityGov {
	res := models.CityGov{
		ID:        g.ID,
		UserID:    g.UserID,
		CityID:    g.CityID,
		Active:    g.Active,
		Role:      g.Role,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.Label.Valid {
		res.Label = &g.Label.String
	}

	return res
}
