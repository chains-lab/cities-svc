package city

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityName(ctx context.Context, req *svc.UpdateCityNameRequest) (*svc.City, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, fmt.Sprint("city_id is invalid"), &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid initiator ID format")

		return nil, problem.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	city, err := s.app.UpdateCityName(ctx, cityID, initiatorID, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city name")
		return nil, err
	}

	return response.City(city), nil
}
