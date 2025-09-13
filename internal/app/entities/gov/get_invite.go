package gov

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (g Gov) GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error) {
	inv, err := g.inv.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.Invite{}, errx.ErrorInternal.Raise(fmt.Errorf("get invite by ID, cause %w", err))
		}
	}

	token, err := g.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		ID:        inv.ID,
		CityID:    inv.CityID,
		Role:      inv.Role,
		ExpiredAt: inv.ExpiresAt,
		CreatedAt: inv.CreatedAt,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	return inviteFromDB(inv, token), nil
}
