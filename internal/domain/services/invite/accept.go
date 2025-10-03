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

func (s Service) Accept(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error) {
	data, err := s.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid token: %w", err),
		)
	}

	now := time.Now().UTC()

	inv, err := s.GetInvite(ctx, data.InviteID)
	if err != nil {
		return models.Invite{}, err
	}
	if inv.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}

	if err = s.jwt.VerifyInviteToken(token, inv.Token); err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invite token mismatch"),
		)
	}

	adm, err := s.db.GetCityAdminByUserAndCityID(ctx, userID, inv.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin by user_id %s and city_id %s, cause: %w", userID, inv.CityID, err),
		)
	}
	if !adm.IsNil() {
		return models.Invite{}, errx.ErrorUserIsAlreadyCityAdmin.Raise(
			fmt.Errorf("city admin with user_id %s already exists in city_id %s", userID, inv.CityID),
		)
	}

	err = s.CityIsOfficialSupport(ctx, inv.CityID)
	if err != nil {
		return models.Invite{}, err
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		adm = models.CityAdmin{
			UserID:    userID,
			CityID:    inv.CityID,
			Role:      data.Role,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = s.db.CreateCityAdmin(ctx, adm)
		if err != nil {
			return err
		}

		if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusAccepted, now); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}
		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = &userID
	inv.AnsweredAt = &now

	return inv, nil
}
