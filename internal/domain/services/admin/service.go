package admin

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db          database
	userGuesser UserGuesser
	event       EventPublisher
}

func NewService(db database, userGuesser UserGuesser, eventPub EventPublisher) Service {
	return Service{
		db:          db,
		userGuesser: userGuesser,
		event:       eventPub,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateAdmin(ctx context.Context, input models.CityAdmin) error
	GetAdmin(ctx context.Context, filters GetFilters) (models.CityAdmin, error)
	UpdateAdmin(ctx context.Context, userID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	DeleteAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	FilterAdmins(ctx context.Context, filter FilterParams, page, size uint64) (models.CityAdminsCollection, error)

	CreateInvite(ctx context.Context, input models.Invite) (models.Invite, error)
	GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)
	AnswerToInvite(ctx context.Context, inviteID uuid.UUID, answer string) error

	ExistsGov(ctx context.Context, userID uuid.UUID) (bool, error)
	GetCity(ctx context.Context, cityID uuid.UUID) (models.City, error)
	GetCityStatus(ctx context.Context, cityID uuid.UUID) (models.CityStatus, error)
}

type UserGuesser interface {
	Guess(ctx context.Context, userIDs ...uuid.UUID) (map[uuid.UUID]models.Profile, error)
}

type EventPublisher interface {
	CityAdminCreated(ctx context.Context, admin models.CityAdmin) error
	CityAdminDeleted(ctx context.Context, userID, cityID uuid.UUID) error
}
