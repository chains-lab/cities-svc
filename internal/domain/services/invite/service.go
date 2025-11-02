package invite

import (
	"context"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db          database
	userGuesser UserGuesser
}

func NewService(db database, userGuesser UserGuesser) Service {
	return Service{
		db:          db,
		userGuesser: userGuesser,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateAdmin(ctx context.Context, input models.CityAdmin) error

	CreateInvite(ctx context.Context, input models.Invite) (models.Invite, error)
	GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)
	AnswerToInvite(ctx context.Context, inviteID uuid.UUID, answer string) error

	ExistsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	GetCity(ctx context.Context, cityID uuid.UUID) (models.City, error)
	GetCityStatus(ctx context.Context, cityID uuid.UUID) (models.CityStatus, error)
}

type UserGuesser interface {
	Guess(ctx context.Context, userIDs ...uuid.UUID) (map[uuid.UUID]models.Profile, error)
}
