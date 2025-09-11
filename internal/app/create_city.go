package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
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
	return a.cities.Create(ctx, city.CreateCityParams{
		Name:      params.Name,
		CountryID: params.CountryID,
		Point:     params.Point,
		Timezone:  params.Timezone,
	})
}
