package country

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) UpdateStatus(ctx context.Context, countryID uuid.UUID, status string) (models.Country, error) {
	err := enum.CheckCountryStatus(status)
	if err != nil {
		return models.Country{}, errx.ErrorInvalidCountryStatus.Raise(err)
	}

	now := time.Now().UTC()

	cou, err := s.GetByID(ctx, countryID)
	if err != nil {
		return models.Country{}, err
	}

	txErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		switch status {

		case enum.CountryStatusSupported:
			//skip all good wee need only update status
		case enum.CountryStatusUnsupported,
			enum.CountryStatusDeprecated:
			err = s.db.UpdateStatusForAllCountryCities(ctx, countryID, status, now)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to update city status, cause: %w", err),
				)
			}
		}

		err = s.db.UpdateCountryStatus(ctx, countryID, status, now)

		cou.Status = status
		cou.UpdatedAt = now

		return nil
	})

	if txErr != nil {
		return models.Country{}, txErr
	}

	return cou, nil
}
