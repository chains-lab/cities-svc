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
	Position *string
	Label    *string
}

func (s Service) Update(ctx context.Context, UserID uuid.UUID, params UpdateParams) (models.CityAdminsWithUserData, error) {
	res, err := s.Get(ctx, GetFilters{
		UserID: &UserID,
	})
	if err != nil {
		return models.CityAdminsWithUserData{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		res.Admin.Label = params.Label
	}

	err = s.db.UpdateAdmin(ctx, UserID, params, now)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city initiator, cause: %w", err),
		)
	}

	return res, nil
}

type UpdateOwnParams struct {
	Position *string
	Label    *string
}

func (s Service) UpdateOwn(ctx context.Context, userID uuid.UUID, params UpdateOwnParams) (models.CityAdminsWithUserData, error) {
	res, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return models.CityAdminsWithUserData{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		res.Admin.Label = params.Label
	}

	err = s.db.UpdateAdmin(ctx, userID, UpdateParams{
		Label:    params.Label,
		Position: params.Position,
	}, now)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city admin, cause: %w", err),
		)
	}

	return res, nil
}
