package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/google/uuid"
)

func (a App) TransferGovMayor(ctx context.Context, initiatorID, UserID uuid.UUID) error {
	initiator, err := a.GetInitiatorGov(ctx, initiatorID)
	if err != nil {
		return err
	}

	_, err = a.Get(ctx, UserID)
	if err != nil && !errors.Is(err, errx.ErrorCityGovNotFound) {
		return err
	}
	if err == nil {
		return errx.ErrorGovAlreadyExists.Raise(
			fmt.Errorf("active mayor already exists in city %s", initiator.CityID),
		)
	}

	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.gov.DeleteOne(txCtx, initiator.CityID)
		if err != nil {
			return err
		}

		_, err = a.gov.Create(txCtx, gov.CreateParams{
			CityID: initiator.CityID,
			UserID: UserID,
			Role:   constant.CityGovRoleMayor,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}
