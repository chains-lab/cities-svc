package app

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (a App) AcceptInvite(ctx context.Context, initiatorID uuid.UUID, token string) (models.Invite, error) {
	var invite models.Invite
	var err error

	txErr := a.transaction(func(ctx context.Context) error {
		invite, err = a.gov.AcceptInvite(ctx, initiatorID, token)
		if err != nil {
			return err
		}

		city, err := a.GetCityByID(ctx, invite.CityID)
		if err != nil {
			return err
		}
		if city.Status != enum.CityStatusOfficial {
			return errx.ErrorAnswerToInviteForInactiveCity.Raise(
				fmt.Errorf("cannot answer to invite for inactive city %s", city.ID),
			)
		}

		return nil
	})
	if txErr != nil {
		return models.Invite{}, txErr
	}

	return invite, err
}
