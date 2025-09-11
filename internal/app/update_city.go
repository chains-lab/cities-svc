package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateCityParams struct {
	Name     *string
	Status   *string
	Point    *orb.Point
	Icon     *string
	Slug     *string
	Timezone *string
}

func (a App) UpdateCity(ctx context.Context, cityID uuid.UUID, params UpdateCityParams) (models.City, error) {
	update := city.UpdateCityParams{}

	if params.Point != nil {
		update.Point = params.Point
	}
	if params.Slug != nil {
		update.Slug = params.Slug
	}
	if params.Name != nil {
		update.Name = params.Name
	}
	if params.Timezone != nil {
		update.Timezone = params.Timezone
	}
	if params.Icon != nil {
		update.Icon = params.Icon
	}
	if params.Status != nil {
		update.Status = params.Status
	}

	return a.cities.UpdateOne(ctx, cityID, update)
}
