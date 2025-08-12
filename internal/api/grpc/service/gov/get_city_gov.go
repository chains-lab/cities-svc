package gov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetCityGov(ctx context.Context, req *svc.GetCityGovRequest) (*svc.CityGov, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "user id is invalid format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	cityAdmin, err := s.app.GetCityGov(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get city admin")

		return nil, err
	}

	return response.CityAdmin(cityAdmin), nil
}
