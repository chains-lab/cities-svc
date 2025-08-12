package govadmin

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/govadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteCityGov(ctx context.Context, req *svc.DeleteCityGovRequest) (*emptypb.Empty, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "delete city government",
		roles.Admin, roles.SuperUser)
	if err != nil {
		return nil, err
	}

	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid city_id", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid user_id format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	err = s.app.DeleteCityGov(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to delete city admin")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
