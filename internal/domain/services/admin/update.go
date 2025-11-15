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

	admin, err := s.Get(ctx, GetFilters{UserID: &userID})
	if err != nil {
		return models.CityAdmin{}, err
	}

	if !enum.RightCityAdminsTechPolitics(initiator.Role, admin.Role) {
		return models.CityAdmin{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator %s has no rights to update admin %s", initiatorID, userID),
		)
	}

	if params.Role != nil && !enum.RightCityAdminsTechPolitics(initiator.Role, *params.Role) {
		return models.CityAdmin{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator %s has no rights to set role %s for admin %s", initiatorID, *params.Role, userID),
		)
	}

	if params.Role != nil && *params.Role == enum.CityAdminRoleTechLead && initiator.UserID != admin.UserID {
		if initiator.Role != enum.CityAdminRoleTechLead {
			return models.CityAdmin{}, errx.ErrorInitiatorHasNoRights.Raise(
				fmt.Errorf("only tech lead can promote another admin to tech lead"),
			)
		}

		city, err := s.getSupportedCity(ctx, admin.CityID)
		if err != nil {
			return models.CityAdmin{}, err
		}

		now := time.Now().UTC()
		moderRole := enum.CityAdminRoleModerator

		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			if err := s.db.UpdateCityAdmin(ctx, initiator.UserID, UpdateParams{
				Role: &moderRole,
			}, now); err != nil {
				return errx.ErrorInternal.Raise(fmt.Errorf("failed to update initiator admin, cause: %w", err))
			}

			if err := s.db.UpdateCityAdmin(ctx, admin.UserID, params, now); err != nil {
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

	return s.applyAdminUpdate(ctx, admin, params)
}

func (s Service) UpdateBySysAdmin(
	ctx context.Context,
	userID uuid.UUID,
	params UpdateParams,
) (models.CityAdmin, error) {
	admin, err := s.Get(ctx, GetFilters{UserID: &userID})
	if err != nil {
		return models.CityAdmin{}, err
	}

	if params.Role != nil && *params.Role == enum.CityAdminRoleTechLead {
		city, err := s.getSupportedCity(ctx, admin.CityID)
		if err != nil {
			return models.CityAdmin{}, err
		}

		now := time.Now().UTC()
		var currentLead models.CityAdmin

		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			moderRole := enum.CityAdminRoleModerator
			leadRole := enum.CityAdminRoleTechLead

			currentLead, err = s.Get(ctx, GetFilters{
				CityID: &admin.CityID,
				Role:   &leadRole,
			})
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to get current tech lead, cause: %w", err),
				)
			}

			if err := s.db.UpdateCityAdmin(ctx, currentLead.UserID, UpdateParams{
				Role: &moderRole,
			}, now); err != nil {
				return errx.ErrorInternal.Raise(fmt.Errorf("failed to update current tech lead, cause: %w", err))
			}

			if err := s.db.UpdateCityAdmin(ctx, admin.UserID, params, now); err != nil {
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

		// шлём ивенты
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

	return s.applyAdminUpdate(ctx, admin, params)
}

type UpdateOwnParams struct {
	Label    *string
	Position *string
}

func (s Service) UpdateOwn(
	ctx context.Context,
	userID uuid.UUID,
	params UpdateOwnParams,
) (models.CityAdmin, error) {
	admin, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	return s.applyAdminUpdate(ctx, admin, UpdateParams{
		Label:    params.Label,
		Position: params.Position,
	})
}

func (s Service) applyAdminUpdate(
	ctx context.Context,
	admin models.CityAdmin,
	params UpdateParams,
) (models.CityAdmin, error) {
	city, err := s.getSupportedCity(ctx, admin.CityID)
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

	if err = s.db.UpdateCityAdmin(ctx, admin.UserID, params, now); err != nil {
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
