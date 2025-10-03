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
			cou, err := s.getCountryByID(ctx, ci.CountryID)
			if err != nil {
				return err
			}
			if cou.Status != enum.CountryStatusSupported {
				return errx.ErrorCountryIsNotSupported.Raise(
					fmt.Errorf("%s is not supported", ci.CountryID),
				)
			}

		case enum.CityStatusCommunity:
			cou, err := s.getCountryByID(ctx, ci.CountryID)
			if err != nil {
				return err
			}
			if cou.Status != enum.CountryStatusSupported {
				return errx.ErrorCountryIsNotSupported.Raise(
					fmt.Errorf("%s is not supported", ci.CountryID),
				)
			}

			err = s.db.DeleteGovForCity(ctx, ci.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status: %w", err),
				)
			}

		case enum.CityStatusDeprecated:
			err = s.db.DeleteGovForCity(ctx, ci.ID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to delete city status: %w", err),
				)
			}
		}

		err = s.db.UpdateCityStatus(ctx, cityID, status, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update city status: %w", err),
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

func (s Service) getCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	country, err := s.db.GetCountryByID(ctx, ID)
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("error get country by id %s, cause: %w", ID, err),
		)
	}

	if country.IsNil() {
		return models.Country{}, errx.ErrorCountryNotFound.Raise(
			fmt.Errorf("country not found %s", ID),
		)
	}

	return country, nil
}
