package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type JwtManager interface {
	CreateInviteToken(
		inviteID uuid.UUID,
		role string,
		cityID uuid.UUID,
		ExpiredAt time.Time,
	) (string, error)

	DecryptInviteToken(tokenStr string) (models.InviteTokenData, error)

	HashInviteToken(tokenStr string) (string, error)
	VerifyInviteToken(tokenStr, hashed string) error
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

	CreateCityAdmin(ctx context.Context, input models.CityAdmin) error
	DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	GetCityAdminByUserAndCityID(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error)
	GetCityAdminByUserID(ctx context.Context, userID uuid.UUID) (models.CityAdmin, error)

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, inviteID, userID uuid.UUID, status string, now time.Time) error

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
