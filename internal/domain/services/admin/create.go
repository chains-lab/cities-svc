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

// Create deprecated -- use CreateViaInvite instead
func (s Service) Create(ctx context.Context, userID, cityID uuid.UUID, role string) (models.CityAdmin, error) {
	err := enum.CheckCityAdminRole(role)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInvalidCityAdminRole.Raise(
			fmt.Errorf("invalid city admin role, cause: %w", err),
		)
	}

	city, err := s.getCity(ctx, cityID)
	if err != nil {
		return models.CityAdmin{}, err
	}
	if city.Status != enum.CityStatusSupported {
		return models.CityAdmin{}, errx.ErrorCityIsNotSupported.Raise(
			fmt.Errorf("city not supported"),
		)
	}

	now := time.Now().UTC()
	res := models.CityAdmin{
		UserID:    userID,
		CityID:    cityID,
		Role:      role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		if role == enum.CityAdminRoleTechLead {
			existingTechLead, err := s.db.GetCityTechLead(ctx, cityID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to get existing tech lead for city %s, cause: %w", cityID, err),
				)
			}

			// Theoretically, this part can be removed, but to avoid bugs, it is better not to touch it.
			// For default city tech-lead must be in active city, but if something went wrong before, we need to handle this case.
			if !existingTechLead.IsNil() {
				err = s.db.DeleteCityAdmin(ctx, existingTechLead.UserID, cityID)
				if err != nil {
					return errx.ErrorInternal.Raise(
						fmt.Errorf("failed to delete existing tech lead for city %s, cause: %w", cityID, err),
					)
				}
			}
		}

		err = s.db.CreateCityAdmin(ctx, res)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to creating city admin, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.CityAdmin{}, err
	}

	admins, err := s.db.GetCityAdmins(ctx, cityID, enum.CityAdminRoleModerator)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", cityID, err),
		)
	}

	err = s.event.PublishCityAdminCreated(ctx, res, city, append(admins.GetUserIDs(), userID)...)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin created events, cause: %w", err),
		)
	}

	return res, nil
}
