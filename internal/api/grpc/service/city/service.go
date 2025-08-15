package city

import (
	"context"
	"errors"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
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
	CreateCity(ctx context.Context, input app.CreateCityInput) (models.City, error)

	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, request pagination.Request) ([]models.City, pagination.Response, error)

	UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)
	UpdateCityName(ctx context.Context, cityID uuid.UUID, name string) (models.City, error)

	GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedCityServiceServer
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

	if gov.Role != enum.CityAdminRoleAdmin {
		logger.Log(ctx).Warnf("user: %s is not a city admin for city %s, but try to do action: '%s'",
			InitiatorID, cityID, action)

		return models.CityGov{}, problems.PermissionDeniedError(ctx, "user is not a city admin")
	}

	return gov, nil
}
