package city

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func (s Service) GetByID(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	city, err := s.db.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city by id: %s, cause: %w", cityID, err),
		)
	}

	if city.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("сity not found by id: %s, cause: %w", cityID, err),
		)
	}

	return city, nil
}

func (s Service) GetByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error) {
	city, err := s.db.GetCityByRadius(ctx, point, radius)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get nearest city, cause: %w", err),
		)
	}

	if city.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("nearest city not found, cause: %w", err),
		)
	}

	return city, nil
}

func (s Service) GetBySlug(ctx context.Context, slug string) (models.City, error) {
	city, err := s.db.GetCityBySlug(ctx, slug)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get nearest city, cause: %w", err),
		)
	}

	if city.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("nearest city not found, cause: %w", err),
		)
	}

	return city, nil
}
