package country

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (c Country) Create(ctx context.Context, name string, status string) (models.Country, error) {
	now := time.Now().UTC()
	ID := uuid.New()

	_, err := c.GetByName(ctx, name)
	if err == nil {
		return models.Country{}, errx.ErrorCountryAlreadyExistsWithThisName.Raise(err)
	}

	err = enum.CheckCountryStatus(status)
	if err != nil {
		return models.Country{}, errx.ErrorInvalidCountryStatus.Raise(
			fmt.Errorf("failed to parse country status: %w", err),
		)
	}

	err = c.countryQ.New().Insert(ctx, dbx.Country{
		ID:        ID,
		Name:      name,
		Status:    enum.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating country: %w", err),
		)
	}

	return models.Country{
		ID:        ID,
		Name:      name,
		Status:    enum.CountryStatusSupported,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
