package invite

import (
	"context"
	"fmt"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db    database
	event EventPublisher
}

func NewService(db database, event EventPublisher) Service {
	return Service{
		db:    db,
		event: event,
	}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateCityAdmin(ctx context.Context, input models.CityAdmin) error

	GetCityAdminByUserAndCityID(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error)
	GetCityAdminByUserID(ctx context.Context, userID uuid.UUID) (models.CityAdmin, error)

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, inviteID, userID uuid.UUID, status string) error

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
}

type EventPublisher interface {
	PublishCityAdminCreated(ctx context.Context, admin models.CityAdmin) error
	PublishInviteCreated(ctx context.Context, invite models.Invite) error
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
