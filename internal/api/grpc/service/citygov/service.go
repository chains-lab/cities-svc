package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type application interface {
	CreateCityOwner(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error)
	DeleteCityOwner(ctx context.Context, cityID, userID uuid.UUID) error

	TransferCityOwnership(ctx context.Context, initiatorID, newOwnerID, cityID uuid.UUID) error

	CreateCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID, input app.CreateCityAdminInput) (models.CityAdmin, error)
	GetCityAdmin(ctx context.Context, cityID, userID uuid.UUID) (models.CityAdmin, error)
	GetUserCitiesAdmins(ctx context.Context, userID uuid.UUID, pag pagination.Request) ([]models.CityAdmin, pagination.Response, error)
	UpdateCityAdminRole(ctx context.Context, initiatorID, cityID, userID uuid.UUID, role string) error

	RefuseOwnAdminRights(ctx context.Context, cityID, userID uuid.UUID) error
	DeleteCityAdmin(ctx context.Context, initiatorID, cityID, userID uuid.UUID) error
	GetCityAdmins(ctx context.Context, cityID uuid.UUID, pag pagination.Request) ([]models.CityAdmin, pagination.Response, error)
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedCityGovServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func RequestID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	requestID, ok := ctx.Value(interceptor.RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return requestID
}
