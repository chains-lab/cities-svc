package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Label    *string
	Position *string
}

func (s Service) Update(
	ctx context.Context,
	userID uuid.UUID,
	initiatorID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	city, err := s.getSupportedCity(ctx, initiator.CityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	admin, err := s.Get(ctx, GetFilters{
		UserID: &userID,
		CityID: &initiator.CityID,
	})
	if err != nil {
		return models.CityAdmin{}, err
	}

	access := enum.RightTechPolitics(initiator.Role, admin.Role)
	if !access {
		return models.CityAdmin{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator %s has no rights to update admin %s", initiatorID, userID),
		)
	}

	if initiator.CityID != city.ID {
		return models.CityAdmin{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator %s is not admin for city %s", initiatorID, city.ID),
		)
	}

	now := time.Now().UTC()
	if params.Label != nil {
		admin.Label = params.Label
	}
	if params.Position != nil {
		admin.Position = params.Position
	}

	err = s.db.UpdateCityAdmin(ctx, userID, params, now)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city initiator, cause: %w", err),
		)
	}

	err = s.event.PublishCityAdminUpdated(ctx, admin, city, []uuid.UUID{userID})
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
		)
	}

	return admin, nil
}

func (s Service) UpdateOwn(
	ctx context.Context,
	userID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	initiator, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	city, err := s.getSupportedCity(ctx, initiator.CityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		initiator.Label = params.Label
	}
	if params.Position != nil {
		initiator.Position = params.Position
	}

	err = s.db.UpdateCityAdmin(ctx, userID, params, now)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city admin, cause: %w", err),
		)
	}

	err = s.event.PublishCityAdminUpdated(ctx, initiator, city, []uuid.UUID{userID})
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
		)
	}

	return initiator, nil
}

func (s Service) UpdateBySysAdmin(
	ctx context.Context,
	userID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	user, err := s.Get(ctx, GetFilters{
		UserID: &userID,
	})
	if err != nil {
		return models.CityAdmin{}, err
	}

	city, err := s.getSupportedCity(ctx, user.CityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	now := time.Now().UTC()
	if params.Label != nil {
		user.Label = params.Label
	}
	if params.Position != nil {
		user.Position = params.Position
	}

	err = s.db.UpdateCityAdmin(ctx, userID, params, now)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city admin, cause: %w", err),
		)
	}

	err = s.event.PublishCityAdminUpdated(ctx, user, city, []uuid.UUID{userID})
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
		)
	}

	return user, nil
}
