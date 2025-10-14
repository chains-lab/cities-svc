package controller

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"
	"github.com/chains-lab/logium"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CityModSvc interface {
	Filter(
		ctx context.Context,
		filters admin.FilterParams,
		page, size uint64,
	) (models.CityAdminsCollection, error)

	Get(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error)
	GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityAdmin, error)

	RefuseOwn(ctx context.Context, userID uuid.UUID) error

	Delete(ctx context.Context, UserID, CityID uuid.UUID) error

	UpdateOther(ctx context.Context, UserID uuid.UUID, params admin.UpdateParams) (models.CityAdmin, error)
	UpdateOwn(ctx context.Context, userID uuid.UUID, params admin.UpdateParams) (models.CityAdmin, error)
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

type inviteSvc interface {
	Create(
		ctx context.Context,
		role string,
		cityID uuid.UUID,
		duration time.Duration,
	) (models.Invite, error)

	Accept(ctx context.Context, userID uuid.UUID, token string) (models.Invite, error)
}

type domain struct {
	moder   CityModSvc
	city    CitySvc
	country CountrySvc
	invite  inviteSvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, country CountrySvc, city CitySvc, cityMod CityModSvc, invSvc inviteSvc) Service {
	return Service{
		log: log,
		domain: domain{
			country: country,
			city:    city,
			moder:   cityMod,
			invite:  invSvc,
		},
	}
}
