package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/domain/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) GetInitiatorGov(ctx context.Context, initiatorID uuid.UUID) (models.Gov, error) {
	return a.gov.GetInitiatorGov(ctx, initiatorID)
}

func (a App) GetGov(ctx context.Context, userID uuid.UUID) (models.Gov, error) {
	return a.gov.GetGov(ctx, gov.GetGovFilters{
		UserID: &userID,
	})
}

type FiltersListGovsParams struct {
	CityID *uuid.UUID
	Roles  []string
}

func (a App) ListGovs(
	ctx context.Context,
	filters FiltersListGovsParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Gov, pagi.Response, error) {
	input := gov.FiltersListParams{}
	if filters.CityID != nil {
		input.CityID = filters.CityID
	}
	if len(filters.Roles) > 0 && filters.Roles != nil {
		input.Role = filters.Roles
	}

	return a.gov.ListGovs(ctx, input, pag, sort)
}

type UpdateOwnGovParams struct {
	Label *string
}

func (a App) UpdateOwnActiveGov(ctx context.Context, userID uuid.UUID, params UpdateOwnGovParams) (models.Gov, error) {
	g, err := a.GetInitiatorGov(ctx, userID)
	if err != nil {
		return models.Gov{}, err
	}

	entitiesParams := gov.UpdateGovParams{}
	if params.Label != nil {
		entitiesParams.Label = params.Label
	}

	return a.gov.UpdateOne(ctx, g.UserID, entitiesParams)
}

func (a App) DeleteGov(ctx context.Context, initiatorID, userID, cityID uuid.UUID) error {
	return a.gov.Delete(ctx, initiatorID, userID, cityID)
}

func (a App) RefuseOwnGov(ctx context.Context, userID uuid.UUID) error {
	err := a.gov.RefuseOwnGov(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
