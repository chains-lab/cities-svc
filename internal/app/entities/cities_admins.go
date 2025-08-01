package entities

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/google/uuid"
)

type citiesAdminsQ interface {
	New() dbx.CitiesAdminsQ

	Insert(ctx context.Context, input dbx.CityAdminModel) error
	Update(ctx context.Context, input dbx.UpdateCityAdmin) error
	Get(ctx context.Context) (dbx.CityAdminModel, error)
	Select(ctx context.Context) ([]dbx.CityAdminModel, error)
	Delete(ctx context.Context) error

	FilterUserID(UserID uuid.UUID) dbx.CitiesAdminsQ
	FilterCityID(cityID uuid.UUID) dbx.CitiesAdminsQ
	FilterRole(role string) dbx.CitiesAdminsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.CitiesAdminsQ
}
