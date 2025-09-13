package gov

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (g Gov) AnsweredInvite(ctx context.Context, userID uuid.UUID, token, status string) (models.Invite, error) {
	data, err := g.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid or expired token: %w", err),
		)
	}

	inv, err := g.GetInvite(ctx, data.JTI)
	if err != nil {
		return models.Invite{}, err
	}

	now := time.Now().UTC()

	if inv.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(fmt.Errorf("invite expired"))
	}
	if data.CityID != inv.CityID {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(fmt.Errorf("token city_id mismatch"))
	}

	upd := dbx.UpdateInviteParams{
		Status:     &status,
		UserID:     &uuid.NullUUID{UUID: userID, Valid: true},
		AnsweredAt: &sql.NullTime{Time: now, Valid: true},
	}
	if err := g.inv.New().FilterID(inv.ID).Update(ctx, upd); err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("update invite status: %w", err),
		)
	}

	switch status {
	case enum.InviteStatusAccepted:
		if data.Role == enum.CityGovRoleMayor {
			err = g.deleteCityMayor(ctx, inv.CityID)
			if err != nil {
				return models.Invite{}, err
			}
		}

		_, err = g.createGov(ctx, createParams{
			UserID: userID,
			CityID: inv.CityID,
			Role:   data.Role,
		})
		if err != nil {
			return models.Invite{}, err
		}
	case enum.InviteStatusRejected:
		// nothing to do
	default:
		return models.Invite{}, errx.ErrorUnexpectedInviteStatus.Raise(
			fmt.Errorf("invalid invite status: %s", status),
		)

	}

	inv.Status = status
	inv.UserID = &userID
	inv.AnsweredAt = &now
	return inv, nil
}
