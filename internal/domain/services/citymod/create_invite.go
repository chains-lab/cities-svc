package citymod

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) CreateInvite(
	ctx context.Context,
	role string,
	cityID uuid.UUID,
	duration time.Duration,
) (models.Invite, models.InviteToken, error) {
	inviteID := uuid.New()

	now := time.Now().UTC()

	token, err := s.jwt.CreateInviteToken(inviteID, role, cityID, now.Add(duration))
	if err != nil {
		return models.Invite{}, "", errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	err = enum.CheckCityGovRole(role)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInvalidGovRole.Raise(err)
	}

	if err = s.CityIsOfficialSupport(ctx, cityID); err != nil {
		return models.Invite{}, "", err
	}

	m := models.Invite{
		ID:        inviteID,
		Status:    enum.InviteStatusSent,
		Role:      role,
		CityID:    cityID,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
	}

	err = s.db.CreateInvite(ctx, m)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	inv, err := s.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInternal.Raise(
			fmt.Errorf("get invite: %w", err),
		)
	}

	return inv, token, nil
}
