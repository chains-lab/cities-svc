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

type CreateParams struct {
	UserID   uuid.UUID
	CityID   uuid.UUID
	Role     string
	Duration time.Duration
}

func (s Service) CreateByCityAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	params CreateParams,
) (models.Invite, error) {
	initiator, err := s.getInitiator(ctx, initiatorID, params.CityID)
	if err != nil {
		return models.Invite{}, err
	}
	if initiator.CityID != params.CityID {
		return models.Invite{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator has no rights to create invite for %s", initiatorID),
		)
	}

	if !enum.RightCityAdminsTechPolitics(initiator.Role, params.Role) {
		return models.Invite{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator has no rights to create invite for %s", initiatorID),
		)
	}

	return s.create(ctx, initiatorID, params)
}

func (s Service) CreateBySysAdmin(
	ctx context.Context,
	initiatorID uuid.UUID,
	params CreateParams,
) (models.Invite, error) {
	return s.create(ctx, initiatorID, params)
}

func (s Service) create(
	ctx context.Context,
	initiatorID uuid.UUID,
	params CreateParams,
) (models.Invite, error) {
	inviteID := uuid.New()
	now := time.Now().UTC()

	err := enum.CheckCityAdminRole(params.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidCityAdminRole.Raise(err)
	}

	city, err := s.getCity(ctx, params.CityID)
	if err != nil {
		return models.Invite{}, err
	}

	employeeAlreadyExist, err := s.db.GetCityAdmin(ctx, params.UserID, params.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get emloyee by user id %w", err),
		)
	}
	if !employeeAlreadyExist.IsNil() {
		return models.Invite{}, errx.ErrorCityAdminAlreadyExists.Raise(
			fmt.Errorf("city admin %s already exists in city %s", params.UserID, params.CityID),
		)
	}

	invite := models.Invite{
		ID:          inviteID,
		CityID:      params.CityID,
		UserID:      params.UserID,
		InitiatorID: initiatorID,
		Status:      enum.InviteStatusSent,
		Role:        params.Role,
		ExpiresAt:   now.Add(params.Duration),
		CreatedAt:   now,
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite, cause: %w", err),
		)
	}

	admins, err := s.db.GetCityAdmins(ctx, params.CityID, enum.CityAdminRoleModerator, enum.CityAdminRoleTechLead)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admins for city %s, cause: %w", params.CityID, err),
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
