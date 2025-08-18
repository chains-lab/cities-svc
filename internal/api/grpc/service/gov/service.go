package gov

import (
	"context"
	"errors"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type application interface {
	CreateCityGovAdmin(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)
	GetCityAdmin(ctx context.Context, cityID uuid.UUID) (models.CityGov, error)
	DeleteCityAdmin(ctx context.Context, cityID uuid.UUID) error

	CreateCityGovModer(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)

	GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)
	GetCityGovs(ctx context.Context, cityID uuid.UUID, pag pagination.Request) ([]models.CityGov, pagination.Response, error)

	RefuseOwnCityGovRights(ctx context.Context, cityID, userID uuid.UUID) error
	TransferCityAdminRight(ctx context.Context, cityID, initiatorID, newOwnerID uuid.UUID) error

	DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedGovServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func (s Service) OnlyGov(ctx context.Context, initiatorID, cityID, action string) (models.CityGov, error) {
	InitiatorID, err := uuid.Parse(initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return models.CityGov{}, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	CityID, err := uuid.Parse(cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return models.CityGov{}, problems.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	gov, err := s.app.GetCityGov(ctx, InitiatorID, CityID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityGovNotFound):
			logger.Log(ctx).Warnf("user: %s is not a city government for city %s, but try to do action: '%s'",
				InitiatorID, cityID, action)

			return models.CityGov{}, problems.PermissionDeniedError(ctx, "you are not a city gov")
		default:
			logger.Log(ctx).WithError(err).Error("failed to get city gov")

			return models.CityGov{}, err
		}
	}

	return gov, nil
}

func (s Service) OnlyCityAdmin(ctx context.Context, initiatorID, cityID, action string) (models.CityGov, error) {
	InitiatorID, err := uuid.Parse(initiatorID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return models.CityGov{}, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	CityID, err := uuid.Parse(cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return models.CityGov{}, problems.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	gov, err := s.app.GetCityGov(ctx, InitiatorID, CityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get city gov")

		return models.CityGov{}, err
	}

	if gov.Role != enum.CityGovRoleAdmin {
		logger.Log(ctx).Warnf("user: %s is not a city admin for city %s, but try to do action: '%s'",
			InitiatorID, cityID, action)

		return models.CityGov{}, problems.PermissionDeniedError(ctx, "user is not a city admin")
	}

	return gov, nil
}
