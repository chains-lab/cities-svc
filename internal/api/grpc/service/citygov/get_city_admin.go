package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetCityAdmin(ctx context.Context, req *svc.GetCityAdminRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	cityAdmin, err := s.methods.GetCityAdmin(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get city admin")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CityAdmin(cityAdmin), nil
}
