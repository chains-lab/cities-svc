package app

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/app/entities/gov"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FiltersListGovsParams struct {
	CityID *uuid.UUID
	Roles  []string
}

func (a App) ListGovs(
	ctx context.Context,
	filters FiltersListGovsParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Gov, pagi.Response, error) {
	input := gov.FiltersListParams{}
	if filters.CityID != nil {
		input.CityID = filters.CityID
	}
	if len(filters.Roles) > 0 && filters.Roles != nil {
		input.Role = filters.Roles
	}

	return a.gov.ListGovs(ctx, input, pag, sort)
}
