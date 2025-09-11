package invites

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

type CreateInviteParams struct {
	InitiatorID uuid.UUID
	CityID      uuid.UUID
	Role        string
	TimeLife    time.Duration
}

func (i Invite) Create(ctx context.Context, params CreateInviteParams) (models.Invite, string, error) {
	exAt := time.Now().UTC().Add(params.TimeLife)

	err := constant.CheckCityGovRole(params.Role)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("check city gov role: %w", err),
		)
	}

	token, id, err := i.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		CityID:    params.CityID,
		Role:      params.Role,
		InvitedBy: params.InitiatorID,
		ExpiredAt: exAt,
	})
	if err != nil {
		return models.Invite{}, "", errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	now := time.Now().UTC()

	stmt := dbx.Invite{
		ID:          id,
		Status:      constant.InviteStatusSent,
		Role:        params.Role,
		CityID:      params.CityID,
		InitiatorID: params.InitiatorID,
		ExpiresAt:   exAt,
		CreatedAt:   now,
	}

	err = i.query.New().Insert(ctx, stmt)
	if err != nil {
		return models.Invite{}, "", errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return modelsFromDB(stmt), token, nil
}
