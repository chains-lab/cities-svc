package country

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := s.db.GetCountryByID(ctx, ID)
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("error get country by id %s, cause: %w", ID, err),
		)
	}

	if country.IsNil() {
		return models.Country{}, errx.ErrorCountryNotFound.Raise(
			fmt.Errorf("country not found %s", ID),
		)
	}

	return country, nil
}

func (s Service) GetByName(ctx context.Context, name string) (models.Country, error) {
	country, err := s.db.GetCountryByName(ctx, name)
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get country by name %s, cause: %w", name, err),
		)
	}

	if country.IsNil() {
		return models.Country{}, errx.ErrorCountryNotFound.Raise(
			fmt.Errorf("country not found %s", name),
		)
	}

	return country, nil
}
