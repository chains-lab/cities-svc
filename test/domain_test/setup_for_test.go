package domain_test

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/chains-lab/cities-svc/internal/data"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"
	"github.com/chains-lab/cities-svc/internal/infra/jwtmanager"
	"github.com/chains-lab/cities-svc/test"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

// TEST DATABASE CONNECTION

type CityModSvc interface {
	Filter(
		ctx context.Context,
		filters citymod.FilterParams,
		page, size uint64,
	) (models.CityModersCollection, error)

	Get(ctx context.Context, filters citymod.GetFilters) (models.CityModer, error)
	GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityModer, error)

	RefuseOwn(ctx context.Context, userID uuid.UUID) error

	Delete(ctx context.Context, UserID, CityID uuid.UUID) error

	CreateInvite(
		ctx context.Context,
		role string,
		cityID uuid.UUID,
		duration time.Duration,
	) (models.Invite, models.InviteToken, error)

	AcceptInvite(ctx context.Context, userID uuid.UUID, token string) (models.CityModer, error)

	UpdateOther(ctx context.Context, UserID uuid.UUID, params citymod.UpdateCityModerParams) (models.CityModer, error)
	UpdateOwn(ctx context.Context, userID uuid.UUID, params citymod.UpdateCityModerParams) (models.CityModer, error)

	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
}

type CitySvc interface {
	Create(ctx context.Context, params city.CreateParams) (models.City, error)

	Filter(
		ctx context.Context,
		filters city.FilterParams,
		page, size uint64,
	) (models.CitiesCollection, error)

	GetByID(ctx context.Context, cityID uuid.UUID) (models.City, error)
	GetByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error)
	GetBySlug(ctx context.Context, slug string) (models.City, error)

	UpdateStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)

	Update(ctx context.Context, cityID uuid.UUID, params city.UpdateParams) (models.City, error)
}

type CountrySvc interface {
	Create(ctx context.Context, name string) (models.Country, error)

	GetByID(ctx context.Context, ID uuid.UUID) (models.Country, error)
	GetByName(ctx context.Context, name string) (models.Country, error)

	Filter(
		ctx context.Context,
		filters country.FilterParams,
		page, size uint64,
	) (models.CountriesCollection, error)

	UpdateStatus(ctx context.Context, countryID uuid.UUID, status string) (models.Country, error)

	Update(ctx context.Context, ID uuid.UUID, params country.UpdateParams) (models.Country, error)
}

type domain struct {
	moder   CityModSvc
	city    CitySvc
	country CountrySvc
}

type Setup struct {
	domain domain
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := internal.Config{
		JWT: internal.JWTConfig{
			Invites: struct {
				SecretKey string `mapstructure:"secret_key"`
			}{
				SecretKey: "invitesuperkey", // тут подставь ключ для тестов
			},
		},
		Database: internal.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: test.TestDatabaseURL,
			},
		},
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(pg)

	jwtInviteManager := jwtmanager.NewManager(cfg)

	citySvc := city.NewService(database)
	countrySvc := country.NewService(database)
	cityModerSvc := citymod.NewService(database, jwtInviteManager)

	return Setup{
		domain: domain{
			country: countrySvc,
			city:    citySvc,
			moder:   cityModerSvc,
		},
	}, nil
}
