package country

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
)

type FilterParams struct {
	Name     *string
	Statuses []string
}

func (s Service) Filter(
	ctx context.Context,
	filters FilterParams,
	page, size uint64,
) (models.CountriesCollection, error) {
	res, err := s.db.FilterCountries(ctx, filters, page, size)
	if err != nil {
		return models.CountriesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter countries, cause: %w", err),
		)
	}

	return res, nil
}
