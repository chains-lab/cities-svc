package cities

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/cities-dir-svc/internal/errs"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) UpdateCityStatusBySysAdmin(ctx context.Context, req *svc.UpdateCityStatusSysAdminRequest) (*svc.City, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, renderer.AppError(ctx, RequestID(ctx), errs.RaiseInternal(err))
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

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "initiator_id",
			Description: "invalid UUID format for initiator ID",
		})
	}

	status, ok := enum.ParseCityStatus(req.Status)
	if !ok {
		logger.Log(ctx, RequestID(ctx)).Error("invalid city status provided")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "status",
			Description: "invalid city status provided",
		})
	}

	city, err := s.methods.UpdateCitiesStatusByOwner(ctx, initiatorID, cityID, status)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to update city status")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.City(city), nil
}
