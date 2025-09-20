package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/domain/country"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (a App) CreateCountry(ctx context.Context, name string) (models.Country, error) {
	return a.country.Create(ctx, name, enum.CountryStatusUnsupported)
}

func (a App) GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error) {
	return a.country.GetByID(ctx, ID)
}

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

type UpdateCountryParams struct {
	Name *string
}

func (a App) UpdateCountry(ctx context.Context, countryID uuid.UUID, params UpdateCountryParams) (models.Country, error) {
	update := country.UpdateCountryParams{}

	if params.Name != nil {
		update.Name = params.Name
	}

	return a.country.Update(ctx, countryID, update)
}
