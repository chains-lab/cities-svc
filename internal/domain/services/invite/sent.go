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

func (s Service) Sent(
	ctx context.Context,
	cityID, userID uuid.UUID,
	role string,
	duration time.Duration,
) (models.Invite, error) {
	inviteID := uuid.New()
	now := time.Now().UTC()

	err := enum.CheckCityAdminRole(role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidCityAdminRole.Raise(err)
	}

	if err = s.CityIsOfficialSupport(ctx, cityID); err != nil {
		return models.Invite{}, err
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

	err = s.event.PublishInviteCreated(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish invite created event, cause: %w", err),
		)
	}

	return invite, nil
}
