package citiesadmins

import (
	"context"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

func (s Service) UpdateCityAdmin(ctx context.Context, req *svccitiesadmins.UpdateCityAdminRequest) (*svccitiesadmins.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "initiator_id",
			Description: "invalid UUID format for initiator ID",
		})
	}

	role, ok := enum.ParseCityAdminRole(req.Role)
	if !ok {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "role",
			Description: "invalid role value",
		})
	}

	err = s.methods.UpdateCityAdminRole(ctx, initiatorID, cityID, userID, role)
	if err != nil {
		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	cityAdmin, err := s.methods.GetCityAdminForCity(ctx, cityID, userID)
	if err != nil {
		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CityAdmin(cityAdmin), nil
}
