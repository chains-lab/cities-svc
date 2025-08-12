package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateCityGov(ctx context.Context, req *svc.CreateCityGovRequest) (*svc.CityGov, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	_, err = s.OnlyCityAdmin(ctx, req.Initiator, cityID, "create city government")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problem.InvalidArgumentError(ctx, "user id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	role, err := enum.ParseCityAdminRole(req.Role)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problem.InvalidArgumentError(ctx, "city admin role is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "role",
			Description: err.Error(),
		})
	}

	cityAdmin, err := s.app.CreateCityGov(ctx, cityID, userID, app.CreateCityGovInput{
		Role: role,
	})

	return response.CityAdmin(cityAdmin), nil
}
