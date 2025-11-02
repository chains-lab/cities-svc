package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Create(ctx context.Context, userID, cityID uuid.UUID, role string) (models.CityAdminsWithUserData, error) {
	err := enum.CheckCityAdminRole(role)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInvalidCityAdminRole.Raise(
			fmt.Errorf("invalid city admin role, cause: %w", err),
		)
	}

	now := time.Now().UTC()

	res := models.CityAdmin{
		UserID:    userID,
		CityID:    cityID,
		Role:      role,
		UpdatedAt: now,
		CreatedAt: now,
	}

	err = s.db.CreateAdmin(ctx, res)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city admin, cause: %w", err),
		)
	}

	err = s.event.CityAdminCreated(ctx, res)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to emit city admin created event, cause: %w", err),
		)
	}

	profiles, err := s.userGuesser.Guess(ctx, res.UserID)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to guess city admin data, cause: %w", err),
		)
	}

	return res.AddProfileData(profiles[res.UserID]), nil
}
