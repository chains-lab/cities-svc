package city

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
)

func (c City) GetBySlug(ctx context.Context, slug string) (models.City, error) {
	city, err := c.citiesQ.New().FilterSlug(slug).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("city not found by slug: %s, cause: %w", slug, err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by slug, cause: %w", err),
			)
		}
	}

	return cityFromDb(city), nil
}
