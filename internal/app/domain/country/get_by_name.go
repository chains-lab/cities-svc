package country

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
)

func (c Country) GetByName(ctx context.Context, name string) (models.Country, error) {
	country, err := c.countryQ.New().FilterName(name).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Country{}, errx.ErrorCountryNotFound.Raise(
				fmt.Errorf("country not found, cause: %w", err),
			)
		default:
			return models.Country{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get country by name, cause: %w", err),
			)
		}
	}

	return countryFromDb(country), nil
}
