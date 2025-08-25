package entities

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config/constant/enum"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Gov struct {
	queries *dbx.Queries
}

func NewCityGov(db *pgxpool.Pool) *Gov {
	return &Gov{queries: dbx.New(db)}
}

// Create methods for citygov

func (g Gov) CreateCityGov(ctx context.Context, cityID, userID uuid.UUID, role string) (models.CityGov, error) {
	r, err := enum.ParseCityGovRole(role)
	if err != nil {
		return models.CityGov{}, err
	}

	gov, err := g.queries.CreateCityGov(ctx, dbx.CreateCityGovParams{
		CityID: cityID,
		UserID: userID,
		Role:   dbx.CityGovRoles(r),
	})
	if err != nil {
		return models.CityGov{}, err
	}

	return models.CityGov{
		UserID:    gov.UserID,
		CityID:    gov.CityID,
		Role:      r,
		CreatedAt: gov.CreatedAt.Time,
		UpdatedAt: gov.UpdatedAt.Time,
	}, nil
}

// Get methods for citygov

func (g Gov) GetCityGovAdmin(ctx context.Context, cityID uuid.UUID) (models.CityGov, error) {
	cityAdmin, err := g.queries.GetCityAdmin(ctx, cityID)
	if err != nil {
		return models.CityGov{}, err
	}

	return models.CityGov{
		UserID:    cityAdmin.UserID,
		CityID:    cityAdmin.CityID,
		Role:      string(cityAdmin.Role),
		CreatedAt: cityAdmin.CreatedAt.Time,
		UpdatedAt: cityAdmin.UpdatedAt.Time,
	}, nil
}

func (g Gov) GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error) {
	cityGov, err := g.queries.GetCityGov(ctx, dbx.GetCityGovParams{
		CityID: cityID,
		UserID: userID,
	})
	if err != nil {
		return models.CityGov{}, err
	}

	return models.CityGov{
		UserID:    cityGov.UserID,
		CityID:    cityGov.CityID,
		Role:      string(cityGov.Role),
		CreatedAt: cityGov.CreatedAt.Time,
		UpdatedAt: cityGov.UpdatedAt.Time,
	}, nil
}

func (g Gov) GetInitiatorCityGov(ctx context.Context, cityID, initiatorID uuid.UUID) (models.CityGov, error) {
	initiator, err := g.queries.GetCityGov(ctx, dbx.GetCityGovParams{
		CityID: cityID,
		UserID: initiatorID,
	})
	if err != nil {
		return models.CityGov{}, err
	}

	return models.CityGov{
		UserID:    initiator.UserID,
		CityID:    initiator.CityID,
		Role:      string(initiator.Role),
		CreatedAt: initiator.CreatedAt.Time,
		UpdatedAt: initiator.UpdatedAt.Time,
	}, nil
}

func (g Gov) SelectCityGovs(
	ctx context.Context,
	cityID uuid.UUID,
	pag *pagi.Request,
) ([]models.CityGov, pagi.Response, error) {
	cityGovs, err := g.queries.SelectCityGovs(ctx, dbx.SelectCityGovsParams{
		PageSize: int64(pag.Size),
		Page:     int64(pag.Page),
		CityID:   cityID,
	})
	if err != nil {
		return nil, pagi.Response{}, err
	}

	res := make([]models.CityGov, 0, len(cityGovs))
	total := 0
	if len(cityGovs) > 0 {
		total = int(cityGovs[0].TotalCount)
	}

	for _, cg := range cityGovs {
		res = append(res, models.CityGov{
			UserID:    cg.UserID,
			CityID:    cg.CityID,
			Role:      string(cg.Role),
			CreatedAt: cg.CreatedAt.Time,
			UpdatedAt: cg.UpdatedAt.Time,
		})
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: uint64(total),
	}, nil
}

// Delete methods for citygov

func (g Gov) DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error {
	err := g.queries.DeleteCityGov(ctx, dbx.DeleteCityGovParams{
		CityID: cityID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	return nil
}
