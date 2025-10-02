package citymod

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateCityModerParams struct {
	Label *string
}

func (s Service) UpdateOther(ctx context.Context, UserID uuid.UUID, params UpdateCityModerParams) (models.CityModer, error) {
	gov, err := s.Get(ctx, GetFilters{
		UserID: &UserID,
	})
	if err != nil {
		return models.CityModer{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		gov.Label = params.Label
	}

	err = s.db.UpdateCityModer(ctx, UserID, params, now)
	if err != nil {
		return models.CityModer{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city initiator, cause: %w", err),
		)
	}

	return gov, nil
}

func (s Service) UpdateOwn(ctx context.Context, userID uuid.UUID, params UpdateCityModerParams) (models.CityModer, error) {
	gov, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return models.CityModer{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		gov.Label = params.Label
	}

	err = s.db.UpdateCityModer(ctx, userID, params, now)
	if err != nil {
		return models.CityModer{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city citymod, cause: %w", err),
		)
	}

	return gov, nil
}
