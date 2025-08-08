package countries

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/countries"
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
	SearchCountries(ctx context.Context, name string, status enum.CountryStatus, limit, offset uint64) ([]models.Country, error)
	DeleteCountry(ctx context.Context, ID uuid.UUID) error

	UpdateCountryStatus(ctx context.Context, ID uuid.UUID, status enum.CountryStatus) (models.Country, error)
	UpdateCountryName(ctx context.Context, ID uuid.UUID, name string) (models.Country, error)
}

type Service struct {
	methods methods
	cfg     config.Config

	svc.CountryServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		methods: app,
		cfg:     cfg,
	}
}

func RequestID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	requestID, ok := ctx.Value(interceptors.RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return requestID
}
