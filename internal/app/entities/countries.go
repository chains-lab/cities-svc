package entities

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/google/uuid"
)

type countriesQ interface {
	New() dbx.CountriesQ

	Insert(ctx context.Context, input dbx.CountryModel) error
	Update(ctx context.Context, input dbx.UpdateCountryInput) error
	Get(ctx context.Context) (dbx.CountryModel, error)
	Select(ctx context.Context) ([]dbx.CountryModel, error)
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.CountriesQ
	FilterName(name string) dbx.CountriesQ
	FilterStatus(status string) dbx.CountriesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CountriesQ
}
