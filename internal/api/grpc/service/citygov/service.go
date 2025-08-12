package citygov

import (
	"context"
	"errors"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-proto/gen/go/common/userdata"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type application interface {
	CreateCityGov(ctx context.Context, cityID, userID uuid.UUID, input app.CreateCityGovInput) (models.CityGov, error)

	GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)
	GetCityGovs(ctx context.Context, cityID uuid.UUID, pag pagination.Request) ([]models.CityGov, pagination.Response, error)

	RefuseOwnCityGovRights(ctx context.Context, cityID, userID uuid.UUID) error
	TransferCityAdminRight(ctx context.Context, cityID, initiatorID, newOwnerID uuid.UUID) error

	DeleteCityGov(ctx context.Context, cityID, userID uuid.UUID) error
}

type Service struct {
	app application
	cfg config.Config

	svc.UnimplementedCityGovServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}

func (s Service) OnlyGov(ctx context.Context, req *userdata.UserData, cityID uuid.UUID, action string) (models.CityGov, error) {
	initiatorID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return models.CityGov{}, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	gov, err := s.app.GetCityGov(ctx, initiatorID, cityID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorCityAdminNotFound):
			logger.Log(ctx).Warnf("user: %s is not a city government for city %s, but try to do action: '%s'",
				initiatorID, cityID, action)

			return models.CityGov{}, problem.PermissionDeniedError(ctx, "you are not a city gov")
		default:
			logger.Log(ctx).WithError(err).Error("failed to get city gov")

			return models.CityGov{}, err
		}
	}

	return gov, nil
}

func (s Service) OnlyCityAdmin(ctx context.Context, req *userdata.UserData, cityID uuid.UUID, action string) (models.CityGov, error) {
	initiatorID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return models.CityGov{}, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	gov, err := s.app.GetCityGov(ctx, initiatorID, cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get city gov")

		return models.CityGov{}, err
	}

	if gov.Role != enum.CityAdminRoleAdmin {
		logger.Log(ctx).Warnf("user: %s is not a city admin for city %s, but try to do action: '%s'",
			initiatorID, cityID, action)

		return models.CityGov{}, problem.PermissionDeniedError(ctx, "user is not a city admin")
	}

	return gov, nil
}
