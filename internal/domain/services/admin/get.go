package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type GetFilters struct {
	UserID *uuid.UUID
	CityID *uuid.UUID
	Role   *string
}

func (s Service) Get(ctx context.Context, filters GetFilters) (models.CityAdmin, error) {
	res, err := s.db.GetCityAdminWithFilter(ctx, filters.UserID, filters.CityID, filters.Role)
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

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityAdmin, error) {
	res, err := s.db.GetCityAdminByUserID(ctx, initiatorID)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin, cause: %w", err),
		)
	}

	if res.IsNil() {
		return models.CityAdmin{}, errx.ErrorInitiatorIsNotCityAdmin.Raise(
			fmt.Errorf("city admin not found"),
		)
	}

	return res, nil
}
