package admin

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
) (models.Invite, error) {
	inviteID := uuid.New()

	now := time.Now().UTC()

	token, err := s.jwt.CreateInviteToken(inviteID, role, cityID, now.Add(duration))
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	err = enum.CheckCityAdminRole(role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidCityAdminRole.Raise(err)
	}

	if err = s.CityIsOfficialSupport(ctx, cityID); err != nil {
		return models.Invite{}, err
	}

	hash, err := s.jwt.HashInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("hash invite token: %w", err),
		)
	}

	invite := models.Invite{
		ID:        inviteID,
		Status:    enum.InviteStatusSent,
		Role:      role,
		CityID:    cityID,
		Token:     hash,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	invite.Token = token

	return invite, nil
}
