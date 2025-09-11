package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/country"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/google/uuid"
)

type UpdateCountryParams struct {
	Name *string
}

func (a App) UpdateCountry(ctx context.Context, countryID uuid.UUID, params UpdateCountryParams) (models.Country, error) {
	update := country.UpdateCountryParams{}

	if params.Name != nil {
		update.Name = params.Name
	}

	return a.country.Update(ctx, countryID, update)
}
