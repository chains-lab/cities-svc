package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/enum"
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

//func (a app) GetForCity(ctx context.Context, cityID, userID uuid.UUID) (models.Gov, error) {
//	return a.gov.GetGov(ctx, gov.GetGovFilters{
//		CityID: &cityID,
//		UserID: &userID,
//	})
//}
//
//func (a app) GetForCityAndRole(ctx context.Context, userID, cityID uuid.UUID, role string) (models.Gov, error) {
//	return a.gov.GetGov(ctx, gov.GetGovFilters{
//		UserID: &userID,
//		CityID: &cityID,
//		Role:   &role,
//	})
//}

func (a App) GetCityMayor(ctx context.Context, cityID uuid.UUID) (models.Gov, error) {
	role := enum.CityGovRoleMayor
	return a.gov.GetGov(ctx, gov.GetGovFilters{
		CityID: &cityID,
		Role:   &role,
	})
}
