package cities

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetCityById(ctx context.Context, req *svc.GetCityByIdRequest) (*svc.City, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "id",
			Description: "invalid UUID format for city ID",
		})
	}

	city, err := s.methods.GetCityByID(ctx, cityID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get city by ID")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.City(city), nil
}
