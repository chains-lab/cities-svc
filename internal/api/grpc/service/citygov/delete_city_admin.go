package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteCityAdmin(ctx context.Context, req *svc.DeleteCityAdminRequest) (*emptypb.Empty, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid city_id", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid user_id format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	err = s.app.DeleteCityAdmin(ctx, initiatorID, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to delete city admin")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
