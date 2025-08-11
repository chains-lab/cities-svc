package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityAdmin(ctx context.Context, req *svc.UpdateCityAdminRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid city id format", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid user id format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		return nil, problems.UnauthenticatedError(ctx, "initiator id is invalid format")
	}

	role, err := enum.ParseCityAdminRole(req.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, problems.InvalidArgumentError(ctx, "invalid city admin role", &errdetails.BadRequest_FieldViolation{
			Field:       "role",
			Description: err.Error(),
		})
	}

	err = s.app.UpdateCityAdminRole(ctx, initiatorID, cityID, userID, role)
	if err != nil {
		return nil, err
	}

	cityAdmin, err := s.app.GetCityAdmin(ctx, cityID, userID)
	if err != nil {
		return nil, err
	}

	return responses.CityAdmin(cityAdmin), nil
}
