package gov

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type SentInviteParams struct {
	InitiatorID uuid.UUID
	CityID      uuid.UUID
	Role        string
}

func (g Gov) SentInvite(ctx context.Context, params SentInviteParams) (models.Invite, error) {
	initiator, err := g.GetInitiatorGov(ctx, params.InitiatorID)
	if err != nil {
		return models.Invite{}, err
	}

	if initiator.CityID != params.CityID {
		return models.Invite{}, errx.ErrorInitiatorIsNotThisCityGov.Raise(
			fmt.Errorf("initiator have not access to city %s", params.CityID),
		)
	}

	access, err := enum.CompareCityGovRoles(params.Role, initiator.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("compare city gov roles: %w", err),
		)
	}
	if access <= 0 && initiator.Role != enum.CityGovRoleMayor {
		return models.Invite{}, errx.ErrorInitiatorGovRoleHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	invID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := g.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		ID:        invID,
		CityID:    initiator.CityID,
		Role:      params.Role,
		ExpiredAt: exAt,
		CreatedAt: now,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	stmt := dbx.Invite{
		ID:        invID,
		Status:    enum.InviteStatusSent,
		Role:      params.Role,
		CityID:    initiator.CityID,
		ExpiresAt: exAt,
		CreatedAt: now,
	}

	err = g.inv.New().Insert(ctx, stmt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return inviteFromDB(stmt, token), nil
}

func (g Gov) CreateMayorInvite(ctx context.Context, cityID uuid.UUID) (models.Invite, error) {
	invID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	token, err := g.jwt.CreateInviteToken(jwtmanager.InvitePayload{
		ID:        invID,
		CityID:    cityID,
		Role:      enum.CityGovRoleMayor,
		ExpiredAt: exAt,
		CreatedAt: now,
	})
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite token: %w", err),
		)
	}

	stmt := dbx.Invite{
		ID:        invID,
		Status:    enum.InviteStatusSent,
		Role:      enum.CityGovRoleMayor,
		CityID:    cityID,
		ExpiresAt: exAt,
		CreatedAt: now,
	}

	err = g.inv.New().Insert(ctx, stmt)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("create invite: %w", err),
		)
	}

	return inviteFromDB(stmt, token), nil
}
