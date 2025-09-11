package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/enum"
)

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	return a.country.Create(ctx, name, enum.CountryStatusUnsupported)
}
