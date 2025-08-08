package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type application interface {
	CreateCity(ctx context.Context, input app.CreateCityInput) (models.City, error)

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, request pagination.Request) ([]models.City, pagination.Response, error)

	UpdateCitiesStatusByOwner(ctx context.Context, initiatorID, cityID uuid.UUID, status string) (models.City, error)
	UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)
	UpdateCityName(ctx context.Context, initiatorID, cityID uuid.UUID, name string) (models.City, error)
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedCityServiceServer
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
