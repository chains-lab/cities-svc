package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) DeleteByCityAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	adminID uuid.UUID,
	cityID uuid.UUID,
) error {
	if initiatorID == adminID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("city admin %s cannot delete himself", initiatorID),
		)
	}
	admin, err := s.Get(ctx, adminID, cityID)
	if err != nil {
		return err
	}

	initiator, err := s.validateInitiator(
		ctx, initiatorID, cityID,
		enum.CityAdminRoleModerator, enum.CityAdminRoleTechLead,
	)
	if err != nil {
		return err
	}

	city, err := s.getCity(ctx, admin.CityID)
	if err != nil {
		return err
	}
	if city.Status != enum.CityStatusSupported {
		return errx.ErrorCityIsNotSupported.Raise(
			fmt.Errorf("city not supported"),
		)
	}

	if admin.Role == enum.CityAdminRoleTechLead {
		return errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("only system admin can delete tech lead for city %s", city.ID),
		)
	}

	if !enum.RightCityAdminsTechPolitics(initiator.Role, admin.Role) {
		return errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("only system admin can delete tech lead for city %s", city.ID),
		)
	}

	return s.delete(ctx, admin, city)
}

func (s Service) DeleteBySysAdmin(
	ctx context.Context,
	userID, cityID uuid.UUID,
) error {
	city, err := s.getCity(ctx, cityID)
	if err != nil {
		return err
	}

	admin, err := s.Get(ctx, userID, cityID)
	if err != nil {
		return err
	}

	return s.delete(ctx, admin, city)
}

func (s Service) DeleteOwn(ctx context.Context, userID, cityID uuid.UUID) error {
	initiator, err := s.Get(ctx, userID, cityID)
	if err != nil {
		return err
	}

	city, err := s.db.GetCityByID(ctx, cityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("get city: %w", err),
		)
	}
	if city.IsNil() {
		return errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city not found"),
		)
	}

	if initiator.Role == enum.CityAdminRoleTechLead {
		return errx.ErrorCityAdminTechLeadCannotRefuseOwn.Raise(
			fmt.Errorf("tech lead for city %s cannot refuse own admin role", initiator.CityID),
		)
	}

	return s.delete(ctx, initiator, city)
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

	err = s.event.PublishCityAdminDeleted(ctx, admin, city, admins.GetUserIDs()...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin deleted events, cause: %w", err),
		)
	}

	return nil
}
