package admin

import (
	"context"
	"fmt"
	"time"

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
	GetCityTechLead(ctx context.Context, cityID uuid.UUID) (models.CityAdmin, error)

	FilterCityAdmins(ctx context.Context, filter FilterParams, page, size uint64) (models.CityAdminsCollection, error)
	UpdateCityAdmin(ctx context.Context, userID, cityID uuid.UUID, params UpdateParams, updateAt time.Time) error
	DeleteCityAdmin(ctx context.Context, userID, cityID uuid.UUID) error

	CreateInvite(ctx context.Context, input models.Invite) error
	GetInvite(ctx context.Context, ID uuid.UUID) (models.Invite, error)
	UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, status string) error

	GetCityAdmins(ctx context.Context, cityID uuid.UUID, roles ...string) (models.CityAdminsCollection, error)

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
}

type EventPublisher interface {
	PublishCityAdminCreated(
		ctx context.Context,
		admin models.CityAdmin,
		city models.City,
		recipients ...uuid.UUID,
	) error

	PublishCityAdminUpdated(
		ctx context.Context,
		admin models.CityAdmin,
		city models.City,
		recipients ...uuid.UUID,
	) error

	PublishCityAdminDeleted(
		ctx context.Context,
		admin models.CityAdmin,
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

func (s Service) validateInitiator(
	ctx context.Context,
	userID, cityID uuid.UUID,
	roles ...string,
) (models.CityAdmin, error) {
	admin, err := s.db.GetCityAdmin(ctx, userID, cityID)
	if err != nil {
		return models.CityAdmin{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get city admin by user ID and city ID, cause: %w", err),
		)
	}
	if admin.IsNil() {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("city admin for user ID %s not found", userID),
		)
	}

	if admin.CityID != cityID {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("city admin for user ID %s does not belong to city %s", userID, cityID),
		)
	}

	hasRole := false
	for _, r := range roles {
		if admin.Role == r {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return models.CityAdmin{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("city admin for user ID %s has not enough rights in city %s", userID, cityID),
		)
	}

	return admin, nil
}
