package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type Gov struct {
	queries dbx.CityGovQ
}

func NewGov(db *sql.DB) Gov {
	return Gov{
		queries: dbx.NewCityGovQ(db),
	}
}

// Create methods for citygov

func (g Gov) CreateCityGov(ctx context.Context, cityID, userID uuid.UUID, role string) (models.CityGov, error) {
	r, err := enum.ParseCityGovRole(role)
	if err != nil {
		return models.CityGov{}, problems.RaiseInvalidCityGovRole(
			fmt.Errorf("invalid city gov role: %w", err),
			fmt.Sprintf("invalid city gov role: %s", role),
		)
	}

	now := time.Now().UTC()

	err = g.queries.New().Insert(ctx, dbx.CityGov{
		UserID:    userID,
		CityID:    cityID,
		Role:      r,
		UpdatedAt: now,
		CreatedAt: now,
	})
	if err != nil {
		return models.CityGov{}, problems.RaiseInternal(
			fmt.Errorf("error creating city gov: %w", err),
		)
	}

	return models.CityGov{
		UserID:    userID,
		CityID:    cityID,
		Role:      r,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Get methods for citygov

func (g Gov) GetCityGovAdmin(ctx context.Context, cityID uuid.UUID) (models.CityGov, error) {
	cityAdmin, err := g.queries.New().FilterCityID(cityID).FilterRole(enum.CityGovRoleAdmin).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, problems.RaiseCityGovNotFound(
				fmt.Errorf("city admin not found: %w", err),
				fmt.Sprintf("city admin not found for cityID: %s", cityID),
			)
		default:
			return models.CityGov{}, problems.RaiseInternal(
				fmt.Errorf("get city admin by cityID: %s, cause: %w", cityID, err),
			)
		}
	}

	return models.CityGov{
		UserID:    cityAdmin.UserID,
		CityID:    cityAdmin.CityID,
		Role:      cityAdmin.Role,
		CreatedAt: cityAdmin.CreatedAt,
		UpdatedAt: cityAdmin.UpdatedAt,
	}, nil
}

func (g Gov) GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	cityGov, err := g.queries.New().FilterCityID(cityID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, problems.RaiseCityGovNotFound(
				fmt.Errorf("city gov not found by cityID: %s, userID: %s, cause: %w", cityID, userID, err),
				fmt.Sprintf("city gov not found for cityID: %s, userID: %s", cityID, userID),
			)
		default:
			return models.CityGov{}, problems.RaiseInternal(
				fmt.Errorf("get city gov by cityID: %s, userID: %s, cause: %w", cityID, userID, err),
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

func (g Gov) GetInitiatorCityGov(ctx context.Context, cityID, initiatorID uuid.UUID) (models.CityGov, error) {
	initiator, err := g.queries.New().FilterUserID(initiatorID).FilterCityID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.CityGov{}, problems.RaiseInitiatorIsNotCityGov(
				fmt.Errorf("initiator: %s, not city gov for city: %s, cause: %w", initiatorID, cityID, err),
				fmt.Sprintf("initiator: %s, not city gov for city: %s", initiatorID, cityID),
			)
		default:
			return models.CityGov{}, problems.RaiseInternal(
				fmt.Errorf("initiator: %s, not city gov: %w", initiatorID, err),
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

func (g Gov) SelectCityGovs(
	ctx context.Context,
	cityID uuid.UUID,
	pag pagi.Request,
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

	rows, err := g.queries.New().FilterCityID(cityID).Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(
			fmt.Errorf("select city gov: %w", err),
		)
	}

	prev := pag.Page > 1
	next := len(rows) > int(pag.Size)
	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
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
		Page: pag.Page,
		Size: pag.Size,
		Next: next,
		Prev: prev,
	}, nil
}

// Delete methods for citygov

func (g Gov) DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error {
	err := g.queries.New().FilterUserID(userID).FilterCityID(cityID).Delete(ctx)
	if err != nil {
		return problems.RaiseInternal(
			fmt.Errorf("delete city gov by cityID: %s, userID: %s, cause: %w", cityID, userID, err),
		)
	}

	return nil
}
