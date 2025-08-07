package cities

import (
	"context"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptors"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

type methods interface {
	CreateCity(ctx context.Context, input app.CreateCityInput) (models.City, error)
	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	DeleteCity(ctx context.Context, ID uuid.UUID) error

	SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, page, limit uint64) ([]models.City, error)

	UpdateCitiesStatusByOwner(ctx context.Context, initiatorID, cityID uuid.UUID, status enum.CityStatus) (models.City, error)
	UpdateCitiesStatusBySysAdmin(ctx context.Context, cityID uuid.UUID, status enum.CityStatus) (models.City, error)
	UpdateStatusForCitiesByCountryID(ctx context.Context, countryID uuid.UUID, status enum.CityStatus) error

	UpdateCityName(ctx context.Context, initiatorID, cityID uuid.UUID, name string) (models.City, error)
}

type Service struct {
	methods methods
	cfg     config.Config
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		methods: app,
		cfg:     cfg,
	}
}

func Meta(ctx context.Context) interceptors.MetaData {
	md, ok := ctx.Value(interceptors.MetaCtxKey).(interceptors.MetaData)
	if !ok {
		return interceptors.MetaData{}
	}
	return md
}
