package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) Create(ctx context.Context, userID, cityID uuid.UUID, role string) (models.CityAdmin, error) {
	err := enum.CheckCityAdminRole(role)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInvalidCityAdminRole.Raise(
			fmt.Errorf("invalid city admin role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	resp := models.CityAdmin{
		UserID:    userID,
		CityID:    cityID,
		Role:      role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	err = s.db.CreateCityAdmin(ctx, resp)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city admin: %w", err),
		)
	}

	return resp, nil
}
