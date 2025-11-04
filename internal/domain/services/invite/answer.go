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

func (s Service) Answer(ctx context.Context, answerID, userID uuid.UUID, answer string) (models.Invite, error) {
	err := enum.CheckInviteStatus(answer)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteAnswer.Raise(
			fmt.Errorf("invalid invite answer: %w", err),
		)
	}

	now := time.Now().UTC()
	inv, err := s.GetInvite(ctx, answerID)
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

	adm, err := s.db.GetCityAdminByUserID(ctx, userID)
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

	switch answer {
	case enum.InviteStatusDeclined:
		if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusDeclined); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}
	case enum.InviteStatusAccepted:
		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			adm = models.CityAdmin{
				UserID:    userID,
				CityID:    inv.CityID,
				Role:      inv.Role,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = s.db.CreateCityAdmin(ctx, adm)
			if err != nil {
				return err
			}

			if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusAccepted); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update invite status, cause: %w", err),
				)
			}
			return nil
		}); err != nil {
			return models.Invite{}, err
		}

		err = s.event.PublishCityAdminCreated(ctx, adm)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin created event, cause: %w", err),
			)
		}
	}

	inv.Status = enum.InviteStatusAccepted
	inv.UserID = userID

	return inv, nil
}
