package city

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) UpdateCityStatusBySysAdmin(ctx context.Context, req *svc.UpdateCityStatusSysAdminRequest) (*svc.City, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, problems.UnauthenticatedError(ctx, "initiator role is invalid format")
	}

	if role != roles.Admin && role != roles.SuperUser {
		logger.Log(ctx, RequestID(ctx)).Warnf("user %s with role %s tried to update a city, but only admins and superusers can update cities",
			req.Initiator.Id, req.Initiator.Role)

		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"user %s with role %s is not allowed to update a city", req.Initiator.Id, req.Initiator.Role),
		)
	}

	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, problems.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	cityStatus, err := enum.ParseCityStatus(req.Status)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, problems.InvalidArgumentError(ctx, "city status is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	city, err := s.app.UpdateCitiesStatusByOwner(ctx, initiatorID, cityID, cityStatus)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to update city status")

		return nil, err
	}

	return responses.City(city), nil
}
