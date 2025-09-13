package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/city"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterListCitiesParams struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Location  *FilterListCityDistance
}

type FilterListCityDistance struct {
	Point   orb.Point
	RadiusM uint64
}

// ListCities searches for cities by name, country ID, and status with pagination and sorting.
// This method for sysadmin
func (a App) ListCities(
	ctx context.Context,
	filters FilterListCitiesParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
	paramsToEntity := city.FilterListParams{}
	if filters.Name != nil {
		paramsToEntity.Name = filters.Name
	}
	if filters.Status != nil {
		paramsToEntity.Status = filters.Status
	}
	if filters.CountryID != nil {
		paramsToEntity.CountryID = filters.CountryID
	}
	if filters.Location != nil {
		paramsToEntity.Location = &city.FilterListDistance{
			Point:   filters.Location.Point,
			RadiusM: filters.Location.RadiusM,
		}
	}
	return a.cities.List(ctx, paramsToEntity, pag, sort)
}
