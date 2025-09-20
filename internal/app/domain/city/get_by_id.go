package city

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (c City) GetByID(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	city, err := c.citiesQ.New().FilterID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("сity not found by id: %s, cause: %w", cityID, err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by id: %s, cause: %w", cityID, err),
			)
		}
	}

	return cityFromDb(city), nil
}
