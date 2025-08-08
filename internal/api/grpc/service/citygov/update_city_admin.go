package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) UpdateCityAdmin(ctx context.Context, req *svc.UpdateCityAdminRequest) (*svc.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "initiator_id",
			Description: "invalid UUID format for initiator ID",
		})
	}

	role, err := enum.ParseCityAdminRole(req.Role)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "role",
			Description: err.Error(),
		})
	}

	err = s.methods.UpdateCityAdminRole(ctx, initiatorID, cityID, userID, role)
	if err != nil {
		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	cityAdmin, err := s.methods.GetCityAdminForCity(ctx, cityID, userID)
	if err != nil {
		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CityAdmin(cityAdmin), nil
}
