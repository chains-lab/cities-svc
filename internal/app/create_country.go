package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
)

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	country, err := a.country.Create(ctx, name, constant.CountryStatusUnsupported)
	if err != nil {
		return models.Country{}, err
	}

	return country, nil
}
