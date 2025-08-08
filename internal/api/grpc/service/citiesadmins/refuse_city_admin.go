package citiesadmins

import (
	"context"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) RefuseCityAdminRights(ctx context.Context, req *svccitiesadmins.RefuseCityAdminRightRequest) (*svccitiesadmins.CityAdmin, error) {
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

	err = s.methods.RefuseOwnAdminRights(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to refuse own admin rights")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	cityAdmin, err := s.methods.GetCityAdminForCity(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get city admin after refusal")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CityAdmin(cityAdmin), nil
}
