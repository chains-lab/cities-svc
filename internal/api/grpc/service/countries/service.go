package countries

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
	CreateCountry(ctx context.Context, name string) (models.Country, error)
	GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error)

	DeleteCountry(ctx context.Context, ID uuid.UUID) error

	UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status enum.CountryStatus) (models.Country, error)
	UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error)
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
