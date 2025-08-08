package citiesadmins

import (
	"context"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteCityAdmin(ctx context.Context, req *svccitiesadmins.DeleteCityAdminRequest) (*emptypb.Empty, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
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

	err = s.methods.DeleteCityAdmin(ctx, initiatorID, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to delete city admin")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return &emptypb.Empty{}, nil
}
