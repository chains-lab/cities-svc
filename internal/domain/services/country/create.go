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

func (s Service) Create(ctx context.Context, name string) (models.Country, error) {
	now := time.Now().UTC()
	ID := uuid.New()

	_, err := s.GetByName(ctx, name)
	if err == nil {
		return models.Country{}, errx.ErrorCountryAlreadyExistsWithThisName.Raise(err)
	}

	res, err := s.db.CreateCountry(ctx, models.Country{
		ID:        ID,
		Name:      name,
		Status:    enum.CountryStatusUnsupported,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.Country{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating country: %w", err),
		)
	}

	return res, nil
}
