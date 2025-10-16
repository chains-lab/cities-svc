package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
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

	CreateCityAdmin(ctx context.Context, input models.CityAdmin) error
	GetCityAdmin(ctx context.Context, filters GetFilters) (models.CityAdmin, error)
	FilterCityAdmins(ctx context.Context, filter FilterParams, page, size uint64) (models.CityAdminsCollection, error)
	UpdateCityAdmin(ctx context.Context, userID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, inviteID, userID uuid.UUID, status string, now time.Time) error

	GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error)
	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
}

type UserGuesser interface {
	Guess(ctx context.Context, userIDs ...uuid.UUID) (map[uuid.UUID]models.Profile, error)
}

func (s Service) CityIsOfficialSupport(ctx context.Context, cityID uuid.UUID) error {
	ci, err := s.db.GetCityByID(ctx, cityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("get city: %w", err),
		)
	}
	if ci.IsNil() {
		return errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city not found"),
		)
	}
	if ci.Status != enum.CityStatusOfficial {
		return errx.ErrorCityIsNotSupported.Raise(
			fmt.Errorf("city not supported"),
		)
	}

	return nil
}
