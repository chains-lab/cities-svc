package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Label *string
}

func (s Service) UpdateOther(ctx context.Context, UserID uuid.UUID, params UpdateParams) (models.CityAdminWithUserData, error) {
	res, err := s.Get(ctx, GetFilters{
		UserID: &UserID,
	})
	if err != nil {
		return models.CityAdminWithUserData{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		res.Label = params.Label
	}

	err = s.db.UpdateCityAdmin(ctx, UserID, params, now)
	if err != nil {
		return models.CityAdminWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city initiator, cause: %w", err),
		)
	}

	return res, nil
}

func (s Service) UpdateOwn(ctx context.Context, userID uuid.UUID, params UpdateParams) (models.CityAdminWithUserData, error) {
	res, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return models.CityAdminWithUserData{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		res.Label = params.Label
	}

	err = s.db.UpdateCityAdmin(ctx, userID, params, now)
	if err != nil {
		return models.CityAdminWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city admin, cause: %w", err),
		)
	}
	
	return res, nil
}
