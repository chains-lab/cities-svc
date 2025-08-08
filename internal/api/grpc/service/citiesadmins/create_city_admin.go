package citiesadmins

import (
	"context"
	"fmt"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) CreateCityAdmin(ctx context.Context, req *svccitiesadmins.CreateCityAdminRequest) (*svccitiesadmins.CityAdmin, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiatorID, err := uuid.Parse(req.Initiator.Id)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid initiator ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "initiator_id",
			Description: "invalid UUID format for initiator ID",
		})
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid user ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	role, ok := enum.ParseCityAdminRole(req.Role)
	if !ok {
		logger.Log(ctx, RequestID(ctx)).Errorf("invalid city admin role: %s", req.Role)

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "role",
			Description: fmt.Sprintf("invalid city admin role: %s", req.Role),
		})
	}

	cityAdmin, err := s.methods.CreateCityAdmin(ctx, initiatorID, cityID, userID, app.CreateCityAdminInput{
		Role: role,
	})

	return renderer.CityAdmin(cityAdmin), nil
}
