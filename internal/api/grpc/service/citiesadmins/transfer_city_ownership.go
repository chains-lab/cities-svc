package citiesadmins

import (
	"context"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) TransferCityOwnership(ctx context.Context, req *svccitiesadmins.TransferOwnershipRequest) (*svccitiesadmins.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	newOwnerID, err := uuid.Parse(req.NewOwnerId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid new owner ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "new_owner_id",
			Description: "invalid UUID format for new owner ID",
		})
	}

	err = s.methods.TransferCityOwnership(ctx, RequestID(ctx), newOwnerID, cityID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to transfer city ownership")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	cityAdmin, err := s.methods.GetCityAdminForCity(ctx, cityID, newOwnerID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get city admin after transfer")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CityAdmin(cityAdmin), nil
}
