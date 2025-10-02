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
	gov, err := s.db.GetCityAdmin(ctx, filters)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("invalid city admin role, cause: %w", err),
		)
	}

	if gov.IsNil() {
		return models.CityAdmin{}, errx.ErrorCityAdminNotFound.Raise(
			fmt.Errorf("city admin not found"),
		)
	}

	return gov, nil
}

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityAdmin, error) {
	gov, err := s.db.GetCityAdmin(ctx, GetFilters{
		UserID: &initiatorID,
	})
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("invalid city admin role, cause: %w", err),
		)
	}

	if gov.IsNil() {
		return models.CityAdmin{}, errx.ErrorInitiatorIsNotCityAdmin.Raise(
			fmt.Errorf("city admin not found"),
		)
	}

	return gov, nil
}
