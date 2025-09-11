package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/country"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/pagi"
)

type FilterCountriesListParams struct {
	Name     string
	Statuses []string
}

func (a App) ListCountries(
	ctx context.Context,
	filters FilterCountriesListParams,
	pagination pagi.Request,
	sort []pagi.SortField,
) ([]models.Country, pagi.Response, error) {
	filterForEntities := country.FilterListParams{}
	if filters.Name != "" {
		filterForEntities.Name = filters.Name
	}
	if len(filters.Statuses) > 0 {
		for _, status := range filters.Statuses {
			filterForEntities.Statuses = append(filterForEntities.Statuses, status)
		}
	}
	return a.country.List(ctx, filterForEntities, pagination, sort)
}
