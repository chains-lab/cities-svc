package citiesadmins

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
	CreateCityOwner(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error)
	TransferCityOwnership(ctx context.Context, initiatorID, newOwnerID, cityID uuid.UUID) error

	CreateCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID, input app.CreateCityAdminInput) (models.CityAdmin, error)
	GetCityAdmin(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error)
	GetCityAdminForCity(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error)
	GetUserCitiesAdmins(ctx context.Context, userID uuid.UUID, limit, page uint64) ([]models.CityAdmin, error)
	UpdateCityAdminRole(ctx context.Context, initiatorID, cityID, userID uuid.UUID, role enum.CityAdminRole) error

	RefuseOwnAdminRights(ctx context.Context, cityID, userID uuid.UUID) error
	DeleteCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID) error
	GetCityAdmins(ctx context.Context, cityID uuid.UUID, limit, page uint64) ([]models.CityAdmin, error)
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
