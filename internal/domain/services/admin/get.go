package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, initiatorID, cityID uuid.UUID) (models.CityAdmin, error) {
	res, err := s.db.GetCityAdmin(ctx, initiatorID, cityID)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin, cause: %w", err),
		)
	}

	if res.IsNil() {
		return models.CityAdmin{}, errx.ErrorCityAdminNotFound.Raise(
			fmt.Errorf("city admin not found"),
		)
	}

	return res, nil
}
