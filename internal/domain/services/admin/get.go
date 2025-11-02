package admin

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type GetFilters struct {
	UserID *uuid.UUID
	CityID *uuid.UUID
	Role   *string
}

func (s Service) Get(ctx context.Context, filters GetFilters) (models.CityAdminsWithUserData, error) {
	res, err := s.db.GetAdmin(ctx, filters)
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin, cause: %w", err),
		)
	}

	if res.IsNil() {
		return models.CityAdminsWithUserData{}, errx.ErrorCityAdminNotFound.Raise(
			fmt.Errorf("city admin not found"),
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

func (s Service) GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityAdminsWithUserData, error) {
	res, err := s.db.GetAdmin(ctx, GetFilters{
		UserID: &initiatorID,
	})
	if err != nil {
		return models.CityAdminsWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin, cause: %w", err),
		)
	}

	if res.IsNil() {
		return models.CityAdminsWithUserData{}, errx.ErrorInitiatorIsNotCityAdmin.Raise(
			fmt.Errorf("city admin not found"),
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
