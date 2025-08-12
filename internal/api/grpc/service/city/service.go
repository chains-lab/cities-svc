package city

import (
	"context"
	"errors"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-proto/gen/go/common/userdata"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptor"
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
	GetCityByID(ctx context.Context, ID uuid.UUID) (models.City, error)
	SearchCityInCountry(ctx context.Context, like string, countryID uuid.UUID, request pagination.Request) ([]models.City, pagination.Response, error)

	UpdateCitiesStatus(ctx context.Context, cityID uuid.UUID, status string) (models.City, error)
	UpdateCityName(ctx context.Context, cityID uuid.UUID, name string) (models.City, error)

	GetCityGov(ctx context.Context, cityID, userID uuid.UUID) (models.CityGov, error)

	CreateForm(ctx context.Context, input app.CreateFormInput) (models.Form, error)
	AcceptForm(ctx context.Context, initiatorID, formID, adminID uuid.UUID) (models.Form, error)
	RejectForm(ctx context.Context, formID uuid.UUID, reason string) (models.Form, error)
	GetForm(ctx context.Context, formID uuid.UUID) (models.Form, error)
	SearchForms(ctx context.Context, input app.SearchFormsInput, pagPar pagination.Request, newFirst bool) ([]models.Form, pagination.Response, error)
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

func RequestID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	requestID, ok := ctx.Value(interceptor.RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return requestID
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
