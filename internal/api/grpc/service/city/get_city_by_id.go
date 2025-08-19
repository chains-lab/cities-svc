package city

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetCityById(ctx context.Context, req *svc.GetCityByIdRequest) (*svc.City, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problems.InvalidArgumentError(ctx, fmt.Sprint("city_id is invalid"), &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: "invalid UUID format for city ID",
		})
	}

	city, err := s.app.GetCityByID(ctx, cityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get city by ID")

		return nil, err
	}

	return responses.City(city), nil
}
