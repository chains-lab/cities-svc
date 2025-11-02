package city

import (
	"context"
	"fmt"
	"time"

	"github.com/pariz/gountries"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	CountryID string
	Name      string
	Timezone  string
	StatusID  string
	Point     orb.Point
}

func (s Service) Create(ctx context.Context, params CreateParams) (models.City, error) {
	err := validateTimezone(params.Timezone)
	if err != nil {
		return models.City{}, err
	}

	err = validatePoint(params.Point)
	if err != nil {
		return models.City{}, err
	}

	err = validateName(params.Name)
	if err != nil {
		return models.City{}, err
	}

	_, err = gountries.New().FindCountryByAlpha(params.CountryID)
	if err != nil {
		return models.City{}, errx.ErrorInvalidCountryISO3ID.Raise(
			fmt.Errorf("invalid country ISO3 ID %s: %w", params.CountryID, err),
		)
	}

	status, err := s.StatusAccessible(ctx, params.StatusID)
	if err != nil {
		return models.City{}, err
	}
	if !status.Accessible {
		return models.City{}, errx.ErrorCityStatusNotAccessible.Raise(
			fmt.Errorf("city status %s is not allowed for mods", params.StatusID),
		)
	}

	cityID := uuid.New()
	now := time.Now().UTC()

	res, err := s.db.CreateCity(ctx, models.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Status:    params.StatusID,
		Name:      params.Name,
		Timezone:  params.Timezone,
		Point:     params.Point,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city, cause: %w", err),
		)
	}

	err = s.event.CityCreated(ctx, res)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish city created event, cause: %w", err),
		)
	}

	return res, nil
}
