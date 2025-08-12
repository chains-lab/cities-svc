package gov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) TransferAdminRight(ctx context.Context, req *svc.TransferAdminRightRequest) (*emptypb.Empty, error) {
	initiator, err := s.OnlyCityAdmin(ctx, req.Initiator.UserId, req.CityId, "transfer city admin right")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.NewOwnerId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "user id is invalid format", &errdetails.BadRequest_FieldViolation{
			Field:       "new_owner_id",
			Description: "invalid UUID format for user ID",
		})
	}

	err = s.app.TransferCityAdminRight(ctx, initiator.CityID, initiator.ID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to refuse own admin rights")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
