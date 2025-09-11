package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := a.country.GetByID(ctx, ID)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}
