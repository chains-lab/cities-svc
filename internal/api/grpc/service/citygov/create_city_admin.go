package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateCityAdmin(ctx context.Context, req *svc.CreateCityAdminRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, problems.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "user id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	role, err := enum.ParseCityAdminRole(req.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, problems.InvalidArgumentError(ctx, "city admin role is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "role",
			Description: err.Error(),
		})
	}

	cityAdmin, err := s.app.CreateCityAdmin(ctx, initiatorID, cityID, userID, app.CreateCityAdminInput{
		Role: role,
	})

	return responses.CityAdmin(cityAdmin), nil
}
