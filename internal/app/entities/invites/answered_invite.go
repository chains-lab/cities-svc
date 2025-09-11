package invites

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

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

	inv, err := i.Get(ctx, jti)
	if err != nil {
		return models.Invite{}, err
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
	inv.UserID = &userID
	inv.AnsweredAt = &now
	return inv, nil
}
