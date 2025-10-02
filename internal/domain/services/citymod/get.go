package citymod

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

func (s Service) Get(ctx context.Context, filters GetFilters) (models.CityModer, error) {
	gov, err := s.db.GetCityModer(ctx, filters)
	if err != nil {
		return models.CityModer{}, errx.ErrorInternal.Raise(
			fmt.Errorf("invalid city citymod role, cause: %w", err),
		)
	}

	if gov.IsNil() {
		return models.CityModer{}, errx.ErrorCityGovNotFound.Raise(
			fmt.Errorf("city citymod not found"),
		)
	}

	return gov, nil
}

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityModer, error) {
	gov, err := s.db.GetCityModer(ctx, GetFilters{
		UserID: &initiatorID,
	})
	if err != nil {
		return models.CityModer{}, errx.ErrorInternal.Raise(
			fmt.Errorf("invalid city citymod role, cause: %w", err),
		)
	}

	if gov.IsNil() {
		return models.CityModer{}, errx.ErrorCityGovNotFound.Raise(
			fmt.Errorf("city citymod not found"),
		)
	}

	return gov, nil
}
