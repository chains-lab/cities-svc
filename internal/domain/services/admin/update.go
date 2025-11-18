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
	Role     *string
}

func (s Service) UpdateByCityAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	cityID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	initiator, err := s.validateInitiator(
		ctx, initiatorID, cityID,
		enum.CityAdminRoleTechLead, enum.CityAdminRoleModerator,
	)
	if err != nil {
		return models.CityAdmin{}, err
	}

	admin, err := s.Get(ctx, userID, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	if !enum.RightCityAdminsTechPolitics(initiator.Role, admin.Role) {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator %s has no rights to update admin %s", initiatorID, userID),
		)
	}

	if params.Role != nil && !enum.RightCityAdminsTechPolitics(initiator.Role, *params.Role) {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator %s has no rights to set role %s for admin %s", initiatorID, *params.Role, userID),
		)
	}

	if params.Role != nil && *params.Role == enum.CityAdminRoleTechLead && initiator.UserID != admin.UserID {
		if initiator.Role != enum.CityAdminRoleTechLead {
			return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
				fmt.Errorf("only tech lead can promote another admin to tech lead"),
			)
		}

		city, err := s.getCity(ctx, admin.CityID)
		if err != nil {
			return models.CityAdmin{}, err
		}
		if city.Status != enum.CityStatusSupported {
			return models.CityAdmin{}, errx.ErrorCityIsNotSupported.Raise(
				fmt.Errorf("city not supported"),
			)
		}

		now := time.Now().UTC()
		moderRole := enum.CityAdminRoleModerator

		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			if err := s.db.UpdateCityAdmin(ctx, initiatorID, cityID, UpdateParams{
				Role: &moderRole,
			}, now); err != nil {
				return errx.ErrorInternal.Raise(fmt.Errorf("failed to update initiator admin, cause: %w", err))
			}

			if err := s.db.UpdateCityAdmin(ctx, userID, cityID, params, now); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update city admin, cause: %w", err),
				)
			}

			return nil
		}); err != nil {
			return models.CityAdmin{}, err
		}

		initiator.Role = moderRole
		initiator.UpdatedAt = now
		admin.Role = *params.Role
		if params.Label != nil {
			admin.Label = params.Label
		}
		if params.Position != nil {
			admin.Position = params.Position
		}

		if err = s.event.PublishCityAdminUpdated(ctx, initiator, city); err != nil {
			return models.CityAdmin{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
			)
		}
		if err = s.event.PublishCityAdminUpdated(ctx, admin, city); err != nil {
			return models.CityAdmin{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
			)
		}

		return admin, nil
	}

	return s.update(ctx, admin, params)
}

func (s Service) UpdateBySysAdmin(
	ctx context.Context,
	userID uuid.UUID,
	cityID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	admin, err := s.Get(ctx, userID, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	if params.Role != nil && *params.Role == enum.CityAdminRoleTechLead {
		city, err := s.getCity(ctx, admin.CityID)
		if err != nil {
			return models.CityAdmin{}, err
		}

		now := time.Now().UTC()
		var currentLead models.CityAdmin

		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			moderRole := enum.CityAdminRoleModerator

			currentLead, err = s.db.GetCityTechLead(ctx, city.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to get current tech lead, cause: %w", err),
				)
			}
			if !currentLead.IsNil() {
				if err := s.db.UpdateCityAdmin(ctx, currentLead.UserID, currentLead.CityID, UpdateParams{
					Role: &moderRole,
				}, now); err != nil {
					return errx.ErrorInternal.Raise(fmt.Errorf("failed to update current tech lead, cause: %w", err))
				}
			}

			if err := s.db.UpdateCityAdmin(ctx, userID, cityID, params, now); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update city admin, cause: %w", err),
				)
			}

			return nil
		}); err != nil {
			return models.CityAdmin{}, err
		}

		now = time.Now().UTC()
		currentLead.Role = enum.CityAdminRoleModerator
		currentLead.UpdatedAt = now
		admin.Role = *params.Role
		if params.Label != nil {
			admin.Label = params.Label
		}
		if params.Position != nil {
			admin.Position = params.Position
		}
		admin.UpdatedAt = now

		if err = s.event.PublishCityAdminUpdated(ctx, currentLead, city); err != nil {
			return models.CityAdmin{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
			)
		}
		if err = s.event.PublishCityAdminUpdated(ctx, admin, city); err != nil {
			return models.CityAdmin{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
			)
		}

		return admin, nil
	}

	return s.update(ctx, admin, params)
}

type UpdateOwnParams struct {
	Label    *string
	Position *string
}

func (s Service) UpdateOwn(
	ctx context.Context,
	userID uuid.UUID,
	cityID uuid.UUID,
	params UpdateOwnParams,
) (models.CityAdmin, error) {
	admin, err := s.validateInitiator(ctx, userID, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	return s.update(ctx, admin, UpdateParams{
		Label:    params.Label,
		Position: params.Position,
	})
}

func (s Service) update(
	ctx context.Context,
	admin models.CityAdmin,
	params UpdateParams,
) (models.CityAdmin, error) {
	city, err := s.getCity(ctx, admin.CityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	now := time.Now().UTC()

	if params.Label != nil {
		admin.Label = params.Label
	}
	if params.Position != nil {
		admin.Position = params.Position
	}
	if params.Role != nil {
		admin.Role = *params.Role
	}

	if err = s.db.UpdateCityAdmin(ctx, admin.UserID, admin.CityID, params, now); err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city admin, cause: %w", err),
		)
	}

	if err = s.event.PublishCityAdminUpdated(ctx, admin, city); err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin updated events, cause: %w", err),
		)
	}

	return admin, nil
}
