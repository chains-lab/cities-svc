package country

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type UpdateCountryParams struct {
	Name   *string
	Status *string
}

func (c Country) Update(ctx context.Context, ID uuid.UUID, params UpdateCountryParams) (models.Country, error) {
	cou, err := c.GetByID(ctx, ID)
	if err != nil {
		return models.Country{}, err
	}

	stmt := dbx.UpdateCountryParams{}

	if params.Name == nil && params.Status == nil {
		return models.Country{}, nil
	}

	if params.Name != nil {
		_, err = c.GetByName(ctx, *params.Name)
		if err != nil && !errors.Is(err, errx.ErrorCountryNotFound) {
			return models.Country{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get country by name, cause: %w", err),
			)
		}
		if err == nil {
			return models.Country{}, errx.ErrorCountryAlreadyExistsWithThisName.Raise(
				fmt.Errorf("country already exists with this name"),
			)
		}

		cou.Name = *params.Name
		stmt.Name = params.Name
	}
	if params.Status != nil {
		err := enum.CheckCountryStatus(*params.Status)
		if err != nil {
			return models.Country{}, errx.ErrorInvalidCountryStatus.Raise(
				fmt.Errorf("failed to invalid country status, cause: %w", err),
			)
		}

		cou.Status = *params.Status
		stmt.Status = params.Status
	}

	stmt.UpdatedAt = time.Now().UTC()
	cou.UpdatedAt = stmt.UpdatedAt

	err = c.countryQ.New().FilterID(ID).Update(ctx, stmt)
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update country, cause: %w", err),
		)
	}

	return cou, nil
}
