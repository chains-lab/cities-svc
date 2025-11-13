package controller

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"

	"github.com/chains-lab/logium"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CityModSvc interface {
	Filter(
		ctx context.Context,
		filters admin.FilterParams,
		page, size uint64,
	) (models.CityAdminCollection, error)

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

type inviteSvc interface {
	Create(
		ctx context.Context,
		cityID, userID uuid.UUID,
		role string,
		duration time.Duration,
	) (models.Invite, error)

	Answer(
		ctx context.Context,
		answerID, userID uuid.UUID,
		answer string,
	) (models.Invite, error)
}

type domain struct {
	moder  CityModSvc
	city   CitySvc
	invite inviteSvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, city CitySvc, cityMod CityModSvc, invSvc inviteSvc) Service {
	return Service{
		log: log,
		domain: domain{
			city:   city,
			moder:  cityMod,
			invite: invSvc,
		},
	}
}
