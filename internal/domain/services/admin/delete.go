package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Delete(
	ctx context.Context,
	userID, cityID uuid.UUID,
) error {
	city, err := s.getSupportedCity(ctx, cityID)
	if err != nil {
		return err
	}

	admin, err := s.Get(ctx, GetFilters{
		UserID: &userID,
		CityID: &cityID,
	})
	if err != nil {
		return err
	}

	if admin.Role == enum.CityAdminRoleTechLead {
		return errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("only system admin can delete tech lead for city %s", cityID),
		)
	}

	return s.delete(ctx, admin, city)
}

func (s Service) DeleteBySysAdmin(
	ctx context.Context,
	userID, cityID uuid.UUID,
) error {
	city, err := s.getSupportedCity(ctx, cityID)
	if err != nil {
		return err
	}

	admin, err := s.Get(ctx, GetFilters{
		UserID: &userID,
		CityID: &cityID,
	})
	if err != nil {
		return err
	}

	return s.delete(ctx, admin, city)
}

func (s Service) delete(
	ctx context.Context,
	admin models.CityAdmin,
	city models.City,
) error {
	err := s.db.DeleteCityAdmin(ctx, admin.UserID, city.ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city admin, cause: %w", err),
		)
	}

	admins, err := s.db.GetCityAdmins(ctx, city.ID, enum.CityAdminRoleModerator, enum.CityAdminRoleTechLead)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", city.ID, err),
		)
	}

	err = s.event.PublishCityAdminDeleted(ctx, admin, city, admins.GetUserIDs())
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin deleted events, cause: %w", err),
		)
	}

	return nil
}
