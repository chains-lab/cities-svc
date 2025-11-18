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

func (s Service) Reply(
	ctx context.Context,
	userID, inviteID uuid.UUID,
	reply string,
) (models.Invite, error) {
	err := enum.CheckInviteStatus(reply)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidInviteReply.Raise(
			fmt.Errorf("invalid invite reply: %w", err),
		)
	}

	now := time.Now().UTC()

	invite, err := s.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyReplied.Raise(
			fmt.Errorf("invite already answered with status=%s", invite.Status),
		)
	}
	if now.After(invite.ExpiresAt) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired"),
		)
	}

	admin, err := s.db.GetCityAdmin(ctx, userID, invite.CityID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin by user_id %s and city_id %s, cause: %w", userID, invite.CityID, err),
		)
	}
	if !admin.IsNil() {
		return models.Invite{}, errx.ErrorCityAdminAlreadyExists.Raise(
			fmt.Errorf("city admin with user_id %s already exists in city_id %s", userID, invite.CityID),
		)
	}

	city, err := s.getCity(ctx, invite.CityID)
	if err != nil {
		return models.Invite{}, err
	}
	if city.Status != enum.CityStatusSupported {
		return models.Invite{}, errx.ErrorCityIsNotSupported.Raise(
			fmt.Errorf("city with id %s is not supported", invite.CityID),
		)
	}

	switch reply {
	case enum.InviteStatusAccepted:
		if err = s.db.Transaction(ctx, func(ctx context.Context) error {
			if invite.Role == enum.CityAdminRoleTechLead {
				existingTechLead, err := s.db.GetCityTechLead(ctx, invite.CityID)
				if err != nil {
					return errx.ErrorInternal.Raise(
						fmt.Errorf("failed to get existing tech lead for city %s, cause: %w", invite.CityID, err),
					)
				}

				// Theoretically, this part can be removed, but to avoid bugs, it is better not to touch it.
				// For default city tech-lead must be in active city, but if something went wrong before, we need to handle this case.
				if !existingTechLead.IsNil() {
					err = s.db.DeleteCityAdmin(ctx, existingTechLead.UserID, invite.CityID)
					if err != nil {
						return errx.ErrorInternal.Raise(
							fmt.Errorf("failed to delete existing tech lead for city %s, cause: %w", invite.CityID, err),
						)
					}
				}
			}

			admin = models.CityAdmin{
				UserID:    userID,
				CityID:    invite.CityID,
				Role:      invite.Role,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = s.db.CreateCityAdmin(ctx, admin)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to create city admin, cause: %w", err),
				)
			}

			if err = s.db.UpdateInviteStatus(ctx, invite.ID, enum.InviteStatusAccepted); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update invite status, cause: %w", err),
				)
			}
			return nil
		}); err != nil {
			return models.Invite{}, err
		}

		admins, err := s.db.GetCityAdmins(ctx, invite.CityID, enum.CityAdminRoleModerator)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city admins for city %s, cause: %w", invite.CityID, err),
			)
		}
		err = s.event.PublishCityAdminCreated(ctx, admin, city, append(admins.GetUserIDs(), userID)...)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city admin created events, cause: %w", err),
			)
		}

		err = s.event.PublishInviteAccepted(ctx, invite, city, admin, invite.InitiatorID)
		if err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite accepted events, cause: %w", err),
			)
		}

	case enum.InviteStatusDeclined:
		if err = s.db.UpdateInviteStatus(ctx, invite.ID, enum.InviteStatusDeclined); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update invite status, cause: %w", err),
			)
		}

		if err = s.event.PublishInviteDeclined(ctx, invite, city, invite.InitiatorID); err != nil {
			return models.Invite{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish invite declined events, cause: %w", err),
			)
		}
	}

	invite.Status = reply
	invite.UserID = userID

	return invite, nil
}
