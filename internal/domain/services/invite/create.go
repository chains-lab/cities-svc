package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) CreateByCityAdmin(
	ctx context.Context,
	userID, cityID, initiatorID uuid.UUID,
	role string,
	duration time.Duration,
) (models.Invite, error) {
	initiator, err := s.getInitiator(ctx, initiatorID)
	if err != nil {
		return models.Invite{}, err
	}

	if initiator.CityID != cityID {
		return models.Invite{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator has no rights to create invite for %s", initiatorID),
		)
	}

	err = enum.CheckCityAdminRole(role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidCityAdminRole.Raise(err)
	}

	employeeAlreadyExist, err := s.db.GetCityAdminByUserID(ctx, userID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get emloyee by user id %w", err),
		)
	}
	if !employeeAlreadyExist.IsNil() {
		return models.Invite{}, errx.ErrorCityAdminAlreadyExists.Raise(
			fmt.Errorf("emloyee already exists %s", userID),
		)
	}

	if !enum.RightCityAdminsTechPolitics(initiator.Role, role) {
		return models.Invite{}, errx.ErrorInitiatorHasNoRights.Raise(
			fmt.Errorf("initiator has no rights to create invite for %s", initiatorID),
		)
	}

	return s.create(ctx, userID, initiatorID, role, duration)
}

func (s Service) CreateBySysAdmin(
	ctx context.Context,
	userID, cityID uuid.UUID,
	role string,
	duration time.Duration,
) (models.Invite, error) {
	return s.create(ctx, userID, cityID, role, duration)
}

func (s Service) create(
	ctx context.Context,
	userID, cityID uuid.UUID,
	role string,
	duration time.Duration,
) (models.Invite, error) {
	inviteID := uuid.New()
	now := time.Now().UTC()

	err := enum.CheckCityAdminRole(role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidCityAdminRole.Raise(err)
	}

	city, err := s.getSupportedCity(ctx, cityID)
	if err != nil {
		return models.Invite{}, err
	}

	employeeAlreadyExist, err := s.db.GetCityAdminByUserID(ctx, userID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get emloyee by user id %w", err),
		)
	}
	if !employeeAlreadyExist.IsNil() {
		return models.Invite{}, errx.ErrorCityAdminAlreadyExists.Raise(
			fmt.Errorf("emloyee already exists %s", userID),
		)
	}

	invite := models.Invite{
		ID:        inviteID,
		Status:    enum.InviteStatusSent,
		Role:      role,
		CityID:    cityID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite, cause: %w", err),
		)
	}

	admins, err := s.db.GetCityAdmins(ctx, cityID, enum.CityAdminRoleModerator)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", cityID, err),
		)
	}

	err = s.event.PublishInviteCreated(ctx, invite, city, admins.GetUserIDs()...)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish invite created events, cause: %w", err),
		)
	}

	return invite, nil
}
