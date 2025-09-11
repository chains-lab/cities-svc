package app

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/entities/country"
	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/google/uuid"
)

func (a App) SetCountryStatusSupported(ctx context.Context, countryID uuid.UUID) (models.Country, error) {
	c, err := a.country.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	countryStatus := constant.CountryStatusSupported

	cou, err := a.country.Update(ctx, c.ID, country.UpdateCountryParams{
		Status: &countryStatus,
	})
	if err != nil {
		return models.Country{}, err
	}

	//TODO in future kafka event about country status change

	return cou, nil
}

func (a App) SetCountryStatusDeprecated(ctx context.Context, countryID uuid.UUID) (models.Country, error) {
	c, err := a.country.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	updatedAt := time.Now().UTC()
	countryStatus := constant.CountryStatusDeprecated
	cityStatus := constant.CountryStatusDeprecated

	txErr := a.transaction(func(ctx context.Context) error {
		if _, err = a.country.Update(ctx, c.ID, country.UpdateCountryParams{
			Status: &countryStatus,
		}); err != nil {
			return err
		}

		err = a.cities.UpdateMany(ctx, city.UpdateCitiesFilters{
			CountryID: &c.ID,
			Status:    []string{constant.CityStatusOfficial, constant.CityStatusCommunity},
		}, city.UpdateCitiesParams{
			Status:    &cityStatus,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}

		err = a.gov.DeleteMany(ctx, gov.DeleteGovsFilters{
			CountryID: &c.ID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.Country{}, txErr
	}

	//TODO in future kafka event about country status change

	return models.Country{
		ID:        c.ID,
		Name:      c.Name,
		Status:    countryStatus,
		CreatedAt: c.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}
