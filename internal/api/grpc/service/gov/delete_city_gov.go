package gov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteCityGov(ctx context.Context, req *svc.DeleteCityGovRequest) (*emptypb.Empty, error) {
	user := meta.User(ctx)

	gov, err := s.OnlyCityAdmin(ctx, user.ID.String(), req.CityId, "delete city government")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "invalid user_id format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	err = s.app.DeleteCityGov(ctx, gov.CityID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to delete city admin")

		return nil, err
	}

	logger.Log(ctx).Infof("city government deleted, city ID: %s, user ID: %s", req.CityId, req.UserId)

	return &emptypb.Empty{}, nil
}
