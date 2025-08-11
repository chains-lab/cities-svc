package citygov

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteCityOwner(ctx context.Context, req *svc.DeleteCityOwnerRequest) (*emptypb.Empty, error) {
	role, err := roles.ParseRole(req.Initiator.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid role in request")

		return nil, problems.UnauthenticatedError(ctx, "initiator role is invalid")
	}

	if role != roles.Admin && role != roles.SuperUser {
		logger.Log(ctx, RequestID(ctx)).Warnf("user %s with role %s tried to delete a city owner, but only admins and superusers can delete city owners",
			req.Initiator.Id, req.Initiator.Role)

		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"user %s with role %s is not allowed to delete a city owner", req.Initiator.Id, req.Initiator.Role),
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

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "user id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	err = s.app.DeleteCityOwner(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to delete city owner")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
