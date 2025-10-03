package city

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	CountryID uuid.UUID
	Name      string
	Timezone  string
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

	cityID := uuid.New()
	now := time.Now().UTC()

	if err = s.CountryIsSupported(ctx, params.CountryID); err != nil {
		return models.City{}, err
	}

	res, err := s.db.CreateCity(ctx, models.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Point:     params.Point,
		Status:    enum.CityStatusCommunity,
		Name:      params.Name,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city, cause: %w", err),
		)
	}

	return res, nil
}
