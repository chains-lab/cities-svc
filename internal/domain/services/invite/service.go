package invite

import (
	"context"
	"fmt"

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
	GetCityAdmin(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error)
	GetCityAdmins(ctx context.Context, cityID uuid.UUID, roles ...string) (models.CityAdminsCollection, error)
	GetCityTechLead(ctx context.Context, cityID uuid.UUID) (models.CityAdmin, error)
	DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, status string) error

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
}

type EventPublisher interface {
	PublishInviteCreated(
		ctx context.Context,
		invite models.Invite,
		city models.City,
		recipients ...uuid.UUID,
	) error

	PublishInviteAccepted(
		ctx context.Context,
		invite models.Invite,
		city models.City,
		cityAdmin models.CityAdmin,
		recipients ...uuid.UUID,
	) error

	PublishInviteDeclined(
		ctx context.Context,
		invite models.Invite,
		city models.City,
		recipients ...uuid.UUID,
	) error

	PublishCityAdminCreated(
		ctx context.Context,
		cityAdmin models.CityAdmin,
		city models.City,
		recipients ...uuid.UUID,
	) error
}

func (s Service) getCity(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	ci, err := s.db.GetCityByID(ctx, cityID)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get city: %w", err),
		)
	}
	if ci.IsNil() {
		return models.City{}, errx.ErrorCityNotFound.Raise(
			fmt.Errorf("city not found"),
		)
	}

	return ci, nil
}

func (s Service) getInitiator(ctx context.Context, userID, cityID uuid.UUID) (models.CityAdmin, error) {
	res, err := s.db.GetCityAdmin(ctx, userID, cityID)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin, cause: %w", err),
		)
	}

	if res.IsNil() {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("city admin not found"),
		)
	}

	return res, nil
}
