package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.CityAdmin, error) {
	data, err := s.jwt.DecryptInviteToken(token)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid token: %w", err),
		)
	}

	now := time.Now().UTC()

	inv, err := s.GetInvite(ctx, data.InviteID)
	if err != nil {
		return models.CityAdmin{}, err
	}
	if inv.Status != enum.InviteStatusSent {
		return models.CityAdmin{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.CityAdmin{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}

	if err = s.jwt.VerifyInviteToken(token, inv.Token); err != nil {
		return models.CityAdmin{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invite token mismatch"),
		)
	}

	_, err = s.Get(ctx, GetFilters{
		UserID: &userID,
	})
	if err == nil {
		return models.CityAdmin{}, errx.ErrorUserIsAlreadyCityAdmin.Raise(
			fmt.Errorf("user is already a city admin"),
		)
	}
	if !errors.Is(err, errx.ErrorCityAdminNotFound) {
		return models.CityAdmin{}, err
	}

	err = s.CityIsOfficialSupport(ctx, inv.CityID)
	if err != nil {
		return models.CityAdmin{}, err
	}

	var admin models.CityAdmin
	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		admin, err = s.Create(ctx, userID, inv.CityID, data.Role)
		if err != nil {
			return err
		}

		if err := s.db.UpdateStatusInvite(ctx, inv.ID, userID, enum.InviteStatusAccepted, now); err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("update invite status: %w", err))
		}
		return nil
	})
	if txErr != nil {
		return models.CityAdmin{}, txErr
	}

	return admin, nil
}
