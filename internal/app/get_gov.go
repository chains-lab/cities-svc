package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetInitiatorGov(ctx context.Context, initiatorID uuid.UUID) (models.Gov, error) {
	return a.gov.GetInitiatorGov(ctx, initiatorID)
}

func (a App) Get(ctx context.Context, userID uuid.UUID) (models.Gov, error) {
	return a.gov.Get(ctx, gov.GetGovFilters{
		UserID: &userID,
	})

}

func (a App) GetForCity(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
	return a.gov.Get(ctx, gov.GetGovFilters{
		CityID: &cityID,
		UserID: &userID,
	})
}

func (a App) GetForCityAndRole(ctx context.Context, userID, cityID uuid.UUID, role string) (models.Gov, error) {
	return a.gov.Get(ctx, gov.GetGovFilters{
		UserID: &userID,
		CityID: &cityID,
		Role:   &role,
	})
}
