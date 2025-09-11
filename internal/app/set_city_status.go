package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (a App) SetCityStatusOfficial(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	status := enum.CityStatusOfficial

	var cou models.City
	var err error
	txErr := a.transaction(func(ctx context.Context) error {
		cou, err = a.cities.UpdateOne(ctx, cityID, city.UpdateCityParams{
			Status: &status,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.City{}, txErr
	}

	return cou, nil
}

func (a App) SetCityStatusCommunity(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	status := enum.CityStatusCommunity

	var cou models.City
	var err error
	txErr := a.transaction(func(ctx context.Context) error {
		cou, err = a.cities.UpdateOne(ctx, cityID, city.UpdateCityParams{
			Status: &status,
		})
		if err != nil {
			return err
		}

		err = a.gov.DeleteMany(ctx, gov.DeleteGovsFilters{
			CityID: &cityID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.City{}, txErr
	}

	return cou, nil
}

func (a App) SetCityStatusDeprecated(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	status := enum.CityStatusDeprecated

	var cou models.City
	var err error
	txErr := a.transaction(func(ctx context.Context) error {
		cou, err = a.cities.UpdateOne(ctx, cityID, city.UpdateCityParams{
			Status: &status,
		})
		if err != nil {
			return err
		}

		err = a.gov.DeleteMany(ctx, gov.DeleteGovsFilters{
			CityID: &cityID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.City{}, txErr
	}

	return cou, nil
}
