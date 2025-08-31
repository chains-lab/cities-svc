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

func CreateCityGovEntity(db *sql.DB) Gov {
	return Gov{
		govQ: dbx.NewCityGovQ(db),
	}
}

type CreateGovParams struct {
	Role  string
	Label *string
}

func (g Gov) CreateGov(ctx context.Context, userID, cityID uuid.UUID, params CreateGovParams) (models.CityGov, error) {
	err := constant.ParseCityGovRole(params.Role)
	if err != nil {
		return models.CityGov{}, errx.ErrorInvalidCityGovRole.Raise(
			fmt.Errorf("invalid city gov role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	stmt := dbx.CityGov{
		UserID:    userID,
		CityID:    cityID,
		Role:      params.Role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	resp := models.CityGov{
		UserID:    userID,
		CityID:    cityID,
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

// Get methods for citygov

func (g Gov) GetForCity(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	cityGov, err := g.govQ.New().FilterCityID(cityID).FilterUserID(userID).Get(ctx)
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

	return models.CityGov{
		UserID:    cityGov.UserID,
		CityID:    cityGov.CityID,
		Role:      cityGov.Role,
		CreatedAt: cityGov.CreatedAt,
		UpdatedAt: cityGov.UpdatedAt,
	}, nil
}

func (g Gov) Get(ctx context.Context, userID uuid.UUID) (models.CityGov, error) {
	cityGov, err := g.govQ.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorCityGovNotFound.Raise(
				fmt.Errorf("city gov not found user_id: %s, cause: %w", userID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov user_id: %s, cause: %w", userID, err),
			)
		}
	}

	return models.CityGov{
		UserID:    cityGov.UserID,
		CityID:    cityGov.CityID,
		Role:      cityGov.Role,
		CreatedAt: cityGov.CreatedAt,
		UpdatedAt: cityGov.UpdatedAt,
	}, nil
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

	return models.CityGov{
		UserID:    initiator.UserID,
		CityID:    initiator.CityID,
		Role:      initiator.Role,
		CreatedAt: initiator.CreatedAt,
		UpdatedAt: initiator.UpdatedAt,
	}, nil
}

func (g Gov) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityGov, error) {
	initiator, err := g.govQ.New().FilterUserID(initiatorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, errx.ErrorInitiatorIsNotCityGov.Raise(
				fmt.Errorf("user_id: %s, not city gov, cause: %w", initiatorID, err),
			)
		default:
			return models.CityGov{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city gov for user_id %s, cause: %w", initiatorID, err),
			)
		}
	}

	return models.CityGov{
		UserID:    initiator.UserID,
		CityID:    initiator.CityID,
		Role:      initiator.Role,
		CreatedAt: initiator.CreatedAt,
		UpdatedAt: initiator.UpdatedAt,
	}, nil
}

type SelectGovsParams struct {
	CityID *uuid.UUID
	Role   []string
	Label  *string
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

	if params.CityID != nil {
		query = query.FilterCityID(*params.CityID)
	}

	if params.Label != nil {
		query = query.FilterLabelLike(*params.Label)
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
		res = append(res, models.CityGov{
			UserID:    cg.UserID,
			CityID:    cg.CityID,
			Role:      cg.Role,
			CreatedAt: cg.CreatedAt,
			UpdatedAt: cg.UpdatedAt,
		})
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

// Delete methods for citygov

func (g Gov) DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error {
	err := g.govQ.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city gov by cityID: %s, userID: %s, cause: %w", cityID, userID, err),
		)
	}

	return nil
}
