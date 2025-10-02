package citymod

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterParams struct {
	CityID *uuid.UUID
	Roles  []string
}

func (s Service) Filter(
	ctx context.Context,
	filters FilterParams,
	page, size uint64,
) (models.CityModersCollection, error) {
	res, err := s.db.FilterCityModers(ctx, filters, page, size)
	if err != nil {
		return models.CityModersCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter govs, cause: %w", err),
		)
	}

	return res, err
}
