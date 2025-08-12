package formadmin

import (
	"context"
	"errors"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/formadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Service struct {
	app *app.App
	cfg config.Config

	svc.UnimplementedFormAdminServiceServer
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

		return models.CityGov{}, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	CityID, err := uuid.Parse(cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return models.CityGov{}, problem.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	gov, err := s.app.GetCityGov(ctx, InitiatorID, CityID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityAdminNotFound):
			logger.Log(ctx).Warnf("user: %s is not a city government for city %s, but try to do action: '%s'",
				InitiatorID, cityID, action)

			return models.CityGov{}, problem.PermissionDeniedError(ctx, "you are not a city gov")
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

		return models.CityGov{}, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	CityID, err := uuid.Parse(cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return models.CityGov{}, problem.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
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

		return models.CityGov{}, problem.PermissionDeniedError(ctx, "user is not a city admin")
	}

	return gov, nil
}
