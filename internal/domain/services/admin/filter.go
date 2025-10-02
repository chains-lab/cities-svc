package admin

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
) (models.CityAdminsCollection, error) {
	res, err := s.db.FilterCityAdmins(ctx, filters, page, size)
	if err != nil {
		return models.CityAdminsCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter city admin, cause: %w", err),
		)
	}

	return res, err
}
