package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) RefuseOwn(ctx context.Context, userID uuid.UUID) error {
	initiator, err := s.GetInitiator(ctx, userID)
	if err != nil {
		return err
	}

	city, err := s.getSupportedCity(ctx, initiator.CityID)
	if err != nil {
		return err
	}

	err = s.db.DeleteCityAdmin(ctx, userID, initiator.CityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete city admin, cause: %w", err),
		)
	}

	if initiator.Role == enum.CityAdminRoleTechLead {
		return errx.ErrorCityAdminTechLeadCannotRefuseOwn.Raise(
			fmt.Errorf("tech lead for city %s cannot refuse own admin role", initiator.CityID),
		)
	}

	recipients, err := s.db.GetCityAdmins(ctx, initiator.CityID, enum.CityAdminRoleModerator, enum.CityAdminRoleTechLead)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", initiator.CityID, err),
		)
	}

	err = s.event.PublishCityAdminDeleted(ctx, initiator, city, recipients.GetUserIDs())
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city admin deleted events, cause: %w", err),
		)
	}

	return nil
}
