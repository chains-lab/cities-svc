package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/middleware"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityStatusBySysAdmin(ctx context.Context, req *svc.UpdateCityStatusSysAdminRequest) (*svc.City, error) {
	_, err := middleware.AllowedRoles(ctx, req.Initiator, "update city status by system admin",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return nil, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	cityStatus, err := enum.ParseCityStatus(req.Status)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problem.InvalidArgumentError(ctx, "city status is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	city, err := s.app.UpdateCitiesStatusByOwner(ctx, initiatorID, cityID, cityStatus)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city status")

		return nil, err
	}

	return response.City(city), nil
}
