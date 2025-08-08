package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) TransferCityOwnership(ctx context.Context, req *svc.TransferOwnershipRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	newOwnerID, err := uuid.Parse(req.NewOwnerId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid new owner ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "new_owner_id",
			Description: "invalid UUID format for new owner ID",
		})
	}

	err = s.methods.TransferCityOwnership(ctx, RequestID(ctx), newOwnerID, cityID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to transfer city ownership")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	cityAdmin, err := s.methods.GetCityAdminForCity(ctx, cityID, newOwnerID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get city admin after transfer")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CityAdmin(cityAdmin), nil
}
