package city

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/pagination"
	"github.com/google/uuid"
)

type application interface {
	CreateCity(ctx context.Context, input app.CreateCityInput) (models.City, error)

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, request pagination.Request) ([]models.City, pagination.Response, error)

	UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)
	UpdateCityName(ctx context.Context, cityID uuid.UUID, name string) (models.City, error)
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

func (s Service) OnlyGov(ctx context.Context, initiatorID, cityID, action string) (any, error) {
	return nil, nil
}

func (s Service) OnlyCityAdmin(ctx context.Context, initiatorID, cityID, action string) (any, error) {
	return nil, nil
}
