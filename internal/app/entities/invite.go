package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

type Invite struct {
	jwt   jwtmanager.Manager
	query dbx.InviteQ
}

func NewInvites(pg *sql.DB, cfg config.Config) Invite {
	return Invite{
		jwt:   jwtmanager.NewManager(cfg),
		query: dbx.NewInviteQ(pg),
	}
}

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

func (i Invite) Get(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	inv, err := i.query.New().FilterID(id).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.Invite{}, errx.ErrorInternal.Raise(fmt.Errorf("get invite by id, cause %w", err))
		}
	}

	return modelsFromDB(inv), nil
}

func (i Invite) Answered(ctx context.Context, userID uuid.UUID, token, status string) (models.Invite, error) {
	data, err := i.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid or expired token: %w", err),
		)
	}

	if data.JTI == "" {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(errors.New("token has empty jti"))
	}
	jti, err := uuid.Parse(data.JTI)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(fmt.Errorf("invalid jti format: %w", err))
	}

	inv, err := i.query.New().FilterID(jti).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, errx.ErrorInviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("get invite: %w", err),
			)
		}
	}

	now := time.Now().UTC()

	if inv.Status != constant.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}

	if now.After(inv.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(errors.New("invite expired"))
	}

	if data.CityID != inv.CityID {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(errors.New("token city_id mismatch"))
	}

	err = constant.CheckInviteStatus(data.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("check invite status: %w", err),
		)
	}

	userNull := uuid.NullUUID{UUID: userID, Valid: true}
	ansNull := sql.NullTime{Time: now, Valid: true}

	upd := dbx.UpdateInviteParams{
		Status:     &status,
		UserID:     &userNull,
		AnsweredAt: &ansNull,
	}

	if err := i.query.New().FilterID(inv.ID).Update(ctx, upd); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("update invite status: %w", err),
		)
	}

	inv.Status = status
	inv.UserID = userNull
	inv.AnsweredAt = ansNull

	return modelsFromDB(inv), nil
}

func modelsFromDB(inv dbx.Invite) models.Invite {
	res := models.Invite{
		ID:          inv.ID,
		Status:      inv.Status,
		Role:        inv.Role,
		CityID:      inv.CityID,
		InitiatorID: inv.InitiatorID,
		ExpiresAt:   inv.ExpiresAt,
		CreatedAt:   inv.CreatedAt,
	}
	if inv.UserID.Valid {
		res.UserID = &inv.UserID.UUID
	}
	if inv.AnsweredAt.Valid {
		res.AnsweredAt = &inv.AnsweredAt.Time
	}

	return res
}
