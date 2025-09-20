package city

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/paulmach/orb"
)

func (c City) GetByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error) {
	city, err := c.citiesQ.New().
		FilterWithinRadiusMeters(point, radius).
		OrderByNearest(point, true).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("nearest city not found, cause: %w", err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get nearest city, cause: %w", err),
			)
		}
	}

	return cityFromDb(city), nil
}
