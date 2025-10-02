package country

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Name *string
}

func (s Service) Update(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.Country, error) {
	cou, err := s.GetByID(ctx, ID)
	if err != nil {
		return models.Country{}, err
	}

	if params.Name != nil {
		_, err = s.GetByName(ctx, *params.Name)
		if err != nil && !errors.Is(err, errx.ErrorCountryNotFound) {
			return models.Country{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get country by name, cause: %w", err),
			)
		}

		if err == nil {
			return models.Country{}, errx.ErrorCountryAlreadyExistsWithThisName.Raise(
				fmt.Errorf("country already exists with this name"),
			)
		}
	}

	now := time.Now().UTC()

	err = s.db.UpdateCountry(ctx, ID, params, now)
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update country, cause: %w", err),
		)
	}

	return cou, nil
}
