package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error) {
	return a.cities.GetByID(ctx, ID)
}

func (a App) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	return a.cities.GetBySlug(ctx, slug)
}
