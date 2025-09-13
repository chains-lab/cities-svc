package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateCityParams struct {
	Name      string
	CountryID uuid.UUID
	Point     orb.Point
	Timezone  string
}

func (a App) CreateCity(ctx context.Context, params CreateCityParams) (models.City, error) {
	country, err := a.GetCountryByID(ctx, params.CountryID)
	if err != nil {
		return models.City{}, err
	}

	if country.Status != enum.CountryStatusSupported {
		return models.City{}, errx.ErrorCountryNotSupported.Raise(
			fmt.Errorf("country status is not 'supported'"),
		)
	}

	return a.cities.Create(ctx, city.CreateCityParams{
		Name:      params.Name,
		CountryID: params.CountryID,
		Point:     params.Point,
		Timezone:  params.Timezone,
	})
}
