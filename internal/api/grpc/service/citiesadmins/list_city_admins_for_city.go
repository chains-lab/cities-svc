package citiesadmins

import (
	"context"

	svccitiesadmins "github.com/chains-lab/cities-dir-proto/gen/go/citiesadmins"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) ListCityAdminsForCity(ctx context.Context, req *svccitiesadmins.ListCityAdminsRequest) (*svccitiesadmins.ListCitiesAdmins, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid city ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	cityAdmins, err := s.methods.GetCityAdmins(ctx, cityID, req.Pagination.Limit, req.Pagination.Page)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to list city admins")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CitiesAdminsList(cityAdmins), nil
}
