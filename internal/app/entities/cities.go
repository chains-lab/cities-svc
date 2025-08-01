package entities

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/google/uuid"
)

type citiesQ interface {
	New() dbx.CitiesQ

	Insert(ctx context.Context, input dbx.CityModels) error
	Update(ctx context.Context, input dbx.CityUpdate) error
	Get(ctx context.Context) (dbx.CityModels, error)
	Select(ctx context.Context) ([]dbx.CityModels, error)
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.CitiesQ
	FilterCountryID(countryID uuid.UUID) dbx.CitiesQ
	FilterStatus(status string) dbx.CitiesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesQ
}
