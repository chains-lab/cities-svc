package citymod

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) Create(ctx context.Context, userID, cityID uuid.UUID, role string) (models.CityModer, error) {
	err := enum.CheckCityGovRole(role)
	if err != nil {
		return models.CityModer{}, errx.ErrorInvalidGovRole.Raise(
			fmt.Errorf("invalid city citymod role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	resp := models.CityModer{
		UserID:    userID,
		CityID:    cityID,
		Role:      role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	err = s.db.CreateCityMod(ctx, resp)
	if err != nil {
		return models.CityModer{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city citymod: %w", err),
		)
	}

	return resp, nil
}
