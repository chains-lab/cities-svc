package cities

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) UpdateCityName(ctx context.Context, req *svc.UpdateCityNameRequest) (*svc.City, error) {
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

	city, err := s.methods.UpdateCityName(ctx, cityID, initiatorID, req.Name)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to update city name")
		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.City(city), nil
}
