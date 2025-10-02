package citymod

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

func (s Service) AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.CityModer, error) {
	data, err := s.jwt.DecryptInviteToken(models.InviteToken(token))
	if err != nil {
		return models.CityModer{}, errx.ErrorInvalidInviteToken.Raise(
			fmt.Errorf("invalid or expired token: %w", err),
		)
	}

	now := time.Now().UTC()

	inv, err := s.GetInvite(ctx, data.InviteID)
	if err != nil {
		return models.CityModer{}, err
	}
	if inv.Status != enum.InviteStatusSent {
		return models.CityModer{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered with status=%s", inv.Status),
		)
	}
	if now.After(inv.ExpiresAt) {
		return models.CityModer{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}

	_, err = s.Get(ctx, GetFilters{
		UserID: &userID,
	})
	if err == nil {
		return models.CityModer{}, errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("user is already a city gov"),
		)
	}
	if !errors.Is(err, errx.ErrorCityGovNotFound) {
		return models.CityModer{}, err
	}

	var gov models.CityModer
	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		if data.Role == enum.CityGovRoleMayor {
			r := enum.CityGovRoleMayor
			mayor, err := s.Get(ctx, GetFilters{CityID: &inv.CityID, Role: &r})
			if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
				return err
			}
			if !mayor.IsNil() {
				if err := s.Delete(ctx, mayor.UserID, mayor.CityID); err != nil {
					return err
				}
			}
		}

		gov, err = s.Create(ctx, userID, inv.CityID, data.Role)
		if err != nil {
			return err
		}

		if err := s.db.UpdateStatusInvite(ctx, inv.ID, userID, enum.InviteStatusAccepted, now); err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("update invite status: %w", err))
		}
		return nil
	})
	if txErr != nil {
		return models.CityModer{}, txErr
	}

	return gov, nil
}
