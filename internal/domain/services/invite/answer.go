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

	inv, err := s.Get(ctx, answerID)
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

	admin, err := s.db.GetCityAdminByUserID(ctx, userID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin by user_id %s and city_id %s, cause: %w", userID, inv.CityID, err),
		)
	}
	if !admin.IsNil() {
		return models.Invite{}, errx.ErrorCityAdminAlreadyExists.Raise(
			fmt.Errorf("city admin with user_id %s already exists in city_id %s", userID, inv.CityID),
		)
	}

	city, err := s.getSupportedCity(ctx, inv.CityID)
	if err != nil {
		return models.Invite{}, err
	}

	switch answer {
	case enum.InviteStatusAccepted:
		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			if inv.Role == enum.CityAdminRoleTechLead {
				existingTechLead, err := s.db.GetCityAdminWithFilter(
					ctx, nil, &inv.CityID, &inv.Role,
				)
				if err != nil {
					return errx.ErrorInternal.Raise(
						fmt.Errorf("failed to get existing tech lead for city %s, cause: %w", inv.CityID, err),
					)
				}

				// Theoretically, it can be removed, but to avoid bugs, it is better not to touch it.
				if !existingTechLead.IsNil() {
					err = s.db.DeleteCityAdmin(ctx, existingTechLead.UserID, inv.CityID)
					if err != nil {
						return errx.ErrorInternal.Raise(
							fmt.Errorf("failed to delete existing tech lead for city %s, cause: %w", inv.CityID, err),
						)
					}
				}
			}

			admin = models.CityAdmin{
				UserID:    userID,
				CityID:    inv.CityID,
				Role:      inv.Role,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = s.db.CreateCityAdmin(ctx, admin)
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

		admins, err := s.db.GetCityAdmins(ctx, inv.CityID, enum.CityAdminRoleModerator)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city admins for city %s, cause: %w", inv.CityID, err),
			)
		}
		err = s.event.PublishCityAdminCreated(ctx, admin, city, append(admins.GetUserIDs(), userID)...)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin created events, cause: %w", err),
			)
		}

		err = s.event.PublishInviteAccepted(ctx, inv, city, admin, inv.InitiatorID)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite accepted events, cause: %w", err),
			)
		}

	case enum.InviteStatusDeclined:
		if err = s.db.UpdateInviteStatus(ctx, inv.ID, userID, enum.InviteStatusDeclined); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}

		if err = s.event.PublishInviteDeclined(ctx, inv, city, inv.InitiatorID); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite declined events, cause: %w", err),
			)
		}
	}

	inv.Status = answer
	inv.UserID = userID

	return inv, nil
}
