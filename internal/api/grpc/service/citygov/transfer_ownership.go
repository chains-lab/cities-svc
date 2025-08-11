package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) TransferOwnership(ctx context.Context, req *svc.TransferOwnershipRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid city id format", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	newOwnerID, err := uuid.Parse(req.NewOwnerId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid new owner ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid city owner id", &errdetails.BadRequest_FieldViolation{
			Field:       "new_owner_id",
			Description: "invalid UUID format for new owner ID",
		})
	}

	err = s.app.TransferCityOwnership(ctx, RequestID(ctx), newOwnerID, cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to transfer city ownership")

		return nil, err
	}

	cityAdmin, err := s.app.GetCityAdmin(ctx, cityID, newOwnerID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get city admin after transfer")

		return nil, err
	}

	return response.CityAdmin(cityAdmin), nil
}
