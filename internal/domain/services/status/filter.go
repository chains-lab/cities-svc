package status

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/domain/models"
)

type FilterStatusParams struct {
	ID           []string
	Name         *string
	AllowedAdmin *bool
	Accessible   *bool
}

func (s Service) FilterStatuses(ctx context.Context, filters FilterStatusParams, page, size uint64) (models.CityStatusesCollection, error) {
	res, err := s.db.FilterStatuses(ctx, filters, page, size)
	if err != nil {
		return models.CityStatusesCollection{}, err
	}

	return res, nil
}
