package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) UpdateCityStatus(ctx context.Context, req *svc.UpdateCityStatusRequest) (*svc.City, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "initiator_id",
			Description: "invalid UUID format for initiator ID",
		})
	}

	status, err := enum.ParseCityStatus(req.Status)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	city, err := s.app.UpdateCitiesStatusByOwner(ctx, initiatorID, cityID, status)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to update city status")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.City(city), nil
}
