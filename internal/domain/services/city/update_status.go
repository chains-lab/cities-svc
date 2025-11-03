package city

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) UpdateStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error) {
	err := enum.CheckCityStatus(status)
	if err != nil {
		return models.City{}, errx.ErrorInvalidCityStatus.Raise(err)
	}

	now := time.Now().UTC()

	ci, err := s.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		switch status {

		case enum.CityStatusOfficial:
			err = s.CountryIsSupported(ctx, ci.CountryID)
			if err != nil {
				return err
			}

		case enum.CityStatusCommunity:
			err = s.CountryIsSupported(ctx, ci.CountryID)
			if err != nil {
				return err
			}

			err = s.db.DeleteGovForCity(ctx, ci.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status, cause: %w", err),
				)
			}

		case enum.CityStatusDeprecated:
			err = s.db.DeleteGovForCity(ctx, ci.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status, cause: %w", err),
				)
			}
		}

		err = s.event.PublishCityUpdated(ctx, ci)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish city updated event, cause: %w", err),
			)
		}

		err = s.db.UpdateCityStatus(ctx, cityID, status, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city status, cause: %w", err),
			)
		}

		ci.Status = status
		ci.UpdatedAt = now

		return nil
	})

	if txErr != nil {
		return models.City{}, txErr
	}

	return ci, nil
}
