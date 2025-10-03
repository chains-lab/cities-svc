package city

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterParams struct {
	Name      *string
	Status    *string
	CountryID *uuid.UUID

	Location *FilterDistance
}

type FilterDistance struct {
	Point   orb.Point
	RadiusM uint64
}

func (s Service) Filter(
	ctx context.Context,
	filters FilterParams,
	page, size uint64,
) (models.CitiesCollection, error) {
	res, err := s.db.FilterCities(ctx, filters, page, size)
	if err != nil {
		return models.CitiesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter cities, cause: %w", err),
		)
	}

	return res, nil
}
