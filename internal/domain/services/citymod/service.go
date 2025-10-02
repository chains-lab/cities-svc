package citymod

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

type JwtManager interface {
	CreateInviteToken(
		inviteID uuid.UUID,
		role string,
		cityID uuid.UUID,
		ExpiredAt time.Time,
	) (models.InviteToken, error)

	DecryptInviteToken(tokenStr models.InviteToken) (models.InviteTokenData, error)
}

type Service struct {
	db  database
	jwt JwtManager
}

func NewService(db database, jwt JwtManager) Service {
	return Service{
		db:  db,
		jwt: jwt,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateCityMod(ctx context.Context, input models.CityModer) error
	GetCityModer(ctx context.Context, filters GetFilters) (models.CityModer, error)
	FilterCityModers(ctx context.Context, filter FilterParams, page, size uint64) (models.CityModersCollection, error)
	UpdateCityModer(ctx context.Context, userID uuid.UUID, params UpdateCityModerParams, updatedAt time.Time) error
	DeleteCityModer(ctx context.Context, userID, cityID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateStatusInvite(ctx context.Context, inviteID, userID uuid.UUID, status string, now time.Time) error

	GetCountryByID(ctx context.Context, ID uuid.UUID) (models.Country, error)
	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
}

func (s Service) CityIsOfficialSupport(ctx context.Context, cityID uuid.UUID) error {
	ci, err := s.db.GetCityByID(ctx, cityID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("get city: %w", err),
		)
	}
	if ci == (models.City{}) {
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
