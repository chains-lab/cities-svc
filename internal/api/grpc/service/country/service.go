package country

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/config"

	"github.com/google/uuid"
)

type application interface {
	CreateCountry(ctx context.Context, name string) (models.Country, error)
	GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error)
	SearchCountries(ctx context.Context, name string, status string, pag pagination.Request) ([]models.Country, pagination.Response, error)

	UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status string) (models.Country, error)
	UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error)
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedCountryServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
