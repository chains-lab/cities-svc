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

func (s Service) AnswerToInvite(ctx context.Context, inviteID uuid.UUID, answer string) (models.Invite, error) {
	if answer != enum.InviteStatusAccepted && answer != enum.InviteStatusDeclined {
		return models.Invite{}, errx.ErrorInvalidInviteAnswer.Raise(
			fmt.Errorf("invalid invite answer: %s", answer),
		)
	}

	invite, err := s.db.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get invite: %w", err),
		)
	}
	if invite.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite not found"),
		)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite already answered"),
		)
	}
	now := time.Now().UTC()
	if !invite.ExpiresAt.After(now) { // expires_at <= now
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s (now %s)", invite.ExpiresAt.Format(time.RFC3339), now.Format(time.RFC3339)),
		)
	}

	exist, err := s.db.ExistsAdmin(ctx, invite.UserID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check if city admin exists: %w", err),
		)
	}
	if exist {
		return models.Invite{}, errx.ErrorUserAlreadyCityAdmin.Raise(
			fmt.Errorf("user is already a city admin"),
		)
	}

	city, err := s.db.GetCity(ctx, invite.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city: %w", err),
		)
	}
	if city.IsNil() {
		return models.Invite{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city not found"),
		)
	}

	cityStatus, err := s.db.GetCityStatus(ctx, invite.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city status: %w", err),
		)
	}

	if !cityStatus.AllowedAdmin {
		return models.Invite{}, errx.ErrorCityAdminNotAllowed.Raise(
			fmt.Errorf("city adminernment not allowed"),
		)
	}

	err = s.db.Transaction(ctx, func(txCtx context.Context) error {
		err = s.db.AnswerToInvite(txCtx, inviteID, answer)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to answer to invite: %w", err),
			)
		}

		if answer == enum.InviteStatusAccepted {
			err = s.db.CreateAdmin(txCtx, models.CityAdmin{
				CityID: invite.CityID,
				UserID: invite.UserID,
				Role:   invite.Role,
			})
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to create city admin: %w", err),
				)
			}
		}

		return nil
	})
	if err != nil {
		return models.Invite{}, err
	}

	invite.Status = answer

	return invite, nil
}
