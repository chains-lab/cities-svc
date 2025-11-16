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

type CityAdminSvc interface {
	Filter(
		ctx context.Context,
		filters admin.FilterParams,
		page, size uint64,
	) (models.CityAdminsCollection, error)

	Get(ctx context.Context, filters admin.GetFilters) (models.CityAdmin, error)
	GetInitiator(ctx context.Context, initiatorID uuid.UUID) (models.CityAdmin, error)

	DeleteOwn(ctx context.Context, userID uuid.UUID) error

	DeleteByCityAdmin(ctx context.Context, userID, cityID, initiatorID uuid.UUID) error
	DeleteBySysAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	UpdateByCityAdmin(
		ctx context.Context,
		userID uuid.UUID,
		initiatorID uuid.UUID,
		params admin.UpdateParams,
	) (models.CityAdmin, error)

	UpdateBySysAdmin(
		ctx context.Context,
		userID uuid.UUID,
		params admin.UpdateParams,
	) (models.CityAdmin, error)

	UpdateOwn(
		ctx context.Context,
		userID uuid.UUID,
		params admin.UpdateOwnParams,
	) (models.CityAdmin, error)
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

	UpdateStatusByCityAdmin(ctx context.Context, cityID, initiatorID uuid.UUID, status string) (models.City, error)
	UpdateStatusBySysAdmin(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)

	UpdateByCityAdmin(ctx context.Context, cityID, initiatorID uuid.UUID, params city.UpdateParams) (models.City, error)
	UpdateByAdmin(ctx context.Context, cityID uuid.UUID, params city.UpdateParams) (models.City, error)
}

type inviteSvc interface {
	CreateBySysAdmin(
		ctx context.Context,
		userID, cityID uuid.UUID,
		role string,
		duration time.Duration,
	) (models.Invite, error)

	CreateByCityAdmin(
		ctx context.Context,
		userID, cityID, initiatorID uuid.UUID,
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
	admin  CityAdminSvc
	city   CitySvc
	invite inviteSvc
}

type Service struct {
	domain domain
	log    logium.Logger
}

func New(log logium.Logger, city CitySvc, cityMod CityAdminSvc, invSvc inviteSvc) Service {
	return Service{
		log: log,
		domain: domain{
			city:   city,
			admin:  cityMod,
			invite: invSvc,
		},
	}
}
