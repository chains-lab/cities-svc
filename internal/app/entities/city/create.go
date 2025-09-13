package city

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateCityParams struct {
	CountryID uuid.UUID
	Name      string
	Timezone  string
	Point     orb.Point
}

func (c City) Create(ctx context.Context, params CreateCityParams) (models.City, error) {
	err := c.validateTimezone(params.Timezone)
	if err != nil {
		return models.City{}, err
	}

	err = c.validatePoint(params.Point)
	if err != nil {
		return models.City{}, err
	}

	err = c.validateName(params.Name)
	if err != nil {
		return models.City{}, err
	}

	cityID := uuid.New()
	now := time.Now().UTC()

	resp := models.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Point:     params.Point,
		Status:    enum.CityStatusCommunity,
		Name:      params.Name,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	stmt := dbx.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Point:     params.Point,
		Status:    enum.CityStatusCommunity,
		Name:      params.Name,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = c.citiesQ.Insert(ctx, stmt)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city: %w", err),
		)
	}

	return resp, nil
}
